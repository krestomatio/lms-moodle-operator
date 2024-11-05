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

	lmsv1alpha1 "github.com/krestomatio/lms-moodle-operator/api/lms/v1alpha1"
)

const (
	LMSMoodleNamePrefix        string = "lms-"
	LMSMoodleFinalizer         string = "lms.krestomat.io/finalizer"
	LMSMoodleTemplateNameIndex string = "spec.lmsMoodleTemplate"
	TruncateCharactersInName   int    = 17
)

type LMSMoodleReconcilerContext struct {
	hasNfs                             bool
	hasKeydb                           bool
	hasPostgres                        bool
	markedToBeDeleted                  bool
	moodleSpecFound                    bool
	nfsSpecFound                       bool
	keydbSpecFound                     bool
	postgresSpecFound                  bool
	lmsMoodleTemplateNfsSpecFound      bool
	lmsMoodleTemplateKeydbSpecFound    bool
	lmsMoodleTemplatePostgresSpecFound bool
	name                               string
	lmsMoodleTemplateName              string
	desiredState                       string
	namespaceName                      string
	networkPolicyBaseName              string
	moodleName                         string
	nfsName                            string
	keydbName                          string
	postgresName                       string
	lmsMoodle                          *unstructured.Unstructured
	lmsMoodleTemplate                  *unstructured.Unstructured
	moodle                             *unstructured.Unstructured
	nfs                                *unstructured.Unstructured
	keydb                              *unstructured.Unstructured
	postgres                           *unstructured.Unstructured
	spec                               map[string]interface{}
	moodleSpec                         map[string]interface{}
	nfsSpec                            map[string]interface{}
	keydbSpec                          map[string]interface{}
	postgresSpec                       map[string]interface{}
	lmsMoodleTemplateSpec              map[string]interface{}
	lmsMoodleTemplateMoodleSpec        map[string]interface{}
	lmsMoodleTemplateNfsSpec           map[string]interface{}
	lmsMoodleTemplateKeydbSpec         map[string]interface{}
	lmsMoodleTemplatePostgresSpec      map[string]interface{}
	combinedMoodleSpec                 map[string]interface{}
	combinedNfsSpec                    map[string]interface{}
	combinedKeydbSpec                  map[string]interface{}
	combinedPostgresSpec               map[string]interface{}
	namespace                          *corev1.Namespace
	lmsMoodleNetpolOmit                bool
	lmsMoodleDefaultNetpol             *networkingv1.NetworkPolicy
}

type LMSMoodleTemplateNotFoundError struct {
	Name string // LMSMoodleTemplate name
}

func (f *LMSMoodleTemplateNotFoundError) Error() string {
	return fmt.Sprintf("LMSMoodleTemplate '%s' not found", f.Name)
}

// LMSMoodleReconciler reconciles a LMSMoodle object
type LMSMoodleReconciler struct {
	client.Client
	Scheme                                   *runtime.Scheme
	MoodleGVK, NfsGVK, KeydbGVK, PostgresGVK schema.GroupVersionKind
	lmsMoodleCtx                             LMSMoodleReconcilerContext
}

