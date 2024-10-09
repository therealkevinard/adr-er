package adr

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
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

// ADR is an architectural decision record.
type ADR struct {
	Sequence     int    `json:"sequence"`
	Title        string `json:"title"`
	Context      string `json:"context"`
	Decision     string `json:"decision"`
	Status       string `json:"status"`
	Consequences string `json:"consequences"`
}

// CompileDocument returns a CompiledDocument instance from the ADR.
// the returned CompiledDocument hold the rendered content and other metadata.
// it can be written to disk using the CompiledDocument.Write() method
func (adr *ADR) CompileDocument(format DocumentFormat) (*CompiledDocument, error) {
	// validate input
	if valid := format.Valid(); !valid {
		return nil, errors.New("invalid format")
	}

	// init the cd object
	cd := &CompiledDocument{
		DocumentID: adr.documentID(),

		Format:    format,
		extension: format.GetExtension(),
		Content:   nil,
	}

	// render the document, capturing the content return
	var err error
	cd.Content, err = adr.renderDocument(format)
	if err != nil {
		return nil, fmt.Errorf("error rendering document: %w", err)
	}

	return cd, nil
}

// documentID returns a complete document id, including padded sequence number. can be used as a filename or index key
func (adr *ADR) documentID() string {
	return strings.Join([]string{adr.padSequence(3), adr.slugifyTitle()}, "_")
}

// renderDocument renders the ADR as a full document.
func (adr *ADR) renderDocument(format DocumentFormat) ([]byte, error) {
	var err error

	// formatTemplateMap maps DocumentFormat to a gotemplate.
	// overkill as long as we only have markdown, but great when we add whatever else.
	var formatTemplateMap = map[DocumentFormat]string{
		DocumentFormatMarkdown: renderTemplateStyleMarkdown,
	}

	// check the format
	if valid := format.Valid(); !valid {
		return nil, errors.New("invalid format")
	}

	// choose a template
	tpl, ok := formatTemplateMap[format]
	if !ok {
		return nil, errors.New("no usable template found")
	}

	// parse the tpl
	tmpl, err := template.New("adr-markdown").Parse(tpl)
	if err != nil {
		return nil, errors.New("error parsing markdown template")
	}

	// execute it
	var result bytes.Buffer
	if err = tmpl.Execute(&result, adr); err != nil {
		return nil, errors.New("error executing template")
	}

	// return the bytes
	return result.Bytes(), nil
}

// padSequence pads adr.Sequence with leading zeros up to width
func (adr *ADR) padSequence(width int) string {
	// TODO: make sure this int-bump is accomodated when we start computing a next based on fs inputs
	if adr.Sequence == 0 {
		adr.Sequence = 1
	}

	return fmt.Sprintf("%0*d", width, adr.Sequence)
}

// slugify converts a string into a slug by lowercasing, replacing spaces with hyphens, and removing special characters.
func (adr *ADR) slugifyTitle() string {
	slug := adr.Title                 // initial value
	slug = strings.ToLower(adr.Title) // tolower

	// reduce single or consecutive whitespace to a hyphen
	spaceRe := regexp.MustCompile(`\s+`)
	slug = spaceRe.ReplaceAllString(slug, "-")

	// keep only letters, numbers, and hyphens
	re := regexp.MustCompile(`[^a-z0-9-]+`)
	slug = re.ReplaceAllString(slug, "")

	// trim for tidiness
	slug = strings.Trim(slug, "-")

	return slug
}
