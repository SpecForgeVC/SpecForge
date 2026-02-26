package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type diffService struct {
}

func NewDiffService() DiffService {
	return &diffService{}
}

func (s *diffService) CompareSnapshots(ctx context.Context, oldSnap, newSnap map[string]interface{}) (*domain.DriftReport, error) {
	report := &domain.DriftReport{
		DriftDetected:   false,
		BreakingChanges: []domain.BreakingChange{},
		RiskScore:       0.0,
	}

	// Simple heuristic: Iterate over old keys. If missing in new, it's a breaking change.
	// This is recursive for nested maps.
	s.detectBreakingChanges(oldSnap, newSnap, "", &report.BreakingChanges)

	if len(report.BreakingChanges) > 0 {
		report.DriftDetected = true
		report.RiskScore = float64(len(report.BreakingChanges)) * 10 // Arbitrary score per breaking change
		if report.RiskScore > 100 {
			report.RiskScore = 100
		}
	} else {
		// Check for additions (non-breaking but drift)
		// For now, we only focus on breaking changes for the report's Critical status
	}

	return report, nil
}

func (s *diffService) CompareProjectSnapshot(ctx context.Context, snapshot domain.ProjectSnapshot, projectID uuid.UUID) (*domain.DriftReport, error) {
	report := &domain.DriftReport{
		DriftDetected:   false,
		BreakingChanges: []domain.BreakingChange{},
		RiskScore:       0.0,
	}

	// For project import drift, we are comparing the imported reality against
	// the intended SpecForge roadmap/contracts.
	// Since we are in the middle of implementation, we will perform a basic
	// structural validation as a form of drift detection.

	if len(snapshot.Architecture.Layers) < 2 {
		report.BreakingChanges = append(report.BreakingChanges, domain.BreakingChange{
			Field: "architecture.layers",
			Issue: "Insufficient architectural layering detected",
		})
	}

	if len(snapshot.Contracts.APIContracts) == 0 {
		report.BreakingChanges = append(report.BreakingChanges, domain.BreakingChange{
			Field: "contracts.api_contracts",
			Issue: "No API contracts found in imported snapshot",
		})
	}

	if len(report.BreakingChanges) > 0 {
		report.DriftDetected = true
		report.RiskScore = float64(len(report.BreakingChanges)) * 15.0
	}

	return report, nil
}

func (s *diffService) detectBreakingChanges(old, new map[string]interface{}, prefix string, changes *[]domain.BreakingChange) {
	for key, oldVal := range old {
		newVal, exists := new[key]
		currentPath := key
		if prefix != "" {
			currentPath = prefix + "." + key
		}

		if !exists {
			*changes = append(*changes, domain.BreakingChange{
				Field: currentPath,
				Issue: "Field removed",
			})
			continue
		}

		// Type check
		// Note: JSON numbers are floats in Go map[string]interface{} usually.
		// We skip detailed type checking for now, assuming basic structure.

		// Recurse if both are maps
		oldMap, okOld := oldVal.(map[string]interface{})
		newMap, okNew := newVal.(map[string]interface{})
		if okOld && okNew {
			s.detectBreakingChanges(oldMap, newMap, currentPath, changes)
		}
	}
}
