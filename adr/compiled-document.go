package adr

import (
	"fmt"
	"os"
)

// CompiledDocument wraps several fields necessary for outputing a complete document based on an ADR.
// Namely, the filename, document format, and (of course) the content
type CompiledDocument struct {
	DocumentID string
	Format     DocumentFormat
	Content    []byte

	filename  string
	extension string
}

// Filename is a getter for the derived filename
func (cd CompiledDocument) Filename() string {
	if cd.DocumentID == "" {
		return "invalid.txt"
	}
	if !cd.Format.Valid() {
		return "invalid.txt"
	}

	// compile documentID and Format extension for full filename
	return cd.DocumentID + "." + cd.Format.GetExtension()
}

// Write flushes the document to filesystem
func (cd *CompiledDocument) Write() error {
	var err error

	// validate the document
	if err = cd.validate(); err != nil {
		return fmt.Errorf("not writing. document validation failed: %w", err)
	}

	// get its filename. validate that, too
	fname := cd.Filename()
	if fname == "" || fname == "invalid.txt" {
		return fmt.Errorf("not writing. invalid document filename %s", fname)
	}

	// make the file
	// TODO: make sure this errors if the file already exists. don't want to force-replace existing files.
	f, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", fname, err)
	}
	defer f.Close()

	// write it
	if _, err = f.Write(cd.Content); err != nil {
		return fmt.Errorf("could not write to file %s: %w", fname, err)
	}

	// donesies
	return nil
}

// validate checks several properties of CompiledDocument, returning errors on failure
func (cd *CompiledDocument) validate() error {
	// validate the format/extension
	if validFmt := cd.Format.Valid(); !validFmt {
		return fmt.Errorf("invalid format: %s", cd.Format)
	}
	// has content?
	if len(cd.Content) == 0 {
		return fmt.Errorf("no content")
	}
	// has document id?
	if len(cd.DocumentID) == 0 {
		return fmt.Errorf("invalid document id")
	}
	// valid filename
	if name := cd.Filename(); name == "" || name == "invalid.txt" {
		return fmt.Errorf("invalid document filename")
	}

	return nil
}
