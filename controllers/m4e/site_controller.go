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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/imdario/mergo"
	m4ev1alpha1 "github.com/krestomatio/kio-operator/apis/m4e/v1alpha1"
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
	var m4eReady bool
	var nfsReady bool
	var keydbReady bool

	// Fetch the Site instance
	site := newUnstructuredObject(m4ev1alpha1.GroupVersion.WithKind("Site"))
	if err := r.Get(ctx, req.NamespacedName, site); err != nil {
		log.V(1).Info(err.Error(), "name", req.NamespacedName.Name)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	siteSpec, _, _ := unstructured.NestedMap(site.UnstructuredContent(), "spec")
	siteFlavor, _, _ := unstructured.NestedString(siteSpec, "flavor")
	siteM4eSpec, siteM4eSpecFound, _ := unstructured.NestedMap(siteSpec, "m4eSpec")
	siteNfsSpec, siteNfsSpecFound, _ := unstructured.NestedMap(siteSpec, "nfsSpec")
	siteKeydbSpec, siteKeydbSpecFound, _ := unstructured.NestedMap(siteSpec, "keydbSpec")

	siteCommonLabels := m4ev1alpha1.GroupVersion.Group + "/site_name: " + req.Name + "\n" + m4ev1alpha1.GroupVersion.Group + "/flavor_name: " + siteFlavor

	// Fetch flavor spec
	flavor := newUnstructuredObject(m4ev1alpha1.GroupVersion.WithKind("Flavor"))
	if err := r.Get(ctx, types.NamespacedName{Name: siteFlavor}, flavor); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Flavor resource not found", "site.Spec.Flavor", siteFlavor)
			return ctrl.Result{Requeue: false}, nil
		}
		log.Error(err, "Failed to get Flavor", "site.Spec.Flavor", siteFlavor)
		return ctrl.Result{Requeue: true}, err
	}

	flavorSpec, _, _ := unstructured.NestedMap(flavor.UnstructuredContent(), "spec")
	flavorM4eSpec, _, _ := unstructured.NestedMap(flavorSpec, "m4eSpec")
	flavorNfsSpec, flavorNfsSpecFound, _ := unstructured.NestedMap(flavorSpec, "nfsSpec")
	flavorKeydbSpec, flavorKeydbSpecFound, _ := unstructured.NestedMap(flavorSpec, "keydbSpec")

	// Site namespace
	ns := &corev1.Namespace{}
	// Set namespace name. It must start with an alphabetic character
	nsName := "site-" + req.Name
	ns.SetName(nsName)
	// Create namespace
	if err := r.reconcileCreate(ctx, site, ns); err != nil {
		return ctrl.Result{}, err
	}

	// Server kind from NFS ansible operator
	if siteNfsSpecFound || flavorNfsSpecFound {
		nfs := newUnstructuredObject(r.NfsGVK)
		// Set NFS Server name. It must start with an alphabetic character
		nfsName := "site-" + req.Name
		nfs.SetName(nfsName)
		nfs.SetNamespace(getEnv("NFSNAMESPACE", NFSNAMESPACE))
		// Set NFS storage class name and access modes when using NFS operator
		nfsRelatedM4eSpec := map[string]interface{}{
			"moodlePvcMoodledataStorageClassName":  nfsName + "-nfs-sc",
			"moodlePvcMoodledataStorageAccessMode": m4ev1alpha1.ReadWriteMany,
		}
		if err := mergo.MapWithOverwrite(&flavorM4eSpec, nfsRelatedM4eSpec); err != nil {
			log.V(1).Info(err.Error(), "name", nfs.GetName())
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		// Merge NFS spec if set on site Spec
		if siteNfsSpecFound {
			if err := mergo.MapWithOverwrite(&flavorNfsSpec, siteNfsSpec); err != nil {
				log.V(1).Info(err.Error(), "name", nfs.GetName())
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}
		}
		// Set site labels to nfs
		flavorNfsSpecCommonLabelsString, flavorNfsSpecCommonLabelsFound, _ := unstructured.NestedString(flavorNfsSpec, "commonLabels")
		if flavorNfsSpecCommonLabelsFound {
			flavorNfsSpec["commonLabels"] = siteCommonLabels + "\n" + flavorNfsSpecCommonLabelsString
		} else {
			flavorNfsSpec["commonLabels"] = siteCommonLabels
		}
		// Save NFS server spec
		nfs.Object["spec"] = flavorNfsSpec
		// Apply Server resource
		if err := r.reconcileApply(ctx, site, nfs); err != nil {
			return ctrl.Result{}, err
		}
		// Update Site status about NFS server
		nfsReady = r.SetNfsReadyCondition(ctx, site, nfs)
	}

	// Keydb kind from Keydb ansible operator
	if siteKeydbSpecFound || flavorKeydbSpecFound {
		keydb := newUnstructuredObject(r.KeydbGVK)
		// Set Keydb name. It must start with an alphabetic character
		keydbName := "site-" + req.Name
		keydb.SetName(keydbName)
		keydb.SetNamespace(ns.GetName())
		// Set Keydb host and secret, if not already present in M4e spec
		keydbRelatedM4eSpec := map[string]interface{}{
			"moodleRedisHost":             keydbName + "-keydb-service",
			"moodleRedisSecretAuthSecret": keydbName + "-keydb-secret",
			"moodleRedisSecretAuthKey":    "keydb_password",
		}
		// Merge M4e related keydb spec with flavor M4e spec
		if err := mergo.MapWithOverwrite(&flavorM4eSpec, keydbRelatedM4eSpec); err != nil {
			log.V(1).Info(err.Error(), "name", keydb.GetName())
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		// Merge Keydb spec if set on site Spec
		if siteKeydbSpecFound {
			if err := mergo.MapWithOverwrite(&flavorKeydbSpec, siteKeydbSpec); err != nil {
				log.V(1).Info(err.Error(), "name", keydb.GetName())
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}
		}
		// Set site labels to keydb
		flavorKeydbSpecCommonLabelsString, flavorKeydbSpecCommonLabelsFound, _ := unstructured.NestedString(flavorKeydbSpec, "commonLabels")
		if flavorKeydbSpecCommonLabelsFound {
			flavorKeydbSpec["commonLabels"] = siteCommonLabels + "\n" + flavorKeydbSpecCommonLabelsString
		} else {
			flavorKeydbSpec["commonLabels"] = siteCommonLabels
		}
		// Save Keydb spec
		keydb.Object["spec"] = flavorKeydbSpec
		// Apply Keydb resource
		if err := r.reconcileApply(ctx, site, keydb); err != nil {
			return ctrl.Result{}, err
		}
		// Update Site status about Keydb server
		keydbReady = r.SetKeydbReadyCondition(ctx, site, keydb)
	}

	// M4e kind from M4e ansible operator
	m4e := newUnstructuredObject(r.M4eGVK)
	m4eName := "site-" + truncate(req.Name, 13)
	m4e.SetName(m4eName)
	m4e.SetNamespace(ns.GetName())
	// Merge M4e spec if set on site Spec
	if siteM4eSpecFound {
		if err := mergo.MapWithOverwrite(&flavorM4eSpec, siteM4eSpec); err != nil {
			log.V(1).Info(err.Error(), "name", m4e.GetName())
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
	}
	// Set site labels to M4e
	flavorM4eSpecCommonLabelsString, flavorM4eSpecCommonLabelsFound, _ := unstructured.NestedString(flavorM4eSpec, "commonLabels")
	if flavorM4eSpecCommonLabelsFound {
		flavorM4eSpec["commonLabels"] = siteCommonLabels + "\n" + flavorM4eSpecCommonLabelsString
	} else {
		flavorM4eSpec["commonLabels"] = siteCommonLabels
	}
	// Save M4e server spec
	m4e.Object["spec"] = flavorM4eSpec
	// Apply M4e resource
	if err := r.reconcileApply(ctx, site, m4e); err != nil {
		return ctrl.Result{}, err
	}
	// Update Site status about M4e
	m4eReady = r.SetM4eReadyCondition(ctx, site, m4e)

	// Set site ready contidion status
	if nfsReady && m4eReady && keydbReady {
		r.SetReadyCondition(ctx, site)
	}

	return ctrl.Result{}, nil
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
