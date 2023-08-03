package m4e

import (
	"context"

	"github.com/imdario/mergo"
	m4ev1alpha1 "github.com/krestomatio/kio-operator/apis/m4e/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"
)

const (
	OPERATORNAME string = "kio-operator"
)

// ReconcileCreate create resource if it does not exists. Otherwise it does nothing
func (r *SiteReconciler) ReconcileCreate(ctx context.Context, parentObj client.Object, obj client.Object) error {
	log := log.FromContext(ctx)

	log.V(1).Info("Creating resource", "Resource", obj.GetObjectKind())

	if err := r.Get(ctx, types.NamespacedName{Name: obj.GetName(), Namespace: obj.GetNamespace()}, obj); !errors.IsNotFound(err) {
		log.V(1).Info("Resource already exists", "Resource", obj.GetObjectKind())
		return nil
	} else if client.IgnoreNotFound(err) != nil {
		log.Error(err, "Failed to get resource", "Resource", obj.GetObjectKind())
		return err
	}

	// Set resource ownership
	if err := r.ReconcileSetOwner(ctx, parentObj, obj); err != nil {
		log.Error(err, "Failed to set owner", "Resource", obj.GetObjectKind())
		return err
	}

	if err := r.Create(ctx, obj); err != nil {
		log.Error(err, "Failed to create resource", "Resource", obj.GetObjectKind())
		return err
	}

	log.Info("Resource created", "Resource", obj.GetObjectKind())
	return nil
}

func (r *SiteReconciler) ReconcileApply(ctx context.Context, parentObj client.Object, obj client.Object) error {
	log := log.FromContext(ctx)

	log.V(1).Info("Applying patch", "Resource", obj.GetObjectKind())

	// Set resource ownership
	if err := r.ReconcileSetOwner(ctx, parentObj, obj); err != nil {
		log.Error(err, "Failed setting owner", "Resource", obj.GetObjectKind())
		return err
	}

	// Apply resource
	force := true
	if err := r.Patch(ctx, obj, client.Apply, &client.PatchOptions{Force: &force, FieldManager: OPERATORNAME}); err != nil {
		log.Error(err, "Failed to attempt patching changes", "Resource", obj.GetObjectKind())
		return err
	}

	return nil
}

func (r *SiteReconciler) ReconcileSetOwner(ctx context.Context, parentObj client.Object, obj client.Object) error {
	log := log.FromContext(ctx)

	// Set owner reference
	log.V(1).Info("Setting owner", "Owner", parentObj.GetUID(), "Dependant", obj.GetObjectKind())
	if err := ctrl.SetControllerReference(parentObj, obj, r.Scheme); err != nil {
		log.Error(err, "Failed to set owner reference")
		return err
	}

	return nil
}

// ReconcileDeleteDependant deletes a resource only if it has the owner reference of its parent
func (r *SiteReconciler) ReconcileDeleteDependant(ctx context.Context, parentObj client.Object, obj client.Object) error {
	log := log.FromContext(ctx)

	log.V(1).Info("Deleting dependant resource", "Resource", obj.GetObjectKind())

	if err := r.Get(ctx, types.NamespacedName{Name: obj.GetName(), Namespace: obj.GetNamespace()}, obj); err != nil {
		log.V(1).Info(err.Error(), "Dependant", obj.GetObjectKind())
		return err
	}

	if isObjMarkedToBeDeleted := obj.GetDeletionTimestamp(); isObjMarkedToBeDeleted != nil {
		log.V(1).Info("Dependant resource marked to be deleted", "Dependant", obj.GetObjectKind())
		return nil
	}

	// Check ownership with parent Object
	objOwner := metav1.GetControllerOf(obj)

	if objOwner == nil {
		log.Info("Dependant resource not deleted. It has no owner", "Dependant", obj.GetObjectKind())
		return nil
	} else if objOwner.UID != parentObj.GetUID() {
		log.Info("Dependant resource not deleted. Its owner does not match parent ownership (uid)", "Dependant", obj.GetObjectKind(), "Parent", parentObj.GetObjectKind())
		return nil
	}

	gracePeriodSeconds := int64(0)
	propagationPolicy := metav1.DeletePropagationBackground
	if err := r.Delete(ctx, obj, client.PropagationPolicy(propagationPolicy), client.GracePeriodSeconds(gracePeriodSeconds)); err != nil {
		log.Error(err, "Failed to delete dependant resource ", "Dependant", obj.GetObjectKind())
		return err
	}

	log.Info("Dependant resource set to be deleted", "Dependant", obj.GetObjectKind())
	return nil
}

