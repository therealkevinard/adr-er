package render

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// asserts fundamental behavior of ListTemplates, parseTemplate, and DefaultTemplateForFormat.
func TestListTemplates(t *testing.T) {
	tpls, err := ListTemplates()
	require.NoError(t, err)
	require.NotNil(t, tpls)

	found, err := DefaultTemplateForFormat(DocumentFormatMarkdown)
	require.NoError(t, err)
	require.NotNil(t, found)

	assert.NotNil(t, found)
	assert.Equal(t, "default.markdown.tpl", found.Name)
	assert.Equal(t, DocumentFormatMarkdown, found.Format)
	assert.Equal(t, "default", found.ID)
}

// ensures embed.FS contains only .tpl files (eg: no .go files).
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
