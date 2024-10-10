package output_templates

import (
	"embed"
	"fmt"
	io_document "github.com/therealkevinard/adr-er/io-document"
	"strings"
)

// Templates embeds the *.tpl files in this directory in an embed.FS
// !!contract: templates should be named `{name}.{format}.tpl`.
// see parseTemplate for details on the naming convention.
//
//go:embed *.tpl
var Templates embed.FS

// ListTemplates returns a map[string][]byte index of the available templates.
// keys are the file paths, and values are the template contents for preview
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

// DefaultTemplateForFormat returns the default template file for a given format. according
// to the naming convention, default template for (eg) markdown would be named literally "default.markdown.tpl"
func DefaultTemplateForFormat(format io_document.DocumentFormat) (*ParsedTemplateFile, error) {
	tpls, err := ListTemplates()
	if err != nil {
		return nil, fmt.Errorf("error listing templates: %w", err)
	}

	defaultName := strings.Join([]string{
		"default",
		string(format),
		"tpl",
	}, ".")

	v, ok := tpls[defaultName]
	if !ok {
		return nil, fmt.Errorf("template %q not found", defaultName)
	}

	return v, nil
}

// ParsedTemplateFile unpacks the name and format from the template filename and joins with its content
type ParsedTemplateFile struct {
	ID     string
	Format io_document.DocumentFormat

	Name    string
	Content []byte
}

// parseTemplate parses a template name according to the naming convention, returning its spec.
// this func returns no errors, only nil.
// TODO: maybe return errors one day
func parseTemplate(filename string) *ParsedTemplateFile {
	parts := strings.Split(filename, ".")
	// exactly 3 parts are expected
	if len(parts) != 3 {
		return nil
	}

	parsed := &ParsedTemplateFile{
		ID:     parts[0],
		Format: io_document.DocumentFormat(parts[1]),
		Name:   filename,
	}

	// invalid, continue
	if !parsed.Format.Valid() || parsed.ID == "" {
		return nil
	}

	var err error
	if parsed.Content, err = Templates.ReadFile(parsed.Name); err != nil {
		// error reading, continue
		return nil
	}

	return parsed
}
