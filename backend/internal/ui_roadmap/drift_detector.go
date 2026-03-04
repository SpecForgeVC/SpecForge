package ui_roadmap

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/SpecForgeVC/SpecForge/internal/domain"
)

type DriftIssue struct {
	Type        string `json:"type"` // MISSING_CONTRACT | FIELD_MISMATCH | FORMAT_ERROR
	Field       string `json:"field"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // WARNING | ERROR
}

// contractEntry holds a contract and the endpoint paths found in its schema.
type contractEntry struct {
	contract domain.ContractDefinition
	paths    []string
}

// DetectUIDrift identifies mismatches between UI specs and backend contracts.
// backendContracts is the full list of ContractDefinitions available for the item's project.
func DetectUIDrift(item *UIRoadmapItem, backendContracts []domain.ContractDefinition) []DriftIssue {
	var issues []DriftIssue

	if len(item.BackendBindings) == 0 {
		return issues
	}

	var bindings []BackendBinding
	if err := json.Unmarshal(item.BackendBindings, &bindings); err != nil {
		return append(issues, DriftIssue{
			Type:        "FORMAT_ERROR",
			Field:       "bindings",
			Description: "Invalid binding schema: " + err.Error(),
			Severity:    "ERROR",
		})
	}

	if len(bindings) == 0 {
		return issues
	}

	// Build lookup of contracts with their known endpoint paths
	entries := buildContractEntries(backendContracts)

	for _, binding := range bindings {
		matched := findMatchingContract(binding.Endpoint, binding.Method, entries)

		if matched == nil {
			issues = append(issues, DriftIssue{
				Type:        "MISSING_CONTRACT",
				Field:       binding.Endpoint,
				Description: fmt.Sprintf("Binding endpoint '%s %s' has no matching backend contract", binding.Method, binding.Endpoint),
				Severity:    "ERROR",
			})
			continue
		}

		// Validate input_map keys exist in the contract's InputSchema
		for uiField := range binding.InputMap {
			if _, exists := matched.InputSchema[uiField]; !exists {
				// Try inside "properties" key (JSON Schema pattern)
				props, _ := matched.InputSchema["properties"].(map[string]interface{})
				if _, existsNested := props[uiField]; !existsNested {
					issues = append(issues, DriftIssue{
						Type:        "FIELD_MISMATCH",
						Field:       binding.Endpoint + "." + uiField,
						Description: fmt.Sprintf("Input field '%s' in binding for '%s' is not in the contract input schema", uiField, binding.Endpoint),
						Severity:    "WARNING",
					})
				}
			}
		}

		// Validate output_map keys exist in the contract's OutputSchema
		for uiField := range binding.OutputMap {
			if _, exists := matched.OutputSchema[uiField]; !exists {
				props, _ := matched.OutputSchema["properties"].(map[string]interface{})
				if _, existsNested := props[uiField]; !existsNested {
					issues = append(issues, DriftIssue{
						Type:        "FIELD_MISMATCH",
						Field:       binding.Endpoint + "." + uiField,
						Description: fmt.Sprintf("Output field '%s' in binding for '%s' is not in the contract output schema", uiField, binding.Endpoint),
						Severity:    "WARNING",
					})
				}
			}
		}
	}

	// Deduplicate issues
	return deduplicateIssues(issues)
}

func deduplicateIssues(issues []DriftIssue) []DriftIssue {
	seen := make(map[string]bool)
	unique := make([]DriftIssue, 0, len(issues))
	for _, issue := range issues {
		key := fmt.Sprintf("%s:%s:%s", issue.Type, issue.Field, issue.Description)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, issue)
		}
	}
	return unique
}

// buildContractEntries extracts endpoint path hints from each contract's schema.
func buildContractEntries(contracts []domain.ContractDefinition) []contractEntry {
	entries := make([]contractEntry, 0, len(contracts))
	for _, c := range contracts {
		var paths []string
		// Check common schema keys for explicit path/endpoint declarations
		for _, key := range []string{"path", "endpoint"} {
			if v, ok := c.InputSchema[key].(string); ok {
				paths = append(paths, v)
			}
		}
		// OpenAPI-style "paths" map
		if pathsObj, ok := c.InputSchema["paths"].(map[string]interface{}); ok {
			for k := range pathsObj {
				paths = append(paths, k)
			}
		}
		entries = append(entries, contractEntry{contract: c, paths: paths})
	}
	return entries
}

// findMatchingContract searches contract entries for one matching the given endpoint+method.
func findMatchingContract(endpoint, method string, entries []contractEntry) *domain.ContractDefinition {
	endpoint = strings.ToLower(strings.TrimSpace(endpoint))

	for i, entry := range entries {
		// Direct path match
		for _, p := range entry.paths {
			if strings.EqualFold(p, endpoint) {
				return &entries[i].contract
			}
		}
		// Fallback: endpoint string appears anywhere in the schema value content
		schemaStr := schemaToString(entry.contract.InputSchema) + schemaToString(entry.contract.OutputSchema)
		if strings.Contains(strings.ToLower(schemaStr), endpoint) {
			return &entries[i].contract
		}
	}
	return nil
}

// schemaToString converts a map schema to a searchable string representation.
func schemaToString(schema map[string]interface{}) string {
	if len(schema) == 0 {
		return ""
	}
	b, _ := json.Marshal(schema)
	return string(b)
}
