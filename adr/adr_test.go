package adr

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therealkevinard/adr-er/globals"
	io_document "github.com/therealkevinard/adr-er/io-document"
	"testing"
)

func TestBuildDocument(t *testing.T) {
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
				doc, err := adr.BuildDocument(io_document.DocumentFormatMarkdown)
				require.NoError(t, err)
				assert.Equal(t, "0001: <title>", doc.Title)
				assert.Equal(t, io_document.DocumentFormatMarkdown, doc.Format)
				assert.NotEmpty(t, doc.Content)
				assert.Equal(t, "0001-title", doc.DocumentID())
				assert.Equal(t, "0001-title.md", doc.Filename())
			},
		},
		{
			name: "invalid format",
			adr: &ADR{
				Sequence: 1,
				Title:    "<title>",
			},
			assertFunc: func(t *testing.T, adr *ADR) {
				doc, err := adr.BuildDocument(io_document.DocumentFormat("invalid"))
				assert.Nil(t, doc)
				assert.Error(t, err)

				// assert against the returned descriptive error
				var validationError globals.InputValidationError
				ok := errors.As(err, &validationError)
				assert.True(t, ok)
				assert.Equal(t, "format", validationError.Field)
				assert.Equal(t, "invalid format provided", validationError.Reason)

			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.assertFunc(t, test.adr)
		})
	}
}
