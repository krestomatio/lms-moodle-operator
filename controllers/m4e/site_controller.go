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
	SiteNamePrefix string = "site-"
	SiteFinalizer  string = "m4e.krestomat.io/finalizer"
)

var (
	siteHasNfs            bool
	siteHasKeydb          bool
	siteName              string
	siteFlavor            string
	siteNamespaceName     string
	siteM4eName           string
	siteNfsName           string
	siteNfsNamespace      string = getEnv("NFSNAMESPACE", NFSNAMESPACE)
	siteKeydbName         string
	site                  *unstructured.Unstructured
	siteM4e               *unstructured.Unstructured
	siteNfs               *unstructured.Unstructured
	siteKeydb             *unstructured.Unstructured
	sitePreparedM4eSpec   map[string]interface{}
	sitePreparedNfsSpec   map[string]interface{}
	sitePreparedKeydbSpec map[string]interface{}
	siteNamespace         *corev1.Namespace
)

// SiteReconciler reconciles a Site object
type SiteReconciler struct {
	client.Client
	Scheme                   *runtime.Scheme
	M4eGVK, NfsGVK, KeydbGVK schema.GroupVersionKind
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
	siteName = req.Name

	// prepare resources
	if err := r.reconcilePrepare(ctx); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// set finalizer
	if finalized, requeue, err := r.reconcileFinalize(ctx); err != nil {
		return ctrl.Result{}, err
	} else if requeue {
		return ctrl.Result{Requeue: true}, nil
	} else if finalized {
		return ctrl.Result{}, nil
	}

	// apply resources
	if err := r.reconcilePersist(ctx); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// reconcilePrepare takes care of initial step during reconcile
func (r *SiteReconciler) reconcilePrepare(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Reconcile preparation")

	// Set namespace name. It must start with an alphabetic character
	siteNamespaceName = SiteNamePrefix + siteName
	// Set M4e name. It must start with an alphabetic character
	siteM4eName = SiteNamePrefix + truncate(siteName, 13)
	// Set NFS Server name and namespace. It must start with an alphabetic character
	siteNfsName = SiteNamePrefix + siteName
	siteNfsNamespace = getEnv("NFSNAMESPACE", NFSNAMESPACE)
	// Set Keydb name. It must start with an alphabetic character
	siteKeydbName = SiteNamePrefix + truncate(siteName, 13)
	// Site namespace
	siteNamespace = &corev1.Namespace{}
	siteNamespace.SetName(siteNamespaceName)

	// Fetch the Site instance
	site = newUnstructuredObject(m4ev1alpha1.GroupVersion.WithKind("Site"))
	if err := r.Get(ctx, types.NamespacedName{Name: siteName}, site); err != nil {
		log.V(1).Info(err.Error())
		return err
	}
	siteSpec, _, _ := unstructured.NestedMap(site.UnstructuredContent(), "spec")
	siteM4eSpec, siteM4eSpecFound, _ := unstructured.NestedMap(siteSpec, "m4eSpec")
	siteNfsSpec, siteNfsSpecFound, _ := unstructured.NestedMap(siteSpec, "nfsSpec")
	siteKeydbSpec, siteKeydbSpecFound, _ := unstructured.NestedMap(siteSpec, "keydbSpec")

	siteFlavor, _, _ = unstructured.NestedString(siteSpec, "flavor")

	siteCommonLabels := m4ev1alpha1.GroupVersion.Group + "/site_name: " + siteName + "\n" + m4ev1alpha1.GroupVersion.Group + "/flavor_name: " + siteFlavor

	// Fetch flavor spec
	flavor := newUnstructuredObject(m4ev1alpha1.GroupVersion.WithKind("Flavor"))
	if err := r.Get(ctx, types.NamespacedName{Name: siteFlavor}, flavor); err != nil {
		log.Info(err.Error())
		return err
	}

	flavorSpec, _, _ := unstructured.NestedMap(flavor.UnstructuredContent(), "spec")
	flavorM4eSpec, _, _ := unstructured.NestedMap(flavorSpec, "m4eSpec")
	flavorNfsSpec, flavorNfsSpecFound, _ := unstructured.NestedMap(flavorSpec, "nfsSpec")
	flavorKeydbSpec, flavorKeydbSpecFound, _ := unstructured.NestedMap(flavorSpec, "keydbSpec")

	// dependant components
	siteM4e = newUnstructuredObject(r.M4eGVK)
	siteNfs = newUnstructuredObject(r.NfsGVK)
	siteKeydb = newUnstructuredObject(r.KeydbGVK)
	siteHasNfs = siteNfsSpecFound || flavorNfsSpecFound
	siteHasKeydb = siteKeydbSpecFound || flavorKeydbSpecFound
	// namespaces and names
	siteM4e.SetName(siteM4eName)
	siteM4e.SetNamespace(siteNamespaceName)
	if siteHasNfs {
		siteNfs.SetName(siteNfsName)
		siteNfs.SetNamespace(siteNfsNamespace)
	}
	if siteHasKeydb {
		siteKeydb.SetName(siteKeydbName)
		siteKeydb.SetNamespace(siteNamespaceName)
	}

	// Server kind from NFS ansible operator
	if siteHasNfs {
		// Set NFS storage class name and access modes when using NFS operator
		nfsRelatedM4eSpec := map[string]interface{}{
			"moodlePvcMoodledataStorageClassName":  siteNfsName + "-nfs-sc",
			"moodlePvcMoodledataStorageAccessMode": m4ev1alpha1.ReadWriteMany,
		}
		if err := mergo.MapWithOverwrite(&flavorM4eSpec, nfsRelatedM4eSpec); err != nil {
			log.Error(err, "Couldn't merge spec")
			return err
		}
		// Merge NFS spec if set on site Spec
		if siteNfsSpecFound {
			if err := mergo.MapWithOverwrite(&flavorNfsSpec, siteNfsSpec); err != nil {
				log.Error(err, "Couldn't merge spec")
				return err
			}
		}
		// Set site labels to nfs
		flavorNfsSpecCommonLabelsString, flavorNfsSpecCommonLabelsFound, _ := unstructured.NestedString(flavorNfsSpec, "commonLabels")
		if flavorNfsSpecCommonLabelsFound {
			flavorNfsSpec["commonLabels"] = siteCommonLabels + "\n" + flavorNfsSpecCommonLabelsString
		} else {
			flavorNfsSpec["commonLabels"] = siteCommonLabels
		}
		// save nfs spec
		sitePreparedNfsSpec = make(map[string]interface{})
		sitePreparedNfsSpec = flavorNfsSpec
	}

	// Keydb kind from Keydb ansible operator
	if siteHasKeydb {
		// Set Keydb host and secret, if not already present in M4e spec
		keydbRelatedM4eSpec := map[string]interface{}{
			"moodleRedisHost":             siteKeydbName + "-keydb-service",
			"moodleRedisSecretAuthSecret": siteKeydbName + "-keydb-secret",
			"moodleRedisSecretAuthKey":    "keydb_password",
		}
		// Merge M4e related keydb spec with flavor M4e spec
		if err := mergo.MapWithOverwrite(&flavorM4eSpec, keydbRelatedM4eSpec); err != nil {
			log.Error(err, "Couldn't merge spec")
			return err
		}
		// Merge Keydb spec if set on site Spec
		if siteKeydbSpecFound {
			if err := mergo.MapWithOverwrite(&flavorKeydbSpec, siteKeydbSpec); err != nil {
				log.Error(err, "Couldn't merge spec")
				return err
			}
		}
		// Set site labels to keydb
		flavorKeydbSpecCommonLabelsString, flavorKeydbSpecCommonLabelsFound, _ := unstructured.NestedString(flavorKeydbSpec, "commonLabels")
		if flavorKeydbSpecCommonLabelsFound {
			flavorKeydbSpec["commonLabels"] = siteCommonLabels + "\n" + flavorKeydbSpecCommonLabelsString
		} else {
			flavorKeydbSpec["commonLabels"] = siteCommonLabels
		}
		// save keydb spec
		sitePreparedKeydbSpec = make(map[string]interface{})
		sitePreparedKeydbSpec = flavorKeydbSpec
	}

	// Merge M4e spec if set on site Spec
	if siteM4eSpecFound {
		if err := mergo.MapWithOverwrite(&flavorM4eSpec, siteM4eSpec); err != nil {
			log.Error(err, "Couldn't merge spec")
			return err
		}
	}
	// Set site labels to M4e
	flavorM4eSpecCommonLabelsString, flavorM4eSpecCommonLabelsFound, _ := unstructured.NestedString(flavorM4eSpec, "commonLabels")
	if flavorM4eSpecCommonLabelsFound {
		flavorM4eSpec["commonLabels"] = siteCommonLabels + "\n" + flavorM4eSpecCommonLabelsString
	} else {
		flavorM4eSpec["commonLabels"] = siteCommonLabels
	}
	// save m4e spec
	sitePreparedM4eSpec = make(map[string]interface{})
	sitePreparedM4eSpec = flavorM4eSpec
	return nil
}

// reconcileFinalize configures finalizer
func (r *SiteReconciler) reconcileFinalize(ctx context.Context) (finalized bool, requeue bool, err error) {
	log := log.FromContext(ctx)
	log.Info("Reconcile finalizer")

	// Check if Site instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	if isSiteMarkedToBeDeleted := site.GetDeletionTimestamp(); isSiteMarkedToBeDeleted != nil {
		if controllerutil.ContainsFinalizer(site, SiteFinalizer) {
			// Run finalization logic for SiteFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if requeue, err := r.finalizeSite(ctx); err != nil {
				return false, false, err
			} else if requeue {
				// if finalizer requires requeue
				return false, true, err
			}

			// Remove SiteFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(site, SiteFinalizer)
			if err := r.Update(ctx, site); err != nil {
				return false, false, err
			}
		}
		return true, false, nil
	}
	// Add finalizer for this CR
	if !controllerutil.ContainsFinalizer(site, SiteFinalizer) {
		controllerutil.AddFinalizer(site, SiteFinalizer)
		if err := r.Update(ctx, site); err != nil {
			return false, false, err
		}
	}
	return false, false, nil
}

