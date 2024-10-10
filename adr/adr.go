package adr

import (
	"bytes"
	"fmt"
	"github.com/therealkevinard/adr-er/globals"
	io_document "github.com/therealkevinard/adr-er/io-document"
	"github.com/therealkevinard/adr-er/utils"
	"strings"
	"text/template"
)

// renderTemplateStyleMarkdown holds the gotemplate to render an ADR as markdown.
// TODO(reminder): to support edit/update, this format needs to be parseable into structured data.
const renderTemplateStyleMarkdown = `
{{.Title}} 
--- 

## Status: {{.Status}}

## Context  
{{.Context}}

## Decision  
{{.Decision}}

## Consequences  
{{.Consequences}}
`

// defaultTemplatesMap maps DocumentFormat to a gotemplate.
// overkill as long as we only have markdown, but great when we add whatever else.
var defaultTemplatesMap = map[io_document.DocumentFormat]string{
	io_document.DocumentFormatMarkdown: renderTemplateStyleMarkdown,
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
	// validate input
	if valid := format.Valid(); !valid {
		return nil, globals.ValidationError("format", "invalid format provided")
	}

	// render the document, capturing the content return
	// ... choose a template
	tpl, ok := defaultTemplatesMap[format]
	if !ok {
		return nil, globals.TemplateNotFoundError{Requested: string(format)}
	}
	// ... render it
	content, err := adr.render(tpl)
	if err != nil {
		return nil, fmt.Errorf("error rendering document: %w", err)
	}

	// build document title: `{padded sequence}: {title}`
	var docTitle strings.Builder
	docTitle.WriteString(utils.PadValue(adr.Sequence, 4))
	docTitle.WriteString(": ")
	docTitle.WriteString(adr.Title)

	// return the writeable document
	return io_document.NewIODocument(format, docTitle.String(), content)
}

// render renders the document using by inputting the adr to provided gotemplate, returning the content []byte
func (adr *ADR) render(tpl string) ([]byte, error) {
	var err error

	// parse the tpl
	tmpl, err := template.New("adr-doc-tpl").Parse(tpl)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %w", err)
	}

	// execute and return
	var result bytes.Buffer
	if err = tmpl.Execute(&result, adr); err != nil {
		return nil, fmt.Errorf("error executing template: %w", err)
	}

	return result.Bytes(), nil
}
