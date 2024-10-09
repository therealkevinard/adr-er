package adr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestDocumentTitle provides coverage of ADR.documentID in many scenarios
func TestDocumentTitle(t *testing.T) {
	tests := []struct {
		name              string
		wantDocumentTitle string
		input             ADR
	}{
		{
			name:              "basic",
			wantDocumentTitle: "001_fizzy-pop",
			input: ADR{
				Sequence: 1,
				Title:    "Fizzy Pop",
			},
		},
		{
			name:              "many spaces",
			wantDocumentTitle: "001_fizzy-pop",
			input: ADR{
				Sequence: 1,
				Title:    "Fizzy                 Pop",
			},
		},
		{
			name:              "leading-trailing whsp",
			wantDocumentTitle: "001_fizzy-pop",
			input: ADR{
				Sequence: 1,
				Title:    "   Fizzy Pop   ",
			},
		},
		{
			name:              "unsafe characters",
			wantDocumentTitle: "001_fizzy-pop",
			input: ADR{
				Sequence: 1,
				Title:    "Fiz*zy P\\o/p",
			},
		},
		{
			name:              "one-indexed",
			wantDocumentTitle: "001_fizzy-pop",
			input: ADR{
				Sequence: 0,
				Title:    "Fizzy Pop",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.documentID()
			assert.Equal(t, test.wantDocumentTitle, got)
		})
	}
}
