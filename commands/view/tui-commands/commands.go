package tui_commands

import tea "github.com/charmbracelet/bubbletea"

// message types.
type SetFilenameMsg string

// SetFilenameCmd emits a setFilenameMsg.
// FileViewerModel responds to this command by displaying the content of filename.
func SetFilenameCmd(filename string) tea.Cmd {
	return func() tea.Msg {
		return SetFilenameMsg(filename)
	}
}
