package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/teacup/markdown"
	"github.com/therealkevinard/adr-er/theme"
)

var _ tea.Model = (*fileViewer)(nil)

// message types
type setFilenameMsg string

// setFilenameCmd emits a setFilenameMsg. this command is issued from other models that would like to update fileViewer
var setFilenameCmd = func(filename string) tea.Cmd {
	return func() tea.Msg {
		return setFilenameMsg(filename)
	}
}

type fileViewer struct {
	markdown markdown.Model
	// presist the last/current opened file. this allows us to load content only when it's _actually_ changed.
	prevSelectedFilename string
}

//nolint:exhaustruct
func newFileViewer() (fileViewer, error) {
	indigo, ok := theme.ApplicationTheme().KeyColors[theme.ThemeColorIndigo].(lipgloss.AdaptiveColor)
	if !ok {
		indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	}

	return fileViewer{
		markdown: markdown.New(false, true, indigo),
	}, nil
}

func (m fileViewer) Init() tea.Cmd { return nil }

//nolint:ireturn // this is the way
func (m fileViewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch message := msg.(type) {
	// handle window resize
	case tea.WindowSizeMsg:
		// layout constants. space to subtract from full window size to account for sibling elems, margins, borders, etc
		const (
			hMinus = listWidth + 4
			vMinus = 1
		)

		cmds = append(cmds, m.markdown.SetSize(message.Width-hMinus, message.Height-vMinus))

	// update viewing file
	case setFilenameMsg:
		// only evaluate if there's a meaningful change.
		if fname := string(message); fname != "" && fname != m.prevSelectedFilename {
			m.prevSelectedFilename = fname
			m.markdown.GotoTop()
			cmds = append(cmds, m.markdown.SetFileName(string(message)))
		}
	}

	// update child/embedded tea.Model
	m.markdown, cmd = m.markdown.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View returns the tui view for this model
func (m fileViewer) View() string {
	return m.markdown.View()
}

// SetIsActive toggles active/focusState state for this model
func (m fileViewer) SetIsActive(active bool) fileViewer {
	m.markdown.SetIsActive(active)
	m.markdown.SetBorderless(!active)

	return m
}