//+kubebuilder:rbac:groups=lms.krestomat.io,resources=lmsmoodles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lms.krestomat.io,resources=lmsmoodles/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lms.krestomat.io,resources=lmsmoodles/finalizers,verbs=update
//+kubebuilder:rbac:groups=m4e.krestomat.io,resources=moodles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nfs.krestomat.io,resources=ganeshas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=keydb.krestomat.io,resources=keydbs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=postgres.krestomat.io,resources=postgres,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create
//+kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LMSMoodle object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *LMSMoodleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Starting reconcile")

	// Vars
	r.lmsMoodleCtx.name = req.Name

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
	if r.lmsMoodleCtx.desiredState == lmsv1alpha1.SuspendedState {
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
func (r *LMSMoodleReconciler) reconcilePrepare(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.V(1).Info("Reconcile set")

	// set base name for dependant resources
	baseName := truncate(LMSMoodleNamePrefix+r.lmsMoodleCtx.name, TruncateCharactersInName)
	baseNamespace := LMSMoodleNamePrefix + r.lmsMoodleCtx.name
	// if lmsMoodle name already include the prefix, do not use it
	if hasPrefix := strings.HasPrefix(r.lmsMoodleCtx.name, LMSMoodleNamePrefix); hasPrefix {
		baseName = truncate(r.lmsMoodleCtx.name, TruncateCharactersInName)
		baseNamespace = r.lmsMoodleCtx.name
	}
	// set namespace name. It must start with an alphabetic character
	r.lmsMoodleCtx.namespaceName = baseNamespace
	// set network policy base name. It must start with an alphabetic character
	r.lmsMoodleCtx.networkPolicyBaseName = baseName
	// set Moodle name. It must start with an alphabetic character
	r.lmsMoodleCtx.moodleName = baseName
	// set Postgres name. It must start with an alphabetic character
	r.lmsMoodleCtx.postgresName = baseName
	// set NFS Ganesha server name and namespace. It must start with an alphabetic character
	r.lmsMoodleCtx.nfsName = baseName
	// set Keydb name. It must start with an alphabetic character
	r.lmsMoodleCtx.keydbName = baseName
	// lmsMoodle namespace
	r.lmsMoodleCtx.namespace = &corev1.Namespace{}
	r.lmsMoodleCtx.namespace.SetName(r.lmsMoodleCtx.namespaceName)
	// dependant components
	r.lmsMoodleCtx.moodle = newUnstructuredObject(r.MoodleGVK)
	r.lmsMoodleCtx.postgres = newUnstructuredObject(r.PostgresGVK)
	r.lmsMoodleCtx.nfs = newUnstructuredObject(r.NfsGVK)
	r.lmsMoodleCtx.keydb = newUnstructuredObject(r.KeydbGVK)
	// namespaces and names
	r.lmsMoodleCtx.moodle.SetName(r.lmsMoodleCtx.moodleName)
	r.lmsMoodleCtx.moodle.SetNamespace(r.lmsMoodleCtx.namespaceName)

	// Fetch LMSMoodle instance
	r.lmsMoodleCtx.lmsMoodle = newUnstructuredObject(lmsv1alpha1.GroupVersion.WithKind("LMSMoodle"))
	if err := r.Get(ctx, types.NamespacedName{Name: r.lmsMoodleCtx.name}, r.lmsMoodleCtx.lmsMoodle); err != nil {
		log.V(1).Info(err.Error())
		return err
	} else {
		// whether lmsMoodle is marked to be deleted
		r.lmsMoodleCtx.markedToBeDeleted = r.lmsMoodleCtx.lmsMoodle.GetDeletionTimestamp() != nil
	}
	r.lmsMoodleCtx.spec, _, _ = unstructured.NestedMap(r.lmsMoodleCtx.lmsMoodle.UnstructuredContent(), "spec")
	r.lmsMoodleCtx.moodleSpec, r.lmsMoodleCtx.moodleSpecFound, _ = unstructured.NestedMap(r.lmsMoodleCtx.spec, "moodleSpec")
	r.lmsMoodleCtx.postgresSpec, r.lmsMoodleCtx.postgresSpecFound, _ = unstructured.NestedMap(r.lmsMoodleCtx.spec, "postgresSpec")
	r.lmsMoodleCtx.nfsSpec, r.lmsMoodleCtx.nfsSpecFound, _ = unstructured.NestedMap(r.lmsMoodleCtx.spec, "nfsSpec")
	r.lmsMoodleCtx.keydbSpec, r.lmsMoodleCtx.keydbSpecFound, _ = unstructured.NestedMap(r.lmsMoodleCtx.spec, "keydbSpec")
	r.lmsMoodleCtx.lmsMoodleTemplateName, _, _ = unstructured.NestedString(r.lmsMoodleCtx.spec, "lmsMoodleTemplateName")
	r.lmsMoodleCtx.lmsMoodleNetpolOmit, _, _ = unstructured.NestedBool(r.lmsMoodleCtx.spec, "lmsMoodleNetpolOmit")
	r.lmsMoodleCtx.desiredState, _, _ = unstructured.NestedString(r.lmsMoodleCtx.spec, "desiredState")

	// Fetch lmsMoodleTemplate spec
	r.lmsMoodleCtx.lmsMoodleTemplate = newUnstructuredObject(lmsv1alpha1.GroupVersion.WithKind("LMSMoodleTemplate"))
	if err := r.Get(ctx, types.NamespacedName{Name: r.lmsMoodleCtx.lmsMoodleTemplateName}, r.lmsMoodleCtx.lmsMoodleTemplate); err != nil {
		log.Error(err, "LMSMoodleTemplate not found")
		return &LMSMoodleTemplateNotFoundError{r.lmsMoodleCtx.lmsMoodleTemplateName}
	}
	r.lmsMoodleCtx.lmsMoodleTemplateSpec, _, _ = unstructured.NestedMap(r.lmsMoodleCtx.lmsMoodleTemplate.UnstructuredContent(), "spec")
	r.lmsMoodleCtx.lmsMoodleTemplateMoodleSpec, _, _ = unstructured.NestedMap(r.lmsMoodleCtx.lmsMoodleTemplateSpec, "moodleSpec")
	r.lmsMoodleCtx.lmsMoodleTemplatePostgresSpec, r.lmsMoodleCtx.lmsMoodleTemplatePostgresSpecFound, _ = unstructured.NestedMap(r.lmsMoodleCtx.lmsMoodleTemplateSpec, "postgresSpec")
	r.lmsMoodleCtx.lmsMoodleTemplateNfsSpec, r.lmsMoodleCtx.lmsMoodleTemplateNfsSpecFound, _ = unstructured.NestedMap(r.lmsMoodleCtx.lmsMoodleTemplateSpec, "nfsSpec")
	r.lmsMoodleCtx.lmsMoodleTemplateKeydbSpec, r.lmsMoodleCtx.lmsMoodleTemplateKeydbSpecFound, _ = unstructured.NestedMap(r.lmsMoodleCtx.lmsMoodleTemplateSpec, "keydbSpec")

	// set labels
	if err := r.setSiteLabels(ctx); err != nil {
		return err
	}
	r.lmsMoodleCtx.postgres.SetLabels(r.lmsMoodleCtx.lmsMoodle.GetLabels())
	r.lmsMoodleCtx.nfs.SetLabels(r.lmsMoodleCtx.lmsMoodle.GetLabels())
	r.lmsMoodleCtx.keydb.SetLabels(r.lmsMoodleCtx.lmsMoodle.GetLabels())
	r.lmsMoodleCtx.moodle.SetLabels(r.lmsMoodleCtx.lmsMoodle.GetLabels())

	// define default network policy
	r.defineLMSMoodleDefaultNetpol()

	// whether LMSMoodle has dependant components
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
func (r *LMSMoodleReconciler) reconcileFinalize(ctx context.Context) (finalized bool, requeue bool, err error) {
	log := log.FromContext(ctx)
	log.V(1).Info("Reconcile finalizer")

	// Check if LMSMoodle instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	if r.lmsMoodleCtx.markedToBeDeleted {
		// update lmsMoodle state (terminating)
		if requeue, err := r.updateLMSMoodleStatus(ctx); err != nil {
			return false, requeue, err
		}
		if controllerutil.ContainsFinalizer(r.lmsMoodleCtx.lmsMoodle, LMSMoodleFinalizer) {
			// Run finalization logic for SiteFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if requeue, err := r.finalizeLMSMoodle(ctx); err != nil || requeue {
				return false, requeue, err
			}

			// Remove SiteFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(r.lmsMoodleCtx.lmsMoodle, LMSMoodleFinalizer)
			if err := r.Update(ctx, r.lmsMoodleCtx.lmsMoodle); err != nil {
				return false, false, err
			}
		}
		// Finalized
		return true, false, nil
	}
	// Add finalizer for this CR
	if !controllerutil.ContainsFinalizer(r.lmsMoodleCtx.lmsMoodle, LMSMoodleFinalizer) {
		controllerutil.AddFinalizer(r.lmsMoodleCtx.lmsMoodle, LMSMoodleFinalizer)
		if err := r.Update(ctx, r.lmsMoodleCtx.lmsMoodle); err != nil {
			return false, false, err
		}
	}
	return false, false, nil
}

