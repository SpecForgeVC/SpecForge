package ui_roadmap

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// UIRoadmapItem represents a production-grade UI feature specification
type UIRoadmapItem struct {
	ID                uuid.UUID       `json:"id"`
	ProjectID         uuid.UUID       `json:"project_id"`
	LinkedFeatureID   *uuid.UUID      `json:"linked_feature_id,omitempty"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	UserPersona       string          `json:"user_persona"`
	UseCase           string          `json:"use_case"`
	ScreenType        string          `json:"screen_type"` // page | modal | component | layout
	LayoutDefinition  json.RawMessage `json:"layout_definition"`
	ComponentTree     json.RawMessage `json:"component_tree"`
	StateMachine      json.RawMessage `json:"state_machine"`
	BackendBindings   json.RawMessage `json:"backend_bindings"`
	AccessibilitySpec json.RawMessage `json:"accessibility_spec"`
	ResponsiveSpec    json.RawMessage `json:"responsive_spec"`
	ValidationRules   json.RawMessage `json:"validation_rules"`
	AnimationRules    json.RawMessage `json:"animation_rules"`
	DesignTokensUsed  pq.StringArray  `json:"design_tokens_used"`
	EdgeCases         json.RawMessage `json:"edge_cases"`
	TestScenarios     json.RawMessage `json:"test_scenarios"`
	IntelligenceScore float64         `json:"intelligence_score"`
	Version           int             `json:"version"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// ComponentNode represents a deterministic node in the UI component tree
type ComponentNode struct {
	Type       string                 `json:"type"`
	Props      map[string]interface{} `json:"props"`
	Binding    string                 `json:"binding,omitempty"`
	Validation []string               `json:"validation,omitempty"`
	Children   []ComponentNode        `json:"children,omitempty"`
}

// StateMachineDef defines the behavior of the UI across interaction states
type StateMachineDef struct {
	States map[string]StateConfig `json:"states"`
}

// StateConfig defines visual and interaction behavior for a specific state
type StateConfig struct {
	VisualChanges      string `json:"visual_changes"`
	InteractionChanges string `json:"interaction_changes"`
	Messaging          string `json:"messaging"`
}

// BackendBinding defines a link to a backend API contract
type BackendBinding struct {
	Endpoint   string            `json:"endpoint"`
	Method     string            `json:"method"`
	InputMap   map[string]string `json:"input_map"`
	OutputMap  map[string]string `json:"output_map"`
	ErrorShape string            `json:"error_shape"`
}

// AccessibilitySpec defines behavior for screen readers and keyboard navigation
type AccessibilitySpec struct {
	Role               string `json:"role"`
	KeyboardNav        string `json:"keyboard_nav"`
	FocusManagement    string `json:"focus_management"`
	ScreenReaderText   string `json:"screen_reader_text"`
	ContrastCompliance bool   `json:"contrast_compliance"`
}

// ResponsiveSpec defines layout mutations across breakpoints
type ResponsiveSpec struct {
	Mobile  LayoutMutation `json:"mobile"`
	Tablet  LayoutMutation `json:"tablet"`
	Desktop LayoutMutation `json:"desktop"`
}

// LayoutMutation defines specific changes to the layout for a breakpoint
type LayoutMutation struct {
	Columns   int                    `json:"columns"`
	Hidden    []string               `json:"hidden"`
	Overrides map[string]interface{} `json:"overrides"`
}
