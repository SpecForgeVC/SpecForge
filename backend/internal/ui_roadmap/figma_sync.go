package ui_roadmap

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FigmaSyncSession tracks an active connection between a Figma file and a Roadmap Item
type FigmaSyncSession struct {
	ID            uuid.UUID `json:"id"`
	RoadmapItemID uuid.UUID `json:"roadmap_item_id"`
	FigmaFileKey  string    `json:"figma_file_key"`
	LastSyncedAt  time.Time `json:"last_synced_at"`
	SyncStatus    string    `json:"sync_status"` // connected | detached | syncing
}

// FigmaSyncPayload represents the data sent from the Figma plugin
type FigmaSyncPayload struct {
	FigmaFileKey string          `json:"figma_file_key"`
	Hierarchy    json.RawMessage `json:"hierarchy"`
	DesignTokens json.RawMessage `json:"design_tokens"`
}

// FigmaPluginManifest defines the structure for Figma's manifest.json
type FigmaPluginManifest struct {
	Name       string   `json:"name"`
	ID         string   `json:"id"`
	API        string   `json:"api"`
	Main       string   `json:"main"`
	UI         string   `json:"ui"`
	EditorType []string `json:"editorType"`
}

// SyncFigmaData processes the incoming hierarchy and tokens from Figma
func (s *service) SyncFigmaData(ctx context.Context, itemID uuid.UUID, payload FigmaSyncPayload) error {
	// 1. Fetch the existing item
	item, err := s.repo.Get(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to fetch item for sync: %w", err)
	}

	// 2. Map Hierarchy to ComponentTree
	// In a real implementation, we would transform Figma layers into ComponentNodes.
	// For now, we will update the ComponentTree field with the provided hierarchy.
	item.ComponentTree = payload.Hierarchy

	// 3. Update the item in the database
	if err := s.repo.Update(ctx, item); err != nil {
		return fmt.Errorf("failed to save synced item: %w", err)
	}

	return nil
}

// GeneratePluginManifest returns a default manifest for the Figma plugin
func GeneratePluginManifest(projectID string) FigmaPluginManifest {
	return FigmaPluginManifest{
		Name:       "SpecForge Sync - " + projectID,
		ID:         uuid.New().String(),
		API:        "1.0.0",
		Main:       "code.js",
		UI:         "ui.html",
		EditorType: []string{"figma"},
	}
}

// GeneratePluginCode returns the JavaScript logic for the Figma plugin
func GeneratePluginCode(apiEndpoint string, itemID string) string {
	return fmt.Sprintf(`
const API_URL = "%s";
const ITEM_ID = "%s";

figma.showUI(__html__);

figma.ui.onmessage = async (msg) => {
  if (msg.type === 'sync-hierarchy') {
    const selection = figma.currentPage.selection;
    if (selection.length === 0) {
      figma.notify("Select a frame to sync hierarchy.");
      return;
    }

    const hierarchy = selection.map(node => ({
      id: node.id,
      name: node.name,
      type: node.type,
      children: node.children ? node.children.map(c => c.name) : []
    }));

    try {
      const response = await fetch(API_URL + "/api/v1/ui-roadmap/" + ITEM_ID + "/sync", {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ hierarchy })
      });
      
      if (response.ok) {
        figma.notify("Hierarchy synced to SpecForge!");
      } else {
        figma.notify("Failed to sync. Check API connectivity.");
      }
    } catch (e) {
      figma.notify("Error: " + e.message);
    }
  }
};
`, apiEndpoint, itemID)
}