// finalizeSite cleans up before deleting Site
func (r *SiteReconciler) finalizeSite(ctx context.Context) (requeue bool, err error) {
	log := log.FromContext(ctx)
	log.Info("Finalizing")

	// Delete moodle and inmediately requeue in order to wait for it to be completely be removed.
	// By doing so, any dependant CR removal will be done after, and removal
	// conflicts will be avoided
	log.Info("Deleting Moodle", "Moodle.Namespace", r.siteCtx.moodle.GetNamespace(), "Moodle.Name", r.siteCtx.moodle.GetName())
	if err := r.ReconcileDeleteDependant(ctx, r.siteCtx.site, r.siteCtx.moodle); err == nil {
		log.V(1).Info("Set for requeue after Moodle deletion", "Moodle.Namespace", r.siteCtx.moodle.GetNamespace(), "Moodle.Name", r.siteCtx.moodle.GetName())
		return true, nil
	} else if !errors.IsNotFound(err) {
		log.Error(err, "Moodle not deleted", "Moodle.Namespace", r.siteCtx.moodle.GetNamespace(), "Moodle.Name", r.siteCtx.moodle.GetName())
		return false, err
	}

	// Delete Keydb and set for later requeuing in order to wait for it to be completely be removed.
	if r.siteCtx.hasKeydb {
		log.Info("Deleting Keydb", "Keydb.Namespace", r.siteCtx.keydb.GetNamespace(), "Keydb.Name", r.siteCtx.keydb.GetName())
		if err := r.ReconcileDeleteDependant(ctx, r.siteCtx.site, r.siteCtx.keydb); err == nil {
			log.V(1).Info("Set for requeue after Keydb deletion", "Keydb.Namespace", r.siteCtx.keydb.GetNamespace(), "Keydb.Name", r.siteCtx.keydb.GetName())
			requeue = true
		} else if !errors.IsNotFound(err) {
			log.Error(err, "Keydb not deleted", "Keydb.Namespace", r.siteCtx.keydb.GetNamespace(), "Keydb.Name", r.siteCtx.keydb.GetName())
			return false, err
		}
	}

	// Delete Postgres and set for later requeuing in order to wait for it to be completely be removed.
	if r.siteCtx.hasPostgres {
		log.Info("Deleting Postgres", "Postgres.Namespace", r.siteCtx.postgres.GetNamespace(), "Postgres.Name", r.siteCtx.postgres.GetName())
		if err := r.ReconcileDeleteDependant(ctx, r.siteCtx.site, r.siteCtx.postgres); err == nil {
			log.V(1).Info("Set for requeue after Postgres deletion", "Postgres.Namespace", r.siteCtx.postgres.GetNamespace(), "Postgres.Name", r.siteCtx.postgres.GetName())
			requeue = true
		} else if !errors.IsNotFound(err) {
			log.Error(err, "Postgres not deleted", "Postgres.Namespace", r.siteCtx.postgres.GetNamespace(), "Postgres.Name", r.siteCtx.postgres.GetName())
			return false, err
		}
	}

	// Delete nfs ganesha server and set for later requeuing in order to wait for it to be completely be removed.
	if r.siteCtx.hasNfs {
		log.Info("Deleting NFS Ganesha", "Ganesha.Namespace", r.siteCtx.nfs.GetNamespace(), "Ganesha.Name", r.siteCtx.nfs.GetName())
		if err := r.ReconcileDeleteDependant(ctx, r.siteCtx.site, r.siteCtx.nfs); err == nil {
			log.V(1).Info("Set for requeue after NFS Ganesha server deletion", "Ganesha.Namespace", r.siteCtx.nfs.GetNamespace(), "Ganesha.Name", r.siteCtx.nfs.GetName())
			requeue = true
		} else if !errors.IsNotFound(err) {
			log.Error(err, "NFS Ganesha server not deleted", "Ganesha.Namespace", r.siteCtx.nfs.GetNamespace(), "Ganesha.Name", r.siteCtx.nfs.GetName())
			return false, err
		}
	}

	if !requeue {
		log.Info("Successfully finalized site")
	}

	return requeue, nil
}

