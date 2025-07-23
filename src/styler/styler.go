package styler

// TODO:
// Block exports and onreadys and local vars in the tokeniser

import (
	"fmt"
	"godot_linter/printer"
	"os"
	"slices"
	"strings"

	tk "godot_linter/styler/tokendef"
	"godot_linter/styler/tokeniser"
)

type TokenizerError struct {
	FilePath string
	Message  string
}

func (terr TokenizerError) Error() string {
	return fmt.Sprintf("Error tokenising file %s: %s", terr.FilePath, terr.Message)
}

func LintFile(path string, ch chan error, verbose bool, dry bool) {
	if verbose {
		printer.PrintNormal("Linting " + path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		ch <- err
	}

	lines := strings.Split(string(data), "\n")

	tokens, err := tokeniser.Tokenize(lines)
	if err != nil {
		terr := TokenizerError{FilePath: path, Message: err.Error()}
		ch <- terr
		return
	}

	if verbose {
		// Print before changes
		println("<== Tokenisation")
		for _, t := range tokens {
			print(tk.BlockTypeToString(t.Type) + ":\n")
			printer.PPrintArray(t.Content)
		}
	}

	// Sort blocks by enum order
	slices.SortStableFunc(tokens, func(a, b tk.Block) int {
		return int(a.Type) - int(b.Type)
	})

	if verbose {
		// After
		println("<== Tokens after sort")
		for _, t := range tokens {
			print(tk.BlockTypeToString(t.Type) + ":\n")
			printer.PPrintArray(t.Content)
		}
	}

	det := Detokenise(tokens)

	if verbose {
		// After
		println("<== Final")
		print(det + "\n")
	}

	// Write edited file
	if !dry {
		err = os.WriteFile(path, []byte(det), 0644)
		if err != nil {
			ch <- err
		}
	}

	printer.PrintSuccess("Finished: " + path)
}

// Order parts
//	Extends
//	Exports
//	Onready
//	Local vars
//	Functions

// Double space functions, 2 spaces from last thing above

// Remove default comments

// Add return typing

func Detokenise(tokens []tk.Block) string {
	file := ""
	for i, token := range tokens {
		file += strings.Join(token.Content, "\n")

		if i+1 == len(tokens) {
			break
		}

		// Start with 1 newline
		newlines := 2

		switch token.Type {
		case tk.ClassName:
			newlines--
		case tk.Function, tk.Ready, tk.Init:
			newlines++
		}

		file += strings.Repeat("\n", newlines)
	}

	return file
}
