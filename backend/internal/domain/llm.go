package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LLMProvider string

const (
	ProviderOpenAI    LLMProvider = "openai"
	ProviderOllama    LLMProvider = "ollama"
	ProviderGemini    LLMProvider = "gemini"
	ProviderAnthropic LLMProvider = "anthropic"
)

type LLMConfiguration struct {
	ID        uuid.UUID   `json:"id"`
	Provider  LLMProvider `json:"provider"`
	APIKey    string      `json:"api_key,omitempty"` // Explicitly handle masking in handlers
	BaseURL   string      `json:"base_url,omitempty"`
	Model     string      `json:"model"`
	IsActive  bool        `json:"is_active"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type LLMClient interface {
	// Generate sends a prompt to the LLM and returns the response.
	// supports streaming via a callback if needed, but for now simple request/response
	Generate(ctx context.Context, prompt string) (string, error)

	// StreamGenerate sends a prompt and streams the response chunks.
	StreamGenerate(ctx context.Context, prompt string, decimals chan<- string) error

	// TestConnection verifies if the client can successfully communicate with the provider.
	TestConnection(ctx context.Context) error

	// ListModels retrieves the list of available models from the provider.
	ListModels(ctx context.Context) ([]string, error)
}
