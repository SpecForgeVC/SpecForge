package app

import (
	"context"
	"fmt"

	"github.com/SpecForgeVC/SpecForge/internal/domain"
	"github.com/google/uuid"
)

type aiProposalService struct {
	repo         AiProposalRepository
	roadmapRepo  RoadmapItemRepository
	snapshotRepo SnapshotRepository
	varRepo      VariableRepository
	contractRepo ContractRepository
	auditLog     AuditLogService
}

func NewAiProposalService(
	repo AiProposalRepository,
	rmRepo RoadmapItemRepository,
	sRepo SnapshotRepository,
	varRepo VariableRepository,
	contractRepo ContractRepository,
	al AuditLogService,
) AiProposalService {
	return &aiProposalService{
		repo:         repo,
		roadmapRepo:  rmRepo,
		snapshotRepo: sRepo,
		varRepo:      varRepo,
		contractRepo: contractRepo,
		auditLog:     al,
	}
}

func (s *aiProposalService) GetProposal(ctx context.Context, id uuid.UUID) (*domain.AiProposal, error) {
	return s.repo.Get(ctx, id)
}

func (s *aiProposalService) ListProposals(ctx context.Context, projectID uuid.UUID) ([]domain.AiProposal, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *aiProposalService) CreateProposal(ctx context.Context, roadmapItemID uuid.UUID, pType domain.ProposalType, diff map[string]interface{}, reasoning string, confidence float64) (*domain.AiProposal, error) {
	p := &domain.AiProposal{
		ID:              uuid.New(),
		RoadmapItemID:   roadmapItemID,
		ProposalType:    pType,
		Diff:            diff,
		Reasoning:       reasoning,
		ConfidenceScore: confidence,
		Status:          domain.Pending,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "ai_proposal", p.ID, "CREATE", uuid.Nil, nil, map[string]interface{}{"type": p.ProposalType})
	return p, nil
}

func (s *aiProposalService) ApproveProposal(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	p, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if p.Status != domain.Pending {
		return fmt.Errorf("proposal is already %s", p.Status)
	}

	// 1. Update proposal status first
	if err := s.repo.UpdateStatus(ctx, id, domain.Approved, userID); err != nil {
		return err
	}
	s.auditLog.Log(ctx, "ai_proposal", id, "APPROVE", userID,
		map[string]interface{}{"status": p.Status},
		map[string]interface{}{"status": domain.Approved},
	)

	// 2. Fetch the roadmap item to apply diff
	rm, err := s.roadmapRepo.Get(ctx, p.RoadmapItemID)
	if err != nil {
		return fmt.Errorf("failed to fetch roadmap item for proposal application: %w", err)
	}

	// 3. Apply diff based on ProposalType
	if err := s.applyProposalDiff(ctx, rm, p, userID); err != nil {
		// Log the failure but don't fail the approval — status is already updated
		s.auditLog.Log(ctx, "ai_proposal", id, "APPLY_DIFF_FAILED", userID, nil,
			map[string]interface{}{"error": err.Error()})
	}

	// 4. Create version snapshot capturing the state at approval time
	snap := &domain.VersionSnapshot{
		ID:            uuid.New(),
		RoadmapItemID: p.RoadmapItemID,
		SnapshotData: map[string]interface{}{
			"proposal_id":   id,
			"proposal_type": p.ProposalType,
			"diff":          p.Diff,
			"roadmap_item":  rm,
		},
		CreatedBy: userID,
	}
	return s.snapshotRepo.Create(ctx, snap)
}

