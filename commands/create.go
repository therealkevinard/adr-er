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
	"slices"
)

var _ CliCommand = (*Create)(nil)

type Create struct {
	// directory to write files into
	outputDir string
	// write to stdout, not file
	outputStdOut bool
	// user set output directory? false if we used LocateADRDirectory
	userDefinedOutputDir bool
	// the next integer sequence for the adrs in this directory
	nextSequence int
}

func NewCreate(outputDir string, userDefinedOutputDir bool, nextSequence int) *Create {
	cmd := &Create{
		outputDir:            outputDir,
		userDefinedOutputDir: userDefinedOutputDir,
		nextSequence:         nextSequence,
	}
	// set stdout flag
	if slices.Contains([]string{"", "-", "/"}, cmd.outputDir) {
		cmd.outputStdOut = true
	}

	return cmd
}

func (n Create) Action(ctx *cli.Context) error {
	var err error

	// form values
	confirmed := false
	record := &adr.ADR{
		Sequence:     n.nextSequence,
		Title:        "",
		Context:      "",
		Decision:     "",
		Status:       "",
		Consequences: "",
	}

	// run with error-or-cancel
	{
		var confirmText string
		if n.outputStdOut {
			confirmText = fmt.Sprintf("this will flush to stderr")
		} else {
			confirmText = fmt.Sprintf("this will create next sequence number %d \nin directory %s", n.nextSequence, n.outputDir)
		}

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
					Title("feeling good about this one?").
					Description(confirmText),
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
		var outputErr error // captures errors inside the spinner closure
		var filename string // captures the document filename after it's built
		var finalMsg string // captures a final message to the user. if writing to stdout, this is the ADR; for file output, it's a friendly status message

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
			if !n.outputStdOut {
				if writeErr := document.Write(n.outputDir); writeErr != nil {
					outputErr = fmt.Errorf("error writing document: %w", writeErr)
					return
				}
				finalMsg = fmt.Sprintf("wrote ADR to %s", filename)
			} else {
				finalMsg = fmt.Sprintf(string(document.Content))
			}
		}).Run()

		if outputErr != nil {
			return fmt.Errorf("error writing adr document: %w", outputErr)
		}

		// we need different styles if we're writing status vs flushing the whole document
		// TODO: writing to stdout is a little awkward, still. refine that some other time.
		if n.outputStdOut {
			fmt.Printf(lipgloss.NewStyle().Render(finalMsg))
		} else {
			fmt.Printf(theme.TitleStyle().Render(lipgloss.JoinVertical(lipgloss.Left, finalMsg)))
		}
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
