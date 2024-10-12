package view

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/therealkevinard/adr-er/commands"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var _ commands.CliCommand = (*Command)(nil)

// Command wraps the cli command for viewing existing ADR documents.
type Command struct {
	// directory holding architecture decision records
	adrDir string
}

func NewCommand(adrDir string) *Command {
	return &Command{adrDir: adrDir}
}

// Action runs the TUI application for viewing Architectural Decision Records.
func (v *Command) Action(_ *cli.Context) error {
	options := []tea.ProgramOption{
		tea.WithAltScreen(),
	}

	// initialize the app models
	model, err := newRootModel(v.adrDir)
	if err != nil {
		return fmt.Errorf("error initializing tui: %w", err)
	}

	// run it
	if _, runErr := tea.NewProgram(model, options...).Run(); runErr != nil {
		return fmt.Errorf("error runing tui: %w", runErr)
	}

	return nil
}

// screenDims returns the terminal width and height.
//
//nolint:mnd // ui layout is all magic
func screenDims() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 120, 80
	}

	// limit width
	if width > 120 {
		width = 120
	}

	if width < 80 {
		width = 80
	}

	// limit height
	if height > 120 {
		height = 120
	}

	if height < 80 {
		height = 80
	}

	return width, height
}
