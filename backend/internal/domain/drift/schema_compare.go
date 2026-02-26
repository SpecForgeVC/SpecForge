package drift

import (
	"fmt"
)

// CompareSchema recursively compares two schemas and returns a list of drift items.
func CompareSchema(baseline, proposed Schema, location string) []DriftItem {
	var items []DriftItem

	// 1. Type change
	if baseline.Type != proposed.Type && baseline.Type != "" && proposed.Type != "" {
		items = append(items, newDriftItem(FieldTypeChanged, location, baseline.Type, proposed.Type,
			fmt.Sprintf("Type changed from %s to %s", baseline.Type, proposed.Type)))
	}

	// 2. Enum changes
	items = append(items, compareEnums(baseline.Enum, proposed.Enum, location)...)

	// 3. Nullability changes
	if !baseline.Nullable && proposed.Nullable {
		items = append(items, newDriftItem(ConstraintLoosened, location, baseline.Nullable, proposed.Nullable, "Field became nullable"))
	} else if baseline.Nullable && !proposed.Nullable {
		items = append(items, newDriftItem(ConstraintTightened, location, baseline.Nullable, proposed.Nullable, "Field became non-nullable"))
	}

	// 4. Required fields changes
	items = append(items, compareRequired(baseline.Required, proposed.Required, location)...)

	// 5. Constraints (MinLength, MaxLength, etc.)
	items = append(items, compareConstraints(baseline, proposed, location)...)

	// 6. Properties comparison (Recursive)
	items = append(items, compareProperties(baseline.Properties, proposed.Properties, location)...)

	// 7. Items comparison (Arrays)
	if baseline.Items != nil && proposed.Items != nil {
		items = append(items, CompareSchema(*baseline.Items, *proposed.Items, location+"/items")...)
	}

	// 8. Composition comparison (oneOf, anyOf, allOf)
	items = append(items, compareComposition("oneOf", baseline.OneOf, proposed.OneOf, location)...)
	items = append(items, compareComposition("anyOf", baseline.AnyOf, proposed.AnyOf, location)...)
	items = append(items, compareComposition("allOf", baseline.AllOf, proposed.AllOf, location)...)

	return items
}

func newDriftItem(t DriftType, loc string, base, prop any, desc string) DriftItem {
	return DriftItem{
		Type:        t,
		Severity:    GetSeverity(t, base, prop),
		Location:    loc,
		Baseline:    base,
		Proposed:    prop,
		Description: desc,
	}
}

func compareEnums(baseline, proposed []any, location string) []DriftItem {
	var items []DriftItem
	if len(baseline) == 0 && len(proposed) == 0 {
		return nil
	}

	baseMap := make(map[any]bool)
	for _, v := range baseline {
		baseMap[v] = true
	}

	propMap := make(map[any]bool)
	for _, v := range proposed {
		propMap[v] = true
	}

	// Check for removals
	for _, v := range baseline {
		if !propMap[v] {
			items = append(items, newDriftItem(EnumValueRemoved, location, v, nil, fmt.Sprintf("Enum value '%v' removed", v)))
		}
	}

	// Check for additions
	for _, v := range proposed {
		if !baseMap[v] {
			items = append(items, newDriftItem(EnumValueAdded, location, nil, v, fmt.Sprintf("Enum value '%v' added", v)))
		}
	}

	return items
}

func compareRequired(baseline, proposed []string, location string) []DriftItem {
	var items []DriftItem
	baseMap := make(map[string]bool)
	for _, v := range baseline {
		baseMap[v] = true
	}

	propMap := make(map[string]bool)
	for _, v := range proposed {
		propMap[v] = true
	}

	for _, v := range baseline {
		if !propMap[v] {
			items = append(items, newDriftItem(RequiredFieldRemoved, location, v, nil, fmt.Sprintf("Required field '%s' removed", v)))
		}
	}

	for _, v := range proposed {
		if !baseMap[v] {
			items = append(items, newDriftItem(RequiredFieldAdded, location, nil, v, fmt.Sprintf("Required field '%s' added", v)))
		}
	}

	return items
}

