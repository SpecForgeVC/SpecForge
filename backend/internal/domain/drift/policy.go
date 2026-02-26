package drift

// DriftPolicy defines the rules for blocking based on drift severity.
type DriftPolicy struct {
	BlockOnBreaking bool
	BlockOnCritical bool
}

// Evaluate applies the policy to a DriftReport and determines if it should be blocked.
func (p DriftPolicy) Evaluate(report *DriftReport) {
	if report == nil {
		return
	}

	blocked := false
	if p.BlockOnCritical && report.CriticalChanges > 0 {
		blocked = true
	}
	if p.BlockOnBreaking && report.BreakingChanges > 0 {
		blocked = true
	}

	report.Blocked = blocked
}
