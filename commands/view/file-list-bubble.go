package view

import (
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
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

	// build the anon-embed list.Model
	listModel := list.New(filesListItems, fileItemDelegate{}, listWidth, listHeight)
	listModel.Title = "ADR Entries"
	listModel.SetShowStatusBar(true)
	listModel.SetFilteringEnabled(true)
	listModel.Styles.Title = titleStyle
	listModel.Styles.PaginationStyle = listPaginationStyle
	listModel.Styles.HelpStyle = helpStyle

	//nolint:exhaustruct
	return fileList{Model: listModel}, nil
}

func (m fileList) Init() tea.Cmd { return nil }

//nolint:ireturn // this is the way
func (m fileList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// handle screen size
	case tea.WindowSizeMsg:
		// TODO: actually, don't. this should be a fixed-width sidebar when we're done-done.
		m.SetWidth(msg.Width)

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

func (m fileList) View() string { return m.Model.View() }

// ... and the FileList items ...
type fileListItem struct {
	name     string
	parent   string
	modified time.Time
}

// FilterValue returns the value to reference when the list is in filter mode.
func (i fileListItem) FilterValue() string { return i.name }

// FullPath returns the absolute path to the file item.
// TODO: this is currently muy naive.
func (i fileListItem) FullPath() string { return path.Join(i.parent, i.name) }

// ... and their delegates ...
type fileItemDelegate struct{}

func (d fileItemDelegate) Height() int                             { return 1 }
func (d fileItemDelegate) Spacing() int                            { return 0 }
func (d fileItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d fileItemDelegate) Render(w io.Writer, listModel list.Model, index int, listItem list.Item) {
	item, ok := listItem.(fileListItem)
	if !ok {
		return
	}

	//
	str := fmt.Sprintf("%s\nmodified %s",
		item.name,
		humanize.RelTime(item.modified, time.Now(), "ago", "from now"),
	)

	// default item renderer
	renderFunc := listItemStyle.Render

	// selected item renderer
	if index == listModel.Index() {
		renderFunc = func(s ...string) string {
			return selectedListItemStyle.Render(strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, renderFunc(str))
}
