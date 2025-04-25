/*
Copyright 2024.

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

package lms

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	lmsv1alpha1 "github.com/krestomatio/lms-moodle-operator/api/lms/v1alpha1"
)

const (
	LMSMoodleTemplateFinalizer string = "lms.krestomat.io/finalizer"
)

type LMSMoodleTemplateReconcilerContext struct {
	markedToBeDeleted bool
	name              string
	lmsMoodleTemplate *unstructured.Unstructured
}

type LMSMoodleTemplateInUsedError struct {
	Name            string // LMSMoodleTemplate name
	LMSMoodleNumber int    // Number of lms moodle using it
}

func (f *LMSMoodleTemplateInUsedError) Error() string {
	return fmt.Sprintf("LMSMoodleTemplate '%s' is in used by %d lms moodle", f.Name, f.LMSMoodleNumber)
}

// LMSMoodleTemplateReconciler reconciles a LMSMoodleTemplate object
type LMSMoodleTemplateReconciler struct {
	client.Client
	Scheme               *runtime.Scheme
	lmsMoodleTemplateCtx LMSMoodleTemplateReconcilerContext
}

// +kubebuilder:rbac:groups=lms.krestomat.io,resources=lmsmoodletemplates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lms.krestomat.io,resources=lmsmoodletemplates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lms.krestomat.io,resources=lmsmoodletemplates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LMSMoodleTemplate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *LMSMoodleTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Starting reconcile")

	// Fetch LMSMoodleTemplate instance
	r.lmsMoodleTemplateCtx.name = req.Name
	r.lmsMoodleTemplateCtx.lmsMoodleTemplate = newUnstructuredObject(lmsv1alpha1.GroupVersion.WithKind("LMSMoodleTemplate"))
	if err := r.Get(ctx, types.NamespacedName{Name: r.lmsMoodleTemplateCtx.name}, r.lmsMoodleTemplateCtx.lmsMoodleTemplate); err != nil {
		log.V(1).Info(err.Error())
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// whether LMSMoodleTemplate is marked to be deleted
	r.lmsMoodleTemplateCtx.markedToBeDeleted = r.lmsMoodleTemplateCtx.lmsMoodleTemplate.GetDeletionTimestamp() != nil

	// Finalize logic
	if finalized, err := r.reconcileFinalize(ctx); err != nil {
		return ctrl.Result{}, err
	} else if finalized {
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, r.updateLMSMoodleTemplateState(ctx)
}

// reconcileFinalize configures finalizer
func (r *LMSMoodleTemplateReconciler) reconcileFinalize(ctx context.Context) (finalized bool, err error) {
	log := log.FromContext(ctx)
	log.Info("Reconcile finalizer")

	// Check if LMSMoodle instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	if r.lmsMoodleTemplateCtx.markedToBeDeleted {
		// update lms moodle state (terminating)
		if err := r.updateLMSMoodleTemplateState(ctx); err != nil {
			return false, err
		}
		if controllerutil.ContainsFinalizer(r.lmsMoodleTemplateCtx.lmsMoodleTemplate, LMSMoodleTemplateFinalizer) {
			// Run finalization logic for LMSMoodleTemplateFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeLMSMoodleTemplate(ctx); err != nil {
				return false, err
			}
			// Remove LMSMoodleTemplateFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(r.lmsMoodleTemplateCtx.lmsMoodleTemplate, LMSMoodleTemplateFinalizer)
			if err := r.Update(ctx, r.lmsMoodleTemplateCtx.lmsMoodleTemplate); err != nil {
				return false, err
			}
		}
		// Finalized
		return true, nil
	}
	// Add finalizer for this CR
	if !controllerutil.ContainsFinalizer(r.lmsMoodleTemplateCtx.lmsMoodleTemplate, LMSMoodleTemplateFinalizer) {
		controllerutil.AddFinalizer(r.lmsMoodleTemplateCtx.lmsMoodleTemplate, LMSMoodleTemplateFinalizer)
		if err := r.Update(ctx, r.lmsMoodleTemplateCtx.lmsMoodleTemplate); err != nil {
			return false, err
		}
	}
	return false, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LMSMoodleTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lmsv1alpha1.LMSMoodleTemplate{}).
		Complete(r)
}
