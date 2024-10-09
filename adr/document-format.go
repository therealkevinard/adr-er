package adr

// DocumentFormat is a string alias that constrains allowed formats and simplifies determining the file extension
type DocumentFormat string

// Valid tests the DocumentFormat is a supportedFormats
func (df DocumentFormat) Valid() bool {
	_, ok := supportedFormats[df]
	return ok
}

// GetExtension returns the file extension.
func (df DocumentFormat) GetExtension() string {
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
