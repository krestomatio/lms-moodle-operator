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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/imdario/mergo"
	m4ev1alpha1 "github.com/krestomatio/kio-operator/apis/m4e/v1alpha1"
)

const (
	SiteNamePrefix  string = "site-"
	SiteFinalizer   string = "m4e.krestomat.io/finalizer"
	FlavorNameIndex string = "spec.flavor"
)

type SiteReconcilerContext struct {
	hasNfs               bool
	hasKeydb             bool
	markedToBeDeleted    bool
	m4eSpecFound         bool
	nfsSpecFound         bool
	keydbSpecFound       bool
	flavorNfsSpecFound   bool
	flavorKeydbSpecFound bool
	name                 string
	flavorName           string
	namespaceName        string
	m4eName              string
	nfsName              string
	nfsNamespaceName     string
	keydbName            string
	commonLabels         string
	site                 *unstructured.Unstructured
	flavor               *unstructured.Unstructured
	m4e                  *unstructured.Unstructured
	nfs                  *unstructured.Unstructured
	keydb                *unstructured.Unstructured
	spec                 map[string]interface{}
	m4eSpec              map[string]interface{}
	nfsSpec              map[string]interface{}
	keydbSpec            map[string]interface{}
	flavorSpec           map[string]interface{}
	flavorM4eSpec        map[string]interface{}
	flavorNfsSpec        map[string]interface{}
	flavorKeydbSpec      map[string]interface{}
	combinedM4eSpec      map[string]interface{}
	combinedNfsSpec      map[string]interface{}
	combinedKeydbSpec    map[string]interface{}
	namespace            *corev1.Namespace
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
	Scheme                   *runtime.Scheme
	M4eGVK, NfsGVK, KeydbGVK schema.GroupVersionKind
	siteCtx                  SiteReconcilerContext
}

//+kubebuilder:rbac:groups=m4e.krestomat.io,resources=sites,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=m4e.krestomat.io,resources=sites/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=m4e.krestomat.io,resources=sites/finalizers,verbs=update
//+kubebuilder:rbac:groups=m4e.krestomat.io,resources=m4es,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nfs.krestomat.io,resources=servers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=keydb.krestomat.io,resources=keydbs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Site object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
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

	// Patch resources
	if err := r.reconcilePersist(ctx); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// reconcilePrepare takes care of initial step during reconcile
