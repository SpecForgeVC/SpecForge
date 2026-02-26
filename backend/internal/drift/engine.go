package drift

import (
	"fmt"
)

type ChangeType string

const (
	Addition ChangeType = "ADDITION"
	Removal  ChangeType = "REMOVAL"
	Modified ChangeType = "MODIFIED"
)

type RiskScore int

const (
	ScoreLow    RiskScore = 1
	ScoreMedium RiskScore = 5
	ScoreHigh   RiskScore = 10
)

type SchemaDiff struct {
	Path        string      `json:"path"`
	Type        ChangeType  `json:"type"`
	OldValue    interface{} `json:"old_value"`
	NewValue    interface{} `json:"new_value"`
	IsBreaking  bool        `json:"is_breaking"`
	Description string      `json:"description"`
	RiskScore   RiskScore   `json:"risk_score"`
}

type DiffEngine interface {
	Compare(oldSchema, newSchema map[string]interface{}) ([]SchemaDiff, error)
}

type diffEngine struct{}

func NewDiffEngine() DiffEngine {
	return &diffEngine{}
}

func (e *diffEngine) Compare(oldSchema, newSchema map[string]interface{}) ([]SchemaDiff, error) {
	var diffs []SchemaDiff
	e.compareRecursive(oldSchema, newSchema, "", &diffs)
	return diffs, nil
}

func (e *diffEngine) compareRecursive(old, new map[string]interface{}, pathPrefix string, diffs *[]SchemaDiff) {
	// Compare Type
	oldType, _ := old["type"].(string)
	newType, _ := new["type"].(string)
	if oldType != "" && newType != "" && oldType != newType {
		typePath := "type"
		if pathPrefix != "" {
			typePath = pathPrefix + ".type"
		}
		*diffs = append(*diffs, SchemaDiff{
			Path:        typePath,
			Type:        Modified,
			OldValue:    oldType,
			NewValue:    newType,
			IsBreaking:  true,
			Description: fmt.Sprintf("Type changed from %s to %s", oldType, newType),
			RiskScore:   ScoreHigh,
		})
		// If type changed, comparison of children might be moot, but let's continue if both are objects
	}

	// Compare Properties
	oldProps, okOld := old["properties"].(map[string]interface{})
	newProps, okNew := new["properties"].(map[string]interface{})

	if okOld {
		if !okNew && newType == "object" {
			// All properties removed? Or properties field removed?
			// If newType is object but no properties, might be valid but effectively removing all specific props
		}

		for name, oldPropVal := range oldProps {
			oldProp, _ := oldPropVal.(map[string]interface{})
			currentPath := name
			if pathPrefix != "" {
				currentPath = fmt.Sprintf("%s.properties.%s", pathPrefix, name)
			} else {
				currentPath = fmt.Sprintf("properties.%s", name)
			}

			newPropVal, exists := newProps[name]
			if !exists {
				*diffs = append(*diffs, SchemaDiff{
					Path:        currentPath,
					Type:        Removal,
					OldValue:    oldProp,
					IsBreaking:  true, // Removal is breaking
					Description: fmt.Sprintf("Property '%s' removed", name),
					RiskScore:   ScoreHigh,
				})
				continue
			}

			newProp, _ := newPropVal.(map[string]interface{})
			e.compareRecursive(oldProp, newProp, currentPath, diffs)
		}
	}

	if okNew {
		for name, newPropVal := range newProps {
			newProp, _ := newPropVal.(map[string]interface{})
			if _, exists := oldProps[name]; !exists {
				currentPath := name
				if pathPrefix != "" {
					currentPath = fmt.Sprintf("%s.properties.%s", pathPrefix, name)
				} else {
					currentPath = fmt.Sprintf("properties.%s", name)
				}

				*diffs = append(*diffs, SchemaDiff{
					Path:        currentPath,
					Type:        Addition,
					NewValue:    newProp,
					IsBreaking:  false, // Addition is generally non-breaking
					Description: fmt.Sprintf("Property '%s' added", name),
					RiskScore:   ScoreLow,
				})
			}
		}
	}

	// Check Required fields changes
	e.checkRequiredFields(old, new, pathPrefix, diffs)

	// Check Array Items
	if newType == "array" && oldType == "array" {
		e.checkArrayItems(old, new, pathPrefix, diffs)
	}
}

func (e *diffEngine) checkRequiredFields(old, new map[string]interface{}, pathPrefix string, diffs *[]SchemaDiff) {
	oldReq, _ := old["required"].([]interface{})
	newReq, _ := new["required"].([]interface{})

	// Convert to map for easy lookup
	oldReqMap := make(map[string]bool)
	for _, r := range oldReq {
		if s, ok := r.(string); ok {
			oldReqMap[s] = true
		}
	}

	for _, r := range newReq {
		s, ok := r.(string)
		if !ok {
			continue
		}
		if !oldReqMap[s] {
			reqPath := "required"
			if pathPrefix != "" {
				reqPath = pathPrefix + ".required"
			}
			// New required field -> Breaking change
			*diffs = append(*diffs, SchemaDiff{
				Path:        reqPath,
				Type:        Modified,
				NewValue:    s,
				IsBreaking:  true,
				Description: fmt.Sprintf("Field '%s' became required", s),
				RiskScore:   ScoreHigh,
			})
		}
	}
}

func (e *diffEngine) checkArrayItems(old, new map[string]interface{}, pathPrefix string, diffs *[]SchemaDiff) {
	oldItems, okOld := old["items"].(map[string]interface{})
	newItems, okNew := new["items"].(map[string]interface{})

	if okOld && okNew {
		currentPath := "items"
		if pathPrefix != "" {
			currentPath = fmt.Sprintf("%s.items", pathPrefix)
		}
		e.compareRecursive(oldItems, newItems, currentPath, diffs)
	}
}
