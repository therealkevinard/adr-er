package theme

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// theme is the huh.Theme used across the application
var theme *huh.Theme

// common string constants used across the application.
const (
	// actionCancelledText is used when user chose to cancel an operation
	actionCancelledText = "❌ Cancelled. No changes were made."
)

// Theme returns the app-wide theme.
// TODO: if we do custom themes one day, this is where it happens. till then... i like ThemeCharm
func Theme() *huh.Theme {
	if theme == nil {
		theme = huh.ThemeCharm()
	}

	return theme
}

// holds handy theme accessors for frequently-used lipgloss styles.
var (
	// TitleStyle is used for chonky title blocks
	TitleStyle = func() lipgloss.Style { return Theme().Focused.Title.MarginTop(1) }

	// SelectedStyle is the style used for selects' focused items
	SelectedStyle = func() lipgloss.Style { return Theme().Focused.SelectedOption }

	// BlockMarginStyle is simply a bit of left-margin. it's used to indent a block of text
	BlockMarginStyle = func() lipgloss.Style { return lipgloss.NewStyle().MarginLeft(1) }
)

// RenderTextBlock renders the provided strings as a line-delimited text block.
func RenderTextBlock(lines ...string) {
	block := lipgloss.JoinVertical(lipgloss.Left, lines...)
	fmt.Printf(BlockMarginStyle().Render(block))
}

// RenderCancelMessage writes the very common "cancelled" message to the user
func RenderCancelMessage() {
	fmt.Println(TitleStyle().Render(
		lipgloss.JoinVertical(lipgloss.Left, actionCancelledText)),
	)
}
