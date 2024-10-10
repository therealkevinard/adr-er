package output_templates

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path"
	"testing"
)

// asserts fundamental behavior of the template list
func TestListTemplates(t *testing.T) {
	tpls, err := ListTemplates()
	require.NoError(t, err)
	require.NotNil(t, tpls)

	defaultTmpl := tpls["adr-markdown-001.md.tpl"]
	assert.NotEmpty(t, defaultTmpl)
}

// ensures embed.FS contains only .tpl files (eg: no .go files)
func TestListTemplates_OnlyTemplates(t *testing.T) {
	tpls, err := ListTemplates()
	require.NoError(t, err)
	require.NotNil(t, tpls)

	// each should explicitly end with .tpl
	for k := range tpls {
		ext := path.Ext(k)
		assert.Equal(t, ".tpl", ext)
	}
}
