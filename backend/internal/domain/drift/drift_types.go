package drift

// DriftType represents the specific kind of drift detected.
type DriftType string

const (
	// Path-Level Drift
	PathRemoved   DriftType = "PATH_REMOVED"
	MethodRemoved DriftType = "METHOD_REMOVED"
	NewPath       DriftType = "NEW_PATH"
	NewMethod     DriftType = "NEW_METHOD"

	// Request/Response Schema Drift
	RequiredFieldRemoved DriftType = "REQUIRED_FIELD_REMOVED"
	RequiredFieldAdded   DriftType = "REQUIRED_FIELD_ADDED"
	FieldTypeChanged     DriftType = "FIELD_TYPE_CHANGED"
	EnumValueRemoved     DriftType = "ENUM_VALUE_REMOVED"
	EnumValueAdded       DriftType = "ENUM_VALUE_ADDED"
	ResponseRemoved      DriftType = "RESPONSE_REMOVED"
	ErrorSchemaChanged   DriftType = "ERROR_SCHEMA_CHANGED"

	// Constraints Drift
	ConstraintTightened DriftType = "CONSTRAINT_TIGHTENED"
	ConstraintLoosened  DriftType = "CONSTRAINT_LOOSENED"
	FieldAdded          DriftType = "FIELD_ADDED"

	// Security Drift
	SecuritySchemeRemoved DriftType = "SECURITY_SCHEME_REMOVED"
	SecurityLoosened      DriftType = "SECURITY_LOOSENED"
	SecurityTightened     DriftType = "SECURITY_TIGHTENED"

	// Metadata Drift
	MetadataChanged DriftType = "METADATA_CHANGED"
)

// DriftInput contains the baseline and proposed documents for comparison.
type DriftInput struct {
	Baseline OpenAPIDocument
	Proposed OpenAPIDocument
}

// OpenAPIDocument is a pre-parsed, normalized, and canonicalized OpenAPI 3.1 document.
// Refs must be resolved before reaching the drift engine.
type OpenAPIDocument struct {
	Paths      map[string]PathItem
	Components Components
}

// PathItem represents the operations available on a single path.
type PathItem struct {
	Operations map[string]Operation
}

// Operation represents an HTTP method on a path.
type Operation struct {
	ID          string
	Summary     string
	Description string
	RequestBody *RequestBody
	Responses   map[string]Response
	Security    []SecurityRequirement
}

// RequestBody represents the content of a request.
type RequestBody struct {
	Content map[string]MediaType
}

// Response represents an HTTP response.
type Response struct {
	Content map[string]MediaType
}

// MediaType represents the schema and examples for a media type.
type MediaType struct {
	Schema Schema
}

// Schema represents an OpenAPI 3.1 Schema object.
type Schema struct {
	Type                 string
	Format               string
	Properties           map[string]*Schema
	Items                *Schema
	Required             []string
	Enum                 []any
	Nullable             bool
	AdditionalProperties *Schema
	MinLength            *int
	MaxLength            *int
	MinItems             *int
	MaxItems             *int
	Minimum              *float64
	Maximum              *float64
	OneOf                []*Schema
	AnyOf                []*Schema
	AllOf                []*Schema
	Description          string
	Default              any
}

// Components represents the components section of an OpenAPI document.
type Components struct {
	Schemas         map[string]Schema
	SecuritySchemes map[string]SecurityScheme
}

// SecurityRequirement represents a security requirement for an operation.
type SecurityRequirement map[string][]string

// SecurityScheme represents a security scheme definition.
type SecurityScheme struct {
	Type string
	In   string
	Name string
}
