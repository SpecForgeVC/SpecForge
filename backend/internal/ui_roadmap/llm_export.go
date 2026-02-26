package ui_roadmap

import (
	"fmt"
)

// GenerateLLMPrompt generates a strict, deterministic prompt for LLM-based UI implementation
func GenerateLLMPrompt(item *UIRoadmapItem) string {
	return fmt.Sprintf(`### 🏗️ UI IMPLEMENTATION SPECIFICATION: %s
**Screen Type**: %s
**Project ID**: %s

---

#### 👤 USER PERSONA & CONTEXT
- **Persona**: %s
- **Use Case**: %s
- **Description**: %s

---

#### 🧩 COMPONENT HIERARCHY (DETERMINISTIC)
%s

---

#### 🔄 STATE MACHINE & INTERACTION
%s

---

#### 🔗 BACKEND BINDINGS & DATA CONTRACTS
%s

---

#### 📱 RESPONSIVE BEHAVIOR
%s

---

#### ♿ ACCESSIBILITY (ARIA & FOCUS)
%s

---

#### ✅ VALIDATION & ERROR HANDLING
%s

---

#### 🛠️ TECHNICAL STACK REQUIREMENTS
- **Framework**: React 19 + Vite
- **Language**: TypeScript (Strict Mode)
- **Styling**: Tailwind CSS + Design Tokens (No raw hex)
- **State Management**: React Hooks (useState, useReducer, useEffect)
- **Data Fetching**: Axios / TanStack Query (linked to provided bindings)

#### 🚫 STRICT RULES
1. **No Assumptions**: Implement exactly as specified.
2. **Deterministic UI**: Align 1:1 with the component tree hierarchy.
3. **Accessibility First**: ARIA roles and keyboard tab indexing are mandatory.
4. **State Completeness**: Implement all mandatory states (idle, loading, success, error, empty, disabled).
5. **Token Compliance**: Use design system tokens for spacing, typography, and color.

#### 🔄 REFINEMENT LOOP INSTRUCTIONS
If the generated code deviates from this spec, prioritize the spec. 
Include rigorous typed props and error boundaries.
`,
		item.Name, item.ScreenType, item.ProjectID,
		item.UserPersona, item.UseCase, item.Description,
		formatJSON(item.ComponentTree),
		formatJSON(item.StateMachine),
		formatJSON(item.BackendBindings),
		formatJSON(item.ResponsiveSpec),
		formatJSON(item.AccessibilitySpec),
		formatJSON(item.ValidationRules),
	)
}

func formatJSON(j []byte) string {
	if len(j) == 0 {
		return "N/A"
	}
	return string(j)
}
