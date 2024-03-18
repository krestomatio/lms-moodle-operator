/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package m4e

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	m4ev1alpha1 "github.com/krestomatio/kio-operator/apis/m4e/v1alpha1"
)

const (
	SiteNamePrefix  string = "site-"
	SiteFinalizer   string = "m4e.krestomat.io/finalizer"
	FlavorNameIndex string = "spec.flavor"
)

type SiteReconcilerContext struct {
	hasNfs                  bool
	hasKeydb                bool
	hasPostgres             bool
	markedToBeDeleted       bool
	moodleSpecFound         bool
	nfsSpecFound            bool
	keydbSpecFound          bool
	postgresSpecFound       bool
	flavorNfsSpecFound      bool
	flavorKeydbSpecFound    bool
	flavorPostgresSpecFound bool
	name                    string
	flavorName              string
	state                   string
	namespaceName           string
	networkPolicyName       string
	moodleName              string
	nfsName                 string
	keydbName               string
	postgresName            string
	site                    *unstructured.Unstructured
	flavor                  *unstructured.Unstructured
	moodle                  *unstructured.Unstructured
	nfs                     *unstructured.Unstructured
	keydb                   *unstructured.Unstructured
	postgres                *unstructured.Unstructured
	spec                    map[string]interface{}
	moodleSpec              map[string]interface{}
	nfsSpec                 map[string]interface{}
	keydbSpec               map[string]interface{}
	postgresSpec            map[string]interface{}
	flavorSpec              map[string]interface{}
	flavorMoodleSpec        map[string]interface{}
	flavorNfsSpec           map[string]interface{}
	flavorKeydbSpec         map[string]interface{}
	flavorPostgresSpec      map[string]interface{}
	combinedMoodleSpec      map[string]interface{}
	combinedNfsSpec         map[string]interface{}
	combinedKeydbSpec       map[string]interface{}
	combinedPostgresSpec    map[string]interface{}
	namespace               *corev1.Namespace
	networkPolicy           *networkingv1.NetworkPolicy
}

type FlavorNotFoundError struct {
	Name string // Flavor name
}

func (f *FlavorNotFoundError) Error() string {
	return fmt.Sprintf("Flavor '%s' not found", f.Name)
}

// SiteReconciler reconciles a Site object
type SiteReconciler struct {
	client.Client
	Scheme                                   *runtime.Scheme
	MoodleGVK, NfsGVK, KeydbGVK, PostgresGVK schema.GroupVersionKind
	siteCtx                                  SiteReconcilerContext
}

//+kubebuilder:rbac:groups=m4e.krestomat.io,resources=sites,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=m4e.krestomat.io,resources=sites/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=m4e.krestomat.io,resources=sites/finalizers,verbs=update
//+kubebuilder:rbac:groups=m4e.krestomat.io,resources=moodles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nfs.krestomat.io,resources=ganeshas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=keydb.krestomat.io,resources=keydbs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=postgres.krestomat.io,resources=postgres,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create
//+kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Site object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *SiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Starting reconcile")

	// Vars
	r.siteCtx.name = req.Name

	// Prepare resource, saved any error for later
	if err := r.reconcilePrepare(ctx); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Finalize logic
	if finalized, requeue, err := r.reconcileFinalize(ctx); err != nil {
		return ctrl.Result{}, err
	} else if finalized || requeue {
		return ctrl.Result{Requeue: requeue}, nil
	}

	// Suspend logic
	if r.siteCtx.state == "suspended" {
		if requeue, err := r.reconcileSuspend(ctx); err != nil {
			return ctrl.Result{}, err
		} else {
			return ctrl.Result{Requeue: requeue}, nil
		}
	}

	// Present resources
	if requeue, err := r.reconcilePresent(ctx); err != nil {
		return ctrl.Result{}, err
	} else {
		return ctrl.Result{Requeue: requeue}, nil
	}
}

