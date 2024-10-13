package view

import (
	"fmt"
	"path"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/therealkevinard/adr-er/theme"
)

var _ tea.Model = (*fileList)(nil)

type fileListKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
}

// fileList is the file list view model
type fileList struct {
	list.Model
	selectedFile fileListItem
	keymap       fileListKeyMap
	active       bool
}

func newFileList(workDirectory string) (fileList, error) {
	// load ADR files from the working directory
	filesListItems, err := getFilesList(workDirectory)
	if err != nil {
		return fileList{}, fmt.Errorf("error listing files: %w", err)
	}

	listModel := list.New(filesListItems, list.NewDefaultDelegate(), 0, 0)
	listModel.Title = "ADR Entries"
	listModel.SetShowStatusBar(true)
	listModel.SetFilteringEnabled(true)
	listModel.Styles.Title = theme.ApplicationTheme().TitleStyle()
	listModel.Styles.HelpStyle = theme.ApplicationTheme().HelpStyle()

	//nolint:exhaustruct
	return fileList{
		Model: listModel,
		keymap: fileListKeyMap{
			Up: newKeyBinding(
				[]string{"up", "w"}, "↑/w", "up",
			),
			Down: newKeyBinding(
				[]string{"down", "s"}, "↓/s", "down",
			),
			Enter: newKeyBinding(
				[]string{"enter"}, "↵/<enter>", "select",
			),
		},
	}, nil
}

func (m fileList) Init() tea.Cmd { return nil }

//nolint:ireturn // this is the way
func (m fileList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// evaluate these only if this model has focusState
	if m.active {
		// update list.model before evaluating keys. we need its result for auto-load
		m.Model, cmd = m.Model.Update(msg)
		cmds = append(cmds, cmd)

		switch message := msg.(type) {
		// handle keys
		case tea.KeyMsg:
			switch {
			// autoload files when selection changes.
			case key.Matches(message, m.keymap.Up), key.Matches(message, m.keymap.Down):
				// fallthrough to load-on-enter case. using fallthrough here will simplify user toggles for this behavior
				fallthrough

			// more conservative load-on-enter behavior
			case key.Matches(message, m.keymap.Enter):
				i, _ := m.SelectedItem().(fileListItem)
				cmds = append(cmds, setFilenameCmd(i.FullPath()))
			}
		}
	}

	// evaluate these regardless of focusState
	switch message := msg.(type) {
	// handle screen size
	case tea.WindowSizeMsg:
		m.SetWidth(listWidth)
		m.SetHeight(message.Height - 4)

		cmds = append(cmds, nil)
	}

	return m, tea.Batch(cmds...)
}

func (m fileList) View() string {
	style := theme.ApplicationTheme().Focused.Base.Border(lipgloss.NormalBorder(), true).Padding(1)

	return style.Render(m.Model.View())
}

// SetIsActive toggles active/focusState state for this model
func (m fileList) SetIsActive(active bool) fileList {
	m.active = active
	return m
}

// ... and the FileList items ...
type fileListItem struct {
	name     string
	parent   string
	modified time.Time
}

// Title is used by list.DefaultDelegate.
func (i fileListItem) Title() string { return i.name }

// Description is used by list.DefaultDelegate.
func (i fileListItem) Description() string {
	return humanize.RelTime(i.modified, time.Now(), "ago", "from now")
}

// FilterValue returns the value to reference when the list is in filter mode.
func (i fileListItem) FilterValue() string {
	return i.name
}

// FullPath returns the absolute path to the file item.
func (i fileListItem) FullPath() string {
	// TODO: this is currently muy naive.
	return path.Join(i.parent, i.name)
}
