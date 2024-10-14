package render

import (
	"embed"
	"fmt"
	"strings"

	"github.com/therealkevinard/adr-er/globals"
)

var _ globals.Validator = (*ParsedTemplateFile)(nil)

type TemplateNotFoundError struct {
	TemplateName string
}

func (e TemplateNotFoundError) Error() string {
	return fmt.Sprintf("template %s not found", e.TemplateName)
}

// TemplateFS embeds the *.tpl files in this directory using an embed.FS.
// Templates should be named `{name}.{format}.tpl` according to the naming convention.
// The `go:embed` directive includes all template files matching the pattern.
//
//go:embed *.tpl
var TemplateFS embed.FS

// ListTemplates reads all embedded template files and returns a map where keys are template file paths and values are
// parsed template specifications.
// The parsed templates contain metadata and content for each file.
func ListTemplates() (map[string]*ParsedTemplateFile, error) {
	tpls, err := TemplateFS.ReadDir(".")
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

	return index, nil
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
		return nil, TemplateNotFoundError{TemplateName: defaultName}
	}

	return v, nil
}
