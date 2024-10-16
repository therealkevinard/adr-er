package io_document

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therealkevinard/adr-er/globals"
	"github.com/therealkevinard/adr-er/render"
)

// TestConstructor guarantees the inline validation behavior of NewIODocument.
func TestConstructor(t *testing.T) {
	validTemplate := testGetDefaultTemplate(t)

	//nolint:thelper // false positive. these aren't helpers, they _are_ the tests
	tests := []struct {
		name   string
		assert func(t *testing.T)
	}{
		{
			name: "happy path",
			assert: func(t *testing.T) {
				doc, err := NewIODocument(validTemplate, "title", []byte("content"))
				require.NoError(t, err)
				assert.NotNil(t, doc)
			},
		},
		{
			name: "missing title",
			assert: func(t *testing.T) {
				doc, err := NewIODocument(validTemplate, "", []byte("content"))
				assert.Nil(t, doc)
				require.Error(t, err)

				var typedErr globals.InputValidationError
				require.ErrorAs(t, err, &typedErr)
				assert.Equal(t, "title", typedErr.Field)
				assert.Equal(t, "title is empty", typedErr.Reason)
			},
		},
		{
			name: "missing content",
			assert: func(t *testing.T) {
				doc, err := NewIODocument(validTemplate, "title", []byte(""))
				assert.Nil(t, doc)
				require.Error(t, err)

				var typedErr globals.InputValidationError
				require.ErrorAs(t, err, &typedErr)
				assert.Equal(t, "content", typedErr.Field)
				assert.Equal(t, "content is empty", typedErr.Reason)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, test.assert)
	}
}

// TestValidate covers the various validations within .Validate().
func TestValidate(t *testing.T) {
	// newDoc replicated NewIODocument, but without inline validation. this allows testing validation directly
	newDoc := func(title string, content string, template *render.ParsedTemplateFile) IODocument {
		return IODocument{
			Title:    title,
			Content:  []byte(content),
			Template: template,
		}
	}

	//nolint:thelper // false positive. these aren't helpers, they _are_ the tests
	tests := []struct {
		name       string
		document   IODocument
		assertFunc func(t *testing.T, err error)
	}{
		{
			name:     "no error",
			document: newDoc("test document", "test content", testGetDefaultTemplate(t)),
			assertFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "invalid title",
			document: newDoc(
				"",
				"test content",
				testGetDefaultTemplate(t),
			),
			assertFunc: func(t *testing.T, err error) {
				require.Error(t, err)

				var ive globals.InputValidationError
				require.ErrorAs(t, err, &ive)
				assert.Equal(t, "title", ive.Field)
				assert.Equal(t, "title is empty", ive.Reason)
			},
		},
		{
			name: "invalid content",
			document: newDoc(
				"test document",
				"",
				testGetDefaultTemplate(t),
			),
			assertFunc: func(t *testing.T, err error) {
				require.Error(t, err)

				var ive globals.InputValidationError
				require.ErrorAs(t, err, &ive)
				assert.Equal(t, "content", ive.Field)
				assert.Equal(t, "content is empty", ive.Reason)
			},
		},
		{
			name: "invalid format",
			document: newDoc(
				"test document",
				"test content",
				testBreakValidTemplate(t, func(tpl *render.ParsedTemplateFile) *render.ParsedTemplateFile {
					tpl.Format = render.DocumentFormat("<invalid>")

					return tpl
				}),
			),
			assertFunc: func(t *testing.T, err error) {
				require.Error(t, err)

				var ive globals.InputValidationError
				require.ErrorAs(t, err, &ive)
				assert.Equal(t, "template", ive.Field)
				assert.Equal(t, "invalid ouput format: format failed validation: unsupported format", ive.Reason)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.document.Validate()
			test.assertFunc(t, err)
		})
	}
}

// testGetDefaultTemplate returns the default markdown template for testing purposes.
func testGetDefaultTemplate(t *testing.T) *render.ParsedTemplateFile {
	t.Helper()

	found, err := render.DefaultTemplateForFormat(render.DocumentFormatMarkdown)
	require.NoError(t, err)
	require.NotNil(t, found)

	return found
}

// testBreakValidTemplate supports testing invalid paths.
// it loads a valid template, mutates it with breakFunc, and returns the now-invalid template.
func testBreakValidTemplate(
	t *testing.T,
	breakFunc func(tpl *render.ParsedTemplateFile) *render.ParsedTemplateFile,
) *render.ParsedTemplateFile {
	t.Helper()

	valid := testGetDefaultTemplate(t)

	return breakFunc(valid)
}
