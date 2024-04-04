package lms

import (
	"context"
	"time"

	lmsv1alpha1 "github.com/krestomatio/lms-moodle-operator/api/lms/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	ReadyConditionType         string = "Ready"
	MoodleReadyConditionType   string = "MoodleReady"
	NfsReadyConditionType      string = "NfsReady"
	KeydbReadyConditionType    string = "KeydbReady"
	PostgresReadyConditionType string = "PostgresReady"
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

// SetSuccessfulReadyCondition set successful ready condition
func (r *LMSMoodleReconciler) SetSuccessfulReadyCondition(ctx context.Context) (changed bool, err error) {
	return r.SetReadyCondition(ctx, "True", lmsv1alpha1.SuccessfulState, "LMSMoodle is ready")
}

// SetFalseReadyCondition set false ready condition
func (r *LMSMoodleReconciler) SetFalseReadyCondition(ctx context.Context, reason string, message string) (changed bool, err error) {
	return r.SetReadyCondition(ctx, "False", reason, message)
}

// SetReadyCondition set ready condition
func (r *LMSMoodleReconciler) SetReadyCondition(ctx context.Context, status string, reason string, message string) (changed bool, err error) {
	log := log.FromContext(ctx)

	readyCondition := map[string]interface{}{
		"type":    "Ready",
		"status":  status,
		"reason":  reason,
		"message": message,
	}

	changed, setConditionErr := SetCondition(r.lmsMoodleCtx.lmsMoodle, readyCondition)
	if setConditionErr != nil {
		log.Error(setConditionErr, "unable to set ready condition")
		return false, setConditionErr
	}

	return changed, err
}

// SetMoodleReadyCondition set ready condition depending on ready status of Moodle
// and returns bool flag which indicates ready condition status of that dependant object
func (r *LMSMoodleReconciler) SetMoodleReadyCondition(ctx context.Context, parentObj *unstructured.Unstructured, dependantObj *unstructured.Unstructured) (changed bool, status bool) {
	return r.SetConditionFromDependantByType(ctx, parentObj, dependantObj, MoodleReadyConditionType, ReadyConditionType)
}

// SetPostgresReadyCondition set ready condition depending on ready status of Postgres
// and returns bool flag which indicates ready condition status of that dependant object
func (r *LMSMoodleReconciler) SetPostgresReadyCondition(ctx context.Context, parentObj *unstructured.Unstructured, dependantObj *unstructured.Unstructured) (changed bool, status bool) {
	return r.SetConditionFromDependantByType(ctx, parentObj, dependantObj, PostgresReadyConditionType, ReadyConditionType)
}

// SetNfsReadyCondition set ready condition depending on ready status of NFS Ganesha
// and returns bool flag which indicates ready condition status of that dependant object
func (r *LMSMoodleReconciler) SetNfsReadyCondition(ctx context.Context, parentObj *unstructured.Unstructured, dependantObj *unstructured.Unstructured) (changed bool, status bool) {
	return r.SetConditionFromDependantByType(ctx, parentObj, dependantObj, NfsReadyConditionType, ReadyConditionType)
}

// SetKeydbReadyCondition set ready condition depending on ready status of Keydb
// and returns bool flag which indicates ready condition status of that dependant object
func (r *LMSMoodleReconciler) SetKeydbReadyCondition(ctx context.Context, parentObj *unstructured.Unstructured, dependantObj *unstructured.Unstructured) (changed bool, status bool) {
	return r.SetConditionFromDependantByType(ctx, parentObj, dependantObj, KeydbReadyConditionType, ReadyConditionType)
}

// SetConditionFromDependantByType set a condition in a parent object from
// a ready type condition of a dependant object based on its status
// Returns bool flag which indicates ready condition status of the dependant object
func (r *LMSMoodleReconciler) SetConditionFromDependantByType(ctx context.Context, parentObj *unstructured.Unstructured, dependantObj *unstructured.Unstructured, parentConditionType string, dependantConditionType string) (changed bool, status bool) {
	log := log.FromContext(ctx)
	dependantConditionStatus := false

	// get dependant condition by type
	dependantCondition, dependantConditionFound, dependantConditionErr := getConditionByType(dependantObj, dependantConditionType)

	if dependantConditionErr != nil {
		log.Error(dependantConditionErr, "unable to get dependant condition")
		return false, false
	}

	if !dependantConditionFound {
		log.V(1).Info("dependant condition not found")
		return false, false
	}

	// rename type to set parent condition
	dependantCondition["type"] = parentConditionType
	_, setConditionErr := SetCondition(parentObj, dependantCondition)
	if setConditionErr != nil {
		log.Error(setConditionErr, "unable to set condition based on dependant")
		return false, false
	}

	// if depedant condition status is false, set parent as not ready either
	if dependantCondition["status"] != "True" {
		parentReadyCondition := map[string]interface{}{
			"type":    "Ready",
			"status":  "False",
			"reason":  "DependantNotReady",
			"message": "Dependant is not ready",
		}
		parentReadyConditionchanged, hasSetParentReadyConditionErr := SetCondition(parentObj, parentReadyCondition)
		if hasSetParentReadyConditionErr != nil {
			log.Error(hasSetParentReadyConditionErr, "unable to set ready condition based on dependant condition status")
			return false, false
		}
		changed = parentReadyConditionchanged
	} else {
		dependantConditionStatus = true
	}

	return changed, dependantConditionStatus
}

// getConditionByType returns a condition by type from a unstructure object,
// a bool flag which indicates whether condition exists, and
// any error getting the condition
func getConditionByType(objU *unstructured.Unstructured, conditionType string) (map[string]interface{}, bool, error) {
	// look for conditions slice in unstructured object
	conditions, conditionsFound, conditionsErr := unstructured.NestedSlice(objU.Object, "status", "conditions")
	if conditionsErr != nil {
		return make(map[string]interface{}), false, conditionsErr
	}

	if !conditionsFound {
		return make(map[string]interface{}), false, nil
	}

	// look for condition type
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
	conditions, _, err := unstructured.NestedSlice(unstructuredObj.Object, "status", "conditions")
	if err != nil {
		return false, err
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
