package drift

import (
	"fmt"
)

// DetectDrift is the main entry point for the drift detection engine.
func DetectDrift(input DriftInput) DriftReport {
	var items []DriftItem

	// 1. Compare Paths
	items = append(items, ComparePaths(input.Baseline.Paths, input.Proposed.Paths)...)

	// 2. Compare Components (Schemas only for now as requested)
	items = append(items, compareComponents(input.Baseline.Components, input.Proposed.Components)...)

	return NewDriftReport(items)
}

func compareComponents(baseline, proposed Components) []DriftItem {
	var items []DriftItem

	// Compare Schemas in Components
	for name, baseSchema := range baseline.Schemas {
		propSchema, exists := proposed.Schemas[name]
		if !exists {
			items = append(items, newDriftItem(RequiredFieldRemoved, "/components/schemas", name, nil, fmt.Sprintf("Component schema '%s' removed", name)))
			continue
		}
		items = append(items, CompareSchema(baseSchema, propSchema, "/components/schemas/"+name)...)
	}

	for name := range proposed.Schemas {
		if _, exists := baseline.Schemas[name]; !exists {
			items = append(items, newDriftItem(FieldAdded, "/components/schemas", nil, name, fmt.Sprintf("New component schema '%s' added", name)))
		}
	}

	return items
}
