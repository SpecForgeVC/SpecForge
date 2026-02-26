package ui_roadmap

import (
	"fmt"
)

// ValidateComponentTree checks for deterministic tree structure
func ValidateComponentTree(node ComponentNode) ValidationResult {
	var errors []string

	// Max depth to prevent infinite recursion/denial of service
	const maxDepth = 100

	depth := checkDepth(node, 0)
	if depth > maxDepth {
		errors = append(errors, fmt.Sprintf("Component tree depth exceeds limit of %d", maxDepth))
	}

	validateNodeRecursive(node, &errors)

	return ValidationResult{
		Valid:  len(errors) == 0,
		Errors: errors,
	}
}

func validateNodeRecursive(node ComponentNode, errors *[]string) {
	if node.Type == "" {
		*errors = append(*errors, "Component type cannot be empty")
	}

	// Enforce typed props (at least ensure props is not nil if expected)
	if node.Props == nil {
		node.Props = make(map[string]interface{})
	}

	// Validation logic for bindings
	if node.Binding != "" {
		// Future: Validate binding against BackendBindings map
	}

	for _, child := range node.Children {
		validateNodeRecursive(child, errors)
	}
}

func checkDepth(node ComponentNode, currentDepth int) int {
	if len(node.Children) == 0 {
		return currentDepth
	}
	max := currentDepth
	for _, child := range node.Children {
		d := checkDepth(child, currentDepth+1)
		if d > max {
			max = d
		}
	}
	return max
}
