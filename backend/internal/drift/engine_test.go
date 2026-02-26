package drift

import (
	"testing"
)

func TestDiffEngine_Compare(t *testing.T) {
	engine := NewDiffEngine()

	oldSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id":   map[string]interface{}{"type": "string"},
			"name": map[string]interface{}{"type": "string"},
		},
	}

	newSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id":   map[string]interface{}{"type": "string"},
			"name": map[string]interface{}{"type": "integer"}, // Type change (breaking)
			"age":  map[string]interface{}{"type": "integer"}, // Addition (non-breaking)
		},
	}

	diffs, err := engine.Compare(oldSchema, newSchema)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}

	foundBreaking := false
	for _, d := range diffs {
		if d.Path == "properties.name.type" && d.IsBreaking {
			foundBreaking = true
		}
	}

	if !foundBreaking {
		t.Error("Expected breaking change for properties.name.type not found")
	}
}

func TestDiffEngine_Compare_Complex(t *testing.T) {
	engine := NewDiffEngine()

	oldSchema := map[string]interface{}{
		"type":     "object",
		"required": []interface{}{"id"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{"type": "string"},
			"tags": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}

	newSchema := map[string]interface{}{
		"type":     "object",
		"required": []interface{}{"id", "tags"}, // Adding required field (breaking)
		"properties": map[string]interface{}{
			"id": map[string]interface{}{"type": "string"},
			"tags": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "integer", // Type change in array items (breaking)
				},
			},
		},
	}

	diffs, err := engine.Compare(oldSchema, newSchema)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}

	foundRequiredChange := false
	foundArrayItemChange := false

	for _, d := range diffs {
		if d.Path == "required" && d.IsBreaking {
			foundRequiredChange = true
		}
		if d.Path == "properties.tags.items.type" && d.IsBreaking {
			foundArrayItemChange = true
		}
	}

	if !foundRequiredChange {
		t.Error("Expected breaking change for required fields not found")
	}
	if !foundArrayItemChange {
		t.Error("Expected breaking change for array items type not found")
	}
}