// finalizeSite cleans up before deleting Flavor
func (r *FlavorReconciler) finalizeFlavor(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Finalizing")

	// Whether any Site is using this flavor
	log.Info("Deleting Flavor", "Flavor.Namespace", r.flavorCtx.flavor.GetNamespace(), "Flavor.Name", r.flavorCtx.flavor.GetName())
	siteList := &m4ev1alpha1.SiteList{}
	if err := r.List(ctx, siteList, client.MatchingFields{"spec.flavor": r.flavorCtx.flavor.GetName()}); err != nil {
		log.Error(err, "Unable to list child sites")
		return err
	}

	sitesUsingFlavor := len(siteList.Items)
	if sitesUsingFlavor > 0 {
		flavorNotFoundError := &FlavorInUsedError{r.flavorCtx.flavor.GetName(), sitesUsingFlavor}
		log.Error(flavorNotFoundError, "Cannot delete flavor")
		return flavorNotFoundError
	}

	log.Info("Successfully finalized flavor")
	return nil
}

// updateSiteState update site state
// return any error
func (r *SiteReconciler) updateSiteState(ctx context.Context) error {
	log := log.FromContext(ctx)

	var state string

	// get ready condition
	readyCondition, readyConditionFound, readyConditionErr := getConditionByType(r.siteCtx.site, ReadyConditionType)

	if readyConditionErr != nil {
		log.Error(readyConditionErr, "unable to update Site '"+r.siteCtx.site.GetName()+"' state")
		return readyConditionErr
	}

	// get Moodle ready condition
	moodleReadyCondition, moodleReadyConditionFound, moodleReadyConditionErr := getConditionByType(r.siteCtx.site, MoodleReadyConditionType)

	if moodleReadyConditionErr != nil {
		log.Error(moodleReadyConditionErr, "unable to update Site '"+r.siteCtx.site.GetName()+"' state")
		return moodleReadyConditionErr
	}

	if readyConditionFound && moodleReadyConditionFound {
		state = r.setSiteState(readyCondition, moodleReadyCondition)
	} else {
		state = string(m4ev1alpha1.SettingUpState)
	}

	// set state in site object
	stateUpdate, err := SetStatusState(r.siteCtx.site, state)
	if err != nil {
		log.Error(err, "unable to update Site '"+r.siteCtx.site.GetName()+"' state")
		return err
	}

	if !stateUpdate {
		log.V(1).Info("Site state not updated")
		return nil
	}

	// save status
	if err := r.Status().Update(ctx, r.siteCtx.site); err != nil {
		log.Error(err, "Unable to update Site '"+r.siteCtx.name+"' state")
		return err
	}

	log.V(1).Info("Site state updated")
	return nil
}

// updateFlavorState update flavor state
// return any error
func (r *FlavorReconciler) updateFlavorState(ctx context.Context) error {
	log := log.FromContext(ctx)

	state := r.setFlavorState()

	// set state in site object
	stateUpdate, err := SetStatusState(r.flavorCtx.flavor, state)
	if err != nil {
		log.Error(err, "unable to update Flavor '"+r.flavorCtx.flavor.GetName()+"' state")
		return err
	}

	if !stateUpdate {
		log.V(1).Info("Flavor state not updated")
		return nil
	}

	// save status
	if err := r.Status().Update(ctx, r.flavorCtx.flavor); err != nil {
		log.Error(err, "Unable to update Flavor '"+r.flavorCtx.flavor.GetName()+"' state")
		return err
	}

	log.V(1).Info("Flavor state updated")
	return nil
}

// setSiteState defines Site state value from ready condition
// return state string
func (r *SiteReconciler) setSiteState(readyCondition map[string]interface{}, moodleReadyCondition map[string]interface{}) string {
	status := readyCondition["status"]
	moodleStatus := moodleReadyCondition["status"]
	moodleReason := moodleReadyCondition["reason"]

	// Terminating
	if r.siteCtx.markedToBeDeleted {
		return string(m4ev1alpha1.TerminatingState)
	}

	if status == "False" || moodleStatus == "False" {
		// Failed
		if moodleReason == "Error" {
			return string(m4ev1alpha1.FailedState)
		}
		// Creating
		if moodleReason == "NotInstantiated" || moodleReason == "Instantiated" || moodleReason == "NotCreated" {
			return string(m4ev1alpha1.CreatingState)
		}
	}
	// Ready
	if moodleStatus == "True" {
		return string(m4ev1alpha1.ReadyState)
	}

	return string(m4ev1alpha1.UnknownState)
}