// reconcilePrepare takes care of initial step during reconcile
func (r *SiteReconciler) reconcilePrepare(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.V(1).Info("Reconcile set")

	// set base name for dependant resources
	baseName := truncate(SiteNamePrefix+r.siteCtx.name, 18)
	baseNamespace := SiteNamePrefix + r.siteCtx.name
	// if site name already include the prefix, do not use it
	if hasPrefix := strings.HasPrefix(r.siteCtx.name, SiteNamePrefix); hasPrefix {
		baseName = truncate(r.siteCtx.name, 18)
		baseNamespace = r.siteCtx.name
	}
	// set namespace name. It must start with an alphabetic character
	r.siteCtx.namespaceName = baseNamespace
	// set network policy name. It must start with an alphabetic character
	r.siteCtx.networkPolicyName = baseName
	// set Moodle name. It must start with an alphabetic character
	r.siteCtx.moodleName = baseName
	// set Postgres name. It must start with an alphabetic character
	r.siteCtx.postgresName = baseName
	// set NFS Ganesha server name and namespace. It must start with an alphabetic character
	r.siteCtx.nfsName = baseName
	// set Keydb name. It must start with an alphabetic character
	r.siteCtx.keydbName = baseName
	// site namespace
	r.siteCtx.namespace = &corev1.Namespace{}
	r.siteCtx.namespace.SetName(r.siteCtx.namespaceName)
	// site network policy
	r.siteNetworkPolicy()
	// dependant components
	r.siteCtx.moodle = newUnstructuredObject(r.MoodleGVK)
	r.siteCtx.postgres = newUnstructuredObject(r.PostgresGVK)
	r.siteCtx.nfs = newUnstructuredObject(r.NfsGVK)
	r.siteCtx.keydb = newUnstructuredObject(r.KeydbGVK)
	// namespaces and names
	r.siteCtx.moodle.SetName(r.siteCtx.moodleName)
	r.siteCtx.moodle.SetNamespace(r.siteCtx.namespaceName)

	// Fetch Site instance
	r.siteCtx.site = newUnstructuredObject(m4ev1alpha1.GroupVersion.WithKind("Site"))
	if err := r.Get(ctx, types.NamespacedName{Name: r.siteCtx.name}, r.siteCtx.site); err != nil {
		log.V(1).Info(err.Error())
		return err
	} else {
		// whether site is marked to be deleted
		r.siteCtx.markedToBeDeleted = r.siteCtx.site.GetDeletionTimestamp() != nil
	}
	r.siteCtx.spec, _, _ = unstructured.NestedMap(r.siteCtx.site.UnstructuredContent(), "spec")
	r.siteCtx.moodleSpec, r.siteCtx.moodleSpecFound, _ = unstructured.NestedMap(r.siteCtx.spec, "moodleSpec")
	r.siteCtx.postgresSpec, r.siteCtx.postgresSpecFound, _ = unstructured.NestedMap(r.siteCtx.spec, "postgresSpec")
	r.siteCtx.nfsSpec, r.siteCtx.nfsSpecFound, _ = unstructured.NestedMap(r.siteCtx.spec, "nfsSpec")
	r.siteCtx.keydbSpec, r.siteCtx.keydbSpecFound, _ = unstructured.NestedMap(r.siteCtx.spec, "keydbSpec")
	r.siteCtx.flavorName, _, _ = unstructured.NestedString(r.siteCtx.spec, "flavor")
	r.siteCtx.state, _, _ = unstructured.NestedString(r.siteCtx.spec, "state")

	// Fetch flavor spec
	r.siteCtx.flavor = newUnstructuredObject(m4ev1alpha1.GroupVersion.WithKind("Flavor"))
	if err := r.Get(ctx, types.NamespacedName{Name: r.siteCtx.flavorName}, r.siteCtx.flavor); err != nil {
		log.Error(err, "Flavor not found")
		return &FlavorNotFoundError{r.siteCtx.flavorName}
	}
	r.siteCtx.flavorSpec, _, _ = unstructured.NestedMap(r.siteCtx.flavor.UnstructuredContent(), "spec")
	r.siteCtx.flavorMoodleSpec, _, _ = unstructured.NestedMap(r.siteCtx.flavorSpec, "moodleSpec")
	r.siteCtx.flavorPostgresSpec, r.siteCtx.flavorPostgresSpecFound, _ = unstructured.NestedMap(r.siteCtx.flavorSpec, "postgresSpec")
	r.siteCtx.flavorNfsSpec, r.siteCtx.flavorNfsSpecFound, _ = unstructured.NestedMap(r.siteCtx.flavorSpec, "nfsSpec")
	r.siteCtx.flavorKeydbSpec, r.siteCtx.flavorKeydbSpecFound, _ = unstructured.NestedMap(r.siteCtx.flavorSpec, "keydbSpec")

	// set labels
	if err := r.setSiteLabels(ctx); err != nil {
		return err
	}
	r.siteCtx.postgres.SetLabels(r.siteCtx.site.GetLabels())
	r.siteCtx.nfs.SetLabels(r.siteCtx.site.GetLabels())
	r.siteCtx.keydb.SetLabels(r.siteCtx.site.GetLabels())
	r.siteCtx.moodle.SetLabels(r.siteCtx.site.GetLabels())

	// whether Site has dependant components
	if err := r.postgresSpec(); err != nil {
		return err
	}
	if err := r.nfsSpec(); err != nil {
		return err
	}
	if err := r.keydbSpec(); err != nil {
		return err
	}

	// moodle spec
	if err := r.moodleSpec(); err != nil {
		return err
	}

	// set UUID when it has to notify status to a url
	if err := r.setNotifyUUID(); err != nil {
		log.Error(err, "Couldn't add status uuid")
		return err
	}

	return nil
}

