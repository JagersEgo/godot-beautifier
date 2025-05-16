package main

import (
	"godot_linter/printer"
	"os"
	"strings"
)

type BlockType int8

const indent = "	"

const (
	Tool BlockType = iota
	ClassName
	Extend
	DocString

	Signals
	Enum
	Constants
	Export
	Onready
	Class
	LocalVar

	Init
	Ready
	Function
	Unknown
)

type Block struct {
	Type    BlockType
	Content []string
}

func lint_file(path string, ch chan error) {

	printer.PrintNormal("Linting " + path)

	data, err := os.ReadFile(path)
	if err != nil {
		ch <- err
	}

	lines := strings.Split(string(data), "\n")

	printer.DebugPrintArray(lines)

	// // Write edited file
	// err = os.WriteFile(path, []byte(content), 0644)
	// if err != nil {
	// 	ch<-err
	// }

	// Order parts
	//	Extends
	//	Exports
	//	Onready
	//	Local vars
	//	Functions

	// Double space functions, 2 spaces from last thing above

	// Remove default comments

	// Add return typing
}

func tokenize(lines []string) []Block {
	linked_above := 0
	linked_below := 0
	//tokens := make([]Block, len(lines)/2)

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "@tool"): // Tool
		case strings.HasPrefix(line, "class_name"): // ClassName
		case strings.HasPrefix(line, "extends"): // Extends
		case strings.HasPrefix(line, "\"\"\""): // DocString
		case strings.HasPrefix(line, "signal"): // Signals
		case strings.HasPrefix(line, "enum"): // Enum
		case strings.HasPrefix(line, "const"): // Constants
		case strings.HasPrefix(line, "@export"): // Export
		case strings.HasPrefix(line, "@onready"): // Onready
		case strings.HasPrefix(line, "class "): // Class
		case strings.HasPrefix(line, "var"): // LocalVar
		case strings.HasPrefix(line, "func _init("): // Init
		case strings.HasPrefix(line, "func _ready("): // Ready
		case strings.HasPrefix(line, "func"): // Function

		case strings.HasPrefix(line, "#"): // Could be DocString or Unknown
			linked_above++

		case strings.HasPrefix(line, indent): // Indented
		case line == "": // Blank line

		default: // Unknown
		}
	}

	return nil
}

func find_function_end(lines []string, idx int) int {
	for ; strings.HasPrefix(lines[idx], indent); idx++ {
	}

	return idx
}

func find_docstring_end(lines []string, idx int) int {
	for ; strings.HasPrefix(lines[idx], "\"\"\""); idx++ {
	}

	return idx
}

func make_block(btype BlockType) Block {
	block := Block{Type: btype}
	return block
}