// setNotifyUUID defines site uuid if notifying status to an endpoint
// Should be used once combinedMoodleSpec is set
// By default, site name is used as UUID
func (r *SiteReconciler) setNotifyUUID() error {
	// whether it has to notify status to a url
	_, moodleSiteRoutineStatusCrNotifyFound, _ := unstructured.NestedMap(r.siteCtx.combinedMoodleSpec, "routineStatusCrNotify")
	if moodleSiteRoutineStatusCrNotifyFound {
		_, moodleSiteRoutineStatusCrNotifyUuidFound, _ := unstructured.NestedMap(r.siteCtx.combinedMoodleSpec, "routineStatusCrNotify", "uuid")
		if !moodleSiteRoutineStatusCrNotifyUuidFound {
			// set uuid to notify about
			if err := unstructured.SetNestedField(r.siteCtx.combinedMoodleSpec, r.siteCtx.name, "routineStatusCrNotify", "uuid"); err != nil {
				return err
			}
		}
	}
	return nil
}

// setFlavorState defines Flavor state value
// return state string
func (r *FlavorReconciler) setFlavorState() string {
	if r.flavorCtx.markedToBeDeleted {
		return string(m4ev1alpha1.TerminatingState)
	}

	return string(m4ev1alpha1.ReadyState)
}

// SetStatusState set status state key in unstructure object
// It returns a bool flag if state was updated, and
// any error
func SetStatusState(objU *unstructured.Unstructured, state string) (bool, error) {
	updateState := false
	objState, objStateFound, objStateErr := unstructured.NestedString(objU.UnstructuredContent(), "status", "state")

	if objStateErr != nil {
		return false, objStateErr
	}

	if !objStateFound || objState != state {
		updateState = true
		if err := unstructured.SetNestedField(objU.Object, state, "status", "state"); err != nil {
			return false, err
		}
	}

	return updateState, nil
}

// Init a new unstructured object with determined GVK
func newUnstructuredObject(gvk schema.GroupVersionKind) *unstructured.Unstructured {
	objU := &unstructured.Unstructured{}
	objU.SetGroupVersionKind(gvk)
	return objU
}

// truncate a string
func truncate(str string, length int) (truncated string) {
	if length <= 0 {
		return
	}
	for i, char := range str {
		if i >= length {
			break
		}
		truncated += string(char)
	}
	return
}

// setSiteLabels set site base labels
func (r *SiteReconciler) setSiteLabels(ctx context.Context) error {
	log := log.FromContext(ctx)
	siteLabels := make(map[string]string)

	// use flavor labels
	for key, value := range r.siteCtx.flavor.GetLabels() {
		siteLabels[key] = value
	}

	// set base labels
	siteLabels[m4ev1alpha1.GroupVersion.Group+"/site-name"] = r.siteCtx.name
	siteLabels[m4ev1alpha1.GroupVersion.Group+"/flavor-name"] = r.siteCtx.flavorName
	siteLabels[m4ev1alpha1.GroupVersion.Group+"/meta-operator-name"] = OPERATORNAME

	r.siteCtx.site.SetLabels(siteLabels)

	if err := r.Patch(ctx, r.siteCtx.site, client.Merge); err != nil {
		log.Error(err, "Failed to attempt patching site labels", "Site", r.siteCtx.site.GetName())
		return err
	}

	return nil
}

// commonLabels set common labels
func (r *SiteReconciler) commonLabels(objSpec map[string]interface{}) (err error) {
	siteLabelsBytes, _ := yaml.Marshal(r.siteCtx.site.GetLabels())
	siteLabelsString := string(siteLabelsBytes)

	objSpecCommonLabelsString, objSpecCommonLabelsFound, _ := unstructured.NestedString(objSpec, "commonLabels")
	if objSpecCommonLabelsFound {
		objSpec["commonLabels"] = siteLabelsString + "\n" + objSpecCommonLabelsString
	} else {
		objSpec["commonLabels"] = siteLabelsString
	}

	return err
}