// reconcileFinalize configures finalizer
func (r *SiteReconciler) reconcileFinalize(ctx context.Context) (finalized bool, requeue bool, err error) {
	log := log.FromContext(ctx)
	log.V(1).Info("Reconcile finalizer")

	// Check if Site instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	if r.siteCtx.markedToBeDeleted {
		// update site state (terminating)
		if requeue, err := r.updateSiteStatus(ctx); err != nil {
			return false, requeue, err
		}
		if controllerutil.ContainsFinalizer(r.siteCtx.site, SiteFinalizer) {
			// Run finalization logic for SiteFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if requeue, err := r.finalizeSite(ctx); err != nil || requeue {
				return false, requeue, err
			}

			// Remove SiteFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(r.siteCtx.site, SiteFinalizer)
			if err := r.Update(ctx, r.siteCtx.site); err != nil {
				return false, false, err
			}
		}
		// Finalized
		return true, false, nil
	}
	// Add finalizer for this CR
	if !controllerutil.ContainsFinalizer(r.siteCtx.site, SiteFinalizer) {
		controllerutil.AddFinalizer(r.siteCtx.site, SiteFinalizer)
		if err := r.Update(ctx, r.siteCtx.site); err != nil {
			return false, false, err
		}
	}
	return false, false, nil
}

// reconcileSuspend take care of suspend state
func (r *SiteReconciler) reconcileSuspend(ctx context.Context) (requeue bool, err error) {
	log := log.FromContext(ctx)
	log.V(1).Info("Reconcile persist")

	// Save Moodle spec
	r.siteCtx.moodle.Object["spec"] = r.siteCtx.combinedMoodleSpec
	// Set suspended
	if err := unstructured.SetNestedField(r.siteCtx.combinedMoodleSpec, "suspended", "cr_state"); err != nil {
		return false, err
	}
	// Update site status about Moodle
	r.SetMoodleReadyCondition(ctx, r.siteCtx.site, r.siteCtx.moodle)
	// Apply Moodle resource
	if err := r.ReconcileApply(ctx, r.siteCtx.site, r.siteCtx.moodle); err != nil {
		return false, err
	}
	// Whether moodle is suspended
	if suspended := r.isDependantSuspended(ctx, r.siteCtx.moodle); !suspended {
		log.Info("Moodle resource is being suspended")
		_, err := r.updateSiteStatus(ctx)
		return true, err
	}

	// Save Keydb spec
	if r.siteCtx.hasKeydb {
		r.siteCtx.keydb.Object["spec"] = r.siteCtx.combinedKeydbSpec
		// Set suspended
		if err := unstructured.SetNestedField(r.siteCtx.combinedKeydbSpec, "suspended", "cr_state"); err != nil {
			return false, err
		}
		// Update Site status about Keydb
		r.SetKeydbReadyCondition(ctx, r.siteCtx.site, r.siteCtx.keydb)
		// Apply Keydb resource
		if err := r.ReconcileApply(ctx, r.siteCtx.site, r.siteCtx.keydb); err != nil {
			return false, err
		}
		// Whether keydb is suspended
		if suspended := r.isDependantSuspended(ctx, r.siteCtx.keydb); !suspended {
			log.Info("Keydb resource is being suspended")
			_, err := r.updateSiteStatus(ctx)
			return true, err
		}
	}

	// Save NFS Ganesha server spec
	if r.siteCtx.hasNfs {
		// Save NFS Ganesha server spec
		r.siteCtx.nfs.Object["spec"] = r.siteCtx.combinedNfsSpec
		// Set suspended
		if err := unstructured.SetNestedField(r.siteCtx.combinedNfsSpec, "suspended", "cr_state"); err != nil {
			return false, err
		}
		// Update Site status about NFS Ganesha
		r.SetNfsReadyCondition(ctx, r.siteCtx.site, r.siteCtx.nfs)
		// Apply NFS Ganesha server resource
		if err := r.ReconcileApply(ctx, r.siteCtx.site, r.siteCtx.nfs); err != nil {
			return false, err
		}
		// Whether nfs is suspended
		if suspended := r.isDependantSuspended(ctx, r.siteCtx.nfs); !suspended {
			log.Info("Nfs resource is being suspended")
			_, err := r.updateSiteStatus(ctx)
			return true, err
		}
	}

	// Save Postgres spec
	if r.siteCtx.hasPostgres {
		r.siteCtx.postgres.Object["spec"] = r.siteCtx.combinedPostgresSpec
		// Set suspended
		if err := unstructured.SetNestedField(r.siteCtx.combinedPostgresSpec, "suspended", "cr_state"); err != nil {
			return false, err
		}
		// Update Site status about Postgres
		r.SetPostgresReadyCondition(ctx, r.siteCtx.site, r.siteCtx.postgres)
		// Apply Postgres resource
		if err := r.ReconcileApply(ctx, r.siteCtx.site, r.siteCtx.postgres); err != nil {
			return false, err
		}
		// Whether postgres is suspended
		if suspended := r.isDependantSuspended(ctx, r.siteCtx.postgres); !suspended {
			log.Info("Postgres resource is being suspended")
			_, err := r.updateSiteStatus(ctx)
			return true, err
		}
	}

	// site is suspended
	return r.updateSiteStatus(ctx)
}

