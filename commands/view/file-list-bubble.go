package view

import (
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

var (
	titleStyle    = lipgloss.NewStyle().MarginLeft(2)
	helpStyle     = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)

	itemStyle         = lipgloss.NewStyle().PaddingLeft(2).PaddingBottom(1)
	selectedItemStyle = itemStyle.Foreground(lipgloss.Color("170"))

	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
)

// the model...
type fileList struct {
	list.Model
	selectedFile fileListItem
	quitting     bool
}

func newFileList(wd string) (fileList, error) {
	const (
		listWidth  = 20
		listHeight = 14
		panelTitle = "ADR Entries"
	)

	// load ADR files from the working directory
	filesListItems, err := getFilesList(wd)
	if err != nil {
		return fileList{}, fmt.Errorf("error listing files: %w", err)
	}

	// build the anon-embed list.Model
	l := list.New(filesListItems, fileItemDelegate{}, listWidth, listHeight)
	l.Title = panelTitle
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return fileList{
		Model:        l,
		selectedFile: fileListItem{},
		quitting:     false,
	}, nil
}

func (m fileList) Init() tea.Cmd { return nil }

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

// FilterValue returns the value to reference when the list is in filter mode
func (i fileListItem) FilterValue() string { return i.name }

// FullPath returns the absolute path to the file item.
// TODO: this is currently muy naive.
func (i fileListItem) FullPath() string { return path.Join(i.parent, i.name) }

// ... and their delegates ...
type fileItemDelegate struct{}

func (d fileItemDelegate) Height() int                             { return 1 }
func (d fileItemDelegate) Spacing() int                            { return 0 }
func (d fileItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d fileItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(fileListItem)
	if !ok {
		return
	}

	//
	str := fmt.Sprintf("%s\nmodified %s",
		i.name,
		humanize.RelTime(i.modified, time.Now(), "ago", "from now"),
	)

	// default item renderer
	renderFunc := itemStyle.Render

	// selected item renderer
	if index == m.Index() {
		renderFunc = func(s ...string) string {
			return selectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, renderFunc(str))
}
