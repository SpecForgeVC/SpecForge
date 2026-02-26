package ui_roadmap

import (
	"encoding/json"
)

type UIScores struct {
	VisualContractScore     float64 `json:"visual_contract_score"`
	StateCoverageScore      float64 `json:"state_coverage_score"`
	AccessibilityScore      float64 `json:"accessibility_score"`
	BackendBindingScore     float64 `json:"backend_binding_score"`
	TokenComplianceScore    float64 `json:"token_compliance_score"`
	ResponsiveScore         float64 `json:"responsive_score"`
	EdgeCaseScore           float64 `json:"edge_case_score"`
	AggregateReadinessScore float64 `json:"aggregate_readiness_score"`
}

// CalculateScores computes the UI Implementation Readiness Score
func CalculateScores(item *UIRoadmapItem) UIScores {
	scores := UIScores{
		VisualContractScore:  calculateVisualScore(item.ComponentTree),
		StateCoverageScore:   calculateStateScore(item.StateMachine),
		AccessibilityScore:   calculateAccessibilityScore(item.AccessibilitySpec),
		BackendBindingScore:  calculateBindingScore(item.BackendBindings),
		TokenComplianceScore: calculateTokenScore(item.DesignTokensUsed),
		ResponsiveScore:      calculateResponsiveScore(item.ResponsiveSpec),
		EdgeCaseScore:        calculateEdgeCaseScore(item.EdgeCases),
	}

	// Simple average for readiness score
	scores.AggregateReadinessScore = (scores.VisualContractScore +
		scores.StateCoverageScore +
		scores.AccessibilityScore +
		scores.BackendBindingScore +
		scores.TokenComplianceScore +
		scores.ResponsiveScore +
		scores.EdgeCaseScore) / 7.0

	return scores
}

func calculateVisualScore(tree json.RawMessage) float64 {
	if isEmptyJSON(tree) {
		return 0
	}
	var node ComponentNode
	if err := json.Unmarshal(tree, &node); err != nil {
		return 0
	}
	// Score based on recursion depth and node count (higher is more detailed)
	count := countNodes(node)
	if count > 5 {
		return 100
	}
	return float64(count) * 20.0
}

func countNodes(node ComponentNode) int {
	count := 1
	for _, child := range node.Children {
		count += countNodes(child)
	}
	return count
}

func calculateStateScore(sm json.RawMessage) float64 {
	if isEmptyJSON(sm) {
		return 0
	}
	var def StateMachineDef
	if err := json.Unmarshal(sm, &def); err != nil {
		return 0
	}
	mandatory := []string{"idle", "loading", "success", "error", "empty", "disabled"}
	covered := 0
	for _, s := range mandatory {
		if _, exists := def.States[s]; exists {
			covered++
		}
	}
	return (float64(covered) / float64(len(mandatory))) * 100.0
}

func calculateAccessibilityScore(acc json.RawMessage) float64 {
	if isEmptyJSON(acc) {
		return 0
	}
	var spec AccessibilitySpec
	if err := json.Unmarshal(acc, &spec); err != nil {
		return 0
	}
	score := 0.0
	if spec.Role != "" {
		score += 20
	}
	if spec.KeyboardNav != "" {
		score += 20
	}
	if spec.FocusManagement != "" {
		score += 20
	}
	if spec.ScreenReaderText != "" {
		score += 20
	}
	if spec.ContrastCompliance {
		score += 20
	}
	return score
}

func calculateBindingScore(bindings json.RawMessage) float64 {
	if isEmptyJSON(bindings) {
		return 0
	}
	// Placeholder: check if bindings are non-empty and well-formed
	return 100.0
}

func calculateTokenScore(tokens []string) float64 {
	if len(tokens) == 0 {
		return 0
	}
	// Score based on number of tokens used (ensuring at least some design system integration)
	if len(tokens) > 3 {
		return 100
	}
	return float64(len(tokens)) * 25.0
}

func calculateResponsiveScore(resp json.RawMessage) float64 {
	if isEmptyJSON(resp) {
		return 0
	}
	var spec ResponsiveSpec
	if err := json.Unmarshal(resp, &spec); err != nil {
		return 0
	}
	score := 0.0
	if spec.Mobile.Columns > 0 {
		score += 33.3
	}
	if spec.Tablet.Columns > 0 {
		score += 33.3
	}
	if spec.Desktop.Columns > 0 {
		score += 33.4
	}
	return score
}

func calculateEdgeCaseScore(ec json.RawMessage) float64 {
	if isEmptyJSON(ec) {
		return 0
	}
	return 100.0
}