// reconcilePresent take care of present state
func (r *SiteReconciler) reconcilePresent(ctx context.Context) (requeue bool, err error) {
	log := log.FromContext(ctx)
	log.V(1).Info("Reconcile persist")

	// Vars
	moodleReady := false
	nfsReady := !r.siteCtx.hasNfs
	keydbReady := !r.siteCtx.hasKeydb
	postgresReady := !r.siteCtx.hasPostgres

	// Create namespace
	if err := r.ReconcileCreate(ctx, r.siteCtx.site, r.siteCtx.namespace); err != nil {
		return false, err
	}

	// Create network policy
	if err := r.ReconcileCreate(ctx, r.siteCtx.site, r.siteCtx.networkPolicy); err != nil {
		return false, err
	}

	// Save Postgres spec
	if r.siteCtx.hasPostgres {
		r.siteCtx.postgres.Object["spec"] = r.siteCtx.combinedPostgresSpec
		// Update Site status about Postgres
		r.SetPostgresReadyCondition(ctx, r.siteCtx.site, r.siteCtx.postgres)
		// Apply Postgres resource
		if err := r.ReconcileApply(ctx, r.siteCtx.site, r.siteCtx.postgres); err != nil {
			return false, err
		}
		// check if postgres ready
		if postgresReady, err = getReadyStatus(ctx, r.siteCtx.postgres); err != nil {
			return false, err
		}
	}

	// Save Keydb spec
	if r.siteCtx.hasKeydb {
		r.siteCtx.keydb.Object["spec"] = r.siteCtx.combinedKeydbSpec
		// Update Site status about Keydb
		r.SetKeydbReadyCondition(ctx, r.siteCtx.site, r.siteCtx.keydb)
		// Apply Keydb resource
		if err := r.ReconcileApply(ctx, r.siteCtx.site, r.siteCtx.keydb); err != nil {
			return false, err
		}
		// check if keydb ready
		if keydbReady, err = getReadyStatus(ctx, r.siteCtx.keydb); err != nil {
			return false, err
		}
	}

	// Save NFS Ganesha server spec
	if r.siteCtx.hasNfs {
		// Save NFS Ganesha server spec
		r.siteCtx.nfs.Object["spec"] = r.siteCtx.combinedNfsSpec
		// Update Site status about NFS Ganesha
		r.SetNfsReadyCondition(ctx, r.siteCtx.site, r.siteCtx.nfs)
		// Apply NFS Ganesha server resource
		if err := r.ReconcileApply(ctx, r.siteCtx.site, r.siteCtx.nfs); err != nil {
			return false, err
		}
		// check if nfs ready
		if nfsReady, err = getReadyStatus(ctx, r.siteCtx.nfs); err != nil {
			return false, err
		}
	}

	// Wait for postgres to be ready; otherwise requeue
	if !postgresReady {
		log.Info("Postgres is not ready, requeueing...", "Postgres.Name", r.siteCtx.postgres.GetName())
		return r.updateSiteStatus(ctx)
	}
	// Wait for Keydb to be ready; otherwise requeue
	if !keydbReady {
		log.Info("Keydb is not ready, requeueing...", "Keydb.Name", r.siteCtx.keydb.GetName())
		return r.updateSiteStatus(ctx)
	}
	// Wait for NFS Ganesha to be ready; otherwise requeue
	// NFS Ganesha server must be ready in order to mount its export as pvc
	if !nfsReady {
		log.Info("(NFS) Ganesha server is not ready, requeueing...", "Ganesha.Name", r.siteCtx.nfs.GetName())
		return r.updateSiteStatus(ctx)
	}

	// Save Moodle spec
	r.siteCtx.moodle.Object["spec"] = r.siteCtx.combinedMoodleSpec
	// Update site status about Moodle
	r.SetMoodleReadyCondition(ctx, r.siteCtx.site, r.siteCtx.moodle)
	// Apply Moodle resource
	if err := r.ReconcileApply(ctx, r.siteCtx.site, r.siteCtx.moodle); err != nil {
		return false, err
	}
	// check if moodle ready
	if moodleReady, err = getReadyStatus(ctx, r.siteCtx.moodle); err != nil {
		return false, err
	}
	// Wait for Moodle to be ready; otherwise requeue
	if !moodleReady {
		log.Info("Moodle is not ready, requeueing...", "Moodle.Name", r.siteCtx.moodle.GetName())
		return r.updateSiteStatus(ctx)
	}

	// site is ready
	return r.updateSiteStatus(ctx)
}

