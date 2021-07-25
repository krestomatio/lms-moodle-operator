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

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/imdario/mergo"
	m4ev1alpha1 "github.com/krestomatio/kio-operator/apis/m4e/v1alpha1"
)

// SiteReconciler reconciles a Site object
type SiteReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	M4eGVK, NfsGVK schema.GroupVersionKind
}

//+kubebuilder:rbac:groups=m4e.app.krestomat.io,resources=sites,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=m4e.app.krestomat.io,resources=sites/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=m4e.app.krestomat.io,resources=sites/finalizers,verbs=update
//+kubebuilder:rbac:groups=m4e.krestomat.io,resources=m4e,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Site object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *SiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// your logic here

	// Fetch the Site instance
	site := newUnstructuredObject(m4ev1alpha1.GroupVersion.WithKind("Site"))
	if err := r.Get(ctx, req.NamespacedName, site); err != nil {
		log.V(1).Info(err.Error(), "name", req.NamespacedName.Name)
		return ctrl.Result{}, ignoreNotFound(err)
	}
	siteSpec, _, _ := unstructured.NestedMap(site.UnstructuredContent(), "spec")
	siteNamespace, _, _ := unstructured.NestedString(siteSpec, "namespace")
	siteFlavor, _, _ := unstructured.NestedString(siteSpec, "flavor")
	siteM4eSpec, siteM4eSpecFound, _ := unstructured.NestedMap(siteSpec, "m4eSpec")
	siteNfsSpec, siteNfsSpecFound, _ := unstructured.NestedMap(siteSpec, "nfsSpec")

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

	// Server kind from NFS ansible operator
	if siteNfsSpecFound || flavorNfsSpecFound {
		nfs := newUnstructuredObject(r.NfsGVK)
		nfs.SetName(req.Name)
		nfs.SetNamespace(getEnv("NFSNAMESPACE", NFSNAMESPACE))
		if siteNfsSpecFound {
			mergo.MapWithOverwrite(&flavorNfsSpec, siteNfsSpec)
		}
		nfs.Object["spec"] = flavorNfsSpec

		if _, err := r.reconcileApply(ctx, site, nfs); err != nil {
			return ctrl.Result{}, err
		}
	}

	// M4e kind from M4e ansible operator
	m4e := newUnstructuredObject(r.M4eGVK)
	m4e.SetName(req.Name)
	m4e.SetNamespace(siteNamespace)
	if siteM4eSpecFound {
		mergo.MapWithOverwrite(&flavorM4eSpec, siteM4eSpec)
	}
	m4e.Object["spec"] = flavorM4eSpec

	if _, err := r.reconcileApply(ctx, site, m4e); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	mgr.GetScheme().AddKnownTypeWithName(r.M4eGVK, &unstructured.Unstructured{})
	metav1.AddToGroupVersion(mgr.GetScheme(), r.M4eGVK.GroupVersion())

	return ctrl.NewControllerManagedBy(mgr).
		For(&m4ev1alpha1.Site{}).
		Owns(newUnstructuredObject(r.M4eGVK)).
		Owns(newUnstructuredObject(r.NfsGVK)).
		Complete(r)
}
