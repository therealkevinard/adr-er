package view

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = (*rootModel)(nil)

// layout constants
const (
	// width of the sidebar files list
	listWidth = 32
)

// rootKeyMap holds the keymap for rootmodel
type rootKeyMap struct {
	Quit key.Binding
	Next key.Binding
	Prev key.Binding
}

// rootModel is the outer tea.Model
type rootModel struct {
	FileList   fileList
	FileViewer fileViewer
	keymap     rootKeyMap

	// tracks focusState state to support cycling child models
	currentFocus focusState
}

func newRootModel(workDirectory string) (*rootModel, error) {
	//nolint:varnamelen // i approve these varnames
	var (
		err error
		fl  fileList
		fv  fileViewer
	)

	// init the fileList
	fl, err = newFileList(workDirectory)
	if err != nil {
		return nil, fmt.Errorf("error initializing filelist: %w", err)
	}

	// init the viewer
	fv, err = newFileViewer()
	if err != nil {
		return nil, fmt.Errorf("error initializing fileviewer: %w", err)
	}

	return &rootModel{
		FileList:   fl,
		FileViewer: fv,
		keymap: rootKeyMap{
			Quit: newKeyBinding(
				[]string{"q", "ctrl+c"}, "q/ctrl+c", "quit application",
			),
			Next: newKeyBinding(
				[]string{"right", "d", "tab"}, "→/d/tab", "next panel",
			),
			Prev: newKeyBinding(
				[]string{"left", "a", "shift+tab"}, "←/a/shift+tab", "prev panel",
			),
		},
		currentFocus: focusList,
	}, nil
}

func (m rootModel) Init() tea.Cmd {
	return tea.Batch(
		m.FileList.Init(),
		m.FileViewer.Init(),
	)
}

// The update function, which processes messages and updates the model.
//
//nolint:ireturn // this is the way
func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch message := msg.(type) {
	// check keys, rootmodel intercepts quit keys for tea.Quit
	case tea.KeyMsg:
		switch {
		// quit command
		case key.Matches(message, m.keymap.Quit):
			return m, tea.Quit

		// cycle next focusState
		case key.Matches(message, m.keymap.Next):
			m.currentFocus = m.currentFocus.Next(m.currentFocus)

		// cycle previous focusState
		case key.Matches(message, m.keymap.Prev):
			m.currentFocus = m.currentFocus.Prev(m.currentFocus)
		}
	}

	// handle focusState change
	switch m.currentFocus {
	case focusList:
		m.FileList = m.FileList.SetIsActive(true)
		m.FileViewer = m.FileViewer.SetIsActive(false)

	case focusViewer:
		m.FileViewer = m.FileViewer.SetIsActive(true)
		m.FileList = m.FileList.SetIsActive(false)
	}

	// update child/embedded tea.Models
	{
		// update filelist first, as the result may affect flows below here
		flm, listCmd := m.FileList.Update(msg)
		m.FileList = flm.(fileList) //nolint:errcheck // fileList.Update can only return fileList
		cmds = append(cmds, listCmd)

		// update fileviewer
		fvm, viewCmd := m.FileViewer.Update(msg)
		m.FileViewer = fvm.(fileViewer) //nolint:errcheck // fileViewer.Update can only return fileViewer
		cmds = append(cmds, viewCmd)
	}

	return m, tea.Batch(cmds...)
}

// The view function, which renders the UI.
func (m rootModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.FileList.View(), m.FileViewer.View())
}
