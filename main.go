package main

import (
	"fmt"
	"godot_linter/printer"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const ROOT = "./"

func main() {
	var project_root string
	if len(os.Args) == 2 {
		project_root = os.Args[1]
	} else if len(os.Args) < 2 {
		project_root = ROOT
	} else {
		printer.PrintError("Too many arguments, only provide the path to the godot project or none for cwd")
		os.Exit(1)
	}

	printer.PrintNormal(fmt.Sprintf("Using godot project at: `%s`", project_root))

	files, err := scan_gd_files(project_root)
	if err != nil {
		printer.PrintError("Not continuing, could not open all files in project at " + project_root)
		os.Exit(1)
	} else if len(files) == 0 {
		printer.PrintError("Not continuing, no GDScript files found." + project_root)
		os.Exit(1)
	}

	printer.PrintNormal("GDScript files found:")
	printer.PPrintArray(files)

	printer.PrintObvious(fmt.Sprintf("Ensure the project root is at `%s`", project_root))
	keep_going := printer.AskConfirmation("Continue to lint?")
	if !keep_going {
		printer.PrintNormal("Exiting")
		os.Exit(0)
	}

	backup_files(project_root, files)

	lint_files_st(files)

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

func lint_files_mt(files []string) {
	var wg sync.WaitGroup

	ch := make(chan error)

	go func() {
		for state := range ch {
			printer.PrintWarning("Error while linting: " + state.Error())
		}
	}()

	for _, file := range files {
		wg.Add(1)

		go func(path string) {
			defer wg.Done()
			lint_file(path, ch)
		}(file)
	}

	wg.Wait()

	close(ch)

	return
}

func lint_files_st(files []string) {
	for _, file := range files {

		lint_file(file, nil)
	}

	return
}
