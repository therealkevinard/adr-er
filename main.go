package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/therealkevinard/adr-er/commands/create"
	"github.com/therealkevinard/adr-er/commands/view"
	"github.com/therealkevinard/adr-er/utils"
	"github.com/urfave/cli/v2"
)

func main() {
	var (
		// root dir to write files into
		adrDirectory string
		// next int sequence. detemined by regex-match on existing filenames in --dir
		nextSequence int
	)

	app := &cli.App{
		Name:  "adr-er",
		Usage: "a friendly little thing for managing architectural decision records",
		// evaluates environment, assigning adrDirectory and nextSequence
		Before: func(ctx *cli.Context) error {
			// TODO: these blocks can hold error-cases, but we need file logging to report them.

			// determine correct output dir
			// don't return on error, just use zero-value (will trigger stdout flag)
			dir, _ := determineADRDirectory(ctx)
			adrDirectory = dir

			// determine next sequence number
			// don't return on error, just increment from zero
			seq, _ := utils.GetHighestSequenceNumber(adrDirectory)
			nextSequence = seq + 1

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "dir",
				Usage: `
root directory to store adr files.

if empty: 
  the application will search for a viable directory according to some conventions.  
  directories in CWD named "architectural-decision-records", "adr", or ".adr" will be checked. 
  we will set --dir to the first in the the list that is  
  a) empty, or b) holds only adr files and optionally subdirectories.
if provided:
  the application will not validate contents - we'll trust your judgement

if something goes wrong: 
in any case, files will be written to stdout if we have a meaningless/dangerous --dir (like /, literal "", or -)  
`,
				Aliases: []string{"d"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "create",
				Aliases:     []string{"c"},
				Usage:       "create a new adr document",
				Description: "new is used to create a brand-spankin-new adr document",
				Action: func(ctx *cli.Context) error {
					return create.NewCommand(adrDirectory, nextSequence).Action(ctx)
				},
			},
			{
				Name:        "view",
				Aliases:     []string{"v"},
				Usage:       "view existing ADR history",
				Description: "runs a tui application for reading historical ADRs",
				Action: func(ctx *cli.Context) error {
					return view.NewCommand(adrDirectory).Action(ctx)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("oops. app exited with an error:", err)
	}
}

// determineADRDirectory determines the correct root/output directory for ADR files
// returns the normalized absolute path.
func determineADRDirectory(ctx *cli.Context) (string, error) {
	var (
		err       error  // an error
		outputDir string // normalized dir
		dir       string // intermediate dir var, from either flag or LocateADRDirectory
	)

	// init dir based on --dir flag: if provided, use it; if not use the conventions codified in utils.LocateADRDirectory
	if userDir := ctx.String("dir"); userDir != "" {
		dir = userDir
	} else {
		dir, err = utils.LocateADRDirectory("")
		if err != nil {
			return "", fmt.Errorf("error evaluating candidate directories: %w", err)
		}
	}

	outputDir, err = filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("error normalizing path %s: %w", dir, err)
	}

	if _, err = os.Stat(outputDir); err != nil {
		return "", fmt.Errorf("error accessing path %s: %w", outputDir, err)
	}

	return outputDir, nil
}