// applyProposalDiff mutates the roadmap item (and related entities) based on the proposal type.
func (s *aiProposalService) applyProposalDiff(ctx context.Context, rm *domain.RoadmapItem, p *domain.AiProposal, userID uuid.UUID) error {
	switch p.ProposalType {
	case domain.EditDescription:
		// Apply description, business context, and/or technical context changes
		if v, ok := p.Diff["description"].(string); ok && v != "" {
			rm.Description = v
		}
		if v, ok := p.Diff["business_context"].(string); ok && v != "" {
			rm.BusinessContext = v
		}
		if v, ok := p.Diff["technical_context"].(string); ok && v != "" {
			rm.TechnicalContext = v
		}
		return s.roadmapRepo.Update(ctx, rm)

	case domain.ModifySchema:
		// Apply schema changes to a specific contract
		contractIDStr, ok := p.Diff["contract_id"].(string)
		if !ok {
			return fmt.Errorf("ModifySchema proposal missing contract_id in diff")
		}
		contractID, err := uuid.Parse(contractIDStr)
		if err != nil {
			return fmt.Errorf("invalid contract_id in diff: %w", err)
		}
		contract, err := s.contractRepo.Get(ctx, contractID)
		if err != nil {
			return fmt.Errorf("failed to fetch contract: %w", err)
		}
		if inputSchema, ok := p.Diff["input_schema"].(map[string]interface{}); ok {
			for k, v := range inputSchema {
				contract.InputSchema[k] = v
			}
		}
		if outputSchema, ok := p.Diff["output_schema"].(map[string]interface{}); ok {
			for k, v := range outputSchema {
				contract.OutputSchema[k] = v
			}
		}
		return s.contractRepo.Update(ctx, contract)

	case domain.AddVariable:
		// Create a new variable from the diff
		varData, ok := p.Diff["variable"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("AddVariable proposal missing 'variable' map in diff")
		}
		contractIDStr, _ := varData["contract_id"].(string)
		contractID, err := uuid.Parse(contractIDStr)
		if err != nil {
			return fmt.Errorf("AddVariable proposal has invalid contract_id: %w", err)
		}
		name, _ := varData["name"].(string)
		vType, _ := varData["type"].(string)
		required, _ := varData["required"].(bool)
		defaultValue, _ := varData["default_value"].(string)
		description, _ := varData["description"].(string)
		validationRules, _ := varData["validation_rules"].(map[string]interface{})

		newVar := &domain.VariableDefinition{
			ID:              uuid.New(),
			ContractID:      contractID,
			Name:            name,
			Type:            vType,
			Required:        required,
			DefaultValue:    defaultValue,
			Description:     description,
			ValidationRules: validationRules,
		}
		return s.varRepo.Create(ctx, newVar)

	case domain.RemoveField:
		// Remove a field from a contract's input or output schema
		contractIDStr, ok := p.Diff["contract_id"].(string)
		if !ok {
			return fmt.Errorf("RemoveField proposal missing contract_id in diff")
		}
		contractID, err := uuid.Parse(contractIDStr)
		if err != nil {
			return fmt.Errorf("invalid contract_id in diff: %w", err)
		}
		contract, err := s.contractRepo.Get(ctx, contractID)
		if err != nil {
			return fmt.Errorf("failed to fetch contract: %w", err)
		}
		fieldName, _ := p.Diff["field"].(string)
		schemaTarget, _ := p.Diff["schema"].(string) // "input" or "output"

		if fieldName != "" {
			if schemaTarget == "output" {
				delete(contract.OutputSchema, fieldName)
				// Also check in properties sub-object
				if props, ok := contract.OutputSchema["properties"].(map[string]interface{}); ok {
					delete(props, fieldName)
				}
			} else {
				delete(contract.InputSchema, fieldName)
				if props, ok := contract.InputSchema["properties"].(map[string]interface{}); ok {
					delete(props, fieldName)
				}
			}
			// Add to deprecated_fields to signal backward compatibility concern
			contract.BackwardCompatible = false
			contract.DeprecatedFields = append(contract.DeprecatedFields, fieldName)
		}
		return s.contractRepo.Update(ctx, contract)

	default:
		return fmt.Errorf("unknown proposal type: %s", p.ProposalType)
	}
}

func (s *aiProposalService) RejectProposal(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	p, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateStatus(ctx, id, domain.Rejected, userID); err != nil {
		return err
	}
	s.auditLog.Log(ctx, "ai_proposal", id, "REJECT", userID,
		map[string]interface{}{"status": p.Status},
		map[string]interface{}{"status": domain.Rejected},
	)
	return nil
}
