package ui_roadmap

import (
	"fmt"
)

// ValidateStateMachine ensures completeness of the UI state machine
func ValidateStateMachine(sm StateMachineDef) ValidationResult {
	var errors []string

	mandatoryStates := []string{"idle", "loading", "success", "error", "empty", "disabled"}
	for _, state := range mandatoryStates {
		config, exists := sm.States[state]
		if !exists {
			errors = append(errors, fmt.Sprintf("Mandatory state '%s' is missing", state))
			continue
		}

		if config.VisualChanges == "" {
			errors = append(errors, fmt.Sprintf("State '%s' must define visual changes", state))
		}
		if config.InteractionChanges == "" {
			errors = append(errors, fmt.Sprintf("State '%s' must define interaction changes", state))
		}

		// Interaction rules
		if state == "loading" {
			// In loading state, we usually expect interactions to be disabled or limited
			// This could be a warning or a strict rule
		}
		if state == "error" {
			if config.Messaging == "" {
				errors = append(errors, "Error state must define messaging behavior")
			}
		}
	}

	return ValidationResult{
		Valid:  len(errors) == 0,
		Errors: errors,
	}
}