// reconcileSuspend take care of suspend state
func (r *LMSMoodleReconciler) reconcileSuspend(ctx context.Context) (requeue bool, err error) {
	log := log.FromContext(ctx)
	log.V(1).Info("Reconcile persist")

	// Save Moodle spec
	r.lmsMoodleCtx.moodle.Object["spec"] = r.lmsMoodleCtx.combinedMoodleSpec
	// Set suspended
	if err := unstructured.SetNestedField(r.lmsMoodleCtx.combinedMoodleSpec, "suspended", "cr_state"); err != nil {
		return false, err
	}
	// Update lmsMoodle status about Moodle
	r.SetMoodleReadyCondition(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.moodle)
	// Apply Moodle resource
	if err := r.ReconcileApply(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.moodle); err != nil {
		return false, err
	}
	// Whether moodle is suspended
	if suspended := r.isDependantSuspended(ctx, r.lmsMoodleCtx.moodle); !suspended {
		log.Info("Moodle resource is being suspended")
		_, err := r.updateLMSMoodleStatus(ctx)
		return true, err
	}

	// Save Keydb spec
	if r.lmsMoodleCtx.hasKeydb {
		r.lmsMoodleCtx.keydb.Object["spec"] = r.lmsMoodleCtx.combinedKeydbSpec
		// Set suspended
		if err := unstructured.SetNestedField(r.lmsMoodleCtx.combinedKeydbSpec, "suspended", "cr_state"); err != nil {
			return false, err
		}
		// Update LMSMoodle status about Keydb
		r.SetKeydbReadyCondition(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.keydb)
		// Apply Keydb resource
		if err := r.ReconcileApply(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.keydb); err != nil {
			return false, err
		}
		// Whether keydb is suspended
		if suspended := r.isDependantSuspended(ctx, r.lmsMoodleCtx.keydb); !suspended {
			log.Info("Keydb resource is being suspended")
			_, err := r.updateLMSMoodleStatus(ctx)
			return true, err
		}
	}

	// Save NFS Ganesha server spec
	if r.lmsMoodleCtx.hasNfs {
		// Save NFS Ganesha server spec
		r.lmsMoodleCtx.nfs.Object["spec"] = r.lmsMoodleCtx.combinedNfsSpec
		// Set suspended
		if err := unstructured.SetNestedField(r.lmsMoodleCtx.combinedNfsSpec, "suspended", "cr_state"); err != nil {
			return false, err
		}
		// Update LMSMoodle status about NFS Ganesha
		r.SetNfsReadyCondition(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.nfs)
		// Apply NFS Ganesha server resource
		if err := r.ReconcileApply(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.nfs); err != nil {
			return false, err
		}
		// Whether nfs is suspended
		if suspended := r.isDependantSuspended(ctx, r.lmsMoodleCtx.nfs); !suspended {
			log.Info("Nfs resource is being suspended")
			_, err := r.updateLMSMoodleStatus(ctx)
			return true, err
		}
	}

	// Save Postgres spec
	if r.lmsMoodleCtx.hasPostgres {
		r.lmsMoodleCtx.postgres.Object["spec"] = r.lmsMoodleCtx.combinedPostgresSpec
		// Set suspended
		if err := unstructured.SetNestedField(r.lmsMoodleCtx.combinedPostgresSpec, "suspended", "cr_state"); err != nil {
			return false, err
		}
		// Update LMSMoodle status about Postgres
		r.SetPostgresReadyCondition(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.postgres)
		// Apply Postgres resource
		if err := r.ReconcileApply(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.postgres); err != nil {
			return false, err
		}
		// Whether postgres is suspended
		if suspended := r.isDependantSuspended(ctx, r.lmsMoodleCtx.postgres); !suspended {
			log.Info("Postgres resource is being suspended")
			_, err := r.updateLMSMoodleStatus(ctx)
			return true, err
		}
	}

	// lmsMoodle is suspended
	return r.updateLMSMoodleStatus(ctx)
}

