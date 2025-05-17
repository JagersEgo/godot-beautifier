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

func LintFile(path string, ch chan error) {

	printer.PrintNormal("Linting " + path)

	data, err := os.ReadFile(path)
	if err != nil {
		ch <- err
	}

	lines := strings.Split(string(data), "\n")

	tokens := tokeniser.Tokenize(lines)

	// Print before changes
	for _, t := range tokens {
		print(tk.BlockTypeToString(t.Type) + ":\n")
		printer.PPrintArray(t.Content)
	}

	slices.SortStableFunc(tokens, func(a, b tk.Block) int {
		return int(a.Type) - int(b.Type)
	})

	// After
	println("---")
	for _, t := range tokens {
		print(tk.BlockTypeToString(t.Type) + ":\n")
		printer.PPrintArray(t.Content)
	}

	det := Detokenise(tokens)
	print(det)

	// // Write edited file
	// err = os.WriteFile(path, []byte(DETOKENISED), 0644)
	// if err != nil {
	// 	ch <- err
	// }

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
		file += "\n\n"
		if token.Type == tk.Function || token.Type == tk.Init || token.Type == tk.Ready {
			file += "\n"
		}

	}
	return file
}
