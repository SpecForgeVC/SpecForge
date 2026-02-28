package app

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockProjectRepo struct{ mock.Mock }

func (m *mockProjectRepo) Get(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Project), args.Error(1)
}
func (m *mockProjectRepo) List(ctx context.Context, workspaceID uuid.UUID) ([]domain.Project, error) {
	args := m.Called(ctx, workspaceID)
	return args.Get(0).([]domain.Project), args.Error(1)
}
func (m *mockProjectRepo) Create(ctx context.Context, p *domain.Project) error {
	return m.Called(ctx, p).Error(0)
}
func (m *mockProjectRepo) Update(ctx context.Context, p *domain.Project) error {
	return m.Called(ctx, p).Error(0)
}
func (m *mockProjectRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type mockAuditLog struct{ mock.Mock }

func (m *mockAuditLog) Log(ctx context.Context, entityType string, entityID uuid.UUID, action string, userID uuid.UUID, oldData, newData map[string]interface{}) error {
	return m.Called(ctx, entityType, entityID, action, userID, oldData, newData).Error(0)
}
func (m *mockAuditLog) GetEntityLogs(ctx context.Context, entityType string, entityID uuid.UUID) ([]domain.AuditLog, error) {
	return nil, nil
}
func (m *mockAuditLog) ListDriftEvents(ctx context.Context) ([]domain.AuditLog, error) {
	return nil, nil
}

type mockLLMService struct{ mock.Mock }

func (m *mockLLMService) GetClient(ctx context.Context) (domain.LLMClient, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(domain.LLMClient), args.Error(1)
}
func (m *mockLLMService) GetActiveConfig(ctx context.Context) (*domain.LLMConfiguration, error) {
	return nil, nil
}
func (m *mockLLMService) UpdateConfig(ctx context.Context, config *domain.LLMConfiguration) error {
	return nil
}
func (m *mockLLMService) TestConfiguration(ctx context.Context, config *domain.LLMConfiguration) error {
	return nil
}
func (m *mockLLMService) ListModels(ctx context.Context, config *domain.LLMConfiguration) ([]string, error) {
	return nil, nil
}

type mockLLMClient struct{ mock.Mock }

func (m *mockLLMClient) Generate(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}
func (m *mockLLMClient) StreamGenerate(ctx context.Context, prompt string, decimals chan<- string) error {
	return nil
}
func (m *mockLLMClient) TestConnection(ctx context.Context) error         { return nil }
func (m *mockLLMClient) ListModels(ctx context.Context) ([]string, error) { return nil, nil }

func TestRecommendStack(t *testing.T) {
	pRepo := new(mockProjectRepo)
	audit := new(mockAuditLog)
	llmSvc := new(mockLLMService)
	llmClient := new(mockLLMClient)

	service := NewProjectService(pRepo, audit, llmSvc)

	t.Run("successful recommendation", func(t *testing.T) {
		purpose := "A simple todo app"
		ctx := context.Background()

		llmSvc.On("GetClient", ctx).Return(llmClient, nil)
		jsonResponse := `{"recommended_stack": {"Frontend": "React"}, "reasoning": "Standard for simple apps"}`
		llmClient.On("Generate", ctx, mock.Anything).Return(jsonResponse, nil)

		res, err := service.RecommendStack(ctx, purpose, "")

		assert.NoError(t, err)
		assert.Equal(t, "React", res.RecommendedStack["Frontend"])
		assert.Equal(t, "Standard for simple apps", res.Reasoning)
	})

	t.Run("fallback on LLM error", func(t *testing.T) {
		purpose := "Complex system"
		ctx := context.Background()

		llmSvc.ExpectedCalls = nil // Reset
		llmSvc.On("GetClient", ctx).Return(llmClient, nil)
		llmClient.ExpectedCalls = nil // Reset
		llmClient.On("Generate", ctx, mock.Anything).Return("", assert.AnError)

		res, err := service.RecommendStack(ctx, purpose, "")

		assert.NoError(t, err)
		assert.NotNil(t, res.RecommendedStack)
		assert.Contains(t, res.Reasoning, "Standard modern stack")
	})
}
