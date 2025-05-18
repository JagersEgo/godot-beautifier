package styler

// TODO:
// Block exports and onreadys and local vars in the tokeniser

import (
	"godot_linter/printer"
	"os"
	"slices"
	"strings"

	tk "godot_linter/styler/tokendef"
	"godot_linter/styler/tokeniser"
)

const VERBOSE = false

func LintFile(path string, ch chan error) {

	printer.PrintNormal("Linting " + path)

	data, err := os.ReadFile(path)
	if err != nil {
		ch <- err
	}

	lines := strings.Split(string(data), "\n")

	tokens := tokeniser.Tokenize(lines)

	if VERBOSE {
		// Print before changes
		for _, t := range tokens {
			print(tk.BlockTypeToString(t.Type) + ":\n")
			printer.PPrintArray(t.Content)
		}
	}

	slices.SortStableFunc(tokens, func(a, b tk.Block) int {
		return int(a.Type) - int(b.Type)
	})

	if VERBOSE {
		// After
		println("---")
		for _, t := range tokens {
			print(tk.BlockTypeToString(t.Type) + ":\n")
			printer.PPrintArray(t.Content)
		}
	}

	det := Detokenise(tokens)

	if VERBOSE {
		print(det)
	}

	// Write edited file
	err = os.WriteFile(path, []byte(det), 0644)
	if err != nil {
		ch <- err
	}

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

		if token.Type == tk.ClassName {
			newlines--
		}

		file += strings.Repeat("\n", newlines)
	}

	return file
}
