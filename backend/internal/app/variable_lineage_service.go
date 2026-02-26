package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type variableLineageService struct {
	repo VariableLineageRepository
}

func NewVariableLineageService(repo VariableLineageRepository) VariableLineageService {
	return &variableLineageService{repo: repo}
}

func (s *variableLineageService) TrackEvent(ctx context.Context, variableID uuid.UUID, eventType domain.LineageEventType, source, description string, userID uuid.UUID, metadata map[string]interface{}) error {
	event := &domain.VariableLineageEvent{
		VariableID:      variableID,
		EventType:       eventType,
		SourceComponent: source,
		Description:     description,
		PerformedBy:     userID,
		Metadata:        metadata,
		CreatedAt:       time.Now(),
	}
	return s.repo.CreateEvent(ctx, event)
}

func (s *variableLineageService) GetLineageEvents(ctx context.Context, variableID uuid.UUID) ([]domain.VariableLineageEvent, error) {
	return s.repo.ListEvents(ctx, variableID)
}

func (s *variableLineageService) GetLineageGraph(ctx context.Context, variableID uuid.UUID) (map[string]interface{}, error) {
	// For now, return a simple graph centered on this variable
	// In the future, traverse dependencies

	deps, err := s.repo.ListDependencies(ctx, variableID)
	if err != nil {
		return nil, err
	}

	nodes := []map[string]interface{}{
		{
			"id":   variableID.String(),
			"type": "variable",
			"data": map[string]string{"label": "Current Variable"},
		},
	}
	edges := []map[string]interface{}{}

	for _, dep := range deps {
		nodes = append(nodes, map[string]interface{}{
			"id":   dep.TargetVariableID.String(),
			"type": "variable",
			"data": map[string]string{"label": "Dependent"},
		})
		edges = append(edges, map[string]interface{}{
			"id":     dep.ID.String(),
			"source": variableID.String(),
			"target": dep.TargetVariableID.String(),
		})
	}

	return map[string]interface{}{
		"nodes": nodes,
		"edges": edges,
	}, nil
}
