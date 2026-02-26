package drift

import (
	"testing"
)

func TestDetectDrift_RequiredFieldRemoval(t *testing.T) {
	baseline := OpenAPIDocument{
		Components: Components{
			Schemas: map[string]Schema{
				"User": {
					Type: "object",
					Properties: map[string]*Schema{
						"id":   {Type: "string"},
						"name": {Type: "string"},
					},
					Required: []string{"id", "name"},
				},
			},
		},
	}

	proposed := OpenAPIDocument{
		Components: Components{
			Schemas: map[string]Schema{
				"User": {
					Type: "object",
					Properties: map[string]*Schema{
						"id":   {Type: "string"},
						"name": {Type: "string"},
					},
					Required: []string{"id"},
				},
			},
		},
	}

	report := DetectDrift(DriftInput{Baseline: baseline, Proposed: proposed})

	if report.CriticalChanges == 0 {
		t.Errorf("Expected critical changes for required field removal, got 0")
	}

	found := false
	for _, item := range report.Items {
		if item.Type == RequiredFieldRemoved && item.Severity == Critical {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Did not find RequiredFieldRemoved with Critical severity")
	}
}

func TestDetectDrift_EnumNarrowing(t *testing.T) {
	baseline := OpenAPIDocument{
		Components: Components{
			Schemas: map[string]Schema{
				"Status": {
					Type: "string",
					Enum: []any{"pending", "active", "deleted"},
				},
			},
		},
	}

	proposed := OpenAPIDocument{
		Components: Components{
			Schemas: map[string]Schema{
				"Status": {
					Type: "string",
					Enum: []any{"pending", "active"},
				},
			},
		},
	}

	report := DetectDrift(DriftInput{Baseline: baseline, Proposed: proposed})

	if report.CriticalChanges == 0 {
		t.Errorf("Expected critical changes for enum narrowing, got 0")
	}

	found := false
	for _, item := range report.Items {
		if item.Type == EnumValueRemoved && item.Severity == Critical {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Did not find EnumValueRemoved with Critical severity")
	}
}

func TestDetectDrift_PathRemoval(t *testing.T) {
	baseline := OpenAPIDocument{
		Paths: map[string]PathItem{
			"/users": {
				Operations: map[string]Operation{
					"get": {},
				},
			},
		},
	}

	proposed := OpenAPIDocument{
		Paths: map[string]PathItem{},
	}

	report := DetectDrift(DriftInput{Baseline: baseline, Proposed: proposed})

	if report.CriticalChanges == 0 {
		t.Errorf("Expected critical changes for path removal, got 0")
	}

	found := false
	for _, item := range report.Items {
		if item.Type == PathRemoved && item.Severity == Critical {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Did not find PathRemoved with Critical severity")
	}
}

func TestDetectDrift_ConstraintTightening(t *testing.T) {
	maxBase := 100
	maxProp := 50
	baseline := OpenAPIDocument{
		Components: Components{
			Schemas: map[string]Schema{
				"Input": {
					Type:      "string",
					MaxLength: &maxBase,
				},
			},
		},
	}

	proposed := OpenAPIDocument{
		Components: Components{
			Schemas: map[string]Schema{
				"Input": {
					Type:      "string",
					MaxLength: &maxProp,
				},
			},
		},
	}

	report := DetectDrift(DriftInput{Baseline: baseline, Proposed: proposed})

	if report.BreakingChanges == 0 {
		t.Errorf("Expected breaking changes for maxLength tightening, got 0")
	}

	found := false
	for _, item := range report.Items {
		if item.Type == ConstraintTightened && item.Severity == Breaking {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Did not find ConstraintTightened with Breaking severity")
	}
}

func TestDriftPolicy_Evaluate(t *testing.T) {
	policy := DriftPolicy{
		BlockOnBreaking: true,
		BlockOnCritical: true,
	}

	report := DriftReport{
		CriticalChanges: 1,
	}

	policy.Evaluate(&report)
	if !report.Blocked {
		t.Errorf("Expected report to be blocked due to critical changes")
	}

	report2 := DriftReport{
		Warnings: 5,
	}
	policy.Evaluate(&report2)
	if report2.Blocked {
		t.Errorf("Expected report NOT to be blocked for warnings only")
	}
}
