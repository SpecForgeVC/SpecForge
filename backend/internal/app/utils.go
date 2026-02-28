package app

import (
	"strings"
)

// CleanJSON strips markdown code blocks, comments, and extracts JSON from mixed content
func CleanJSON(input string) string {
	// First strip comments (simple line-by-line approach)
	lines := strings.Split(input, "\n")
	var cleanedLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") {
			continue
		}
		// Check for inline comments (very basic, might break if // is in a string)
		if idx := strings.Index(line, " //"); idx != -1 {
			line = line[:idx]
		}
		cleanedLines = append(cleanedLines, line)
	}
	input = strings.Join(cleanedLines, "\n")
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
