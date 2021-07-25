package m4e

import (
	"context"
	"os"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	OPERATORNAME string = "kio-operator"
	NFSNAMESPACE string = "rook-nfs"
)

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}

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

func (r *SiteReconciler) reconcileApply(ctx context.Context, parentObj client.Object, obj client.Object) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Set owner reference
	log.V(1).Info("Setting owner", "Owner", parentObj.GetUID(), "Dependant.Kind", obj.GetObjectKind(), "Dependant.Namespace", obj.GetNamespace(), "Dependant.Name", obj.GetName())
	if err := ctrl.SetControllerReference(parentObj, obj, r.Scheme); err != nil {
		log.Error(err, "Failed to set owner reference")
		return ctrl.Result{}, err
	}

	// Apply resource
	log.V(1).Info("Applying changes", "Resource.Kind", obj.GetObjectKind(), "Resource.Namespace", obj.GetNamespace(), "Resource.Name", obj.GetName())
	force := true
	if err := r.Patch(ctx, obj, client.Apply, &client.PatchOptions{Force: &force, FieldManager: OPERATORNAME}); err != nil {
		log.Error(err, "Failed to apply changes", "Resource.Kind", obj.GetObjectKind(), "Resource.Namespace", obj.GetNamespace(), "Resource.Name", obj.GetName())
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
