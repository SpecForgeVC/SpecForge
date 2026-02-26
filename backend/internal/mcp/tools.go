package mcp

// Tool represents an MCP tool definition
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// GetToolDefinitions returns the set of tools supported by the RAE MCP server
func GetToolDefinitions() []Tool {
	return []Tool{
		{
			Name:        "create_snapshot",
			Description: "Initiates a pre-implementation environment extraction request.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_id":      map[string]interface{}{"type": "string", "description": "The unique identifier for the project."},
					"roadmap_item_id": map[string]interface{}{"type": "string", "description": "The unique identifier for the roadmap item."},
					"mode":            map[string]interface{}{"type": "string", "description": "The mode of extraction (e.g., pre-implementation)."},
				},
				"required": []string{"project_id", "roadmap_item_id", "mode"},
			},
		},
		{
			Name:        "post_snapshot",
			Description: "Submits structured environment snapshot back to SpecForge.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"snapshot_id":          map[string]interface{}{"type": "string", "description": "The ID of the snapshot initiated via create_snapshot."},
					"environment_snapshot": map[string]interface{}{"type": "object", "description": "The structured environment data."},
				},
				"required": []string{"snapshot_id", "environment_snapshot"},
			},
		},
		{
			Name:        "get_snapshot_status",
			Description: "Returns the current lifecycle state of a snapshot.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"snapshot_id": map[string]interface{}{"type": "string", "description": "The ID of the snapshot to query."},
				},
				"required": []string{"snapshot_id"},
			},
		},
		{
			Name:        "list_active_snapshots",
			Description: "Lists all snapshots currently in an active state (not completed or failed).",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_id": map[string]interface{}{"type": "string", "description": "Optional filter by project ID."},
				},
			},
		},
		{
			Name:        "init_project_import",
			Description: "Initializes a structured import of an existing project. This tool sets up the extraction requirements and alignment rules for the initial reality snapshot. It returns a session ID to be used in subsequent snapshot submissions.",
			InputSchema: ToolInputInitProjectImportSchema,
		},
		{
			Name:        "submit_project_snapshot",
			Description: "Submits a comprehensive reality snapshot of the project, including architecture, contracts, variables, security, and tests. SpecForge will analyze this snapshot against its intelligence models to identify alignment gaps and architectural risks.",
			InputSchema: ToolInputSubmitProjectSnapshotSchema,
		},
		{
			Name:        "get_import_alignment_rules",
			Description: "Retrieves the specific alignment and enforcement rules applicable to the project import session. These rules define the strictness of contract locking and architectural isolation required during the ingestion phase.",
			InputSchema: ToolInputGetImportAlignmentRulesSchema,
		},
		{
			Name:        "submit_post_import_snapshot",
			Description: "Submits a final reality snapshot after changes have been applied during the import process. SpecForge uses this to detect drift and verify that the final state aligns with the intended roadmap and alignment rules.",
			InputSchema: ToolInputSubmitPostImportSnapshotSchema,
		},
		{
			Name:        "finalize_project_import",
			Description: "Signals that the project cataloguing and documentation is 100% complete and requests finalization of the import process. This will trigger a transition to the project dashboard.",
			InputSchema: ToolInputFinalizeProjectImportSchema,
		},
		{
			Name:        "help",
			Description: "Returns tool descriptions, required usage order, and JSON schema examples.",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
	}
}

var ToolInputInitProjectImportSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"project_id":      map[string]interface{}{"type": "string", "description": "The unique identifier for the project."},
		"language_stack":  map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
		"repository_type": map[string]interface{}{"type": "string", "enum": []string{"monorepo", "polyrepo", "single"}},
		"estimated_size":  map[string]interface{}{"type": "string", "enum": []string{"small", "medium", "large"}},
		"import_mode":     map[string]interface{}{"type": "string", "enum": []string{"light", "full"}},
	},
	"required": []string{"project_id", "language_stack", "repository_type", "estimated_size", "import_mode"},
}

var ToolInputSubmitProjectSnapshotSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"project_id":       map[string]interface{}{"type": "string", "description": "The unique identifier for the project."},
		"snapshot_version": map[string]interface{}{"type": "string", "description": "The version of the snapshot."},
		"snapshot_payload": map[string]interface{}{
			"type":        "object",
			"description": "The structured project model. MUST use exactly the defined schema with no deviations. This is critical for system ingestion.",
			"properties": map[string]interface{}{
				"project_overview": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name":                 map[string]interface{}{"type": "string"},
						"description":          map[string]interface{}{"type": "string"},
						"domain":               map[string]interface{}{"type": "string"},
						"primary_language":     map[string]interface{}{"type": "string"},
						"architecture_pattern": map[string]interface{}{"type": "string"},
					},
					"required": []string{"name", "description", "domain", "primary_language", "architecture_pattern"},
				},
				"tech_stack": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"languages":      map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
						"frameworks":     map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
						"databases":      map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
						"infrastructure": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
						"build_tools":    map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
					},
				},
				"modules": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name":               map[string]interface{}{"type": "string"},
							"description":        map[string]interface{}{"type": "string"},
							"responsibilities":   map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
							"risk_level":         map[string]interface{}{"type": "string", "enum": []string{"LOW", "MEDIUM", "HIGH"}},
							"change_sensitivity": map[string]interface{}{"type": "string", "enum": []string{"LOW", "MEDIUM", "HIGH"}},
							"dependencies":       map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
						},
						"required": []string{"name", "description", "responsibilities", "risk_level"},
					},
				},
				"apis": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"endpoint":        map[string]interface{}{"type": "string"},
							"method":          map[string]interface{}{"type": "string", "enum": []string{"GET", "POST", "PUT", "PATCH", "DELETE"}},
							"auth_type":       map[string]interface{}{"type": "string"},
							"request_schema":  map[string]interface{}{"type": "object"},
							"response_schema": map[string]interface{}{"type": "object"},
						},
						"required": []string{"endpoint", "method"},
					},
				},
				"data_models": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name": map[string]interface{}{"type": "string"},
							"relationships": map[string]interface{}{
								"type": "array",
								"items": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"target": map[string]interface{}{"type": "string"},
										"type":   map[string]interface{}{"type": "string"},
									},
								},
							},
						},
						"required": []string{"name"},
					},
				},
				"contracts": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name":            map[string]interface{}{"type": "string"},
							"contract_type":   map[string]interface{}{"type": "string", "enum": []string{"REST", "GRAPHQL", "EVENT", "INTERNAL_FUNCTION"}},
							"source_module":   map[string]interface{}{"type": "string"},
							"stability_score": map[string]interface{}{"type": "number"},
						},
						"required": []string{"name", "contract_type", "source_module"},
					},
				},
				"risks": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"area":        map[string]interface{}{"type": "string"},
							"severity":    map[string]interface{}{"type": "string", "enum": []string{"LOW", "MEDIUM", "HIGH", "CRITICAL"}},
							"description": map[string]interface{}{"type": "string"},
						},
						"required": []string{"area", "severity", "description"},
					},
				},
				"change_sensitivity": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"module":      map[string]interface{}{"type": "string"},
							"sensitivity": map[string]interface{}{"type": "string", "enum": []string{"LOW", "MEDIUM", "HIGH"}},
							"reason":      map[string]interface{}{"type": "string"},
						},
						"required": []string{"module", "sensitivity", "reason"},
					},
				},
			},
			"required": []string{"project_overview", "tech_stack", "modules", "apis", "data_models", "contracts", "risks", "change_sensitivity"},
		},
		"final_submission": map[string]interface{}{
			"type":        "boolean",
			"description": "Set to true only when you have documented all requested categories and are ready to finalize the import.",
		},
	},
	"required": []string{"project_id", "snapshot_version", "snapshot_payload"},
}

var ToolInputGetImportAlignmentRulesSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"project_id": map[string]interface{}{"type": "string", "description": "The unique identifier for the project."},
	},
	"required": []string{"project_id"},
}

var ToolInputSubmitPostImportSnapshotSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"project_id":       map[string]interface{}{"type": "string", "description": "The unique identifier for the project."},
		"snapshot_version": map[string]interface{}{"type": "string", "description": "The version of the snapshot."},
		"snapshot_payload": map[string]interface{}{"type": "object", "description": "The structured project model."},
	},
	"required": []string{"project_id", "snapshot_version", "snapshot_payload"},
}

var ToolInputFinalizeProjectImportSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"project_id": map[string]interface{}{"type": "string", "description": "The unique identifier for the project."},
	},
	"required": []string{"project_id"},
}
