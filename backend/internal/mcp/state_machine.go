package mcp

import (
	"errors"
	"fmt"

	"github.com/scott/specforge/internal/domain"
)

var (
	ErrInvalidStateTransition = errors.New("invalid state transition")
	ErrSnapshotNotFound       = errors.New("snapshot not found")
)

// SnapshotStateMachine handles the lifecycle of a reality snapshot
type SnapshotStateMachine struct{}

func NewSnapshotStateMachine() *SnapshotStateMachine {
	return &SnapshotStateMachine{}
}

// CanTransition checks if a transition from current to next is allowed
func (sm *SnapshotStateMachine) CanTransition(current, next domain.SnapshotState) error {
	switch current {
	case domain.StateInitiated:
		if next == domain.StateAwaitingPost {
			return nil
		}
	case domain.StateAwaitingPost:
		if next == domain.StateAnalyzing {
			return nil
		}
	case domain.StateAnalyzing:
		if next == domain.StateCompleted || next == domain.StateFailed {
			return nil
		}
	}

	return fmt.Errorf("%w: %s -> %s", ErrInvalidStateTransition, current, next)
}

// GetRequiredNextStep returns the user-facing description of the next step
func (sm *SnapshotStateMachine) GetRequiredNextStep(state domain.SnapshotState) string {
	switch state {
	case domain.StateInitiated:
		return "collect_and_post_snapshot"
	case domain.StateAwaitingPost:
		return "post_snapshot"
	case domain.StateAnalyzing:
		return "wait_for_analysis"
	case domain.StateCompleted, domain.StateFailed:
		return "none"
	default:
		return "unknown"
	}
}
