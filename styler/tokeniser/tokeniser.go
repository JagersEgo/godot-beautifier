package tokeniser

import (
	"errors"
	"slices"
	"strings"

	tk "godot_linter/styler/tokendef"
)

const indent = "	"

type HandlerFunc func(line string, lines []string, i *int, blocks *[]tk.Block, linkedAbove *[]string)

var handlers = map[string]HandlerFunc{
	"@tool":        handleTool,
	"class_name":   handleClassName,
	"extends":      handleExtend,
	`"""`:          handleDocString,
	"signal":       handleSignals,
	"enum":         handleEnum,
	"const":        handleConstants,
	"@export":      handleExport,
	"@onready":     handleOnReady,
	"class":        handleClass,
	"static":       handleStatic,
	"var":          handleVar,
	"func _init(":  handleInit,
	"func _ready(": handleReady,
	"func":         handleFunction,
}

func Tokenize(lines []string) ([]tk.Block, error) {
	lines = ConvertSpaceIndentsToTabs(lines)

	var blocks []tk.Block
	blocks = make([]tk.Block, 0, len(lines)/2)

	var linked_above []string

	unknown_component := false

	flush_above := func() {
		linked_above = nil
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		fn, ok := handlers[strings.Split(line, " ")[0]]
		if ok {
			fn(line, lines, &i)
			flush_above()
		} else {
			handleUnknown(line, lines, &i, blocks, linked_above)
			unknown_component = true
		}
	}

	var e error
	if unknown_component {
		e = errors.New("Unknown component in script")
	} else {
		e = nil
	}

	return blocks, e
}

// countIndent counts how many indent tabs at line start.
func countIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, indent))
}

func consumeWithAbove(linked_above *[]string, content ...string) []string {
	trimmedAbove := trimBlankLines(*linked_above)
	*linked_above = nil

	if len(trimmedAbove) == 0 {
		return *&content
	}
	out := slices.Concat(trimmedAbove, content)
	return out

}

// findBlockEnd finds the last line index of a func/class by indent level.
func findBlockEnd(lines []string, idx int) int {
	baseIndent := countIndent(lines[idx])
	i := idx + 1
	for ; i < len(lines); i++ {
		if countIndent(lines[i]) <= baseIndent && strings.TrimSpace(lines[i]) != "" {
			break
		}
	}
	return i - 1
}

// findBlockEnd finds the last line index of a func/class by finding the start of the next block
func findImplicitBlockEnd(lines []string, idx int) int {
	i := idx + 1

outer:
	for ; i < len(lines); i++ {
		for _, pre := range tk.Prefixes {
			if strings.HasPrefix(lines[i], pre) {
				break outer
			}
		}
	}
	return i - 1
}

// findBlockEnd finds the last line index of a func/class by finding the start of the next block, extends the original type
func findImplicitExtendedBlockEnd(lines []string, idx int, exception string) int {
	i := idx + 1

outer:
	for ; i < len(lines); i++ {
		for _, pre := range tk.Prefixes {
			if !strings.HasPrefix(lines[i], exception) && strings.HasPrefix(lines[i], pre) {
				break outer
			}
		}
	}
	return i - 1
}

// Find when a segment of lines that start with `id` ends
func findSegmentEnd(lines []string, idx int, id string) int {
	i := idx + 1
	for ; i < len(lines); i++ {
		// stop when indent back to <= base or blank
		if !strings.HasPrefix(lines[i], id) {
			return i - 1
		}
	}
	return i
}

func find_docstring_end(lines []string, idx int) int {
	i := idx + 1
	for ; i < len(lines); i++ {
		if strings.HasSuffix(strings.TrimSpace(lines[i]), "\"\"\"") {
			return i
		}
	}
	// didn't find closer → just return start
	return idx
}

func makeBlock(btype tk.BlockType, lines []string) tk.Block {
	block := tk.Block{Type: btype, Content: lines}
	return block
}

func trimBlankLines(lines []string) []string {
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}

	end := len(lines)
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}

	return lines[start:end]
}

// ---- Handler implementations ----

func handleTool(line string, _ []string, _ *int, blocks *[]tk.Block, linkedAbove *[]string) {
	*blocks = append(*blocks, makeBlock(tk.Tool,
		consumeWithAbove(linkedAbove, line),
	))
}
func handleClassName(line string, _ []string, _ *int, blocks *[]tk.Block, linkedAbove *[]string) {
	*blocks = append(*blocks, makeBlock(tk.ClassName,
		consumeWithAbove(linkedAbove, line),
	))
}
func handleExtend(line string, _ []string, _ *int, blocks *[]tk.Block, linkedAbove *[]string) {
	*blocks = append(*blocks, makeBlock(tk.Extend,
		consumeWithAbove(linkedAbove, line),
	))
}
func handleDocString(_ string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := find_docstring_end(lines, *idx)
	docLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.DocString,
		consumeWithAbove(linkedAbove, docLines...),
	))
	*idx = end
}
func handleSignals(line string, _ []string, _ *int, blocks *[]tk.Block, linkedAbove *[]string) {
	*blocks = append(*blocks, makeBlock(tk.Signals,
		consumeWithAbove(linkedAbove, line),
	))
}
func handleEnum(_ string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := findImplicitBlockEnd(lines, *idx)
	enumLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.Enum,
		trimBlankLines(consumeWithAbove(linkedAbove, enumLines...)),
	))
	*idx = end
}
func handleConstants(_ string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := findImplicitBlockEnd(lines, *idx)
	constLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.Constants,
		trimBlankLines(consumeWithAbove(linkedAbove, constLines...)),
	))
	*idx = end
}
func handleExport(_ string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := findImplicitExtendedBlockEnd(lines, *idx, "@export")
	fnLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.Export,
		trimBlankLines(consumeWithAbove(linkedAbove, fnLines...)),
	))
	*idx = end
}
func handleOnReady(_ string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := findImplicitExtendedBlockEnd(lines, *idx, "@onready")
	fnLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.Ready,
		trimBlankLines(consumeWithAbove(linkedAbove, fnLines...)),
	))
	*idx = end
}
func handleClass(_ string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := findBlockEnd(lines, *idx)
	classLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.Class,
		consumeWithAbove(linkedAbove, classLines...),
	))
	*idx = end
}
func handleStaticVar(_ string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := findImplicitExtendedBlockEnd(lines, *idx, "static var")
	varLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.LocalVar,
		trimBlankLines(consumeWithAbove(linkedAbove, varLines...)),
	))
	*idx = end
}
func handleVar(_ string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := findImplicitExtendedBlockEnd(lines, *idx, "var")
	varLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.LocalVar,
		trimBlankLines(consumeWithAbove(linkedAbove, varLines...)),
	))
	*idx = end
}
func handleInit(_ string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := findBlockEnd(lines, *idx)
	initLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.Init,
		trimBlankLines(consumeWithAbove(linkedAbove, initLines...)),
	))
	*idx = end
}
func handleReady(_ string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := findBlockEnd(lines, *idx)
	readyLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.Ready,
		trimBlankLines(consumeWithAbove(linkedAbove, readyLines...)),
	))
	*idx = end
}
func handleFunction(_ string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := findBlockEnd(lines, *idx)
	fnLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.Function,
		trimBlankLines(consumeWithAbove(linkedAbove, fnLines...)),
	))
	*idx = end
}

func handleUnknown() (_ string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	return
}

// flushAbove only resets linkedAbove slice
func flushAbove(linkedAbove *[]string) {
	*linkedAbove = nil
}

func isIndentOnly(s string) bool {
	return strings.Trim(s, " \t") == ""
}