// ignoreDeletionPredicate filters Delete events on resources that have been confirmed deleted
func ignoreDeletionPredicate() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: func(e event.DeleteEvent) bool {
			// Evaluates to false if the object has been confirmed deleted.
			return !e.DeleteStateUnknown
		},
	}
}

// sitesByFlavor select sites that are using a flavor
// It returns a list of reconcile.Request
func (r *SiteReconciler) sitesByFlavor(flavor client.Object) []reconcile.Request {
	SiteList := &m4ev1alpha1.SiteList{}

	// Filter the list of sites by the ones using the flavor name
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(FlavorNameIndex, flavor.GetName()),
	}

	err := r.Client.List(context.Background(), SiteList, listOps)
	if err != nil {
		return []reconcile.Request{}
	}

	reconcileRequests := make([]reconcile.Request, len(SiteList.Items))
	for i, site := range SiteList.Items {
		reconcileRequests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name: site.Name,
			},
		}
	}
	return reconcileRequests
}

// SetupWithManager sets up the controller with the Manager.
func (r *SiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	mgr.GetScheme().AddKnownTypeWithName(r.MoodleGVK, &unstructured.Unstructured{})
	metav1.AddToGroupVersion(mgr.GetScheme(), r.MoodleGVK.GroupVersion())

	// Add spec.flavor index
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &m4ev1alpha1.Site{}, FlavorNameIndex, func(obj client.Object) []string {
		flavorName := obj.(*m4ev1alpha1.Site).Spec.Flavor
		if flavorName == "" {
			return nil
		}
		return []string{flavorName}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&m4ev1alpha1.Site{}).
		WithEventFilter(ignoreDeletionPredicate()).
		Owns(newUnstructuredObject(r.MoodleGVK)).
		Owns(newUnstructuredObject(r.NfsGVK)).
		Owns(newUnstructuredObject(r.KeydbGVK)).
		Owns(newUnstructuredObject(r.PostgresGVK)).
		Watches(&source.Kind{Type: &m4ev1alpha1.Flavor{}}, handler.EnqueueRequestsFromMapFunc(r.sitesByFlavor)).
		Complete(r)
}
