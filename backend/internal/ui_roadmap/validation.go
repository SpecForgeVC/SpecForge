package ui_roadmap

import (
	"encoding/json"
)

type ValidationResult struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors"`
}

// ValidateUIRoadmapItem performs strict validation on a UI Roadmap Item
func ValidateUIRoadmapItem(item *UIRoadmapItem) ValidationResult {
	var errors []string

	// Mandatory Section Checks
	if isEmptyJSON(item.LayoutDefinition) {
		errors = append(errors, "Layout definition is mandatory")
	}
	if isEmptyJSON(item.ComponentTree) {
		errors = append(errors, "Component tree is mandatory")
	}
	if isEmptyJSON(item.StateMachine) {
		errors = append(errors, "State machine is mandatory")
	}
	if isEmptyJSON(item.BackendBindings) {
		errors = append(errors, "Backend bindings are mandatory")
	}
	if isEmptyJSON(item.ValidationRules) {
		errors = append(errors, "Validation rules are mandatory")
	}
	if isEmptyJSON(item.AccessibilitySpec) {
		errors = append(errors, "Accessibility rules are mandatory")
	}
	if isEmptyJSON(item.ResponsiveSpec) {
		errors = append(errors, "Responsive rules are mandatory")
	}
	if isEmptyJSON(item.TestScenarios) {
		errors = append(errors, "Test scenarios are mandatory")
	}

	// Deep Validation: Component Tree
	var tree ComponentNode
	if err := json.Unmarshal(item.ComponentTree, &tree); err != nil {
		errors = append(errors, "Invalid component tree format: "+err.Error())
	} else {
		treeRes := ValidateComponentTree(tree)
		if !treeRes.Valid {
			errors = append(errors, treeRes.Errors...)
		}
	}

	// Deep Validation: State Machine
	var sm StateMachineDef
	if err := json.Unmarshal(item.StateMachine, &sm); err != nil {
		errors = append(errors, "Invalid state machine format: "+err.Error())
	} else {
		smRes := ValidateStateMachine(sm)
		if !smRes.Valid {
			errors = append(errors, smRes.Errors...)
		}
	}

	// Deep Validation: Accessibility
	var acc AccessibilitySpec
	if err := json.Unmarshal(item.AccessibilitySpec, &acc); err != nil {
		errors = append(errors, "Invalid accessibility spec format: "+err.Error())
	} else {
		accRes := ValidateAccessibility(acc)
		if !accRes.Valid {
			errors = append(errors, accRes.Errors...)
		}
	}

	return ValidationResult{
		Valid:  len(errors) == 0,
		Errors: errors,
	}
}

func isEmptyJSON(j json.RawMessage) bool {
	if len(j) == 0 {
		return true
	}
	s := string(j)
	return s == "{}" || s == "[]" || s == "null"
}
