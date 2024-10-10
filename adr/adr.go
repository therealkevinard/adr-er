package adr

import (
	"bytes"
	"fmt"
	"github.com/therealkevinard/adr-er/globals"
	io_document "github.com/therealkevinard/adr-er/io-document"
	output_templates "github.com/therealkevinard/adr-er/output-templates"
	"github.com/therealkevinard/adr-er/utils"
	"strings"
	"text/template"
)

// defaultTemplatesMap maps DocumentFormat to a gotemplate.
// overkill as long as we only have markdown, but great when we add whatever else.
var defaultTemplatesMap = map[io_document.DocumentFormat]string{
	io_document.DocumentFormatMarkdown: "default.markdown.tpl",
}

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
func (adr *ADR) BuildDocument(format io_document.DocumentFormat) (*io_document.IODocument, error) {
	// render the document, capturing the content return
	content, err := adr.render(format)
	if err != nil {
		return nil, fmt.Errorf("error rendering adr: %w", err)
	}

	// return the writeable document
	return io_document.NewIODocument(format, adr.SequencedTitle(), content)
}

func (adr *ADR) SequencedTitle() string {
	var docTitle strings.Builder
	docTitle.WriteString(utils.PadValue(adr.Sequence, 4))
	docTitle.WriteString(": ")
	docTitle.WriteString(adr.Title)

	return docTitle.String()
}

// render renders the document using by inputting the adr to provided gotemplate, returning the content []byte
func (adr *ADR) render(format io_document.DocumentFormat) ([]byte, error) {
	if valid := format.Valid(); !valid {
		return nil, globals.ValidationError("format", "invalid format provided")
	}

	tplPath, ok := defaultTemplatesMap[format]
	if !ok {
		return nil, globals.TemplateNotFoundError{Requested: string(format)}
	}

	tpl, err := template.ParseFS(output_templates.Templates, tplPath)
	if err != nil {
		return nil, fmt.Errorf("error parsing embedded template: %w", err)
	}

	var content bytes.Buffer
	if err = tpl.Execute(&content, adr); err != nil {
		return nil, fmt.Errorf("error rendering document: %w", err)
	}

	return content.Bytes(), nil
}
