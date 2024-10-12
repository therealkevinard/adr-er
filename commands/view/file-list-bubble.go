package view

import (
	"fmt"
	"path"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/therealkevinard/adr-er/theme"
)

// the model...
type fileList struct {
	list.Model
	selectedFile fileListItem
	quitting     bool
}

func newFileList(workDirectory string) (fileList, error) {
	// load ADR files from the working directory
	filesListItems, err := getFilesList(workDirectory)
	if err != nil {
		return fileList{}, fmt.Errorf("error listing files: %w", err)
	}

	listModel := list.New(filesListItems, list.NewDefaultDelegate(), listWidth, listHeight)
	listModel.Title = "ADR Entries"
	listModel.SetShowStatusBar(true)
	listModel.SetFilteringEnabled(true)
	listModel.Styles.Title = titleStyle
	listModel.Styles.HelpStyle = helpStyle

	listModel.Styles.PaginationStyle = theme.ApplicationTheme().ListPaginationStyle
	listModel.Paginator.PerPage = 10

	//nolint:exhaustruct
	return fileList{Model: listModel}, nil
}

func (m fileList) Init() tea.Cmd { return nil }

//nolint:ireturn // this is the way
func (m fileList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// handle screen size
	case tea.WindowSizeMsg:
		m.SetWidth(listWidth)
		m.SetHeight(msg.Height - 6)

		return m, nil

	// handle keys
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true

			return m, tea.Quit

		case "enter":
			i, ok := m.SelectedItem().(fileListItem)
			if ok {
				m.selectedFile = i
			}

			return m, nil
		}
	}

	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg)

	return m, cmd
}

func (m fileList) View() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(theme.ApplicationTheme().PrimaryColor).
		Padding(1)
	return style.Render(m.Model.View())
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
