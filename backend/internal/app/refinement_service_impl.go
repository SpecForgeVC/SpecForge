package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type refinementService struct {
	repo        RefinementRepository
	projectRepo ProjectRepository
	llmService  LLMService
	eventBus    map[uuid.UUID]chan domain.RefinementEvent // Simple in-memory event bus for now
}

func NewRefinementService(repo RefinementRepository, projectRepo ProjectRepository, llm LLMService) RefinementService {
	return &refinementService{
		repo:        repo,
		projectRepo: projectRepo,
		llmService:  llm,
		eventBus:    make(map[uuid.UUID]chan domain.RefinementEvent),
	}
}

func (s *refinementService) StartSession(ctx context.Context, artifactType string, targetType string, prompt string, contextData map[string]any, maxIterations int) (*domain.RefinementSession, error) {
	session := &domain.RefinementSession{
		ID:            uuid.New(),
		ArtifactType:  artifactType,
		TargetType:    targetType,
		InitialPrompt: prompt,
		ContextData:   contextData,
		MaxIterations: maxIterations,
		Status:        domain.RefinementStatusInProgress,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	// Create event channel
	s.eventBus[session.ID] = make(chan domain.RefinementEvent, 100)

	// Start async orchestrator
	go s.runOrchestrator(session)

	return session, nil
}

func (s *refinementService) GetSession(ctx context.Context, id uuid.UUID) (*domain.RefinementSession, error) {
	return s.repo.GetSession(ctx, id)
}

func (s *refinementService) GetSessionEvents(ctx context.Context, id uuid.UUID) (<-chan domain.RefinementEvent, error) {
	ch, ok := s.eventBus[id]
	if !ok {
		return nil, fmt.Errorf("session not active or found")
	}
	return ch, nil
}

func (s *refinementService) ApproveSession(ctx context.Context, id uuid.UUID) error {
	session, err := s.repo.GetSession(ctx, id)
	if err != nil {
		return err
	}
	session.Status = domain.RefinementStatusApproved
	session.UpdatedAt = time.Now()
	return s.repo.UpdateSession(ctx, session)
}

func (s *refinementService) runOrchestrator(session *domain.RefinementSession) {
	ctx := context.Background()
	eventCh := s.eventBus[session.ID]
	defer close(eventCh)

	publish := func(eventType, msg string, payload map[string]any) {
		eventCh <- domain.RefinementEvent{Type: eventType, Message: msg, Payload: payload}
	}

	publish("INFO", fmt.Sprintf("Starting refinement session for %s (%s)...", session.ArtifactType, session.TargetType), nil)

	llmClient, err := s.llmService.GetClient(ctx)
	if err != nil {
		publish("ERROR", fmt.Sprintf("Failed to get LLM client: %v", err), nil)
		session.Status = domain.RefinementStatusFailed
		_ = s.repo.UpdateSession(ctx, session)
		return
	}

	// Construct System Prompt based on TargetType
	var systemPrompt string
	switch session.TargetType {
	case "contract":
		systemPrompt = "You are an expert API architect. Generate a complete, valid OpenAPI 3.1.0 specification in JSON format that can be directly converted to YAML and saved as an openapi.spec.yaml file. The spec MUST include the top-level 'openapi', 'info', 'paths', and 'components' fields. Define all request/response schemas under 'components/schemas' and reference them in path operations. ALSO identify any necessary environment variables or secrets required (e.g., API keys, database URLs). Return a JSON object with two fields: 'contract' (the complete OpenAPI 3.1.0 spec object) and 'variables' (an array of variable definitions with 'name', 'description', 'required' fields)."
	case "variable":
		systemPrompt = "You are a DevOps engineer. Identify necessary environment variables, secrets, and configuration flags based on the requirements. Return a JSON object with a single field 'variables' containing an array of variable definitions (with 'name', 'description', 'required', 'default_value'). Return ONLY valid JSON, no markdown formatting, no explanations."
	case "context":
		systemPrompt = "You are a Product Manager and Technical Architect. Analyze the input and generate detailed Business Context and Technical Context. Return JSON with 'business_context' and 'technical_context' fields."
	case "roadmap_item":
		systemPrompt = "You are a Product Manager. Generate a comprehensive roadmap item based on the input. Return JSON with 'title', 'description' (detailed), 'business_context', 'technical_context', 'type' (EPIC/FEATURE/TASK/BUGFIX/REFACTOR), and 'priority' (LOW/MEDIUM/HIGH/CRITICAL) fields."
	case "requirement":
		systemPrompt = "You are an expert Technical Lead. Generate a list of detailed technical requirements based on the input. Return a JSON object with two fields: 'requirements' (an array of objects with 'title', 'acceptance_criteria', 'testable' (boolean), and 'priority' (LOW/MEDIUM/HIGH)) and 'variables' (an array of necessary environment variables with 'name', 'description', 'required', 'default_value')."
	case "schema_suggestion":
		systemPrompt = "You are an expert API Architect. Analyze the Roadmap Item context (Title, Description) and the existing partial Contract (Type, other schemas). Generate the requested JSON Schema for `{target_field}`. It must be distinct and appropriate for the specific role (e.g., Error schema should define error codes/messages, not copy Input). Return JSON with a single field 'schema' containing the JSON schema object."
	case "validation_rule":
		systemPrompt = "You are an expert Security and Quality Engineer. Analyze the project context and generate appropriate validation rules. These rules protect variables, contracts, and business logic. Return a JSON object with a single field 'rules', which is an array of objects. Each object must have: 'name' (descriptive), 'rule_type' (e.g., 'REGEX', 'RANGE', 'ENUM', 'CUSTOM'), 'description' (clear explanation), and 'rule_config' (a JSON object with specific parameters for the type, e.g., {'pattern': '^v.+'} or {'min': 0, 'max': 100}). Return ONLY valid JSON."
	default:
		systemPrompt = "You are a helpful AI assistant. Generate the requested artifact in JSON format."
	}

	currentPrompt := fmt.Sprintf("%s\n\nTask: %s", systemPrompt, session.InitialPrompt)
	if session.ContextData != nil {
		contextFunc, _ := json.Marshal(session.ContextData)
		currentPrompt += fmt.Sprintf("\n\nContext:\n%s", string(contextFunc))
	}

	for i := 1; i <= session.MaxIterations; i++ {
		session.CurrentIteration = i
		_ = s.repo.UpdateSession(ctx, session)
		publish("ITERATION_START", fmt.Sprintf("Starting iteration %d/%d", i, session.MaxIterations), nil)

		validationErrors := []string{}

		// 1. Generate
		publish("STEP", "Generating artifact...", nil)
		resp, err := llmClient.Generate(ctx, currentPrompt)
		if err != nil {
			publish("ERROR", fmt.Sprintf("LLM generation failed: %v", err), nil)
			break
		}

		// Log response for debugging
		if err := logResponse("Refinement Iteration "+fmt.Sprint(i), currentPrompt, resp); err != nil {
			fmt.Printf("Failed to log response: %v\n", err)
		}

		// Clean JSON response (strip markdown)
		cleanedResp := cleanJSON(resp)

		// Parse response
		var artifact map[string]any
		if err := json.Unmarshal([]byte(cleanedResp), &artifact); err != nil {
			// Log the failed parse attempt
			publish("WARN", "Failed to parse JSON response. Retrying with format instruction...", nil)
			currentPrompt += "\n\nCRITICAL ERROR: The previous response was not valid JSON. Please return ONLY the raw JSON object, no markdown formatting."
			continue
		}

		// 2. Self-Evaluation (if enabled)
		var selfEval *domain.SelfEvaluationResult
		projectIDStr, ok := session.ContextData["project_id"].(string)
		if ok {
			projectID, err := uuid.Parse(projectIDStr)
			if err == nil {
				project, err := s.projectRepo.Get(ctx, projectID)
				if err == nil && project.Settings != nil {
					if selfEvalEnabled(project.Settings) {
						publish("STEP", "AI Self-Critique phase...", nil)
						eval, err := s.generateSelfEvaluation(ctx, llmClient, artifact)
						if err == nil {
							selfEval = eval
							publish("INFO", fmt.Sprintf("AI Score: %d/10. %s", eval.Score, eval.ImprovementSuggestions[0]), nil)

							// If score is too low, treat as validation failure
							if eval.Score < 7 {
								validationErrors = append(validationErrors, fmt.Sprintf("AI Self-Critique flagged low score (%d/10). Issues: %v", eval.Score, eval.ImprovementSuggestions))
							}
						} else {
							publish("WARN", fmt.Sprintf("Self-evaluation failed: %v", err), nil)
						}
					}
				}
			}
		}

		iteration := &domain.RefinementIteration{
			ID:             uuid.New(),
			SessionID:      session.ID,
			Iteration:      i,
			Prompt:         currentPrompt,
			Response:       resp,
			Artifact:       artifact,
			SelfEvaluation: selfEval,
			CreatedAt:      time.Now(),
		}
		_ = s.repo.CreateIteration(ctx, iteration)

		if len(validationErrors) == 0 {
			session.Status = domain.RefinementStatusValidated
			session.Result = artifact
			session.ValidationErrors = nil
			if selfEval != nil {
				session.ConfidenceScore = float64(selfEval.Score) / 10.0
			} else {
				session.ConfidenceScore = 1.0 // Mock
			}
			_ = s.repo.UpdateSession(ctx, session)
			publish("SUCCESS", "Validation passed!", map[string]any{"artifact": artifact, "evaluation": selfEval})
			return
		}

		// 3. Feedback Loop
		publish("INFO", fmt.Sprintf("Validation failed with %d errors. Refining...", len(validationErrors)), nil)
		feedbackPrompt := fmt.Sprintf("%s\n\nPrevious attempt failed validation or self-critique:\nErrors: %v\n\n", systemPrompt, validationErrors)
		if selfEval != nil {
			feedbackPrompt += "### AI Self-Critique Findings:\n"
			if len(selfEval.AmbiguityFlags) > 0 {
				feedbackPrompt += fmt.Sprintf("- Ambiguity: %v\n", selfEval.AmbiguityFlags)
			}
			if len(selfEval.MissingConstraints) > 0 {
				feedbackPrompt += fmt.Sprintf("- Missing Constraints: %v\n", selfEval.MissingConstraints)
			}
			if len(selfEval.WeakValidations) > 0 {
				feedbackPrompt += fmt.Sprintf("- Weak Validations: %v\n", selfEval.WeakValidations)
			}
			if len(selfEval.SecurityConcerns) > 0 {
				feedbackPrompt += fmt.Sprintf("- Security Concerns: %v\n", selfEval.SecurityConcerns)
			}
			if len(selfEval.ImprovementSuggestions) > 0 {
				feedbackPrompt += fmt.Sprintf("- Improvement Suggestions: %v\n", selfEval.ImprovementSuggestions)
			}
			feedbackPrompt += "\n"
		}
		feedbackPrompt += "Please address all the issues listed above and regenerate the JSON artifact."
		currentPrompt = feedbackPrompt
	}

	session.Status = domain.RefinementStatusFailed
	_ = s.repo.UpdateSession(ctx, session)
	publish("ERROR", "Max iterations reached without validation success.", nil)
}

func (s *refinementService) generateSelfEvaluation(ctx context.Context, client domain.LLMClient, artifact map[string]any) (*domain.SelfEvaluationResult, error) {
	artifactJSON, _ := json.MarshalIndent(artifact, "", "  ")

	prompt := fmt.Sprintf(`You are a lead security architect. Critically evaluate the following generated artifact:

<artifact>
%s
</artifact>

Identify:
- Missing required schema elements or constraints
- Ambiguous definitions or descriptions
- Weak validation rules
- Security risks (e.g., hardcoded values, missing auth scopes, injection points)
- Logical inconsistencies

Return a structured JSON object only. The JSON must follow this exact format:
{
  "score": (integer 1-10, where 10 is perfect),
  "ambiguity_flags": ["list of items"],
  "missing_constraints": ["list of items"],
  "weak_validations": ["list of items"],
  "security_concerns": ["list of items"],
  "improvement_suggestions": ["list of items"]
}

The "score" must be a number from 1 to 10. You MUST justify the score through the lists provided. Provide ONLY raw JSON, no markdown formatting.`, string(artifactJSON))

	resp, err := client.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	cleanedResp := cleanJSON(resp)
	var result domain.SelfEvaluationResult
	if err := json.Unmarshal([]byte(cleanedResp), &result); err != nil {
		return nil, fmt.Errorf("failed to parse self-evaluation JSON: %w", err)
	}

	return &result, nil
}

func selfEvalEnabled(settings map[string]any) bool {
	if settings == nil {
		return false
	}
	v, ok := settings["enable_self_evaluation"]
	if !ok {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return strings.ToLower(val) == "true"
	case float64:
		return val != 0
	case int:
		return val != 0
	default:
		return false
	}
}

// Helper to log LLM responses to a file for debugging
func logResponse(step, prompt, response string) error {
	f, err := os.OpenFile("llm-response.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := time.Now().Format(time.RFC3339)
	entry := fmt.Sprintf("\n# %s - %s\n\n## Prompt\n```\n%s\n```\n\n## Response\n```json\n%s\n```\n\n---\n", timestamp, step, prompt, response)

	if _, err := f.WriteString(entry); err != nil {
		return err
	}
	return nil
}

// Helper to strip markdown code blocks and extract JSON from mixed content
func cleanJSON(input string) string {
	input = strings.TrimSpace(input)

	// First strip markdown blocks if present
	if strings.Contains(input, "```") {
		start := strings.Index(input, "```json")
		if start == -1 {
			start = strings.Index(input, "```")
		}
		if start != -1 {
			// Find the end of the block
			end := strings.Index(input[start+3:], "```")
			if end != -1 {
				// Extract the content inside the block
				blockContent := input[start:]
				if strings.HasPrefix(blockContent, "```json") {
					blockContent = strings.TrimPrefix(blockContent, "```json")
				} else {
					blockContent = strings.TrimPrefix(blockContent, "```")
				}
				// Trim the end block
				if idx := strings.Index(blockContent, "```"); idx != -1 {
					blockContent = blockContent[:idx]
				}
				input = blockContent
			}
		}
	}

	// Then find the outer braces/brackets to be sure
	firstBrace := strings.Index(input, "{")
	firstBracket := strings.Index(input, "[")

	start := -1
	if firstBrace != -1 && firstBracket != -1 {
		if firstBrace < firstBracket {
			start = firstBrace
		} else {
			start = firstBracket
		}
	} else if firstBrace != -1 {
		start = firstBrace
	} else if firstBracket != -1 {
		start = firstBracket
	}

	if start != -1 {
		input = input[start:]
		// Find last brace or bracket
		lastBrace := strings.LastIndex(input, "}")
		lastBracket := strings.LastIndex(input, "]")
		end := -1

		if lastBrace > lastBracket {
			end = lastBrace
		} else {
			end = lastBracket
		}

		if end != -1 {
			input = input[:end+1]
		}
	}

	return strings.TrimSpace(input)
}
