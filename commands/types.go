package commands

import (
	"github.com/urfave/cli/v2"
)

// CliCommand exports an Action func that can be run in the tui
type CliCommand interface {
	// Action satisfies
	Action(cliCtx *cli.Context) error
}

// Validator is used for form validation
type Validator func(string) error
