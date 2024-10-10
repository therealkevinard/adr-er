package io_document

import (
	"fmt"
	"github.com/therealkevinard/adr-er/globals"
	output_templates "github.com/therealkevinard/adr-er/output-templates"
	"github.com/therealkevinard/adr-er/utils"
	"os"
)

var _ globals.Validator = (*IODocument)(nil)

// IODocument supports filesystem io.
// It exports methods for deriving the filename, validation, and actual filesystem IO.
type IODocument struct {
	// Title is the document title. it's mostly used for presentation. business logic relies on DocumentID, which is derived from Title.
	Title string
	// Template holds the parsed template, containing file metadata
	Template *output_templates.ParsedTemplateFile
	// Content is the literal content of the document
	Content []byte
}

// NewIODocument returns a IODocument using provided format, title, and content.
// the created document is validated before returning. only a valid document is returned.
func NewIODocument(parsedTemplate *output_templates.ParsedTemplateFile, title string, content []byte) (*IODocument, error) {
	if err := parsedTemplate.Validate(); err != nil {
		return nil, fmt.Errorf("invalid template. refusing IODocument: %w", err)
	}

	cd := &IODocument{
		Title:    title,
		Template: parsedTemplate,
		Content:  content,
	}
	if err := cd.Validate(); err != nil {
		return nil, fmt.Errorf("refusing to create invalid document: %w", err)
	}

	return cd, nil
}

// Validate checks several properties of IODocument, returning errors on failure.
// many validations are consolidated here, so many operations can be gated behind this one validator
func (cd *IODocument) Validate() error {
	// base template
	if err := cd.Template.Validate(); err != nil {
		return globals.ValidationError("template", err.Error())
	}

	// file body validations
	{
		// has title?
		if len(cd.Title) == 0 {
			return globals.ValidationError("title", "title is empty")
		}
		// has content?
		if len(cd.Content) == 0 {
			return globals.ValidationError("content", "content is empty")
		}
	}

	// file name validations
	// these are extreme edge-cases, as they're derived downstream from format and/or title, both of which have already been checked.
	{
		if cd.DocumentID() == "" {
			return globals.ValidationError("documentID", "can't create valid document id from title")
		}
		if cd.Template.Format.Extension() == "" {
			return globals.ValidationError("extension", "can't create valid document extension")
		}
		if cd.Filename() == "" {
			return globals.ValidationError("filename", "can't create valid document filename")
		}
	}

	return nil
}

// Filename is a getter for the derived filename
// filename is built from sluggified title and format's extension.
func (cd IODocument) Filename() string {
	return cd.DocumentID() + "." + cd.Template.Format.Extension()
}

// DocumentID is a getter for the derived document id. returns the slugified title
func (cd IODocument) DocumentID() string {
	return utils.Slugify(cd.Title)
}

// Write flushes the document to filesystem
func (cd *IODocument) Write() error {
	var err error

	// validate the document before running any io
	if err = cd.Validate(); err != nil {
		return fmt.Errorf("not writing. document validation failed: %w", err)
	}

	// make the file
	// TODO: make sure this errors if the file already exists. don't want to force-replace existing files.
	f, err := os.Create(cd.Filename())
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", cd.Filename(), err)
	}
	defer f.Close()

	// write it
	if _, err = f.Write(cd.Content); err != nil {
		return fmt.Errorf("could not write to file %s: %w", cd.Filename(), err)
	}

	// donesies
	return nil
}
