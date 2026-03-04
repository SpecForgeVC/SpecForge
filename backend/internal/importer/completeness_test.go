package importer

import (
	"testing"
)

func TestCompletenessScoring(t *testing.T) {
	scorer := NewCompletenessScorer()

	t.Run("Empty Payload", func(t *testing.T) {
		docs := map[string]interface{}{}
		res := scorer.ScoreSubmission(docs)

		if res.Score != 0 {
			t.Errorf("Expected score 0, got %d", res.Score)
		}
		if len(res.MissingCategories) != 8 {
			t.Errorf("Expected 8 missing categories, got %d", len(res.MissingCategories))
		}
	})

	t.Run("Partial Payload", func(t *testing.T) {
		docs := map[string]interface{}{
			"contracts":   []interface{}{map[string]interface{}{"id": "c1"}},
			"data_models": []interface{}{map[string]interface{}{"id": "d1"}},
		}
		res := scorer.ScoreSubmission(docs)

		if res.Score <= 0 || res.Score >= 100 {
			t.Errorf("Expected partial score, got %d", res.Score)
		}
		if len(res.MissingCategories) != 6 {
			t.Errorf("Expected 6 missing categories, got %d", len(res.MissingCategories))
		}
	})

	t.Run("Complete Payload", func(t *testing.T) {
		docs := map[string]interface{}{
			"project_overview":   map[string]interface{}{"id": "po"},
			"tech_stack":         []interface{}{map[string]interface{}{"id": "ts"}},
			"modules":            []interface{}{map[string]interface{}{"id": "m1"}},
			"apis":               []interface{}{map[string]interface{}{"id": "a1"}},
			"data_models":        []interface{}{map[string]interface{}{"id": "dm1"}},
			"contracts":          []interface{}{map[string]interface{}{"id": "c1"}},
			"risks":              []interface{}{map[string]interface{}{"id": "r1"}},
			"change_sensitivity": map[string]interface{}{"id": "cs1"},
		}
		res := scorer.ScoreSubmission(docs)

		if res.Score != 100 {
			t.Errorf("Expected perfect score 100, got %d", res.Score)
		}
		if len(res.MissingCategories) != 0 {
			t.Errorf("Expected 0 missing categories, got %v", res.MissingCategories)
		}
	})
}
