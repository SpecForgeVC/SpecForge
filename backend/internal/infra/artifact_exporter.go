package infra

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/SpecForgeVC/SpecForge/internal/domain"
)

type ArtifactExporter interface {
	Export(pkg *domain.BuildArtifactPackage, format domain.ExportFormat) ([]byte, string, error)
}

type artifactExporter struct{}

func NewArtifactExporter() ArtifactExporter {
	return &artifactExporter{}
}

func (e *artifactExporter) Export(pkg *domain.BuildArtifactPackage, format domain.ExportFormat) ([]byte, string, error) {
	switch format {
	case domain.ExportFormatJSON:
		data, err := json.MarshalIndent(pkg, "", "  ")
		if err != nil {
			return nil, "", err
		}
		return data, "application/json", nil

	case domain.ExportFormatMarkdown:
		return e.exportMarkdown(pkg)

	case domain.ExportFormatZip:
		return e.exportZip(pkg)

	default:
		return nil, "", fmt.Errorf("unsupported format: %s", format)
	}
}

func (e *artifactExporter) exportMarkdown(pkg *domain.BuildArtifactPackage) ([]byte, string, error) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("# Build Artifact: %s\n\n", pkg.RoadmapContext.Title))
	buf.WriteString(fmt.Sprintf("## Metadata\n- **Artifact ID**: %s\n- **Roadmap Item**: %s\n- **Exported At**: %s\n- **Governance Mode**: %s\n- **Integrity Hash**: `%s`\n\n",
		pkg.Metadata.ArtifactID, pkg.Metadata.RoadmapItemID, pkg.Metadata.ExportedAt.Format("2006-01-02 15:04:05"), pkg.Metadata.GovernanceMode, pkg.Metadata.IntegrityHash))

	buf.WriteString("## Readiness\n")
	buf.WriteString(fmt.Sprintf("- **Score**: %d%%\n", pkg.RoadmapContext.ReadinessScore))
	buf.WriteString(fmt.Sprintf("- **Level**: %s\n", pkg.RoadmapContext.ReadinessLevel))
	buf.WriteString(fmt.Sprintf("- **Priority**: %s\n", pkg.RoadmapContext.Priority))
	buf.WriteString(fmt.Sprintf("- **Risk Level**: %s\n\n", pkg.RoadmapContext.RiskLevel))

	buf.WriteString("## Context\n")
	buf.WriteString(pkg.RoadmapContext.Description + "\n\n")
	if pkg.RoadmapContext.BusinessContext != "" {
		buf.WriteString("### Business Context\n" + pkg.RoadmapContext.BusinessContext + "\n\n")
	}
	if pkg.RoadmapContext.TechnicalContext != "" {
		buf.WriteString("### Technical Context\n" + pkg.RoadmapContext.TechnicalContext + "\n\n")
	}

	buf.WriteString("## Implementation Prompt\n")
	buf.WriteString("```markdown\n")
	buf.WriteString(pkg.BuildPrompts.Implementation)
	buf.WriteString("\n```\n\n")

	buf.WriteString("## Verification Prompt\n")
	buf.WriteString("```markdown\n")
	buf.WriteString(pkg.BuildPrompts.Verification)
	buf.WriteString("\n```\n\n")

	buf.WriteString("## Refinement Instructions\n")
	buf.WriteString(pkg.RefinementLoopPrompts.Instructions + "\n")

	return buf.Bytes(), "text/markdown", nil
}

func (e *artifactExporter) exportZip(pkg *domain.BuildArtifactPackage) ([]byte, string, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// metadata.json
	metaData, _ := json.MarshalIndent(pkg.Metadata, "", "  ")
	e.addToZip(w, "metadata.json", metaData)

	// roadmap-context.md
	contextData := fmt.Sprintf("# %s\n\n%s\n", pkg.RoadmapContext.Title, pkg.RoadmapContext.Description)
	e.addToZip(w, "roadmap-context.md", []byte(contextData))

	// Prompts
	e.addToZip(w, "prompts/implementation.md", []byte(pkg.BuildPrompts.Implementation))
	e.addToZip(w, "prompts/verification.md", []byte(pkg.BuildPrompts.Verification))
	e.addToZip(w, "prompts/refinement.md", []byte(pkg.RefinementLoopPrompts.Instructions))

	// Contracts
	for _, c := range pkg.Contracts {
		data, _ := json.MarshalIndent(c, "", "  ")
		e.addToZip(w, fmt.Sprintf("contracts/%s.json", c.ID), data)
	}

	// Schemas (as part of package, or separate files if needed)
	// For now just include the main JSON as well
	fullPkgData, _ := json.MarshalIndent(pkg, "", "  ")
	e.addToZip(w, "build-artifact.json", fullPkgData)

	if err := w.Close(); err != nil {
		return nil, "", err
	}

	return buf.Bytes(), "application/zip", nil
}

func (e *artifactExporter) addToZip(w *zip.Writer, filename string, content []byte) {
	f, err := w.Create(filename)
	if err != nil {
		return
	}
	f.Write(content)
}
