//nolint:gochecknoglobals // theme package is global by design
package theme

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// common string constants used across the application.
const (
	// actionCancelledText is used when user chose to cancel an operation.
	actionCancelledText = "❌ Cancelled. No changes were made."
)

// theme is a singular instance of the Theme used across the application.
// consumers can use the ApplicationTheme getter to reference it.
var theme *Theme

type Theme struct {
	*huh.Theme
}

// ApplicationTheme returns the singular theme instance.
func ApplicationTheme() *Theme {
	if theme == nil {
		theme = &Theme{
			Theme: huh.ThemeCharm(),
		}

		// resets/overrides
		theme.Focused.Title = theme.Focused.Title.MarginTop(1)
	}

	return theme
}

func (t *Theme) TitleStyle() lipgloss.Style       { return t.Focused.Title }
func (t *Theme) SelectedStyle() lipgloss.Style    { return t.Focused.SelectedOption }
func (t *Theme) BlockMarginStyle() lipgloss.Style { return lipgloss.NewStyle().Margin(1) }

// RenderTextBlock renders the provided strings as a line-delimited text block.
func (t *Theme) RenderTextBlock(lines ...string) {
	fmt.Print(t.BlockMarginStyle().Render(
		lipgloss.JoinVertical(lipgloss.Left, lines...),
	))
}

// RenderCancelMessage writes the very common "cancelled" message to the user.
func (t *Theme) RenderCancelMessage() {
	fmt.Println(t.TitleStyle().Render(
		lipgloss.JoinVertical(lipgloss.Left, actionCancelledText)),
	)
}
