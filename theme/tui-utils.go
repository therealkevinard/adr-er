package theme

import (
	"os"

	"golang.org/x/term"
)

// ScreenDimensions returns the terminal width and height.
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

// Mid is a tiny-tiny helper to return int value constrained by min and max.
func Mid(value int, min, max int) int {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}
