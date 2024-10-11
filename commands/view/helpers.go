package view

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
)

// getFilesList reads a directory of files, returning []list.Item.
// the returned sliced is suitable for pupulating the fileList model
// TODO: this muddies concerns. should return a generic slice that fileList coerces to list.Item
// TODO: this should leverage the regex file filter used elsewhere to only show ADR files (as determined by naming convention)
func getFilesList(wd string) ([]list.Item, error) {
	items, err := os.ReadDir(wd)
	if err != nil {
		return nil, fmt.Errorf("error reading dir %s: %w", wd, err)
	}

	filesList := make([]list.Item, 0)

	for _, item := range items {
		// don't FileList dirs
		if item.IsDir() {
			continue
		}

		info, err := os.Stat(filepath.Join(wd, item.Name()))
		// don't FileList unreadable files
		if err != nil {
			continue
		}

		filesList = append(filesList, fileListItem{
			name:     info.Name(),
			parent:   wd,
			modified: info.ModTime(),
		})
	}

	return filesList, nil
}

// getFileContent reads <file>, returning its string content
func getFileContent(file string) (string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("error reading %s: %w", file, err)
	}

	return string(content), nil
}
