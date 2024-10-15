package file_list

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/therealkevinard/adr-er/commands"
	tui_commands "github.com/therealkevinard/adr-er/commands/view/tui-commands"
	"github.com/therealkevinard/adr-er/globals"
	"github.com/therealkevinard/adr-er/theme"
)

var _ tea.Model = (*FileListModel)(nil)

// FileListModel is the file list view model.
// It renders the list of files found in its bound workDirectory.
type FileListModel struct {
	list.Model
	active bool
	keymap fileListKeyMap
}

// New creates a new FileListModel, bound ot the provided workDirectory.
func New(workDirectory string) (FileListModel, error) {
	// load ADR files from the working directory
	filesListItems, err := getFilesList(workDirectory)
	if err != nil {
		return FileListModel{}, fmt.Errorf("error listing files: %w", err)
	}

	listModel := list.New(filesListItems, list.NewDefaultDelegate(), 0, 0)
	listModel.Title = "ADR Entries"
	listModel.SetShowStatusBar(true)
	listModel.SetFilteringEnabled(true)
	listModel.Styles.Title = theme.ApplicationTheme().TitleStyle()
	listModel.Styles.HelpStyle = theme.ApplicationTheme().HelpStyle()

	return FileListModel{
		Model:  listModel,
		active: false,
		keymap: fileListKeyMap{
			Up: commands.NewKeybinding(
				[]string{"up", "w"}, "↑/w", "up",
			),
			Down: commands.NewKeybinding(
				[]string{"down", "s"}, "↓/s", "down",
			),
			Enter: commands.NewKeybinding(
				[]string{"enter"}, "↵/<enter>", "select",
			),
		},
	}, nil
}

// Init ...
func (m FileListModel) Init() tea.Cmd { return nil }

// Update ...
//
//nolint:ireturn
func (m FileListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// evaluate these only if this model has focusState
	if m.active {
		// update list.model before evaluating keys. we need its result for auto-load
		m.Model, cmd = m.Model.Update(msg)
		cmds = append(cmds, cmd)

		// handle keys
		//nolint:gocritic // keeping singleCaseSwitch for convention
		switch message := msg.(type) {
		case tea.KeyMsg:
			switch {
			// autoload files when selection changes.
			case key.Matches(message, m.keymap.Up), key.Matches(message, m.keymap.Down):
				// fallthrough to load-on-enter case. using fallthrough here will simplify user toggles for this behavior
				fallthrough

			// more conservative load-on-enter behavior
			case key.Matches(message, m.keymap.Enter):
				i, _ := m.SelectedItem().(Item)
				cmds = append(cmds, tui_commands.SetFilenameCmd(i.FullPath()))
			}
		}
	}

	// evaluate these regardless of focusState
	//nolint:gocritic // keeping singleCaseSwitch for convention
	switch message := msg.(type) {
	// handle screen size
	case tea.WindowSizeMsg:
		// layout constants. space to subtract from full window size to account for sibling elems, margins, borders, etc
		const (
			hMinus = 0
			vMinus = 3
		)

		m.Model.SetWidth(globals.ListModelWidth - hMinus)
		m.Model.SetHeight(message.Height - vMinus)

		cmds = append(cmds, nil)
	}

	return m, tea.Batch(cmds...)
}

// View ...
func (m FileListModel) View() string {
	focusedBorderColor := theme.ApplicationTheme().KeyColors[theme.ThemeColorIndigo]
	style := lipgloss.NewStyle().BorderForeground(focusedBorderColor)

	// toggle border visible based on active/focus state
	if m.active {
		style = style.BorderStyle(lipgloss.NormalBorder())
	} else {
		style = style.BorderStyle(lipgloss.HiddenBorder())
	}

	return style.Render(m.Model.View())
}

// SetIsActive toggles active/focusState state for this model.
func (m FileListModel) SetIsActive(active bool) FileListModel {
	m.active = active

	return m
}

// getFilesList reads a directory of files, returning []list.Item.
// the returned sliced is suitable for pupulating the fileList model
// TODO: this should leverage the regex file filter used elsewhere to only show ADR files (per naming convention)
func getFilesList(workDirectory string) ([]list.Item, error) {
	items, err := os.ReadDir(workDirectory)
	if err != nil {
		return nil, fmt.Errorf("error reading dir %s: %w", workDirectory, err)
	}

	filesList := make([]list.Item, 0)

	for _, item := range items {
		// don't list dirs
		if item.IsDir() {
			continue
		}

		info, statErr := os.Stat(filepath.Join(workDirectory, item.Name()))
		// don't FileList unreadable files
		if statErr != nil {
			continue
		}

		filesList = append(filesList, NewItem(
			info.Name(),
			workDirectory,
			info.ModTime(),
		))
	}

	return filesList, nil
}

// fileListKeyMap holds the keys this model responds to.
type fileListKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
}
