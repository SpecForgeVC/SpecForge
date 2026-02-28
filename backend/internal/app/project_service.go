package app

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type projectService struct {
	repo       ProjectRepository
	auditLog   AuditLogService
	llmService LLMService
}

func NewProjectService(repo ProjectRepository, auditLog AuditLogService, llmService LLMService) ProjectService {
	return &projectService{
		repo:       repo,
		auditLog:   auditLog,
		llmService: llmService,
	}
}

func (s *projectService) RecommendStack(ctx context.Context, purpose, constraints string) (*domain.TechStackRecommendation, error) {
	client, err := s.llmService.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	prompt := "Recommend a modern tech stack for a project with the following purpose: " + purpose
	if constraints != "" {
		prompt += "\nConstraints: " + constraints
	}
	prompt += "\n\nReturn ONLY a JSON object with 'recommended_stack' (a map of categories to technologies) and 'reasoning' (a string) fields. Do not include markdown formatting or explanations outside the JSON."

	resp, err := client.Generate(ctx, prompt)
	if err != nil {
		// Fallback for demo/simple cases if LLM fails or for testing
		return &domain.TechStackRecommendation{
			RecommendedStack: map[string]interface{}{
				"Frontend": "React + TailwindCSS + Vite",
				"Backend":  "Go + Echo",
				"Database": "PostgreSQL",
				"Auth":     "JWT",
			},
			Reasoning: "Standard modern stack for high performance and developer productivity.",
		}, nil
	}

	// Clean JSON response (strip markdown)
	cleanedResp := CleanJSON(resp)

	var recommendation domain.TechStackRecommendation
	if err := json.Unmarshal([]byte(cleanedResp), &recommendation); err != nil {
		return nil, err
	}

	return &recommendation, nil
}

func (s *projectService) GetProject(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	return s.repo.Get(ctx, id)
}

func (s *projectService) ListProjects(ctx context.Context, workspaceID uuid.UUID) ([]domain.Project, error) {
	return s.repo.List(ctx, workspaceID)
}

func (s *projectService) CreateProject(ctx context.Context, workspaceID uuid.UUID, name, description string, techStack map[string]interface{}, settings map[string]interface{}, mcpSettings domain.MCPSettings, repositoryURL string, userID uuid.UUID) (*domain.Project, error) {
	p := &domain.Project{
		WorkspaceID:   workspaceID,
		Name:          name,
		Description:   description,
		TechStack:     techStack,
		Settings:      settings,
		MCPSettings:   mcpSettings,
		RepositoryURL: repositoryURL,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "project", p.ID, "CREATE", userID, nil, map[string]interface{}{"name": p.Name})
	return p, nil
}

func (s *projectService) UpdateProject(ctx context.Context, id uuid.UUID, name, description string, techStack map[string]interface{}, settings map[string]interface{}, mcpSettings domain.MCPSettings, repositoryURL string, userID uuid.UUID) (*domain.Project, error) {
	oldP, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	p := *oldP
	p.Name = name
	p.Description = description
	p.TechStack = techStack
	p.Settings = settings
	p.MCPSettings = mcpSettings
	p.RepositoryURL = repositoryURL
	if err := s.repo.Update(ctx, &p); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "project", id, "UPDATE", userID, map[string]interface{}{"name": oldP.Name}, map[string]interface{}{"name": p.Name})
	return &p, nil
}

func (s *projectService) DeleteProject(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	p, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.auditLog.Log(ctx, "project", id, "DELETE", userID, map[string]interface{}{"name": p.Name}, nil)
	return nil
}