// reconcilePresent take care of present state
func (r *LMSMoodleReconciler) reconcilePresent(ctx context.Context) (requeue bool, err error) {
	log := log.FromContext(ctx)
	log.V(1).Info("Reconcile persist")

	// Vars
	moodleReady := false
	nfsReady := !r.lmsMoodleCtx.hasNfs
	keydbReady := !r.lmsMoodleCtx.hasKeydb
	postgresReady := !r.lmsMoodleCtx.hasPostgres

	// Create namespace
	if err := r.ReconcileCreate(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.namespace); err != nil {
		return false, err
	}

	// Whether default network policy should be present
	if r.lmsMoodleCtx.lmsMoodleNetpolOmit {
		if err := r.ReconcileDeleteDependant(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.lmsMoodleDefaultNetpol); client.IgnoreNotFound(err) != nil {
			return false, err
		}
	} else {
		if err := r.ReconcileCreate(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.lmsMoodleDefaultNetpol); err != nil {
			return false, err
		}
	}

	// Save Postgres spec
	if r.lmsMoodleCtx.hasPostgres {
		r.lmsMoodleCtx.postgres.Object["spec"] = r.lmsMoodleCtx.combinedPostgresSpec
		// Update LMSMoodle status about Postgres
		r.SetPostgresReadyCondition(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.postgres)
		// Apply Postgres resource
		if err := r.ReconcileApply(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.postgres); err != nil {
			return false, err
		}
		// check if postgres ready
		if postgresReady, err = getReadyStatus(ctx, r.lmsMoodleCtx.postgres); err != nil {
			return false, err
		}
	}

	// Save Keydb spec
	if r.lmsMoodleCtx.hasKeydb {
		r.lmsMoodleCtx.keydb.Object["spec"] = r.lmsMoodleCtx.combinedKeydbSpec
		// Update LMSMoodle status about Keydb
		r.SetKeydbReadyCondition(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.keydb)
		// Apply Keydb resource
		if err := r.ReconcileApply(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.keydb); err != nil {
			return false, err
		}
		// check if keydb ready
		if keydbReady, err = getReadyStatus(ctx, r.lmsMoodleCtx.keydb); err != nil {
			return false, err
		}
	}

	// Save NFS Ganesha server spec
	if r.lmsMoodleCtx.hasNfs {
		// Save NFS Ganesha server spec
		r.lmsMoodleCtx.nfs.Object["spec"] = r.lmsMoodleCtx.combinedNfsSpec
		// Update LMSMoodle status about NFS Ganesha
		r.SetNfsReadyCondition(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.nfs)
		// Apply NFS Ganesha server resource
		if err := r.ReconcileApply(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.nfs); err != nil {
			return false, err
		}
		// check if nfs ready
		if nfsReady, err = getReadyStatus(ctx, r.lmsMoodleCtx.nfs); err != nil {
			return false, err
		}
	}

	// Wait for postgres to be ready; otherwise requeue
	if !postgresReady {
		log.Info("Postgres is not ready, requeueing...", "Postgres.Name", r.lmsMoodleCtx.postgres.GetName())
		return r.updateLMSMoodleStatus(ctx)
	}
	// Wait for Keydb to be ready; otherwise requeue
	if !keydbReady {
		log.Info("Keydb is not ready, requeueing...", "Keydb.Name", r.lmsMoodleCtx.keydb.GetName())
		return r.updateLMSMoodleStatus(ctx)
	}
	// Wait for NFS Ganesha to be ready; otherwise requeue
	// NFS Ganesha server must be ready in order to mount its export as pvc
	if !nfsReady {
		log.Info("(NFS) Ganesha server is not ready, requeueing...", "Ganesha.Name", r.lmsMoodleCtx.nfs.GetName())
		return r.updateLMSMoodleStatus(ctx)
	}

	// Save Moodle spec
	r.lmsMoodleCtx.moodle.Object["spec"] = r.lmsMoodleCtx.combinedMoodleSpec
	// Update lmsMoodle status about Moodle
	r.SetMoodleReadyCondition(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.moodle)
	// Apply Moodle resource
	if err := r.ReconcileApply(ctx, r.lmsMoodleCtx.lmsMoodle, r.lmsMoodleCtx.moodle); err != nil {
		return false, err
	}
	// check if moodle ready
	if moodleReady, err = getReadyStatus(ctx, r.lmsMoodleCtx.moodle); err != nil {
		return false, err
	}
	// Wait for Moodle to be ready; otherwise requeue
	if !moodleReady {
		log.Info("Moodle is not ready, requeueing...", "Moodle.Name", r.lmsMoodleCtx.moodle.GetName())
		return r.updateLMSMoodleStatus(ctx)
	}

	// lmsMoodle is ready
	return r.updateLMSMoodleStatus(ctx)
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

// lmsMoodlesByLMSMoodleTemplate select lmsmoodles that are using a lmsMoodleTemplate
// It returns a list of reconcile.Request
func (r *LMSMoodleReconciler) lmsMoodlesByLMSMoodleTemplate(ctx context.Context, lmsMoodleTemplate client.Object) []reconcile.Request {
	SiteList := &lmsv1alpha1.LMSMoodleList{}

	// Filter the list of lmsmoodles by the ones using the lmsMoodleTemplate name
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(LMSMoodleTemplateNameIndex, lmsMoodleTemplate.GetName()),
	}

	err := r.Client.List(context.Background(), SiteList, listOps)
	if err != nil {
		return []reconcile.Request{}
	}

	reconcileRequests := make([]reconcile.Request, len(SiteList.Items))
	for i, lmsmoodle := range SiteList.Items {
		reconcileRequests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name: lmsmoodle.Name,
			},
		}
	}
	return reconcileRequests
}

