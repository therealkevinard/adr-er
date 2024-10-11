package adr

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	io_document "github.com/therealkevinard/adr-er/io-document"
	output_templates "github.com/therealkevinard/adr-er/output-templates"
	"github.com/therealkevinard/adr-er/utils"
)

const (
	numericPadWidth = 4
)

// ADR represents an Architectural Decision Record (ADR).
// It stores details about decisions made during software architecture design.
type ADR struct {
	Sequence     int
	Title        string
	Context      string
	Decision     string
	Status       string
	Consequences string
}

// BuildDocument creates an IODocument from the ADR using the provided template.
// It renders the ADR content into the template and returns a document that can be written to disk.
// Returns an error if rendering fails or if the template is invalid.
func (adr *ADR) BuildDocument(parsedTemplate *output_templates.ParsedTemplateFile) (*io_document.IODocument, error) {
	// render the document, capturing the content return
	content, err := adr.render(parsedTemplate)
	if err != nil {
		return nil, fmt.Errorf("error rendering ADR: %w", err)
	}

	// return the writeable document
	return io_document.NewIODocument(parsedTemplate, adr.SequencedTitle(), content)
}

// SequencedTitle returns the ADR's title prefixed with its sequence number.
// This is used for display purposes to distinguish between different ADRs.
func (adr *ADR) SequencedTitle() string {
	var docTitle strings.Builder
	docTitle.WriteString(utils.PadValue(adr.Sequence, numericPadWidth))
	docTitle.WriteString(": ")
	docTitle.WriteString(adr.Title)

	return docTitle.String()
}

// render processes the ADR through the provided ParsedTemplateFile, generating the rendered content.
// Returns the rendered content as a byte slice, or an error if the template is invalid or rendering fails.
func (adr *ADR) render(parsedTemplate *output_templates.ParsedTemplateFile) ([]byte, error) {
	if err := parsedTemplate.Validate(); err != nil {
		return nil, fmt.Errorf("refusing to render invalid template: %w", err)
	}

	tpl, err := template.New(parsedTemplate.ID).Parse(string(parsedTemplate.Content))
	if err != nil {
		return nil, fmt.Errorf("error preparing template: %w", err)
	}

	var content bytes.Buffer
	if err = tpl.Execute(&content, adr); err != nil {
		return nil, fmt.Errorf("error rendering document: %w", err)
	}

	return content.Bytes(), nil
}