func (r *SiteReconciler) reconcilePrepare(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.V(1).Info("Reconcile set")

	// set namespace name. It must start with an alphabetic character
	r.siteCtx.namespaceName = SiteNamePrefix + r.siteCtx.name
	// set M4e name. It must start with an alphabetic character
	r.siteCtx.m4eName = SiteNamePrefix + truncate(r.siteCtx.name, 13)
	// set NFS Server name and namespace. It must start with an alphabetic character
	r.siteCtx.nfsName = SiteNamePrefix + r.siteCtx.name
	r.siteCtx.nfsNamespaceName = getEnv("NFSNAMESPACE", NFSNAMESPACE)
	// set Keydb name. It must start with an alphabetic character
	r.siteCtx.keydbName = SiteNamePrefix + truncate(r.siteCtx.name, 13)
	// site namespace
	r.siteCtx.namespace = &corev1.Namespace{}
	r.siteCtx.namespace.SetName(r.siteCtx.namespaceName)
	// dependant components
	r.siteCtx.m4e = newUnstructuredObject(r.M4eGVK)
	r.siteCtx.nfs = newUnstructuredObject(r.NfsGVK)
	r.siteCtx.keydb = newUnstructuredObject(r.KeydbGVK)
	// namespaces and names
	r.siteCtx.m4e.SetName(r.siteCtx.m4eName)
	r.siteCtx.m4e.SetNamespace(r.siteCtx.namespaceName)

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
	r.siteCtx.m4eSpec, r.siteCtx.m4eSpecFound, _ = unstructured.NestedMap(r.siteCtx.spec, "m4eSpec")
	r.siteCtx.nfsSpec, r.siteCtx.nfsSpecFound, _ = unstructured.NestedMap(r.siteCtx.spec, "nfsSpec")
	r.siteCtx.keydbSpec, r.siteCtx.keydbSpecFound, _ = unstructured.NestedMap(r.siteCtx.spec, "keydbSpec")

	r.siteCtx.flavorName, _, _ = unstructured.NestedString(r.siteCtx.spec, "flavor")

	r.siteCtx.commonLabels = m4ev1alpha1.GroupVersion.Group + "/site_name: " + r.siteCtx.name + "\n" + m4ev1alpha1.GroupVersion.Group + "/flavor_name: " + r.siteCtx.flavorName

	// Fetch flavor spec
	r.siteCtx.flavor = newUnstructuredObject(m4ev1alpha1.GroupVersion.WithKind("Flavor"))
	if err := r.Get(ctx, types.NamespacedName{Name: r.siteCtx.flavorName}, r.siteCtx.flavor); err != nil {
		log.Error(err, "Flavor not found")
		return &FlavorNotFoundError{r.siteCtx.flavorName}
	}

	r.siteCtx.flavorSpec, _, _ = unstructured.NestedMap(r.siteCtx.flavor.UnstructuredContent(), "spec")
	r.siteCtx.flavorM4eSpec, _, _ = unstructured.NestedMap(r.siteCtx.flavorSpec, "m4eSpec")
	r.siteCtx.flavorNfsSpec, r.siteCtx.flavorNfsSpecFound, _ = unstructured.NestedMap(r.siteCtx.flavorSpec, "nfsSpec")
	r.siteCtx.flavorKeydbSpec, r.siteCtx.flavorKeydbSpecFound, _ = unstructured.NestedMap(r.siteCtx.flavorSpec, "keydbSpec")

	// whether Site has dependant components
	r.siteCtx.hasNfs = r.siteCtx.nfsSpecFound || r.siteCtx.flavorNfsSpecFound
	r.siteCtx.hasKeydb = r.siteCtx.keydbSpecFound || r.siteCtx.flavorKeydbSpecFound
	if r.siteCtx.hasNfs {
		r.siteCtx.nfs.SetName(r.siteCtx.nfsName)
		r.siteCtx.nfs.SetNamespace(r.siteCtx.nfsNamespaceName)
	}
	if r.siteCtx.hasKeydb {
		r.siteCtx.keydb.SetName(r.siteCtx.keydbName)
		r.siteCtx.keydb.SetNamespace(r.siteCtx.namespaceName)
	}

	// Server kind from NFS ansible operator
	if r.siteCtx.hasNfs {
		// Set NFS storage class name and access modes when using NFS operator
		nfsRelatedM4eSpec := map[string]interface{}{
			"moodlePvcMoodledataStorageClassName":  r.siteCtx.nfsName + "-nfs-sc",
			"moodlePvcMoodledataStorageAccessMode": m4ev1alpha1.ReadWriteMany,
		}
		if err := mergo.MapWithOverwrite(&r.siteCtx.flavorM4eSpec, nfsRelatedM4eSpec); err != nil {
			log.Error(err, "Couldn't merge spec")
			return err
		}
		// Merge NFS spec if set on site Spec
		if r.siteCtx.nfsSpecFound {
			if err := mergo.MapWithOverwrite(&r.siteCtx.flavorNfsSpec, r.siteCtx.nfsSpec); err != nil {
				log.Error(err, "Couldn't merge spec")
				return err
			}
		}
		// Set site labels to nfs
		flavorNfsSpecCommonLabelsString, flavorNfsSpecCommonLabelsFound, _ := unstructured.NestedString(r.siteCtx.flavorNfsSpec, "commonLabels")
		if flavorNfsSpecCommonLabelsFound {
			r.siteCtx.flavorNfsSpec["commonLabels"] = r.siteCtx.commonLabels + "\n" + flavorNfsSpecCommonLabelsString
		} else {
			r.siteCtx.flavorNfsSpec["commonLabels"] = r.siteCtx.commonLabels
		}
		// save nfs spec
		r.siteCtx.combinedNfsSpec = make(map[string]interface{})
		r.siteCtx.combinedNfsSpec = r.siteCtx.flavorNfsSpec
	}

	// Keydb kind from Keydb ansible operator
	if r.siteCtx.hasKeydb {
		// Set Keydb host and secret, if not already present in M4e spec
		keydbRelatedM4eSpec := map[string]interface{}{
			"moodleRedisHost":             r.siteCtx.keydbName + "-keydb-service",
			"moodleRedisSecretAuthSecret": r.siteCtx.keydbName + "-keydb-secret",
			"moodleRedisSecretAuthKey":    "keydb_password",
		}
		// Merge M4e related keydb spec with flavor M4e spec
		if err := mergo.MapWithOverwrite(&r.siteCtx.flavorM4eSpec, keydbRelatedM4eSpec); err != nil {
			log.Error(err, "Couldn't merge spec")
			return err
		}
		// Merge Keydb spec if set on site Spec
		if r.siteCtx.keydbSpecFound {
			if err := mergo.MapWithOverwrite(&r.siteCtx.flavorKeydbSpec, r.siteCtx.keydbSpec); err != nil {
				log.Error(err, "Couldn't merge spec")
				return err
			}
		}
		// Set site labels to keydb
		flavorKeydbSpecCommonLabelsString, flavorKeydbSpecCommonLabelsFound, _ := unstructured.NestedString(r.siteCtx.flavorKeydbSpec, "commonLabels")
		if flavorKeydbSpecCommonLabelsFound {
			r.siteCtx.flavorKeydbSpec["commonLabels"] = r.siteCtx.commonLabels + "\n" + flavorKeydbSpecCommonLabelsString
		} else {
			r.siteCtx.flavorKeydbSpec["commonLabels"] = r.siteCtx.commonLabels
		}
		// save keydb spec
		r.siteCtx.combinedKeydbSpec = make(map[string]interface{})
		r.siteCtx.combinedKeydbSpec = r.siteCtx.flavorKeydbSpec
	}

	// Merge M4e spec if set on site Spec
	if r.siteCtx.m4eSpecFound {
		if err := mergo.MapWithOverwrite(&r.siteCtx.flavorM4eSpec, r.siteCtx.m4eSpec); err != nil {
			log.Error(err, "Couldn't merge spec")
			return err
		}
	}
	// Set site labels to M4e
	flavorM4eSpecCommonLabelsString, flavorM4eSpecCommonLabelsFound, _ := unstructured.NestedString(r.siteCtx.flavorM4eSpec, "commonLabels")
	if flavorM4eSpecCommonLabelsFound {
		r.siteCtx.flavorM4eSpec["commonLabels"] = r.siteCtx.commonLabels + "\n" + flavorM4eSpecCommonLabelsString
	} else {
		r.siteCtx.flavorM4eSpec["commonLabels"] = r.siteCtx.commonLabels
	}
	// save m4e spec
	r.siteCtx.combinedM4eSpec = make(map[string]interface{})
	r.siteCtx.combinedM4eSpec = r.siteCtx.flavorM4eSpec
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
		if err := r.updateSiteState(ctx); err != nil {
			return false, false, err
		}
		if controllerutil.ContainsFinalizer(r.siteCtx.site, SiteFinalizer) {
			// Run finalization logic for SiteFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if requeue, err := r.finalizeSite(ctx); err != nil {
				return false, false, err
			} else if requeue {
				// finalizer requires requeue
				return false, true, err
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

// reconcilePersist take care of applying and persisting changes
func (r *SiteReconciler) reconcilePersist(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.V(1).Info("Reconcile persist")

	// Vars
	m4eReady := false
	nfsReady := !r.siteCtx.hasNfs
	keydbReady := !r.siteCtx.hasKeydb

	// Create namespace
	if err := r.ReconcileCreate(ctx, r.siteCtx.site, r.siteCtx.namespace); err != nil {
		return err
	}

	// Save NFS Server spec
	if r.siteCtx.hasNfs {
		// Save NFS Server spec
		r.siteCtx.nfs.Object["spec"] = r.siteCtx.combinedNfsSpec
		// Apply NFS Server resource
		if err := r.ReconcileApply(ctx, r.siteCtx.site, r.siteCtx.nfs); err != nil {
			return err
		}
		// Update Site status about NFS Server
		nfsReady = r.SetNfsReadyCondition(ctx, r.siteCtx.site, r.siteCtx.nfs)
	}

	// Save Keydb spec
	if r.siteCtx.hasKeydb {
		r.siteCtx.keydb.Object["spec"] = r.siteCtx.combinedKeydbSpec
		// Apply Keydb resource
		if err := r.ReconcileApply(ctx, r.siteCtx.site, r.siteCtx.keydb); err != nil {
			return err
		}
		// Update Site status about Keydb
		keydbReady = r.SetKeydbReadyCondition(ctx, r.siteCtx.site, r.siteCtx.keydb)
	}

	// Save M4e spec
	r.siteCtx.m4e.Object["spec"] = r.siteCtx.combinedM4eSpec
	// Apply M4e resource
	if err := r.ReconcileApply(ctx, r.siteCtx.site, r.siteCtx.m4e); err != nil {
		return err
	}
	// Update site status about M4e
	m4eReady = r.SetM4eReadyCondition(ctx, r.siteCtx.site, r.siteCtx.m4e)

	// Set site ready contidion status and state
	if nfsReady && m4eReady && keydbReady {
		r.SetReadyCondition(ctx, r.siteCtx.site)
	}

	return r.updateSiteState(ctx)
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

// SetupWithManager sets up the controller with the Manager.
func (r *SiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	mgr.GetScheme().AddKnownTypeWithName(r.M4eGVK, &unstructured.Unstructured{})
	metav1.AddToGroupVersion(mgr.GetScheme(), r.M4eGVK.GroupVersion())

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
		Owns(newUnstructuredObject(r.M4eGVK)).
		Owns(newUnstructuredObject(r.NfsGVK)).
		Owns(newUnstructuredObject(r.KeydbGVK)).
		Complete(r)
}
