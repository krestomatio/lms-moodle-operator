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

	m4ev1alpha1 "github.com/krestomatio/kio-operator/apis/m4e/v1alpha1"
)

const (
	siteControllerName string = "site_controller"
)

// SiteReconciler reconciles a Site object
type SiteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	M4eGVK schema.GroupVersionKind
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
	// Fetch the Memcached instance
	site := &m4ev1alpha1.Site{}
	err := r.Get(ctx, req.NamespacedName, site)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Site resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Site")
		return ctrl.Result{}, err
	}

	flavor := &m4ev1alpha1.Flavor{}
	err = r.Get(ctx, types.NamespacedName{Name: site.Spec.Flavor}, flavor)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Flavor resource not found. Ignoring since object must be deleted", "site.Spec.Flavor", site.Spec.Flavor)
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Flavor")
		return ctrl.Result{}, err
	}

	// M4e kind from ansible operator
	m4e := r.newM4eObject()
	m4e.SetName(req.Name)
	m4e.SetNamespace(site.Spec.Namespace)
	m4e.Object["spec"] = flavor.Spec
	log.V(1).Info("Setting M4e owner", "Owner", site.GetUID(), "M4e.Namespace", m4e.GetNamespace(), "M4e.Name", m4e.GetName())
	err = ctrl.SetControllerReference(site, m4e, r.Scheme)
	if err != nil {
		log.Error(err, "Failed to set owner reference")
		return ctrl.Result{}, err
	}
	// Patch M4e
	log.V(1).Info("Applying M4e changes", "M4e.Namespace", m4e.GetNamespace(), "M4e.Name", m4e.GetName())
	force := true
	err = r.Patch(ctx, m4e, client.Apply, &client.PatchOptions{Force: &force, FieldManager: siteControllerName})
	if err != nil {
		log.Error(err, "Failed to apply M4e changes", "M4e.Namespace", m4e.GetNamespace(), "M4e.Name", m4e.GetName())
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
		Owns(r.newM4eObject()).
		Complete(r)
}

func (r *SiteReconciler) newM4eObject() *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(r.M4eGVK)
	return obj
}
