package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"

	"github.com/therealkevinard/adr-er/adr"
	io_document "github.com/therealkevinard/adr-er/output-templates"
	"github.com/therealkevinard/adr-er/theme"
	"github.com/therealkevinard/adr-er/validators"
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
		outputStdOut:         false,
	}
	// set stdout flag
	if slices.Contains([]string{"", "-", "/"}, cmd.outputDir) {
		cmd.outputStdOut = true
	}

	return cmd
}

//nolint:funlen // tui apps are long by nature
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
			confirmText = fmt.Sprint("this will flush to stderr")
		} else {
			// a simple inline to dompute a relative path. if it errors, just show the abs path
			// TODO: this absolutely belongs in a util, and should also be used in the final message.
			displayPath := func() string {
				cwd, err := os.Getwd()
				if err != nil {
					return ""
				}
				rel, err := filepath.Rel(cwd, n.outputDir)
				if err != nil {
					return ""
				}

				return rel
			}()

			if displayPath == "" {
				displayPath = n.outputDir
			}

			confirmText = fmt.Sprintf("this will create next sequence number %d \nin %s", n.nextSequence, displayPath)
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
					OptionsFunc(n.statusOptions, nil).
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
		// captures errors inside the spinner closure
		var outputErr error
		// captures the document filename after it's built
		var filename string
		// captures a final message to the user.
		// if writing to stdout, this is the ADR string; for file output, it's a friendly status message
		var finalMsg string

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
				finalMsg = fmt.Sprint(string(document.Content))
			}
		}).Run()

		if outputErr != nil {
			return fmt.Errorf("error writing adr document: %w", outputErr)
		}

		// we need different styles if we're writing status vs flushing the whole document
		// TODO: writing to stdout is a little awkward, still. refine that some other time.
		if n.outputStdOut {
			fmt.Print(lipgloss.NewStyle().Render(finalMsg))
		} else {
			fmt.Print(theme.TitleStyle().Render(lipgloss.JoinVertical(lipgloss.Left, finalMsg)))
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
