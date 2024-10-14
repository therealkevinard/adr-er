package render

import (
	"fmt"
	"strings"

	"github.com/therealkevinard/adr-er/globals"
)

// ParsedTemplateFile unpacks the name and format from the template filename and joins with its content.
// once constructed, downstream business logic is entirely decoupled from the filesystem.
type ParsedTemplateFile struct {
	// ID is the template id.
	// This is the sluggified version of the human name, used primarily to create the physical filename.
	ID string
	// Name is the human name of the document
	Name string

	// Format is the DocumentFormat for this template.
	Format DocumentFormat
	// Content holds the templated content bytes, ready for writing
	Content []byte
}

// Validate validates the object, returning any error.
func (t *ParsedTemplateFile) Validate() error {
	// nested validators
	if err := t.Format.Validate(); err != nil {
		return fmt.Errorf("invalid ouput format: %w", err)
	}

	// self validators
	if t.ID == "" {
		return globals.ValidationError("id", "empty template id")
	}

	if t.Name == "" {
		return globals.ValidationError("name", "empty name")
	}

	if len(t.Content) == 0 {
		return globals.ValidationError("content", "empty content")
	}

	return nil
}

// parseTemplate parses a template file name according to the `{name}.{format}.tpl` naming convention.
// Returns a ParsedTemplateFile with extracted metadata and content, or nil if parsing fails.
func parseTemplate(filename string) *ParsedTemplateFile {
	parts := strings.Split(filename, ".")
	// exactly 3 parts are expected
	//nolint:mnd // this isn't magic, it's from the regex capture
	if len(parts) != 3 {
		return nil
	}

	parsed := &ParsedTemplateFile{
		ID:      parts[0],
		Format:  DocumentFormat(parts[1]),
		Name:    filename,
		Content: nil,
	}

	var err error
	if parsed.Content, err = TemplateFS.ReadFile(parsed.Name); err != nil {
		// error reading, continue
		return nil
	}

	// invalid, continue
	if err = parsed.Validate(); err != nil {
		return nil
	}

	return parsed
}
