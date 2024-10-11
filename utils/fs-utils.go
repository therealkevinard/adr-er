package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/therealkevinard/adr-er/globals"
)

// this regex will match existing ADR output files.
// it follows the naming conventions outlined for output_templates.ParsedTemplateFile.
// example matches: 0001-fizzy-pop.md, 0002-bubble-gupp.md, 0003-thing-two.txt.
var adrFileNamePattern = regexp.MustCompile(`^(\d+)-.+\.\w+$`)

// LocateADRDirectory attempts to locate the correct directory to store ADRs, starting at root
// root defaults to os.Getwd if empty.
func LocateADRDirectory(root string) (string, error) {
	// default to os.Getwd()
	if root == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("error getting current working directory: %w", err)
		}

		root = cwd
	}

	// candidate directory names. these directory names, within root, are considered candidates
	// these are ordered by preference: the first one that passes all rules is used.
	candidates := []string{
		"architectural-decision-records",
		"adr",
		".adr",
	}

	// evaluate each candidate directory. the first one that passes the rules is used
	for _, dir := range candidates {
		candidatePath := filepath.Join(root, dir)
		ok, err := evaluateCandidate(candidatePath)
		if err != nil {
			continue
		}
		// this one wins
		if ok {
			return candidatePath, nil
		}

		continue
	}

	return "", fmt.Errorf("no ADR directory found in %s", root)
}

// GetHighestSequenceNumber reads the filenames in root, extracting the ADR sequence number.
// The highest existing value is returned.
func GetHighestSequenceNumber(root string) (int, error) {
	highest := 0
	if root == "" {
		return 0, globals.ValidationError("directory", "directory path is empty")
	}

	// read the filename entries
	entries, err := os.ReadDir(root)
	if err != nil {
		return 0, fmt.Errorf("error reading directory: %w", err)
	}

	// iterate entries, capturing the sequence prefix.
	for _, entry := range entries {
		// ignore dirs
		if entry.IsDir() {
			continue
		}

		// search filename using the ADR convention match pattern
		matches := adrFileNamePattern.FindStringSubmatch(entry.Name())
		if len(matches) < 2 {
			continue // No match or no captured group.
		}

		// convert to int
		sequence, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, fmt.Errorf("error parsing sequence number from %s: %w", entry.Name(), err)
		}

		// update highest, if higher.
		if sequence > highest {
			highest = sequence
		}
	}

	return highest, nil
}

// DisplayShortpath creates a relative path from absolute.
// this is used primarily for display, as absolute paths can _easily_ over-wrap.
// for error cases, the absolute path is returned. this guarantees a usable return value.
func DisplayShortpath(absolutePath string) (string, error) {
	var (
		cwd          string
		relativePath string

		err error
	)

	if cwd, err = os.Getwd(); err != nil {
		return absolutePath, fmt.Errorf("error getting current working directory: %w", err)
	}
	if relativePath, err = filepath.Rel(cwd, absolutePath); err != nil {
		return absolutePath, fmt.Errorf("error getting relative path: %w", err)
	}

	return "./" + relativePath, nil
}

// evaluateCandidate checks an os directory as a viable store for ADR files.
// returns true if the directory is a valid candidate, otherwise false
// a viable store must be either empty, or hold only ADR-named files. subdirectories are allowed.
func evaluateCandidate(fullpath string) (bool, error) {
	// read the contents
	entries, err := os.ReadDir(fullpath)
	if err != nil {
		return false, fmt.Errorf("error reading directory: %w", err)
	}

	// vacant? winner!
	if len(entries) == 0 {
		return true, nil
	}

	// iterate the contents of this directory
	for _, entry := range entries {
		// skip directories as we're only checking files.
		// we explicitly want to allow subdirectories, as these _could_ be used for evidence or other docs.
		if entry.IsDir() {
			continue
		}

		// if any file fails the regex match, the directory isn't a candidate.
		if !adrFileNamePattern.MatchString(entry.Name()) {
			return false, nil
		}
	}

	// If all files match the regex, return true.
	return true, nil
}