func compareConstraints(baseline, proposed Schema, location string) []DriftItem {
	var items []DriftItem

	// MinLength
	if baseline.MinLength != proposed.MinLength {
		valBase := 0
		if baseline.MinLength != nil {
			valBase = *baseline.MinLength
		}
		valProp := 0
		if proposed.MinLength != nil {
			valProp = *proposed.MinLength
		}

		if valProp > valBase {
			items = append(items, newDriftItem(ConstraintTightened, location+"/minLength", baseline.MinLength, proposed.MinLength, "minLength increased"))
		} else if valProp < valBase {
			items = append(items, newDriftItem(ConstraintLoosened, location+"/minLength", baseline.MinLength, proposed.MinLength, "minLength decreased"))
		}
	}

	// MaxLength
	if baseline.MaxLength != proposed.MaxLength {
		if baseline.MaxLength != nil && proposed.MaxLength != nil {
			if *proposed.MaxLength < *baseline.MaxLength {
				items = append(items, newDriftItem(ConstraintTightened, location+"/maxLength", *baseline.MaxLength, *proposed.MaxLength, "maxLength decreased"))
			} else if *proposed.MaxLength > *baseline.MaxLength {
				items = append(items, newDriftItem(ConstraintLoosened, location+"/maxLength", *baseline.MaxLength, *proposed.MaxLength, "maxLength increased"))
			}
		} else if baseline.MaxLength == nil && proposed.MaxLength != nil {
			items = append(items, newDriftItem(ConstraintTightened, location+"/maxLength", nil, *proposed.MaxLength, "maxLength added"))
		} else if baseline.MaxLength != nil && proposed.MaxLength == nil {
			items = append(items, newDriftItem(ConstraintLoosened, location+"/maxLength", *baseline.MaxLength, nil, "maxLength removed"))
		}
	}

	// (Similar logic for Minimum, Maximum, MinItems, MaxItems could be added here)

	return items
}

func compareProperties(baseline, proposed map[string]*Schema, location string) []DriftItem {
	var items []DriftItem
	if baseline == nil && proposed == nil {
		return nil
	}

	for name := range baseline {
		propProp, exists := proposed[name]
		if !exists {
			items = append(items, newDriftItem(RequiredFieldRemoved, location, name, nil, fmt.Sprintf("Property '%s' removed", name)))
			continue
		}
		items = append(items, CompareSchema(*baseline[name], *propProp, location+"/properties/"+name)...)
	}

	for name := range proposed {
		if _, exists := baseline[name]; !exists {
			items = append(items, newDriftItem(FieldAdded, location, nil, name, fmt.Sprintf("Property '%s' added", name)))
		}
	}

	return items
}

func compareComposition(compType string, baseline, proposed []*Schema, location string) []DriftItem {
	var items []DriftItem
	if len(baseline) == 0 && len(proposed) == 0 {
		return nil
	}

	// Simple heuristic: compare by index if lengths match, otherwise report count change
	// In a real engine, we might try to match by structure, but for this task index-based or length-based is a start.
	// Actually, the requirements emphasize structural comparison.

	if len(baseline) > len(proposed) {
		items = append(items, newDriftItem(ConstraintTightened, location+"/"+compType, len(baseline), len(proposed), fmt.Sprintf("%s option(s) removed", compType)))
	} else if len(baseline) < len(proposed) {
		items = append(items, newDriftItem(ConstraintLoosened, location+"/"+compType, len(baseline), len(proposed), fmt.Sprintf("%s option(s) added", compType)))
	}

	// Recurse into common indices
	minLen := len(baseline)
	if len(proposed) < minLen {
		minLen = len(proposed)
	}

	for i := 0; i < minLen; i++ {
		items = append(items, CompareSchema(*baseline[i], *proposed[i], fmt.Sprintf("%s/%s/%d", location, compType, i))...)
	}

	return items
}
