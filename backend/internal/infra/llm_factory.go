package infra

import (
	"context"
	"fmt"

	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/llm"
)

type LLMFactory struct{}

func NewLLMFactory() *LLMFactory {
	return &LLMFactory{}
}

func (f *LLMFactory) CreateClient(config *domain.LLMConfiguration) (domain.LLMClient, error) {
	switch config.Provider {
	case domain.ProviderOpenAI:
		return llm.NewOpenAIClient(config.APIKey, config.Model), nil
	case domain.ProviderOllama:
		return llm.NewOllamaClient(config.BaseURL, config.Model), nil
	case domain.ProviderGemini:
		return llm.NewGeminiClient(context.Background(), config.APIKey, config.Model)
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", config.Provider)
	}
}
