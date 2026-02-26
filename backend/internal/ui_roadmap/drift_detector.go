package ui_roadmap

import (
	"encoding/json"
	"fmt"
)

type DriftIssue struct {
	Type        string `json:"type"` // MISMATCH | MISSING_STATE | DEPRECATED_ENDPOINT
	Field       string `json:"field"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // WARNING | ERROR
}

// DetectUIDrift identifies mismatches between UI specs and backend contracts
func DetectUIDrift(item *UIRoadmapItem, backendContracts []json.RawMessage) []DriftIssue {
	var issues []DriftIssue

	var bindings []BackendBinding
	if err := json.Unmarshal(item.BackendBindings, &bindings); err != nil {
		return append(issues, DriftIssue{
			Type:        "FORMAT_ERROR",
			Field:       "bindings",
			Description: "Invalid binding schema",
			Severity:    "ERROR",
		})
	}

	for _, binding := range bindings {
		// Verify endpoint existence and field surface
		endpointFound := false
		for range backendContracts {
			// In a real scenario, we'd unmarshal the contract and compare types/fields
			// Placeholder logic for now
			if binding.Endpoint == "exists" { // Simulating check
				endpointFound = true
			}
		}

		if !endpointFound {
			issues = append(issues, DriftIssue{
				Type:        "MISSING_CONTRACT",
				Field:       binding.Endpoint,
				Description: fmt.Sprintf("Binding refers to non-existent endpoint: %s", binding.Endpoint),
				Severity:    "ERROR",
			})
		}
	}

	// State machine completeness check (already in validation, but here as "drift" if requirements change)

	return issues
}
