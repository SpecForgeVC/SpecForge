package ui_roadmap

import (
	"fmt"
	"strings"
)

// GenerateFigmaMakePrompt generates instructions for the Figma Make AI tool
func GenerateFigmaMakePrompt(item *UIRoadmapItem) string {
	return fmt.Sprintf(`### FIGMA MAKE: DETERMINISTIC UI
**Component**: %s

---

#### 📐 LAYOUT CONFIGURATION
- **Auto-Layout**: %s
- **Responsive Frames**: Mobile (390px), Tablet (834px), Desktop (1440px)
- **Spacing Scale**: %s

#### 🎨 DESIGN TOKENS (STRICT)
- **Colors**: %s
- **Typography**: %s

#### 🎭 COMPONENT VARIANTS
%s

#### 🔗 INTERACTION FLOWS
- Transition from 'idle' to 'loading' on primary button click.
- Show 'error' variant on validation failure.
- Automate state transitions based on defined State Machine.
`,
		item.Name,
		formatJSON(item.LayoutDefinition),
		strings.Join(item.DesignTokensUsed, ", "),
		"Linked to Design System Registry", // Future integration
		"Linked to Design System Registry",
		formatJSON(item.StateMachine),
	)
}

// GenerateClaudeInFigmaPrompt generates instructions specifically optimized for Claude in Figma
func GenerateClaudeInFigmaPrompt(item *UIRoadmapItem) string {
	return fmt.Sprintf(`### CLAUDE IN FIGMA SPECIFICATION
**Task**: Build a %s based on the following deterministic spec.

---

#### 🌳 COMPONENT TREE (LAYER HIERARCHY)
%s

#### ⚙️ VARIANT CONFIGURATION
Generate Figma variants for the following states:
%s

#### 📐 AUTO-LAYOUT COMMANDS
- Set all containers to 'Hug' or 'Fill' as appropriate.
- Use spacing tokens: %s

#### 🎨 STYLING
- Naming Convention: %s
- Design Tokens: %s

[INSTRUCTION]: Ensure all layers are named according to the component tree hierarchy. Use Figma Auto-Layout 5.0 features.
`,
		item.Name,
		formatJSON(item.ComponentTree),
		formatJSON(item.StateMachine),
		strings.Join(item.DesignTokensUsed, ", "),
		"BEM-style layer naming",
		"Integrated with Figma Design System variables",
	)
}