// DefaultAffinity set the default affinity for a site
func (r *SiteReconciler) defaultAffinityYaml(objSpec map[string]interface{}, fieldName string) (err error) {
	var defaultAffinityYamlBytes []byte

	defaultAffinity := corev1.Affinity{
		PodAffinity: &corev1.PodAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
				{
					Weight: int32(100),
					PodAffinityTerm: corev1.PodAffinityTerm{
						TopologyKey: "kubernetes.io/hostname",
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      m4ev1alpha1.GroupVersion.Group + "/site-name",
									Operator: metav1.LabelSelectorOpIn,
									Values: []string{
										r.siteCtx.name,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	defaultAffinityYamlBytes, err = yaml.Marshal(defaultAffinity)

	if objSpecDefaultAffinityString, objDefaultAffinityFound, err := unstructured.NestedString(objSpec, fieldName); err != nil {
		return err
	} else if objDefaultAffinityFound {
		objSpec[fieldName] = string(defaultAffinityYamlBytes) + "\n" + objSpecDefaultAffinityString
	} else {
		objSpec[fieldName] = string(defaultAffinityYamlBytes)
	}

	return err
}

// Merge value in nested string present in both objects into the first object
// string + '\n' + string
func (r *SiteReconciler) mergeNestedString(firstObjSpec map[string]interface{}, secondObjSpec map[string]interface{}, fields ...string) (err error) {

	firstObjSpecNestedField, firstObjSpecNestedFieldFound, err := unstructured.NestedString(firstObjSpec, fields...)
	if err != nil {
		return err
	}

	secondObjSpecNestedField, secondObjSpecFieldNestedFound, err := unstructured.NestedString(secondObjSpec, fields...)
	if err != nil {
		return err
	}

	if firstObjSpecNestedFieldFound && secondObjSpecFieldNestedFound {
		mergeNestedField := firstObjSpecNestedField + "\n" + secondObjSpecNestedField
		if err := unstructured.SetNestedField(firstObjSpec, mergeNestedField, fields...); err != nil {
			return err
		}
	}

	return err
}

// moodleDefaultAffinityYaml set the default affinity for Moodle
func (r *SiteReconciler) moodleDefaultAffinityYaml() (err error) {
	if err = r.defaultAffinityYaml(r.siteCtx.flavorMoodleSpec, "moodleCronjobAffinity"); err != nil {
		return err
	}
	if err = r.defaultAffinityYaml(r.siteCtx.flavorMoodleSpec, "moodleUpdateJobAffinity"); err != nil {
		return err
	}
	if err = r.defaultAffinityYaml(r.siteCtx.flavorMoodleSpec, "moodleNewInstanceJobAffinity"); err != nil {
		return err
	}
	if err = r.defaultAffinityYaml(r.siteCtx.flavorMoodleSpec, "phpFpmAffinity"); err != nil {
		return err
	}
	if err = r.defaultAffinityYaml(r.siteCtx.flavorMoodleSpec, "nginxAffinity"); err != nil {
		return err
	}
	return err
}

// postgresSpec handle any postgres spec
func (r *SiteReconciler) postgresSpec() (err error) {
	r.siteCtx.hasPostgres = r.siteCtx.postgresSpecFound || r.siteCtx.flavorPostgresSpecFound

	if r.siteCtx.hasPostgres {
		r.siteCtx.postgres.SetName(r.siteCtx.postgresName)
		r.siteCtx.postgres.SetNamespace(r.siteCtx.namespaceName)
	}

	// Postgres kind from Postgres ansible operator
	if r.siteCtx.hasPostgres {
		// Set Postgres host and secret, if not already present in Moodle spec
		postgresRelatedMoodleSpec := map[string]interface{}{
			"moodlePostgresMetaName": r.siteCtx.postgresName,
		}
		// Merge Moodle related postgres spec with flavor Moodle spec
		if err := mergo.MapWithOverwrite(&r.siteCtx.flavorMoodleSpec, postgresRelatedMoodleSpec); err != nil {
			return err
		}
		// Merge Postgres spec if set on site Spec
		if r.siteCtx.postgresSpecFound {
			if err := mergo.MapWithOverwrite(&r.siteCtx.flavorPostgresSpec, r.siteCtx.postgresSpec); err != nil {
				return err
			}
		}
		// Set site labels to postgres
		if err := r.commonLabels(r.siteCtx.flavorPostgresSpec); err != nil {
			return err
		}
		// set default affinity
		if err := r.defaultAffinityYaml(r.siteCtx.flavorPostgresSpec, "postgresAffinity"); err != nil {
			return err
		}
		// save postgres spec
		r.siteCtx.combinedPostgresSpec = make(map[string]interface{})
		r.siteCtx.combinedPostgresSpec = r.siteCtx.flavorPostgresSpec
	}

	return err
}

// nfsSpec handle any nfs spec
func (r *SiteReconciler) nfsSpec() (err error) {
	r.siteCtx.hasNfs = r.siteCtx.nfsSpecFound || r.siteCtx.flavorNfsSpecFound

	if r.siteCtx.hasNfs {
		r.siteCtx.nfs.SetName(r.siteCtx.nfsName)
		r.siteCtx.nfs.SetNamespace(r.siteCtx.namespaceName)
	}

	// Ganesha server kind from NFS ansible operator
	if r.siteCtx.hasNfs {
		// Set NFS storage class name and access modes when using NFS operator
		nfsRelatedMoodleSpec := map[string]interface{}{
			"moodleNfsMetaName": r.siteCtx.nfsName,
		}
		// Merge Moodle related nfs spec with flavor Moodle spec
		if err := mergo.MapWithOverwrite(&r.siteCtx.flavorMoodleSpec, nfsRelatedMoodleSpec); err != nil {
			return err
		}
		// Merge NFS spec if set on site Spec
		if r.siteCtx.nfsSpecFound {
			if err := mergo.MapWithOverwrite(&r.siteCtx.flavorNfsSpec, r.siteCtx.nfsSpec); err != nil {
				return err
			}
		}
		// Set site labels to nfs
		if err := r.commonLabels(r.siteCtx.flavorNfsSpec); err != nil {
			return err
		}
		// set default affinity
		if err := r.defaultAffinityYaml(r.siteCtx.flavorNfsSpec, "ganeshaAffinity"); err != nil {
			return err
		}
		// save nfs spec
		r.siteCtx.combinedNfsSpec = make(map[string]interface{})
		r.siteCtx.combinedNfsSpec = r.siteCtx.flavorNfsSpec
	}

	return err
}

// keydbSpec handle any keydb spec
func (r *SiteReconciler) keydbSpec() (err error) {
	r.siteCtx.hasKeydb = r.siteCtx.keydbSpecFound || r.siteCtx.flavorKeydbSpecFound

	if r.siteCtx.hasKeydb {
		r.siteCtx.keydb.SetName(r.siteCtx.keydbName)
		r.siteCtx.keydb.SetNamespace(r.siteCtx.namespaceName)
	}

	// Keydb kind from Keydb ansible operator
	if r.siteCtx.hasKeydb {
		// Set Keydb host and secret, if not already present in Moodle spec
		keydbRelatedMoodleSpec := map[string]interface{}{
			"moodleKeydbMetaName": r.siteCtx.keydbName,
		}
		// Merge Moodle related keydb spec with flavor Moodle spec
		if err := mergo.MapWithOverwrite(&r.siteCtx.flavorMoodleSpec, keydbRelatedMoodleSpec); err != nil {
			return err
		}
		// Merge Keydb spec if set on site Spec
		if r.siteCtx.keydbSpecFound {
			if err := mergo.MapWithOverwrite(&r.siteCtx.flavorKeydbSpec, r.siteCtx.keydbSpec); err != nil {
				return err
			}
		}
		// Set site labels to keydb
		if err := r.commonLabels(r.siteCtx.flavorKeydbSpec); err != nil {
			return err
		}
		// set default affinity
		if err := r.defaultAffinityYaml(r.siteCtx.flavorKeydbSpec, "keydbAffinity"); err != nil {
			return err
		}
		// save keydb spec
		r.siteCtx.combinedKeydbSpec = make(map[string]interface{})
		r.siteCtx.combinedKeydbSpec = r.siteCtx.flavorKeydbSpec
	}

	return err
}

// moodleSpec handle any keydb spec
func (r *SiteReconciler) moodleSpec() (err error) {
	// Merge ingress annotations
	if err := r.mergeNestedString(r.siteCtx.moodleSpec, r.siteCtx.flavorMoodleSpec, "nginxIngressAnnotations"); err != nil {
		return err
	}

	// Merge Moodle spec if set on site Spec
	if r.siteCtx.moodleSpecFound {
		if err := mergo.MapWithOverwrite(&r.siteCtx.flavorMoodleSpec, r.siteCtx.moodleSpec); err != nil {
			return err
		}
	}
	// Set site labels to Moodle
	if err := r.commonLabels(r.siteCtx.flavorMoodleSpec); err != nil {
		return err
	}
	// set moodle default affinity
	if err := r.moodleDefaultAffinityYaml(); err != nil {
		return err
	}
	// save moodle spec
	r.siteCtx.combinedMoodleSpec = make(map[string]interface{})
	r.siteCtx.combinedMoodleSpec = r.siteCtx.flavorMoodleSpec

	return err
}
