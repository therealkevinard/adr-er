package io_document

import (
	"fmt"
	"github.com/therealkevinard/adr-er/globals"
	output_templates "github.com/therealkevinard/adr-er/output-templates"
	"github.com/therealkevinard/adr-er/utils"
	"os"
	"path"
)

var _ globals.Validator = (*IODocument)(nil)

// IODocument represents a document that is templated and prepared for filesystem operations.
// It provides methods for validation, deriving filenames, and writing content to disk.
type IODocument struct {
	// Title is the document title. it's mostly used for presentation. business logic relies on DocumentID, which is derived from Title.
	Title string
	// Content is the literal content of the document
	Content []byte
	// Template holds the parsed template, containing file metadata
	Template *output_templates.ParsedTemplateFile
}

// NewIODocument creates a new IODocument with the given parsed template, title, and content.
// It validates the provided template and document before returning.
// Returns an error if validation fails.
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

// Validate checks the properties of an IODocument, including title, content, and template metadata.
// Returns an error if any of the required properties are missing or invalid.
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

// DocumentID is a getter for the derived document id. returns the slugified title.
func (cd IODocument) DocumentID() string {
	return utils.Slugify(cd.Title)
}

// Write attempts to write the document content to a file on disk within the directory <inDir>.
// It first validates the document before creating the file. Returns an error if validation or writing fails.
// TODO: Ensure the method does not overwrite existing files.
func (cd *IODocument) Write(inDir string) error {
	var err error
	if inDir == "" {
		return globals.ValidationError("directory", "directory is empty")
	}

	// validate the document before running any io
	if err = cd.Validate(); err != nil {
		return fmt.Errorf("not writing. document validation failed: %w", err)
	}

	// make the file
	// TODO: make sure this errors if the file already exists. don't want to force-replace existing files.
	f, err := os.Create(path.Join(inDir, cd.Filename()))
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
