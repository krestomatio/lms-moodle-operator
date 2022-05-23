package m4e

import (
	"context"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	m4ev1alpha1 "github.com/krestomatio/kio-operator/apis/m4e/v1alpha1"
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

	// Delete m4e and requeue in order to wait for it to be completely be removed.
	// By doing so, any NFS Server or Keydb removal will be done after, and removal
	// conflicts will be avoided
	log.Info("Deleting M4e", "M4e.Namespace", r.siteCtx.m4e.GetNamespace(), "M4e.Name", r.siteCtx.m4e.GetName())
	if err := r.ReconcileDeleteDependant(ctx, r.siteCtx.site, r.siteCtx.m4e); err == nil {
		log.V(1).Info("Requeueing after M4e deletion", "M4e.Namespace", r.siteCtx.m4e.GetNamespace(), "M4e.Name", r.siteCtx.m4e.GetName())
		return true, nil
	} else if !errors.IsNotFound(err) {
		log.Error(err, "M4e not deleted", "M4e.Namespace", r.siteCtx.m4e.GetNamespace(), "M4e.Name", r.siteCtx.m4e.GetName())
		return false, err
	}

	// Delete Keydb and requeue in order to wait for it to be completely be removed.
	if r.siteCtx.hasKeydb {
		log.Info("Deleting Keydb", "Keydb.Namespace", r.siteCtx.keydb.GetNamespace(), "Keydb.Name", r.siteCtx.keydb.GetName())
		if err := r.ReconcileDeleteDependant(ctx, r.siteCtx.site, r.siteCtx.keydb); err == nil {
			log.V(1).Info("Requeueing after Keydb deletion", "Keydb.Namespace", r.siteCtx.keydb.GetNamespace(), "Keydb.Name", r.siteCtx.keydb.GetName())
			return true, nil
		} else if !errors.IsNotFound(err) {
			log.Error(err, "Keydb not deleted", "Keydb.Namespace", r.siteCtx.keydb.GetNamespace(), "Keydb.Name", r.siteCtx.keydb.GetName())
			return false, err
		}
	}

	// Delete nfs server and requeue in order to wait for it to be completely be removed.
	if r.siteCtx.hasNfs {
		log.Info("Deleting NFS Server", "Server.Namespace", r.siteCtx.nfs.GetNamespace(), "Server.Name", r.siteCtx.nfs.GetName())
		if err := r.ReconcileDeleteDependant(ctx, r.siteCtx.site, r.siteCtx.nfs); err == nil {
			log.V(1).Info("Requeueing after NFS Server deletion", "Server.Namespace", r.siteCtx.nfs.GetNamespace(), "Server.Name", r.siteCtx.nfs.GetName())
			return true, nil
		} else if !errors.IsNotFound(err) {
			log.Error(err, "NFS Server not deleted", "Server.Namespace", r.siteCtx.nfs.GetNamespace(), "Server.Name", r.siteCtx.nfs.GetName())
			return false, err
		}
	}

	// Delete Postgres and requeue in order to wait for it to be completely be removed.
	if r.siteCtx.hasPostgres {
		log.Info("Deleting Postgres", "Postgres.Namespace", r.siteCtx.postgres.GetNamespace(), "Postgres.Name", r.siteCtx.postgres.GetName())
		if err := r.ReconcileDeleteDependant(ctx, r.siteCtx.site, r.siteCtx.postgres); err == nil {
			log.V(1).Info("Requeueing after Postgres deletion", "Postgres.Namespace", r.siteCtx.postgres.GetNamespace(), "Postgres.Name", r.siteCtx.postgres.GetName())
			return true, nil
		} else if !errors.IsNotFound(err) {
			log.Error(err, "Postgres not deleted", "Postgres.Namespace", r.siteCtx.postgres.GetNamespace(), "Postgres.Name", r.siteCtx.postgres.GetName())
			return false, err
		}
	}

	log.Info("Successfully finalized site")
	return false, nil
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

	// get M4e ready condition
	m4eReadyCondition, m4eReadyConditionFound, m4eReadyConditionErr := getConditionByType(r.siteCtx.site, M4eReadyConditionType)

	if m4eReadyConditionErr != nil {
		log.Error(m4eReadyConditionErr, "unable to update Site '"+r.siteCtx.site.GetName()+"' state")
		return m4eReadyConditionErr
	}

	if readyConditionFound && m4eReadyConditionFound {
		state = r.setSiteState(readyCondition, m4eReadyCondition)
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
	stateUpdate, err := SetStatusState(r.flavorCtx.flavor, string(state))
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
func (r *SiteReconciler) setSiteState(readyCondition map[string]interface{}, m4eReadyCondition map[string]interface{}) string {
	status := readyCondition["status"]
	m4eStatus := m4eReadyCondition["status"]
	m4eReason := m4eReadyCondition["reason"]

	// Terminating
	if r.siteCtx.markedToBeDeleted {
		return string(m4ev1alpha1.TerminatingState)
	}

	if status == "False" || m4eStatus == "False" {
		// Failed
		if m4eReason == "Error" {
			return string(m4ev1alpha1.FailedState)
		}
		// Creating
		if m4eReason == "NotInstantiated" || m4eReason == "Instantiated" || m4eReason == "NotCreated" {
			return string(m4ev1alpha1.CreatingState)
		}
	}
	// Ready
	if m4eStatus == "True" {
		return string(m4ev1alpha1.ReadyState)
	}

	return string(m4ev1alpha1.UnknownState)
}

// setNotifyUUID defines site uuid if notifying status to an endpoint
// Should be used once combinedM4eSpec is set
// By default, site name is used as UUID
func (r *SiteReconciler) setNotifyUUID() error {
	// whether it has to notify status to a url
	_, m4eSiteNotifyStatusFound, _ := unstructured.NestedMap(r.siteCtx.combinedM4eSpec, "notifyStatus")
	if m4eSiteNotifyStatusFound {
		_, m4eSiteNotifyStatusUuidFound, _ := unstructured.NestedMap(r.siteCtx.combinedM4eSpec, "notifyStatus", "uuid")
		if !m4eSiteNotifyStatusUuidFound {
			// set uuid to notify about
			if err := unstructured.SetNestedField(r.siteCtx.combinedM4eSpec, r.siteCtx.name, "notifyStatus", "uuid"); err != nil {
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

// get an environmental variable
func getEnv(envVar string, defaultVal string) string {
	val, ok := os.LookupEnv(envVar)
	if !ok {
		return defaultVal
	} else {
		return val
	}
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
