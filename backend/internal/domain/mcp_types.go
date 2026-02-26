package domain

// --- MCP Tool Structs ---

type ToolInputInitProjectImport struct {
	ProjectID      string   `json:"project_id"`
	LanguageStack  []string `json:"language_stack"`
	RepositoryType string   `json:"repository_type"`
	EstimatedSize  string   `json:"estimated_size"`
	ImportMode     string   `json:"import_mode"`
}

type ToolOutputInitProjectImport struct {
	Status               string                 `json:"status"`
	RequiredDocuments    []string               `json:"required_documents"`
	DocumentSchema       map[string]interface{} `json:"document_schema"`
	ScaffoldInstructions string                 `json:"scaffold_instructions"`
	SubmissionProtocol   map[string]interface{} `json:"submission_protocol"`
	SessionID            string                 `json:"session_id"`
}

type ToolInputSubmitProjectSnapshot struct {
	ProjectID       string                 `json:"project_id"`
	SnapshotVersion string                 `json:"snapshot_version"`
	SnapshotPayload map[string]interface{} `json:"snapshot_payload"`
	FinalSubmission bool                   `json:"final_submission"`
}

type ToolOutputSubmitProjectSnapshot struct {
	Status                       string   `json:"status"`
	CompletenessScore            int      `json:"completeness_score"`
	MissingCategories            []string `json:"missing_categories"`
	UnresolvedReferences         []string `json:"unresolved_references"`
	SelfAssessmentPrompt         string   `json:"self_assessment_prompt"`
	RequiresAdditionalSubmission bool     `json:"requires_additional_submission"`
	CatalogState                 string   `json:"catalog_state,omitempty"`
	ProjectImportStatus          string   `json:"project_import_status,omitempty"`
}

type ToolInputGetImportAlignmentRules struct {
	ProjectID string `json:"project_id"`
}

type ToolOutputGetImportAlignmentRules struct {
	StrictRules            []string               `json:"strict_rules"`
	ForbiddenActions       []string               `json:"forbidden_actions"`
	RequiredSnapshotPolicy SnapshotPolicy         `json:"required_snapshot_policy"`
	AlignmentConstraints   map[string]interface{} `json:"alignment_constraints"`
}

type SnapshotPolicy struct {
	MustCreateSnapshotBeforeChanges bool `json:"must_create_snapshot_before_changes"`
	MustPostSnapshotAfterChanges    bool `json:"must_post_snapshot_after_changes"`
}

type ToolInputSubmitPostImportSnapshot struct {
	ProjectID       string          `json:"project_id"`
	SnapshotVersion string          `json:"snapshot_version"`
	SnapshotPayload ProjectSnapshot `json:"snapshot_payload"`
}

type ToolOutputSubmitPostImportSnapshot struct {
	DriftScore               int     `json:"drift_score"`
	ContractBreakageDetected bool    `json:"contract_breakage_detected"`
	RiskDelta                float64 `json:"risk_delta"`
	AlignmentDelta           float64 `json:"alignment_delta"`
	Recommendation           string  `json:"recommendation"`
}

type ToolInputFinalizeProjectImport struct {
	ProjectID string `json:"project_id"`
}

type ToolOutputFinalizeProjectImport struct {
	Status      string `json:"status"`
	RedirectURL string `json:"redirect_url,omitempty"`
}
