package adr

import (
	"bytes"
	"fmt"
	io_document "github.com/therealkevinard/adr-er/io-document"
	output_templates "github.com/therealkevinard/adr-er/output-templates"
	"github.com/therealkevinard/adr-er/utils"
	"strings"
	"text/template"
)

// ADR is an architectural decision record.
type ADR struct {
	Sequence     int
	Title        string
	Context      string
	Decision     string
	Status       string
	Consequences string
}

// BuildDocument returns a CompiledDocument instance from the ADR.
// the returned CompiledDocument hold the rendered content and other metadata.
// it can be written to disk using the CompiledDocument.Write() method
func (adr *ADR) BuildDocument(parsedTemplate *output_templates.ParsedTemplateFile) (*io_document.IODocument, error) {
	// render the document, capturing the content return
	content, err := adr.render(parsedTemplate)
	if err != nil {
		return nil, fmt.Errorf("error rendering ADR: %w", err)
	}

	// return the writeable document
	return io_document.NewIODocument(parsedTemplate, adr.SequencedTitle(), content)
}

// SequencedTitle prefixes title with sequence for display
func (adr *ADR) SequencedTitle() string {
	var docTitle strings.Builder
	docTitle.WriteString(utils.PadValue(adr.Sequence, 4))
	docTitle.WriteString(": ")
	docTitle.WriteString(adr.Title)

	return docTitle.String()
}

// render renders the document using by inputting the ADR to provided gotemplate, returning the content []byte
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