// reconcilePersist take care of applying and persisting changes
func (r *SiteReconciler) reconcilePersist(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Reconcile persistance")

	// Vars
	m4eReady := false
	nfsReady := false
	keydbReady := false

	// Create namespace
	if err := r.ReconcileCreate(ctx, site, siteNamespace); err != nil {
		return err
	}

	// Save NFS Server spec
	if siteHasNfs {
		// Save NFS Server spec
		siteNfs.Object["spec"] = sitePreparedNfsSpec
		// Apply NFS Server resource
		if err := r.ReconcileApply(ctx, site, siteNfs); err != nil {
			return err
		}
		// Update Site status about NFS Server
		nfsReady = r.SetNfsReadyCondition(ctx, site, siteNfs)
	}

	// Save Keydb spec
	if siteHasKeydb {
		siteKeydb.Object["spec"] = sitePreparedKeydbSpec
		// Apply Keydb resource
		if err := r.ReconcileApply(ctx, site, siteKeydb); err != nil {
			return err
		}
		// Update Site status about Keydb
		keydbReady = r.SetKeydbReadyCondition(ctx, site, siteKeydb)
	}

	// Save M4e spec
	siteM4e.Object["spec"] = sitePreparedM4eSpec
	// Apply M4e resource
	if err := r.ReconcileApply(ctx, site, siteM4e); err != nil {
		return err
	}
	// Update site status about M4e
	m4eReady = r.SetM4eReadyCondition(ctx, site, siteM4e)

	// Set site ready contidion status
	if nfsReady && m4eReady && keydbReady {
		r.SetReadyCondition(ctx, site)
	}

	return nil
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

	return ctrl.NewControllerManagedBy(mgr).
		For(&m4ev1alpha1.Site{}).
		WithEventFilter(ignoreDeletionPredicate()).
		Owns(newUnstructuredObject(r.M4eGVK)).
		Owns(newUnstructuredObject(r.NfsGVK)).
		Owns(newUnstructuredObject(r.KeydbGVK)).
		Complete(r)
}
