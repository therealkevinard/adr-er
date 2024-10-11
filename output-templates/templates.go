package output_templates

import (
	"embed"
	"fmt"
	"github.com/therealkevinard/adr-er/globals"
	"strings"
)

var _ globals.Validator = (*ParsedTemplateFile)(nil)

// Templates embeds the *.tpl files in this directory using an embed.FS.
// Templates should be named `{name}.{format}.tpl` according to the naming convention.
// The `go:embed` directive includes all template files matching the pattern.
//
//go:embed *.tpl
var Templates embed.FS

// ListTemplates reads all embedded template files and returns a map where keys are template file paths and values are
// parsed template specifications.
// The parsed templates contain metadata and content for each file.
func ListTemplates() (map[string]*ParsedTemplateFile, error) {
	tpls, err := Templates.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("error listing templates: %w", err)
	}

	index := make(map[string]*ParsedTemplateFile, len(tpls))
	for _, tpl := range tpls {
		// no nested files are allowed
		if tpl.IsDir() {
			continue
		}

		// parse the name, getting important metadata and its content
		parsed := parseTemplate(tpl.Name())
		if parsed == nil {
			continue
		}

		index[tpl.Name()] = parsed
	}

	return index, err
}

// DefaultTemplateForFormat retrieves the default template for a given DocumentFormat.
// Default templates are expected to follow the naming pattern "default.{format}.tpl".
// Returns an error if no template matching the default pattern is found.
func DefaultTemplateForFormat(format DocumentFormat) (*ParsedTemplateFile, error) {
	tpls, err := ListTemplates()
	if err != nil {
		return nil, fmt.Errorf("error listing templates: %w", err)
	}

	defaultName := strings.Join([]string{"default", string(format), "tpl"}, ".")

	v, ok := tpls[defaultName]
	if !ok {
		return nil, fmt.Errorf("template %q not found", defaultName)
	}

	return v, nil
}

// ParsedTemplateFile unpacks the name and format from the template filename and joins with its content.
// once constructed, downstream business logic is entirely decoupled from the filesystem.
type ParsedTemplateFile struct {
	ID     string
	Format DocumentFormat

	Name    string
	Content []byte
}

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
	if len(parts) != 3 {
		return nil
	}

	parsed := &ParsedTemplateFile{
		ID:     parts[0],
		Format: DocumentFormat(parts[1]),
		Name:   filename,
	}

	var err error
	if parsed.Content, err = Templates.ReadFile(parsed.Name); err != nil {
		// error reading, continue
		return nil
	}

	// invalid, continue
	if err = parsed.Validate(); err != nil {
		return nil
	}

	return parsed
}
