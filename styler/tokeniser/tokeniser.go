package tokeniser

import (
	"errors"
	"slices"
	"strings"

	"godot_linter/printer"
	tk "godot_linter/styler/tokendef"
)

const indent = "	"

type HandlerFunc func(line string, lines []string, i *int, blocks *[]tk.Block, linkedAbove *[]string)

var handlers = map[string]HandlerFunc{
	"@tool":      handleTool,
	"class_name": handleClassName,
	"extends":    handleExtend,
	`"""`:        handleDocString,
	"signal":     handleSignals,
	"enum":       handleEnum,
	"const":      handleConstants,
	"@export":    handleExport,
	"@onready":   handleOnReady,
	"class":      handleClass,
	"static":     handleStatic,
	"var":        handleVar,
	"func":       handleFunction,
	"#":          handleComment,
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
			fn(line, lines, &i, &blocks, &linked_above)
			flush_above()
		} else {
			// Cases not recognised by first character
			switch {
			case isIndentOnly(line):
				continue
			case line[0] == '#', stripIndents(&line)[0] == '#':
				handleComment(line, lines, &i, &blocks, &linked_above)
			default:
				handleUnknown(line, lines, &i, &blocks, &linked_above)
				unknown_component = true
			}
		}
	}

	// Scan for unknown component
	if !unknown_component {
		for _, block := range blocks {
			if block.Type == tk.Unknown {
				unknown_component = true
				break
			}
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
	hadBlank := false

outer:
	for ; i < len(lines); i++ {
		for _, pre := range tk.Prefixes {
			if isIndentOnly(lines[i]) {
				hadBlank = true
				goto nl
			}

			if strings.HasPrefix(lines[i], pre) {
				if strings.HasPrefix(lines[i], exception) && hadBlank {
					// exception prefix after blank line(s)
					break outer
				} else if strings.HasPrefix(lines[i], exception) {
					// exception prefix without blank line
					goto nl
				} else {
					// non-exception prefix
					break outer
				}
			}
		}

	nl:
		continue
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
func handleStatic(line string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	if line[7] == 'v' {
		handleStaticVar_(line, lines, idx, blocks, linkedAbove)
	} else if line[7] == 'f' {
		handleStaticFunction_(line, lines, idx, blocks, linkedAbove)
	} else {
		handleUnknown(line, lines, idx, blocks, linkedAbove)
	}
}
func handleStaticVar_(_ string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := findImplicitExtendedBlockEnd(lines, *idx, "static var")
	varLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.LocalVar,
		trimBlankLines(consumeWithAbove(linkedAbove, varLines...)),
	))
	*idx = end
}
func handleStaticFunction_(_ string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := findBlockEnd(lines, *idx)
	fnLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.Function,
		trimBlankLines(consumeWithAbove(linkedAbove, fnLines...)),
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
func handleFunction(line string, lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	if strings.HasPrefix(line, "func _init(") {
		handleInit_(lines, idx, blocks, linkedAbove)
		return
	} else if strings.HasPrefix(line, "func _ready(") {
		handleReady_(lines, idx, blocks, linkedAbove)
		return
	}

	end := findBlockEnd(lines, *idx)
	fnLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.Function,
		trimBlankLines(consumeWithAbove(linkedAbove, fnLines...)),
	))
	*idx = end
}
func handleInit_(lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := findBlockEnd(lines, *idx)
	initLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.Init,
		trimBlankLines(consumeWithAbove(linkedAbove, initLines...)),
	))
	*idx = end
}
func handleReady_(lines []string, idx *int, blocks *[]tk.Block, linkedAbove *[]string) {
	end := findBlockEnd(lines, *idx)
	readyLines := lines[*idx : end+1]
	*blocks = append(*blocks, makeBlock(tk.Ready,
		trimBlankLines(consumeWithAbove(linkedAbove, readyLines...)),
	))
	*idx = end
}
func handleComment(line string, _ []string, _ *int, blocks *[]tk.Block, linkedAbove *[]string) {
	*linkedAbove = append(*linkedAbove, line)
}
func handleUnknown(line string, _ []string, _ *int, blocks *[]tk.Block, linkedAbove *[]string) {
	printer.PrintWarning("Unknown line parsed: " + line)
	*blocks = append(*blocks, makeBlock(tk.Unknown,
		trimBlankLines(consumeWithAbove(linkedAbove, line)),
	))
	flushAbove(linkedAbove)
}

// flushAbove only resets linkedAbove slice
func flushAbove(linkedAbove *[]string) {
	*linkedAbove = nil
}

func isIndentOnly(s string) bool {
	return strings.Trim(s, " \t") == ""
}

func stripIndents(line *string) string {
	return strings.TrimLeft(*line, indent)

}
