package main

import (
	"context"
	"fmt"
	"godot_linter/printer"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"godot_linter/styler"

	"github.com/urfave/cli/v3"
)

const ROOT = "./"

func main() {
	cmd := &cli.Command{
		Name:      "Godot Beautifier",
		Usage:     "Beautify/format GDScript code!",
		UsageText: "godot-beautifier [path to project/file] [args...]",

		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Value:   false,
				Usage:   "verbose logging of editing process,",
			},
			&cli.BoolFlag{
				Name:    "dry",
				Aliases: []string{"d"},
				Value:   false,
				Usage:   "don't write changed files, use with verbose for testing",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var input_path string
			if cmd.NArg() == 1 {
				input_path = cmd.Args().Get(0)
			} else if cmd.NArg() == 0 {
				input_path = ROOT
			} else {
				printer.PrintError("Too many arguments, only provide the path to the godot project or none for cwd")
				os.Exit(1)
			}

			printer.PrintNormal(fmt.Sprintf("Using godot project at: `%s`", input_path))

			var files []string
			if strings.HasSuffix(input_path, ".gd") {
				// Is single file
				files = append(files, input_path)
				printer.PrintNormal("GDScript file provided: " + input_path)
			} else {
				var err error

				// Is dir
				files, err = scan_gd_files(input_path)
				if err != nil {
					printer.PrintError("Not continuing, could not open all files in project at " + input_path)
					os.Exit(1)
				} else if len(files) == 0 {
					printer.PrintError("Not continuing, no GDScript files found." + input_path)
					os.Exit(1)
				}

				printer.PrintNormal("GDScript files found:")
				printer.PPrintArray(files)

				keep_going := printer.AskConfirmation("Continue to process?")
				if !keep_going {
					printer.PrintNormal("Exiting")
					os.Exit(0)
				}
			}

			backup_files(input_path, files)

			start := time.Now() // Before line
			total, errored := lint_files_mt(files, cmd.Bool("v"), cmd.Bool("d"))
			elapsed := time.Since(start) // After line
			printer.PrintNormal(fmt.Sprintf("Execution took %s for %d files (%d failed)", elapsed, total, errored))

			return nil
		},
	}

	cmd.Run(context.Background(), os.Args)
}

func scan_gd_files(local_root string) ([]string, error) {
	var matches []string
	err := filepath.WalkDir(local_root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// Log and skip errored files/directories
			printer.PrintWarning(fmt.Sprintf("Error accessing %s: %v\n", path, err))
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".gd") {
			matches = append(matches, path)
		}
		return nil
	})
	return matches, err
}

func backup_files(local_root string, locations []string) error {
	path, err := NewBackup(local_root, locations)

	if err != nil {
		printer.PrintError("Failure to create backup, exiting now without changes.")
		printer.PrintError("Raw: " + err.Error())
		os.Exit(1)
	} else {
		printer.PrintInfo("Backup of all *.gd files saved to " + path)
	}

	return nil
}

func lint_files_mt(files []string, verbose bool, dry bool) (total int, errored int) {
	var wg sync.WaitGroup

	ch := make(chan error)
	not_completed := 0

	go func() {
		for state := range ch {
			printer.PrintWarning(state.Error())
			not_completed++
		}
	}()

	for _, file := range files {
		wg.Add(1)

		go func(path string) {
			defer wg.Done()
			styler.LintFile(path, ch, verbose, dry)
		}(file)
	}

	wg.Wait()

	close(ch)

	return len(files), not_completed
}

func makePathLocal(path string, local_root string) string {
	return filepath.Base(local_root) + "/" + strings.TrimPrefix(path, local_root)
}
