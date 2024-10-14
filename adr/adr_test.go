package adr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therealkevinard/adr-er/globals"
	"github.com/therealkevinard/adr-er/render"
)

func TestBuildDocument(t *testing.T) {
	defaultTemplate, tplErr := render.DefaultTemplateForFormat(render.DocumentFormatMarkdown)
	require.NoError(t, tplErr)
	require.NotNil(t, defaultTemplate)

	//nolint:thelper // false positive. these aren't helpers, they _are_ the tests
	tests := []struct {
		name       string
		adr        *ADR
		assertFunc func(t *testing.T, adr *ADR)
	}{
		{
			name: "happy path",
			adr: &ADR{
				Sequence:     1,
				Title:        "<title>",
				Context:      "<context>",
				Decision:     "<decision>",
				Status:       "<status>",
				Consequences: "<consequences>",
			},
			assertFunc: func(t *testing.T, adr *ADR) {
				doc, err := adr.BuildDocument(defaultTemplate)
				require.NoError(t, err)

				assert.Equal(t, "0001: <title>", doc.Title)
				assert.Equal(t, "0001-title", doc.DocumentID())
				assert.Equal(t, "0001-title.md", doc.Filename())
				assert.Equal(t, render.DocumentFormatMarkdown, doc.Template.Format)
				assert.NotEmpty(t, doc.Content)
			},
		},
		{
			name: "invalid format",
			adr: &ADR{
				Sequence: 1,
				Title:    "<title>",
			},
			assertFunc: func(t *testing.T, adr *ADR) {
				doc, err := adr.BuildDocument(&render.ParsedTemplateFile{
					ID:      "default",
					Format:  render.DocumentFormat("invalid"),
					Name:    "default.invalid.tpl",
					Content: []byte("content"),
				})
				assert.Nil(t, doc)
				require.Error(t, err)

				// assert against the returned descriptive error
				var validationError globals.InputValidationError
				ok := errors.As(err, &validationError)
				assert.True(t, ok)
				assert.Equal(t, "format", validationError.Field)
				assert.Equal(t, "unsupported format", validationError.Reason)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.assertFunc(t, test.adr)
		})
	}
}
