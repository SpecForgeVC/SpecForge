package infra

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/scott/specforge/internal/domain"
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
	buf.WriteString(fmt.Sprintf("## Metadata\n- ID: %s\n- Exported At: %s\n- Governance Mode: %s\n\n",
		pkg.Metadata.ArtifactID, pkg.Metadata.ExportedAt.Format("2006-01-02 15:04:05"), pkg.Metadata.GovernanceMode))

	buf.WriteString("## Context\n")
	buf.WriteString(pkg.RoadmapContext.Description + "\n\n")

	buf.WriteString("## Implementation Prompt\n")
	buf.WriteString("```markdown\n")
	buf.WriteString(pkg.BuildPrompts.Implementation)
	buf.WriteString("\n```\n\n")

	buf.WriteString("## Verification Prompt\n")
	buf.WriteString("```markdown\n")
	buf.WriteString(pkg.BuildPrompts.Verification)
	buf.WriteString("\n```\n")

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
