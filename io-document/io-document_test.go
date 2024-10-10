package io_document

import (
	"github.com/stretchr/testify/assert"
	"github.com/therealkevinard/adr-er/globals"
	"testing"
)

// TestConstructor guarantees the inline validation behavior of NewIODocument
func TestConstructor(t *testing.T) {
	tests := []struct {
		name   string
		assert func(t *testing.T)
	}{
		{
			name: "happy path",
			assert: func(t *testing.T) {
				doc, err := NewIODocument(DocumentFormatMarkdown, "title", []byte("content"))
				assert.Nil(t, err)
				assert.NotNil(t, doc)
			},
		},
		{
			name: "missing title",
			assert: func(t *testing.T) {
				doc, err := NewIODocument(DocumentFormatMarkdown, "", []byte("content"))
				assert.Nil(t, doc)
				assert.NotNil(t, err)

				var genericerr globals.GenericInputValidationError
				assert.ErrorAs(t, err, &genericerr)
			},
		},
		{
			name: "missing content",
			assert: func(t *testing.T) {
				doc, err := NewIODocument(DocumentFormatMarkdown, "title", []byte(""))
				assert.Nil(t, doc)
				assert.NotNil(t, err)

				var genericerr globals.GenericInputValidationError
				assert.ErrorAs(t, err, &genericerr)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, test.assert)
	}
}

// TestValidate covers the various validations within .Validate()
func TestValidate(t *testing.T) {
	// newDoc replicated NewIODocument, but without inline validation. this allows testing validation directly
	newDoc := func(format DocumentFormat, title string, content string) IODocument {
		return IODocument{Title: title, Format: format, Content: []byte(content)}
	}

	tests := []struct {
		name       string
		document   IODocument
		assertFunc func(t *testing.T, err error)
	}{
		{
			name:     "no error",
			document: newDoc(DocumentFormatMarkdown, "test document", "test content"),
			assertFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:     "invalid format",
			document: newDoc(DocumentFormat("invalid"), "test document", "test content"),
			assertFunc: func(t *testing.T, err error) {
				assert.Error(t, err)

				var ive globals.InputValidationError
				assert.ErrorAs(t, err, &ive)
				assert.Equal(t, "format", ive.Field)
				assert.Equal(t, "invalid format provided", ive.Reason)
			},
		},
		{
			name:     "invalid title",
			document: newDoc(DocumentFormatMarkdown, "", "test content"),
			assertFunc: func(t *testing.T, err error) {
				assert.Error(t, err)

				var ive globals.InputValidationError
				assert.ErrorAs(t, err, &ive)
				assert.Equal(t, "title", ive.Field)
				assert.Equal(t, "title is empty", ive.Reason)
			},
		},
		{
			name:     "invalid content",
			document: newDoc(DocumentFormatMarkdown, "test document", ""),
			assertFunc: func(t *testing.T, err error) {
				assert.Error(t, err)

				var ive globals.InputValidationError
				assert.ErrorAs(t, err, &ive)
				assert.Equal(t, "content", ive.Field)
				assert.Equal(t, "content is empty", ive.Reason)
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
