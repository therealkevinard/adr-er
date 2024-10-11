package output_templates

import "github.com/therealkevinard/adr-er/globals"

var _ globals.Validator = (*DocumentFormat)(nil)

// DocumentFormat defines a supported document format and provides methods for validation and retrieving file
// extensions.
type DocumentFormat string

// Validate checks if the DocumentFormat is supported.
// It returns an error if the format is not recognized.
func (df DocumentFormat) Validate() error {
	if _, ok := supportedFormats[df]; !ok {
		return globals.ValidationError("format", "unsupported format")
	}

	return nil
}

// Extension returns the file extension associated with the DocumentFormat.
// If the format is unsupported, it returns "txt" as a safe default.
func (df DocumentFormat) Extension() string {
	if ext, ok := supportedFormats[df]; ok {
		return ext
	}

	// safe default
	return "txt"
}

// Available DocumentFormat constants for supported document types.
const (
	// DocumentFormatMarkdown represents a markdown document format.
	DocumentFormatMarkdown DocumentFormat = "markdown"
)

// supportedFormats registers a map of supported formats to their fs extension.
var supportedFormats = map[DocumentFormat]string{
	DocumentFormatMarkdown: "md",
}