// SetupWithManager sets up the controller with the Manager.
func (r *LMSMoodleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	mgr.GetScheme().AddKnownTypeWithName(r.MoodleGVK, &unstructured.Unstructured{})
	metav1.AddToGroupVersion(mgr.GetScheme(), r.MoodleGVK.GroupVersion())

	// Add spec.lmsMoodleTemplate index
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &lmsv1alpha1.LMSMoodle{}, LMSMoodleTemplateNameIndex, func(obj client.Object) []string {
		lmsMoodleTemplateName := obj.(*lmsv1alpha1.LMSMoodle).Spec.LMSMoodleTemplateName
		if lmsMoodleTemplateName == "" {
			return nil
		}
		return []string{lmsMoodleTemplateName}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&lmsv1alpha1.LMSMoodle{}).
		WithEventFilter(ignoreDeletionPredicate()).
		Owns(newUnstructuredObject(r.MoodleGVK)).
		Owns(newUnstructuredObject(r.NfsGVK)).
		Owns(newUnstructuredObject(r.KeydbGVK)).
		Owns(newUnstructuredObject(r.PostgresGVK)).
		Watches(&lmsv1alpha1.LMSMoodleTemplate{}, handler.EnqueueRequestsFromMapFunc(r.lmsMoodlesByLMSMoodleTemplate)).
		Complete(r)
}
