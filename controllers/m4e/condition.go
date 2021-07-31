package m4e

import (
	"context"
	"errors"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// FindConditionUnstructuredByType returns first Condition with given conditionType
// along with bool flag which indicates if the Condition is found or not
func FindConditionUnstructuredByType(conditionsUnstructured []interface{}, conditionType string) (map[string]interface{}, bool) {
	for _, conditionUnstructured := range conditionsUnstructured {
		if conditionAsMap, ok := conditionUnstructured.(map[string]interface{}); ok {
			if typeString, ok := conditionAsMap["type"]; ok && typeString == conditionType {
				return conditionAsMap, true
			}
		}
	}
	return make(map[string]interface{}), false
}

// SetReadyCondition set ready condition
func (r *SiteReconciler) SetReadyCondition(ctx context.Context, parentObj *unstructured.Unstructured) {
	log := log.FromContext(ctx)

	parentReadyCondition := map[string]interface{}{
		"type":    "Ready",
		"status":  "True",
		"reason":  "SiteReady",
		"message": "Site is ready",
	}

	hasSetCondition, setConditionErr := SetCondition(parentObj, parentReadyCondition)
	if setConditionErr != nil {
		log.Error(setConditionErr, "unable to set ready condition based on dependant condition status")
	}

	// update parent status conditions if conditions changed
	if hasSetCondition {
		if err := r.Status().Update(ctx, parentObj); err != nil {
			log.Error(err, "unable to update resource status")
		}
	}
}

// SetM4eReadyCondition set ready condition depending on ready status of M4e Server
// and returns bool flag which indicates ready condition status of that dependant object
// along with any error setting the condition
func (r *SiteReconciler) SetM4eReadyCondition(ctx context.Context, parentObj *unstructured.Unstructured, dependantObj *unstructured.Unstructured) bool {
	return r.SetConditionFromDependantByType(ctx, parentObj, dependantObj, "M4eReady", "Ready", "True")
}

// SetNfsReadyCondition set ready condition depending on ready status of NFS Server
// and returns bool flag which indicates ready condition status of that dependant object
// along with any error setting the condition
func (r *SiteReconciler) SetNfsReadyCondition(ctx context.Context, parentObj *unstructured.Unstructured, dependantObj *unstructured.Unstructured) bool {
	return r.SetConditionFromDependantByType(ctx, parentObj, dependantObj, "NfsReady", "Ready", "True")
}

// SetConditionFromDependantByType set ready condition depending on ready status of dependant object filter by type
// and returns bool flag which indicates ready condition status of that dependant object
// along with any error setting the condition
func (r *SiteReconciler) SetConditionFromDependantByType(ctx context.Context, parentObj *unstructured.Unstructured, dependantObj *unstructured.Unstructured, parentConditionType string, dependantConditionType string, stringToEvalConditionStatus string) bool {
	log := log.FromContext(ctx)
	parentReady := false

	// get dependant condition by type
	dependantCondition, dependantConditionFound, dependantConditionErr := getConditionFromDependantByType(dependantObj, dependantConditionType)

	if dependantConditionErr != nil {
		log.Error(dependantConditionErr, "unable to get dependant condition")
		return false
	}

	if !dependantConditionFound {
		log.V(1).Info("dependant condition not found")
		return false
	}

	// rename type to set parent condition
	dependantCondition["type"] = parentConditionType
	hasSetCondition, setConditionErr := SetCondition(parentObj, dependantCondition)
	if setConditionErr != nil {
		log.Error(setConditionErr, "unable to set condition based on dependant")
		return false
	}

	// if depedant condition status is false, set parent as not ready either
	if dependantCondition["status"] != "True" {
		parentReadyCondition := map[string]interface{}{
			"type":    "Ready",
			"status":  "False",
			"reason":  "DependantNotReady",
			"message": "Dependant is not ready",
		}
		hasSetParentReadyCondition, hasSetParentReadyConditionErr := SetCondition(parentObj, parentReadyCondition)
		if hasSetParentReadyConditionErr != nil {
			log.Error(hasSetParentReadyConditionErr, "unable to set ready condition based on dependant condition status")
			return false
		}
		hasSetCondition = hasSetParentReadyCondition
	} else {
		parentReady = true
	}

	// update parent status conditions
	if hasSetCondition {
		if err := r.Status().Update(ctx, parentObj); err != nil {
			log.Error(err, "unable to update resource status")
			return false
		}
	}

	return parentReady
}

// getConditionFromDependantByType returns condition from dependant resourse by type,
// a bool flag which indicates whether condition exists, and
// any error getting the condition
func getConditionFromDependantByType(dependant *unstructured.Unstructured, conditionType string) (map[string]interface{}, bool, error) {
	// look for conditions slice in unstructured object
	conditions, conditionsFound, conditionsErr := unstructured.NestedSlice(dependant.Object, "status", "conditions")
	if conditionsErr != nil {
		return make(map[string]interface{}), false, conditionsErr
	}

	if !conditionsFound {
		return make(map[string]interface{}), false, errors.New("dependant conditions not found")
	}

	// look for ready condition type
	condition, conditionFound := FindConditionUnstructuredByType(conditions, conditionType)

	return condition, conditionFound, nil
}

// SetCondition update or append a condition if needed
// It returns a bool flag if condition was appended or updated, and
// any error
func SetCondition(unstructuredObj *unstructured.Unstructured, condition map[string]interface{}) (bool, error) {
	appendCondition := true
	updateConditions := false

	// set lastTransitionTime if not set
	if condition["lastTransitionTime"] == nil {
		condition["lastTransitionTime"] = metav1.Now().UTC().Format(time.RFC3339)
	}

	// get resource conditions
	conditions, conditionsFound, err := unstructured.NestedSlice(unstructuredObj.Object, "status", "conditions")
	if err != nil {
		return false, err
	}

	if !conditionsFound {
		return false, errors.New("conditions not found")
	}

	// update or append, depending if condition is present or not
	for i, item := range conditions {
		if conditionObj, ok := item.(map[string]interface{}); ok {
			if conditionType, ok := conditionObj["type"].(string); ok && conditionType == condition["type"] {
				// do not append
				appendCondition = false
				// update if transitioned
				if transitionedCondition, hasTransitioned := HasTransitioned(conditionObj, condition); hasTransitioned {
					conditions[i] = transitionedCondition
					updateConditions = true
				}
			}
		}
	}
	// append
	if appendCondition {
		conditions = append(conditions, condition)
		updateConditions = true
	}

	// set conditions slice in status if condition set
	if updateConditions {
		if err := unstructured.SetNestedField(unstructuredObj.Object, conditions, "status", "conditions"); err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

// HasTransitioned returns the version of the condition and a bool flag. The values depends on whether
// its state has transitioned or not
func HasTransitioned(oldCondition map[string]interface{}, newCondition map[string]interface{}) (map[string]interface{}, bool) {
	if oldCondition["status"] != newCondition["status"] || oldCondition["reason"] != newCondition["reason"] || oldCondition["message"] != newCondition["message"] {
		return newCondition, true
	}
	return oldCondition, false
}
