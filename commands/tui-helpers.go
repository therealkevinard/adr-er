package commands

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/therealkevinard/adr-er/globals"
	"golang.org/x/term"
)

// StrLenValidator returns a func that ensures a string's len is within range.
func StrLenValidator(fieldLabel string, min, max int) func(string) error {
	return func(s string) error {
		if len(s) < min {
			return globals.ValidationError(fieldLabel, fmt.Sprintf("must have at least %d characters", min))
		}

		if len(s) > max {
			return globals.ValidationError(fieldLabel, fmt.Sprintf("must have at most %d characters", min))
		}

		return nil
	}
}

// ScreenDimensions returns the terminal width and height, constrained to `80 <= x <= 120`
//
//nolint:mnd // ui layout is all magic
func ScreenDimensions() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 120, 80
	}

	return Mid(width, 80, 120),
		Mid(height, 80, 120)
}

// Mid is a tiny-tiny helper to return int value constrained by `min <= value <= max`.
func Mid(value int, min, max int) int {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}

// NewKeybinding is a simple constructor that creates a key.Binding from provided strokes and help input.
func NewKeybinding(strokes []string, helpKeys, helpDesc string) key.Binding {
	return key.NewBinding(
		key.WithKeys(strokes...),
		key.WithHelp(helpKeys, helpDesc),
	)
}
