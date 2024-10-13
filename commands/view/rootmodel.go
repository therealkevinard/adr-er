package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	listWidth  = 32
	listHeight = 0
)

type rootModel struct {
	FileList   fileList
	FileViewer fileViewer

	// presist the last/current opened file. this allows us to load content only when it's _actually_ changed.
	prevSelectedFilename string
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
		FileList:             fl,
		FileViewer:           fv,
		prevSelectedFilename: "",
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
	var listCmd, viewCmd tea.Cmd

	// update filelist first
	var flm tea.Model
	flm, listCmd = m.FileList.Update(msg)
	m.FileList = flm.(fileList) //nolint:errcheck // fileList.Update literally can't return anything other than a fileList

	// handle flm selection change.
	if sel := m.FileList.selectedFile.FullPath(); sel != "" && sel != m.prevSelectedFilename {
		// capture the selected file here to ref on later iterations
		m.prevSelectedFilename = sel

		// load the selected file content
		content, err := getFileContent(sel)
		if err != nil {
			content = err.Error()
		}

		// pass content to viewer for display
		m.FileViewer = m.FileViewer.Show(content)
	}

	// update fileviewer
	var fvm tea.Model
	fvm, viewCmd = m.FileViewer.Update(msg)
	//nolint:errcheck // fileViewer.Update literally can't return anything other than a fileViewer
	m.FileViewer = fvm.(fileViewer)

	return m, tea.Batch(listCmd, viewCmd)
}

// The view function, which renders the UI.
func (m rootModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.FileList.View(), m.FileViewer.View())
}
