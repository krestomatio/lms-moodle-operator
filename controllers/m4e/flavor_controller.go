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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	m4ev1alpha1 "github.com/krestomatio/kio-operator/apis/m4e/v1alpha1"
)

const (
	FlavorFinalizer string = "m4e.krestomat.io/finalizer"
)

type FlavorReconcilerContext struct {
	markedToBeDeleted bool
	flavor            *m4ev1alpha1.Flavor
}

type FlavorInUsedError struct {
	Name       string // Flavor name
	SiteNumber int    // Number of site using it
}

func (f *FlavorInUsedError) Error() string {
	return fmt.Sprintf("Flavor '%s' is in used by %d sites", f.Name, f.SiteNumber)
}

// FlavorReconciler reconciles a Flavor object
type FlavorReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	flavorCtx FlavorReconcilerContext
}

//+kubebuilder:rbac:groups=m4e.krestomat.io,resources=flavors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=m4e.krestomat.io,resources=flavors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=m4e.krestomat.io,resources=flavors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Flavor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *FlavorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Starting reconcile")

	// your logic here
	// Fetch the Memcached instance
	r.flavorCtx.flavor = &m4ev1alpha1.Flavor{}
	if err := r.Get(ctx, req.NamespacedName, r.flavorCtx.flavor); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Flavor resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Flavor")
		return ctrl.Result{}, err
	}

	// Whether it is marked to be deleted
	r.flavorCtx.markedToBeDeleted = r.flavorCtx.flavor.GetDeletionTimestamp() != nil

	// Finalize logic
	if finalized, err := r.reconcileFinalize(ctx); err != nil {
		return ctrl.Result{}, err
	} else if finalized {
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// reconcileFinalize configures finalizer
func (r *FlavorReconciler) reconcileFinalize(ctx context.Context) (finalized bool, err error) {
	log := log.FromContext(ctx)
	log.Info("Reconcile finalizer")

	// Check if Site instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	if r.flavorCtx.markedToBeDeleted {
		if controllerutil.ContainsFinalizer(r.flavorCtx.flavor, FlavorFinalizer) {
			// Run finalization logic for FlavorFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeFlavor(ctx); err != nil {
				return false, err
			}
			// Remove FlavorFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(r.flavorCtx.flavor, FlavorFinalizer)
			if err := r.Update(ctx, r.flavorCtx.flavor); err != nil {
				return false, err
			}
		}
		// Finalized
		return true, nil
	}
	// Add finalizer for this CR
	if !controllerutil.ContainsFinalizer(r.flavorCtx.flavor, FlavorFinalizer) {
		controllerutil.AddFinalizer(r.flavorCtx.flavor, FlavorFinalizer)
		if err := r.Update(ctx, r.flavorCtx.flavor); err != nil {
			return false, err
		}
	}
	return false, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FlavorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&m4ev1alpha1.Flavor{}).
		Complete(r)
}
