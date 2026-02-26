package drift

// DriftSeverity represents the severity level of a detected drift.
type DriftSeverity string

const (
	Info     DriftSeverity = "INFO"
	Warning  DriftSeverity = "WARNING"
	Breaking DriftSeverity = "BREAKING"
	Critical DriftSeverity = "CRITICAL"
)

// GetSeverity returns the severity for a given DriftType based on deterministic rules.
func GetSeverity(driftType DriftType, baseline, proposed any) DriftSeverity {
	switch driftType {
	case PathRemoved,
		MethodRemoved,
		RequiredFieldRemoved,
		EnumValueRemoved,
		ResponseRemoved:
		return Critical

	case FieldTypeChanged:
		// Type changed incompatibly is Critical.
		// For simplicity in this engine, we treat all type changes as Critical
		// as they usually break consumers.
		return Critical

	case ConstraintTightened,
		RequiredFieldAdded:
		// Nullable -> non-nullable is also constraint tightening in our logic
		return Breaking

	case FieldAdded,
		NewPath,
		NewMethod,
		ConstraintLoosened,
		EnumValueAdded:
		return Warning

	case MetadataChanged:
		return Info

	default:
		return Info
	}
}
