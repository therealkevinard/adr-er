package output_templates

import (
	"embed"
	"fmt"
)

//go:embed *.tpl
var Templates embed.FS

// ListTemplates returns a map[string][]byte index of the available templates.
// keys are the file paths, and values are the template contents for preview
func ListTemplates() (map[string][]byte, error) {
	tpls, err := Templates.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("error listing templates: %w", err)
	}

	index := make(map[string][]byte, len(tpls))
	for _, tpl := range tpls {
		// no nested files
		if tpl.IsDir() {
			continue
		}

		// read the content
		content, err := Templates.ReadFile(tpl.Name())
		if err != nil {
			// TODO: report error
			continue
		}

		index[tpl.Name()] = content
	}

	return index, err
}
