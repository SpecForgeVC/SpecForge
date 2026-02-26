package drift

import (
	"fmt"
)

// ComparePaths compares the paths and operations between baseline and proposed documents.
func ComparePaths(baseline, proposed map[string]PathItem) []DriftItem {
	var items []DriftItem

	for path, baseItem := range baseline {
		propItem, exists := proposed[path]
		if !exists {
			items = append(items, newDriftItem(PathRemoved, "/paths", path, nil, fmt.Sprintf("Path '%s' removed", path)))
			continue
		}
		items = append(items, compareOperations(path, baseItem.Operations, propItem.Operations)...)
	}

	for path := range proposed {
		if _, exists := baseline[path]; !exists {
			items = append(items, newDriftItem(NewPath, "/paths", nil, path, fmt.Sprintf("New path '%s' added", path)))
		}
	}

	return items
}

func compareOperations(path string, baseline, proposed map[string]Operation) []DriftItem {
	var items []DriftItem
	location := "/paths:" + path

	for method, baseOp := range baseline {
		propOp, exists := proposed[method]
		if !exists {
			items = append(items, newDriftItem(MethodRemoved, location, method, nil, fmt.Sprintf("Method '%s' removed from path '%s'", method, path)))
			continue
		}
		items = append(items, compareOperation(location+":"+method, baseOp, propOp)...)
	}

	for method := range proposed {
		if _, exists := baseline[method]; !exists {
			items = append(items, newDriftItem(NewMethod, location, nil, method, fmt.Sprintf("New method '%s' added to path '%s'", method, path)))
		}
	}

	return items
}

func compareOperation(location string, baseline, proposed Operation) []DriftItem {
	var items []DriftItem

	// 1. Request Body Comparison
	if baseline.RequestBody != nil && proposed.RequestBody != nil {
		items = append(items, compareRequestBody(location+":requestBody", *baseline.RequestBody, *proposed.RequestBody)...)
	} else if baseline.RequestBody != nil && proposed.RequestBody == nil {
		items = append(items, newDriftItem(RequiredFieldRemoved, location, "requestBody", nil, "Request body removed"))
	}

	// 2. Responses Comparison
	for code, baseResp := range baseline.Responses {
		propResp, exists := proposed.Responses[code]
		if !exists {
			items = append(items, newDriftItem(ResponseRemoved, location+":responses", code, nil, fmt.Sprintf("Response '%s' removed", code)))
			continue
		}
		items = append(items, compareResponse(location+":responses:"+code, baseResp, propResp)...)
	}

	// 3. Security Comparison
	items = append(items, compareSecurity(location+":security", baseline.Security, proposed.Security)...)

	return items
}

func compareRequestBody(location string, baseline, proposed RequestBody) []DriftItem {
	var items []DriftItem
	for contentType, baseMT := range baseline.Content {
		propMT, exists := proposed.Content[contentType]
		if !exists {
			items = append(items, newDriftItem(RequiredFieldRemoved, location, contentType, nil, fmt.Sprintf("Content type '%s' removed from request body", contentType)))
			continue
		}
		items = append(items, CompareSchema(baseMT.Schema, propMT.Schema, location+":"+contentType)...)
	}
	return items
}

func compareResponse(location string, baseline, proposed Response) []DriftItem {
	var items []DriftItem
	for contentType, baseMT := range baseline.Content {
		propMT, exists := proposed.Content[contentType]
		if !exists {
			items = append(items, newDriftItem(RequiredFieldRemoved, location, contentType, nil, fmt.Sprintf("Content type '%s' removed from response", contentType)))
			continue
		}
		items = append(items, CompareSchema(baseMT.Schema, propMT.Schema, location+":"+contentType)...)
	}
	return items
}

func compareSecurity(location string, baseline, proposed []SecurityRequirement) []DriftItem {
	var items []DriftItem
	// Simplified security comparison: check if security requirements were loosened or tightened.
	if len(baseline) > 0 && len(proposed) == 0 {
		items = append(items, newDriftItem(SecurityLoosened, location, baseline, nil, "Security requirements removed"))
	} else if len(baseline) == 0 && len(proposed) > 0 {
		items = append(items, newDriftItem(SecurityTightened, location, nil, proposed, "Security requirements added"))
	} else if len(baseline) != len(proposed) {
		// More complex logic could be added here to compare specific requirements
		items = append(items, newDriftItem(SecurityLoosened, location, baseline, proposed, "Security requirements changed"))
	}
	return items
}
