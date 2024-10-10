package output_templates

import "github.com/therealkevinard/adr-er/globals"

var _ globals.Validator = (*DocumentFormat)(nil)

// DocumentFormat is a string alias that constrains allowed formats and simplifies determining the file extension.
type DocumentFormat string

// Validate validates the DocumentFormat, returning errors
func (df DocumentFormat) Validate() error {
	if _, ok := supportedFormats[df]; !ok {
		return globals.ValidationError("format", "unsupported format")
	}

	return nil
}

// Extension returns the file extension for this format.
func (df DocumentFormat) Extension() string {
	if ext, ok := supportedFormats[df]; ok {
		return ext
	}

	// safe default
	return "txt"
}

// this block holds DocumentFormat constants
const (
	// DocumentFormatMarkdown is a markdown document
	DocumentFormatMarkdown DocumentFormat = "markdown"
)

// supportedFormats registers a map of supported formats to their fs extension
var supportedFormats = map[DocumentFormat]string{
	DocumentFormatMarkdown: "md",
}
