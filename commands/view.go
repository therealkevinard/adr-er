package commands

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour/styles"
	glow "github.com/charmbracelet/glow/v2/ui"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var _ CliCommand = (*View)(nil)

// View wraps the cli command for viewing existing ADR documents.
type View struct {
	// directory holding architecture decision records
	adrDir string
}

func NewView(adrDir string) *View {
	return &View{adrDir: adrDir}
}

func (v *View) Action(_ *cli.Context) error {
	var (
		style    = styles.TokyoNightStyle
		width, _ = screenDims()
	)

	glowApp := glow.NewProgram(glow.Config{
		ShowAllFiles:         false,
		ShowLineNumbers:      true,
		Gopath:               "",
		HomeDir:              "",
		EnableMouse:          false,
		PreserveNewLines:     false,
		HighPerformancePager: false,

		WorkingDirectory: v.adrDir,
		GlamourMaxWidth:  uint(width),
		GlamourStyle:     style,
		GlamourEnabled:   true,
	})
	if _, err := glowApp.Run(); err != nil {
		return fmt.Errorf("well, THAT didn't work: %w", err)
	}

	return nil
}

// screenDims returns the terminal width and height.
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
