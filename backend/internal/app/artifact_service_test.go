package app

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type mockRoadmapRepo struct{ mock.Mock }

func (m *mockRoadmapRepo) Get(ctx context.Context, id uuid.UUID) (*domain.RoadmapItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.RoadmapItem), args.Error(1)
}
func (m *mockRoadmapRepo) List(ctx context.Context, projectID uuid.UUID) ([]domain.RoadmapItem, error) {
	return nil, nil
}
func (m *mockRoadmapRepo) Create(ctx context.Context, item *domain.RoadmapItem) error {
	return nil
}
func (m *mockRoadmapRepo) Update(ctx context.Context, item *domain.RoadmapItem) error {
	return nil
}
func (m *mockRoadmapRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

type mockContractRepo struct{ mock.Mock }

func (m *mockContractRepo) Get(ctx context.Context, id uuid.UUID) (*domain.ContractDefinition, error) {
	return nil, nil
}
func (m *mockContractRepo) List(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.ContractDefinition, error) {
	args := m.Called(ctx, roadmapItemID)
	return args.Get(0).([]domain.ContractDefinition), args.Error(1)
}
func (m *mockContractRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.ContractDefinition, error) {
	return nil, nil
}
func (m *mockContractRepo) Create(ctx context.Context, c *domain.ContractDefinition) error {
	return nil
}
func (m *mockContractRepo) Update(ctx context.Context, c *domain.ContractDefinition) error {
	return nil
}
func (m *mockContractRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

type mockVariableRepo struct{ mock.Mock }

func (m *mockVariableRepo) Get(ctx context.Context, id uuid.UUID) (*domain.VariableDefinition, error) {
	return nil, nil
}
func (m *mockVariableRepo) List(ctx context.Context, contractID uuid.UUID) ([]domain.VariableDefinition, error) {
	args := m.Called(ctx, contractID)
	return args.Get(0).([]domain.VariableDefinition), args.Error(1)
}
func (m *mockVariableRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.VariableDefinition, error) {
	return nil, nil
}
func (m *mockVariableRepo) Create(ctx context.Context, v *domain.VariableDefinition) error {
	return nil
}
func (m *mockVariableRepo) Update(ctx context.Context, v *domain.VariableDefinition) error {
	return nil
}
func (m *mockVariableRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

type mockRequirementRepo struct{ mock.Mock }

func (m *mockRequirementRepo) Get(ctx context.Context, id uuid.UUID) (*domain.Requirement, error) {
	return nil, nil
}
func (m *mockRequirementRepo) List(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.Requirement, error) {
	args := m.Called(ctx, roadmapItemID)
	return args.Get(0).([]domain.Requirement), args.Error(1)
}
func (m *mockRequirementRepo) Create(ctx context.Context, r *domain.Requirement) error {
	return nil
}
func (m *mockRequirementRepo) Update(ctx context.Context, r *domain.Requirement) error {
	return nil
}
func (m *mockRequirementRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

type mockValidationRepo struct{ mock.Mock }

func (m *mockValidationRepo) Get(ctx context.Context, id uuid.UUID) (*domain.ValidationRule, error) {
	return nil, nil
}
func (m *mockValidationRepo) List(ctx context.Context, projectID uuid.UUID) ([]domain.ValidationRule, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]domain.ValidationRule), args.Error(1)
}
func (m *mockValidationRepo) Create(ctx context.Context, r *domain.ValidationRule) error {
	return nil
}
func (m *mockValidationRepo) Update(ctx context.Context, r *domain.ValidationRule) error {
	return nil
}
func (m *mockValidationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

type mockGovService struct{ mock.Mock }

func (m *mockGovService) CanBuildFeature(ctx context.Context, featureID uuid.UUID) (bool, []string, error) {
	args := m.Called(ctx, featureID)
	return args.Bool(0), args.Get(1).([]string), args.Error(2)
}
func (m *mockGovService) CanDeployFeature(ctx context.Context, featureID uuid.UUID) (bool, []string, error) {
	return true, nil, nil
}
func (m *mockGovService) CanUpdateContract(ctx context.Context, contractID uuid.UUID) (bool, []string, error) {
	return true, nil, nil
}

func TestArtifactService_GenerateArtifact(t *testing.T) {
	ctx := context.Background()
	roadmapItemID := uuid.New()
	userID := uuid.New()
	projectID := uuid.New()

	item := &domain.RoadmapItem{
		ID:        roadmapItemID,
		ProjectID: projectID,
		Title:     "Test Feature",
		Status:    domain.StatusApproved,
	}

	rmRepo := new(mockRoadmapRepo)
	rmRepo.On("Get", ctx, roadmapItemID).Return(item, nil)

	cRepo := new(mockContractRepo)
	cRepo.On("List", ctx, roadmapItemID).Return([]domain.ContractDefinition{}, nil)

	vRepo := new(mockVariableRepo)
	// Not called if contracts are empty

	reqRepo := new(mockRequirementRepo)
	reqRepo.On("List", ctx, roadmapItemID).Return([]domain.Requirement{}, nil)

	valRepo := new(mockValidationRepo)
	valRepo.On("List", ctx, projectID).Return([]domain.ValidationRule{}, nil)

	govSvc := new(mockGovService)
	govSvc.On("CanBuildFeature", ctx, roadmapItemID).Return(true, []string{"Check Passed"}, nil)

	service := NewBuildArtifactService(rmRepo, cRepo, vRepo, reqRepo, valRepo, govSvc)

	options := ExportOptions{
		IncludeDependencies: true,
		IncludeGovernance:   true,
	}

	pkg, err := service.GenerateArtifact(ctx, roadmapItemID, domain.ExportFormatZip, options, userID)

	assert.NoError(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, roadmapItemID, pkg.Metadata.RoadmapItemID)
	assert.Equal(t, "Test Feature", pkg.RoadmapContext.Title)
	assert.Equal(t, "READY", pkg.GovernanceConstraints.ComplianceStatus)

	rmRepo.AssertExpectations(t)
	cRepo.AssertExpectations(t)
	reqRepo.AssertExpectations(t)
	valRepo.AssertExpectations(t)
	govSvc.AssertExpectations(t)
}
