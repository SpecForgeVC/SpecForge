package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

// LLMFactory interface to avoid cyclic dependency if we used the concrete factory struct directly
type LLMFactory interface {
	CreateClient(config *domain.LLMConfiguration) (domain.LLMClient, error)
}

type llmService struct {
	repo    LLMRepository
	factory LLMFactory
}

func NewLLMService(repo LLMRepository, factory LLMFactory) LLMService {
	return &llmService{
		repo:    repo,
		factory: factory,
	}
}

func (s *llmService) GetActiveConfig(ctx context.Context) (*domain.LLMConfiguration, error) {
	return s.repo.GetActive(ctx)
}

func (s *llmService) UpdateConfig(ctx context.Context, config *domain.LLMConfiguration) error {
	// If API key is masked or empty, try to retrieve the existing one if we have an ID
	if (config.APIKey == "********" || config.APIKey == "") && config.ID != uuid.Nil {
		existing, err := s.repo.GetActive(ctx)
		// We only use existing key if the ID matches or if we are just updating the active config
		// Simpler logic: GetActive returns the *single* active config or row.
		// If our repo supports multiple configs, we'd need GetByID.
		// Assuming GetActive is what we want for now since we seem to support single active provider.
		if err == nil && existing != nil {
			config.APIKey = existing.APIKey
		}
	}

	return s.repo.Upsert(ctx, config)
}

func (s *llmService) GetClient(ctx context.Context) (domain.LLMClient, error) {
	config, err := s.repo.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active llm config: %w", err)
	}
	if config == nil {
		return nil, fmt.Errorf("no active llm configuration found")
	}

	return s.factory.CreateClient(config)
}

func (s *llmService) TestConfiguration(ctx context.Context, config *domain.LLMConfiguration) error {
	// Handle masked key for testing an existing config
	if (config.APIKey == "********" || config.APIKey == "") && config.ID != uuid.Nil {
		existing, err := s.repo.GetActive(ctx)
		if err == nil && existing != nil && existing.ID == config.ID {
			// Use existing key but keep other fields from the request (candidate changes)
			config.APIKey = existing.APIKey
		}
	}

	client, err := s.factory.CreateClient(config)
	if err != nil {
		return err
	}
	return client.TestConnection(ctx)
}

func (s *llmService) ListModels(ctx context.Context, config *domain.LLMConfiguration) ([]string, error) {
	// Handle masked key
	if (config.APIKey == "********" || config.APIKey == "") && config.ID != uuid.Nil {
		existing, err := s.repo.GetActive(ctx)
		if err == nil && existing != nil && existing.ID == config.ID {
			config.APIKey = existing.APIKey
		}
	}

	client, err := s.factory.CreateClient(config)
	if err != nil {
		return nil, err
	}
	return client.ListModels(ctx)
}
