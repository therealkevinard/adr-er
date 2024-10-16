package file_viewer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/teacup/markdown"
	tui_commands "github.com/therealkevinard/adr-er/commands/view/tui-commands"
	"github.com/therealkevinard/adr-er/globals"
	"github.com/therealkevinard/adr-er/theme"
)

var _ tea.Model = (*FileViewerModel)(nil)

// FileViewerModel is the file viewer model.
// it renders a selected file's contents with pretty formatting.
type FileViewerModel struct {
	markdown markdown.Model
	// presist the last/current opened file. this allows us to load content only when it's _actually_ changed.
	prevSelectedFilename string
}

func New() FileViewerModel {
	indigo, ok := theme.ApplicationTheme().KeyColors[theme.ThemeColorIndigo].(lipgloss.AdaptiveColor)
	if !ok {
		indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	}

	return FileViewerModel{
		markdown:             markdown.New(false, true, indigo),
		prevSelectedFilename: "",
	}
}

// Init ...
func (m FileViewerModel) Init() tea.Cmd { return nil }

// Update ...
//
//nolint:ireturn
func (m FileViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch message := msg.(type) {
	// handle window resize
	case tea.WindowSizeMsg:
		// layout constants. space to subtract from full window size to account for sibling elems, margins, borders, etc
		const (
			hMinus = globals.ListModelWidth + 4
			vMinus = 1
		)

		cmds = append(cmds, m.markdown.SetSize(message.Width-hMinus, message.Height-vMinus))

	// update viewing file
	case tui_commands.SetFilenameMsg:
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

// View ...
func (m FileViewerModel) View() string {
	return m.markdown.View()
}

// SetIsActive toggles active/focusState state for this model.
func (m FileViewerModel) SetIsActive(active bool) FileViewerModel {
	m.markdown.SetIsActive(active)
	m.markdown.SetBorderless(!active)

	return m
}
