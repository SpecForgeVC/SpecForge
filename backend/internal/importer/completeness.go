package importer

import (
	"math"
)

type CompletenessScorer interface {
	ScoreSubmission(documents map[string]interface{}) ScoringResult
}

type ScoringResult struct {
	Score                int
	MissingCategories    []string
	UnresolvedReferences []string
	Prompt               string
}

type defaultCompletenessScorer struct{}

func NewCompletenessScorer() CompletenessScorer {
	return &defaultCompletenessScorer{}
}

func (s *defaultCompletenessScorer) ScoreSubmission(documents map[string]interface{}) ScoringResult {
	expectedCategories := []string{
		"project_overview",
		"tech_stack",
		"modules",
		"apis",
		"data_models",
		"contracts",
		"risks",
		"change_sensitivity",
	}

	var missing []string
	var unresolved []string
	totalWeight := 100.0
	earned := 0.0

	categoryWeight := totalWeight / float64(len(expectedCategories))

	for _, cat := range expectedCategories {
		if val, exists := documents[cat]; exists && val != nil {
			// Check if the array or object is empty
			if arr, ok := val.([]interface{}); ok && len(arr) > 0 {
				earned += categoryWeight
			} else if obj, ok := val.(map[string]interface{}); ok && len(obj) > 0 {
				earned += categoryWeight
			} else {
				missing = append(missing, cat)
			}
		} else {
			missing = append(missing, cat)
		}
	}

	// Calculate cross-references
	// E.g., if contract refers to a feature that doesn't exist
	// Since we are mocking the deeply nested schema behavior for the first iteration:
	if len(missing) > 0 {
		unresolved = append(unresolved, "Cannot resolve cross-references for missing categories")
	}

	score := int(math.Round(earned))

	prompt := "All categories are documented."
	if len(missing) > 0 {
		prompt = "Are there any additional undocumented items for the missing categories?"
	}

	return ScoringResult{
		Score:                score,
		MissingCategories:    missing,
		UnresolvedReferences: unresolved,
		Prompt:               prompt,
	}
}
