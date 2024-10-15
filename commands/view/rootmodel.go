package view

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/therealkevinard/adr-er/commands"
	file_list "github.com/therealkevinard/adr-er/commands/view/file-list"
	file_viewer "github.com/therealkevinard/adr-er/commands/view/file-viewer"
)

var _ tea.Model = (*rootModel)(nil)

// rootModel is the outer tea.Model.
type rootModel struct {
	// FileList is the child model that renders the left file list
	FileList file_list.FileListModel
	// FileViewer is the child model the renders the selected file content
	FileViewer file_viewer.FileViewerModel

	// tracks focusState state to support cycling child models
	currentFocus focusState

	// keymap holds the keybindings this model responds to. this feeds the help model to render help text.
	keymap rootKeyMap
	help   help.Model

	// track screen dimensions for layout reasons
	screenW int
	screenH int
}

func newRootModel(workDirectory string) (*rootModel, error) {
	//nolint:varnamelen // i approve these varnames
	var (
		err error
		fl  file_list.FileListModel
		fv  file_viewer.FileViewerModel
		hv  help.Model
	)

	// init the fileList
	fl, err = file_list.New(workDirectory)
	if err != nil {
		return nil, fmt.Errorf("error initializing filelist: %w", err)
	}

	// init the viewer
	fv = file_viewer.New()

	// init help
	hv = help.New()

	return &rootModel{
		FileList:     fl,
		FileViewer:   fv,
		help:         hv,
		currentFocus: focusList,
		screenW:      0,
		screenH:      0,
		keymap: rootKeyMap{
			Quit: commands.NewKeybinding(
				[]string{"q", "ctrl+c"}, "q/ctrl+c", "quit application",
			),
			Next: commands.NewKeybinding(
				[]string{"right", "tab"}, "→/tab", "next panel",
			),
			Prev: commands.NewKeybinding(
				[]string{"left", "shift+tab"}, "←/shift+tab", "prev panel",
			),
		},
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
//nolint:ireturn
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

	case tea.WindowSizeMsg:
		m = m.SetScreenDimensions(message.Width, message.Height)
	}

	// update m.currentFocus
	//nolint:exhaustive // iota case focusMax is computation-only
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
		cmds = append(cmds, listCmd)
		m.FileList = flm.(file_list.FileListModel) //nolint:errcheck // fileList.Update can only return fileList

		// update fileviewer
		fvm, viewCmd := m.FileViewer.Update(msg)
		cmds = append(cmds, viewCmd)
		m.FileViewer = fvm.(file_viewer.FileViewerModel) //nolint:errcheck // fileViewer.Update can only return fileViewer
	}

	return m, tea.Batch(cmds...)
}

// The view function, which renders the UI.
func (m rootModel) View() string {
	mainView := lipgloss.JoinHorizontal(lipgloss.Bottom, m.FileList.View(), m.FileViewer.View())
	helpView := lipgloss.NewStyle().Padding(0, 1).Render(m.help.View(m.keymap))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainView,
		helpView,
	)
}

// SetScreenDimensions updates the outer screen dimensions.
func (m rootModel) SetScreenDimensions(width, height int) rootModel {
	m.screenW = width
	m.screenH = height

	return m
}

// rootKeyMap holds the keymap for rootmodel
// implements help.KeyMap for help panel support.
type rootKeyMap struct {
	Quit key.Binding
	Next key.Binding
	Prev key.Binding
}

func (r rootKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{r.Next, r.Prev, r.Quit}
}

func (r rootKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{r.Next, r.Prev},
		{r.Quit},
	}
}
