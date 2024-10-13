//nolint:gochecknoglobals // theme package is global by design
package theme

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// theme is a read-only singular instance of the Theme used across the application.
// consumers can use the ApplicationTheme getter to reference it, which will create on demand
var theme *Theme

// ThemeColor is a typed map-key const used for Theme.KeyColors
type ThemeColor string

// holds the valid ThemeColor constants
const (
	ThemeColorNormalFG ThemeColor = "color-normal-fg"
	ThemeColorIndigo   ThemeColor = "color-indigo"
	ThemeColorCream    ThemeColor = "color-cream"
	ThemeColorFuchsia  ThemeColor = "color-fuchsia"
	ThemeColorGreen    ThemeColor = "color-green"
	ThemeColorRed      ThemeColor = "color-red"
)

// Theme is an adr-er theme. it anon-embed a *huh.Theme, but also hoists key colors for ad-hoc use.
// this is effectively a huh.Theme that also plays well with bubbletea/bubbles
type Theme struct {
	*huh.Theme
	PrimaryColor lipgloss.TerminalColor
	AccentColor  lipgloss.TerminalColor

	KeyColors map[ThemeColor]lipgloss.TerminalColor
}

// ApplicationTheme returns the singular theme instance.
func ApplicationTheme() *Theme {
	if theme == nil {
		// init from huh.ThemeCharm.
		// color keys ripped from ThemeCharm() constructor and hoisted to _our_ theme for re-use
		var (
			normalFg = lipgloss.AdaptiveColor{Light: "235", Dark: "252"}
			indigo   = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
			cream    = lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}
			fuchsia  = lipgloss.Color("#F780E2")
			green    = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
			red      = lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"}
		)

		//
		theme = &Theme{
			Theme:        huh.ThemeCharm(),
			PrimaryColor: indigo,
			AccentColor:  fuchsia,
			KeyColors: map[ThemeColor]lipgloss.TerminalColor{
				ThemeColorNormalFG: normalFg,
				ThemeColorIndigo:   indigo,
				ThemeColorCream:    cream,
				ThemeColorFuchsia:  fuchsia,
				ThemeColorGreen:    green,
				ThemeColorRed:      red,
			},
		}
	}

	return theme
}

func (t *Theme) TitleStyle() lipgloss.Style { return t.Focused.Title }

func (t *Theme) HelpStyle() lipgloss.Style { return list.DefaultStyles().HelpStyle }

// RenderCancelMessage writes the very common "cancelled" message to the user.
func (t *Theme) RenderCancelMessage() {
	fmt.Println(t.TitleStyle().Render(
		lipgloss.JoinVertical(lipgloss.Left, "❌ Cancelled. No changes were made.")),
	)
}
