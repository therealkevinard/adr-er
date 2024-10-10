package io_document

// DocumentFormat is a string alias that constrains allowed formats and simplifies determining the file extension.
type DocumentFormat string

// Valid tests the DocumentFormat is a supportedFormats
func (df DocumentFormat) Valid() bool {
	_, ok := supportedFormats[df]
	return ok
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
