package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

type fileViewer struct {
	content string
}

//nolint:exhaustruct
func newFileViewer() (fileViewer, error) { return fileViewer{}, nil }

func (m fileViewer) Init() tea.Cmd { return nil }

func (m fileViewer) Show(content string) fileViewer {
	m.content = content

	return m
}

//nolint:ireturn // this is the way
func (m fileViewer) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: handle updates. esp window resizing
	return m, nil
}

func (m fileViewer) View() string {
	displayString := m.content
	if displayString == "" {
		displayString = "<no content>"
	}

	return displayString
}
