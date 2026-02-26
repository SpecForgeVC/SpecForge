package ui_roadmap

import (
	"encoding/json"
)

type ExportBundle struct {
	JSONSpec      []byte
	LLMPrompt     string
	FigmaMake     string
	ClaudeFigma   string
	StorybookSpec string
}

// GenerateExportBundle compiles all build artifacts for a UI Roadmap Item
func GenerateExportBundle(item *UIRoadmapItem) (ExportBundle, error) {
	var bundle ExportBundle

	// 1. JSON Specification
	spec, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return bundle, err
	}
	bundle.JSONSpec = spec

	// 2. LLM Implementation Prompt
	bundle.LLMPrompt = GenerateLLMPrompt(item)

	// 3. Figma Make Prompt
	bundle.FigmaMake = GenerateFigmaMakePrompt(item)

	// 4. Claude in Figma Prompt
	bundle.ClaudeFigma = GenerateClaudeInFigmaPrompt(item)

	// 5. Storybook Scaffold
	bundle.StorybookSpec = GenerateStorybookScaffold(item)

	return bundle, nil
}
