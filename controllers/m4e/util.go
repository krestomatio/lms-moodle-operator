package m4e

import (
	"context"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
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

func (r *SiteReconciler) reconcileCreate(ctx context.Context, parentObj client.Object, obj client.Object) error {
	log := log.FromContext(ctx)

	if err := r.Get(ctx, types.NamespacedName{Name: obj.GetName(), Namespace: obj.GetNamespace()}, obj); err != nil && errors.IsNotFound(err) {
		log.V(1).Info("Create resource", "Resource.Kind", obj.GetObjectKind(), "Resource.Name", obj.GetName())

		// Set resource ownership
		if err := r.reconcileSetOwner(ctx, parentObj, obj); err != nil {
			return err
		}

		if err := r.Create(ctx, obj); err != nil {
			log.Error(err, "Failed to create resource", "Resource.Kind", obj.GetObjectKind(), "Resource.Name", obj.GetName())
			return err
		}
	} else if err != nil {
		log.Error(err, "Failed to get resource", "Resource.Kind", obj.GetObjectKind(), "Resource.Name", obj.GetName())
		return err
	}
	return nil
}

func (r *SiteReconciler) reconcileApply(ctx context.Context, parentObj client.Object, obj client.Object) error {
	log := log.FromContext(ctx)

	// Set resource ownership
	if err := r.reconcileSetOwner(ctx, parentObj, obj); err != nil {
		return err
	}

	// Apply resource
	log.V(1).Info("Applying changes", "Resource.Kind", obj.GetObjectKind(), "Resource.Name", obj.GetName())
	force := true
	if err := r.Patch(ctx, obj, client.Apply, &client.PatchOptions{Force: &force, FieldManager: OPERATORNAME}); err != nil {
		log.Error(err, "Failed to apply changes", "Resource.Kind", obj.GetObjectKind(), "Resource.Namespace", obj.GetNamespace(), "Resource.Name", obj.GetName())
		return err
	}

	return nil
}

func (r *SiteReconciler) reconcileSetOwner(ctx context.Context, parentObj client.Object, obj client.Object) error {
	log := log.FromContext(ctx)

	// Set owner reference
	log.V(1).Info("Setting owner", "Owner", parentObj.GetUID(), "Dependant.Kind", obj.GetObjectKind(), "Dependant.Name", obj.GetName())
	if err := ctrl.SetControllerReference(parentObj, obj, r.Scheme); err != nil {
		log.Error(err, "Failed to set owner reference")
		return err
	}

	return nil
}
