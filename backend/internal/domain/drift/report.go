package drift

import "sort"

// DriftReport represents the outcome of a drift detection process.
type DriftReport struct {
	BreakingChanges int         `json:"breaking_changes"`
	CriticalChanges int         `json:"critical_changes"`
	Warnings        int         `json:"warnings"`
	Infos           int         `json:"infos"`
	Items           []DriftItem `json:"items"`
	Blocked         bool        `json:"blocked"`
}

// DriftItem represents a single instance of detected drift.
type DriftItem struct {
	Type        DriftType     `json:"type"`
	Severity    DriftSeverity `json:"severity"`
	Location    string        `json:"location"`
	Baseline    any           `json:"baseline"`
	Proposed    any           `json:"proposed"`
	Description string        `json:"description"`
}

// NewDriftReport initializes a report and sorts items for stability.
func NewDriftReport(items []DriftItem) DriftReport {
	report := DriftReport{
		Items: items,
	}

	// Sort items by location and then by type for stability
	sort.Slice(report.Items, func(i, j int) bool {
		if report.Items[i].Location != report.Items[j].Location {
			return report.Items[i].Location < report.Items[j].Location
		}
		return report.Items[i].Type < report.Items[j].Type
	})

	for _, item := range report.Items {
		switch item.Severity {
		case Critical:
			report.CriticalChanges++
		case Breaking:
			report.BreakingChanges++
		case Warning:
			report.Warnings++
		case Info:
			report.Infos++
		}
	}

	return report
}
