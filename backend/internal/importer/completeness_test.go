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
		if len(res.MissingCategories) != 6 {
			t.Errorf("Expected 6 missing categories, got %d", len(res.MissingCategories))
		}
	})

	t.Run("Partial Payload", func(t *testing.T) {
		docs := map[string]interface{}{
			"contracts": []interface{}{map[string]interface{}{"id": "c1"}},
			"features":  []interface{}{map[string]interface{}{"id": "f1"}},
		}
		res := scorer.ScoreSubmission(docs)

		if res.Score <= 0 || res.Score >= 100 {
			t.Errorf("Expected partial score, got %d", res.Score)
		}
		if len(res.MissingCategories) != 4 {
			t.Errorf("Expected 4 missing categories, got %d", len(res.MissingCategories))
		}
	})

	t.Run("Complete Payload", func(t *testing.T) {
		docs := map[string]interface{}{
			"contracts":    []interface{}{map[string]interface{}{"id": "c1"}},
			"variables":    []interface{}{map[string]interface{}{"id": "v1"}},
			"features":     []interface{}{map[string]interface{}{"id": "f1"}},
			"architecture": []interface{}{map[string]interface{}{"id": "a1"}},
			"security":     []interface{}{map[string]interface{}{"id": "s1"}},
			"integrations": []interface{}{map[string]interface{}{"id": "i1"}},
		}
		res := scorer.ScoreSubmission(docs)

		// With 6 items mapped perfectly, math.Round of 100.0 is 100
		if res.Score != 100 {
			t.Errorf("Expected perfect score 100, got %d", res.Score)
		}
		if len(res.MissingCategories) != 0 {
			t.Errorf("Expected 0 missing categories, got %v", res.MissingCategories)
		}
	})
}
