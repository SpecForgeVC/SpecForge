package app

import (
	"context"
	"time"

	"github.com/SpecForgeVC/SpecForge/internal/domain"
	"github.com/google/uuid"
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
	nodes := []map[string]interface{}{}
	edges := []map[string]interface{}{}
	visited := make(map[uuid.UUID]bool)
	queue := []uuid.UUID{variableID}

	// Root node
	nodes = append(nodes, map[string]interface{}{
		"id":   variableID.String(),
		"type": "variable",
		"data": map[string]string{"label": "Current Variable"},
	})
	visited[variableID] = true

	// Recursive traversal (BFS) with depth limit
	maxDepth := 5
	currentDepth := 0

	for len(queue) > 0 && currentDepth < maxDepth {
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			currID := queue[0]
			queue = queue[1:]

			deps, err := s.repo.ListDependencies(ctx, currID)
			if err != nil {
				return nil, err
			}

			for _, dep := range deps {
				if !visited[dep.TargetVariableID] {
					visited[dep.TargetVariableID] = true
					nodes = append(nodes, map[string]interface{}{
						"id":   dep.TargetVariableID.String(),
						"type": "variable",
						"data": map[string]string{"label": "Dependent"},
					})
					queue = append(queue, dep.TargetVariableID)
				}

				edges = append(edges, map[string]interface{}{
					"id":     dep.ID.String(),
					"source": currID.String(),
					"target": dep.TargetVariableID.String(),
				})
			}
		}
		currentDepth++
	}

	return map[string]interface{}{
		"nodes": nodes,
		"edges": edges,
	}, nil
}
