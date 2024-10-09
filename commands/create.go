package commands

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/therealkevinard/adr-er/adr"
	"github.com/therealkevinard/adr-er/theme"
	"github.com/therealkevinard/adr-er/validators"
	"github.com/urfave/cli/v2"
)

var _ CliCommand = (*Create)(nil)

type Create struct{}

func (n Create) Action(cliCtx *cli.Context) error {
	var err error

	// form values
	confirmed := false
	record := new(adr.ADR)

	// run with error-or-cancel
	{
		form := huh.NewForm(
			huh.NewGroup(
				// title
				huh.NewInput().
					Value(&record.Title).
					Title("Title").
					Description("name your decision").
					CharLimit(128).
					Inline(false).
					Validate(validators.StrLenValidator("title", 3, 128)),
				// context
				huh.NewText().
					Value(&record.Context).
					Title("Context").
					Description("add relevant context"),
				// decision
				huh.NewText().
					Value(&record.Decision).
					Title("Decision").
					Description("what did you folks decide to do"),
				// status
				huh.NewSelect[string]().
					Value(&record.Status).
					Title("Status").
					OptionsFunc(func() []huh.Option[string] {
						return n.statusOptions()
					}, nil).
					Description("what's the current status?"),
				// consequences
				huh.NewText().
					Value(&record.Consequences).
					Title("Consequences").
					Description("what are the consequences of this decision?"),

				// confirmation
				huh.NewConfirm().
					Value(&confirmed).
					Title("feeling good about this one?"),
			).Title("The Decision"),
		).WithTheme(theme.Theme())

		if err = form.Run(); err != nil {
			return fmt.Errorf("error running form: %w", err)
		}
		if !confirmed {
			theme.RenderCancelMessage()
		}
	}

	// commit the input
	{
		var outputErr error

		// compile the doc. this renders the md template and returns a CompiledDocument instance that can be flushed to disk
		document, outputErr := record.CompileDocument(adr.DocumentFormatMarkdown)
		if outputErr != nil {
			return fmt.Errorf("error rendering document: %w", outputErr)
		}

		// write the file
		_ = spinner.New().Title("saving the file").Action(func() {
			outputErr = document.Write()
		}).Run()
		if outputErr != nil {
			return fmt.Errorf("error writing adr document: %w", outputErr)
		}

		msg := fmt.Sprintf("wrote ADR to %s", document.Filename())
		fmt.Printf(
			theme.TitleStyle().Render(lipgloss.JoinVertical(lipgloss.Left, msg)),
		)

	}

	return nil
}

// statusOptions returns valid options for status selection.
// TODO: this is overkill rn, but the plan is for this func to read from an active ADR record and do things.
func (n Create) statusOptions() []huh.Option[string] {
	return []huh.Option[string]{
		huh.NewOption("proposed", "proposed"),
		huh.NewOption("accepted", "accepted"),
		huh.NewOption("rejected", "rejected"),
		huh.NewOption("deprecated", "deprecated"),
		huh.NewOption("superceded", "superceded"),
	}
}
