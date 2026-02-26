package ui_roadmap

// ValidateAccessibility ensures UI meets governance standards
func ValidateAccessibility(spec AccessibilitySpec) ValidationResult {
	var errors []string

	if spec.Role == "" {
		errors = append(errors, "ARIA role must be defined for the top-level element")
	}
	if spec.KeyboardNav == "" {
		errors = append(errors, "Keyboard navigation logic must be documented")
	}
	if spec.FocusManagement == "" {
		errors = append(errors, "Focus management behavior must be defined")
	}
	if spec.ScreenReaderText == "" {
		errors = append(errors, "Screen reader behavior (alt text, labels) must be defined")
	}
	if !spec.ContrastCompliance {
		errors = append(errors, "Color contrast rules must be enabled and compliant")
	}

	return ValidationResult{
		Valid:  len(errors) == 0,
		Errors: errors,
	}
}
