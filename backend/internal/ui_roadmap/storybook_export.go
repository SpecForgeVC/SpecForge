package ui_roadmap

import (
	"fmt"
)

// GenerateStorybookScaffold generates a Storybook template with stories for each state
func GenerateStorybookScaffold(item *UIRoadmapItem) string {
	return fmt.Sprintf(`### 📚 STORYBOOK SCAFFOLD: %s

#### 🎭 COMPONENT STORIES (STATES)
- **Idle**: Default baseline visualization.
- **Loading**: With loading indicators and disabled interactions.
- **Success**: Post-action successful state.
- **Error**: With validation messaging and error boundaries.
- **Empty**: Zero-state representation.
- **Disabled**: Read-only/Inactive state.

#### 🧪 INTERACTION TESTS (PLAY FUNCTIONS)
%s

#### ♿ ACCESSIBILITY AUDIT HINTS
- ARIA Role: %s
- Screen Reader: %s
- Visual Regression: Test all responsive breakpoints.

#### 🛠️ BOOTSTRAP COMMAND
"npx storybook@latest add @storybook/addon-interactions"
`,
		item.Name,
		formatJSON(item.StateMachine),
		"Check Accessibility Spec",
		"Verify Alt Text and Labels",
	)
}
