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
)

const (
	OPERATORNAME string = "kio-operator"
	NFSNAMESPACE string = "rook-nfs"
)

func newUnstructuredObject(gvk schema.GroupVersionKind) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)
	return obj
}

func getEnv(envVar string, defaultVal string) string {
	val, ok := os.LookupEnv(envVar)
	if !ok {
		return defaultVal
	} else {
		return val
	}
}

// ReconcileCreate create resource if it does not exists. Otherwise it does nothing
func (r *SiteReconciler) ReconcileCreate(ctx context.Context, parentObj client.Object, obj client.Object) error {
	log := log.FromContext(ctx)

	log.V(1).Info("Creating resource", "Resource", obj.GetObjectKind())

	if err := r.Get(ctx, types.NamespacedName{Name: obj.GetName(), Namespace: obj.GetNamespace()}, obj); !errors.IsNotFound(err) {
		log.V(1).Info("Resource already exists", "Resource", obj.GetObjectKind())
		return nil
	} else if err != nil {
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

	log.Info("Owner set", "Owner", parentObj.GetUID(), "Dependant", obj.GetObjectKind())
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

	log.Info("Dependant resource deleted", "Dependant", obj.GetObjectKind())
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

	log.Info("Successfully finalized site")
	return false, nil
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
