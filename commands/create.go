package commands

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/therealkevinard/adr-er/adr"
	io_document "github.com/therealkevinard/adr-er/output-templates"
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
		var filename string

		// run load-compile-write under a spinner
		_ = spinner.New().Title("saving the file").Action(func() {
			// load the template
			tpl, tplErr := io_document.DefaultTemplateForFormat(io_document.DocumentFormatMarkdown)
			if tplErr != nil {
				outputErr = fmt.Errorf("error finding template: %w", tplErr)
				return
			}

			// build the document
			document, buildErr := record.BuildDocument(tpl)
			if buildErr != nil {
				outputErr = fmt.Errorf("error rendering document: %w", buildErr)
				return
			}
			filename = document.Filename()

			// write the document
			if writeErr := document.Write(); writeErr != nil {
				outputErr = fmt.Errorf("error writing document: %w", writeErr)
				return
			}
		}).Run()

		if outputErr != nil {
			return fmt.Errorf("error writing adr document: %w", outputErr)
		}

		style := theme.TitleStyle()
		msg := fmt.Sprintf("wrote ADR to %s", filename)
		fmt.Printf(style.Render(lipgloss.JoinVertical(lipgloss.Left, msg)))
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
