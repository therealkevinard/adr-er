package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/therealkevinard/adr-er/commands"
	"github.com/therealkevinard/adr-er/utils"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "adr-er",
		Usage: "a friendly little thing for managing architectural decision records",
		// this is a good place to evaluate environment and set initial flags
		Before: func(_ *cli.Context) error { return nil },
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
  the special value "-" can be used to indicate stdout  
`,
				Aliases: []string{"d"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "create",
				Aliases:     []string{"new"},
				Usage:       "create a new adr document",
				Description: "new is used to create a brand-spankin-new adr document",
				// this is a good place to evaluate environment and set initial flags
				Before: func(_ *cli.Context) error { return nil },
				Action: func(ctx *cli.Context) error {
					var (
						// root dir to write files into
						outputDir string
						// flag if the user provided the directory. this allows us to ignore some safeguards.
						userDefinedOutputDir bool
						// next int sequence.
						// when bootstrapping, this is detemined by regex-match on existing filenames in --dir
						nextSequence int
					)

					// user provided --dir flag
					if userDir := ctx.String("dir"); userDir != "" {
						outputDir = userDir
						userDefinedOutputDir = true
					}
					// if no user-provided --dir arg was supplied, attempt utils.LocateADRDirectory
					if ctx.String("dir") == "" {
						found, _ := utils.LocateADRDirectory("")
						if found != "" {
							outputDir = found // yay! use the one we found
						} else {
							outputDir = "-" // fallback to stdout
						}
					}

					// for sout, we can finish early
					if outputDir == "-" {
						return commands.NewCreate(outputDir, userDefinedOutputDir, 0).Action(ctx)
					}

					// normalize absolute path
					if !path.IsAbs(outputDir) {
						cwd, _ := os.Getwd()
						outputDir = path.Join(cwd, outputDir)
					}
					// check dir exists
					_, err := os.Stat(outputDir)
					if err != nil {
						return fmt.Errorf("determined output directory %s is inaccessible: %w", outputDir, err)
					}

					// identify next sequence number
					// TODO: this _can_ error, but we need file logging before we can report it
					currentSequence, _ := utils.GetHighestSequenceNumber(outputDir)
					nextSequence = currentSequence + 1

					return commands.NewCreate(outputDir, userDefinedOutputDir, nextSequence).Action(ctx)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("oops. app exited with an error:", err)
	}
}
