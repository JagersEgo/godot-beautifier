package tokeniser

import (
	"slices"
	"strings"

	tk "godot_linter/styler/tokendef"
)

const indent = "	"

func Tokenize(lines []string) []tk.Block {
	lines = ConvertSpaceIndentsToTabs(lines)

	var blocks []tk.Block
	var linked_above []string
	//tokens := make([]Block, len(lines)/2)

	flush_above := func() {
		linked_above = nil
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		//trim := trimBlankLines(line)

		switch {
		// ---- CORE ITEMS ----
		case strings.HasPrefix(line, "@tool"):
			blocks = append(blocks, make_block(tk.Tool, consume_with_above(&linked_above, []string{line}...)))
			flush_above()

		case strings.HasPrefix(line, "class_name"):
			blocks = append(blocks, make_block(tk.ClassName, consume_with_above(&linked_above, []string{line}...)))
			flush_above()

		case strings.HasPrefix(line, "extends"):
			blocks = append(blocks, make_block(tk.Extend, consume_with_above(&linked_above, []string{line}...)))
			flush_above()

		case strings.HasPrefix(line, "\"\"\""):
			// doc‐string might be multi-line
			end := find_docstring_end(lines, i)
			docLines := lines[i : end+1]
			blocks = append(blocks, make_block(tk.DocString, consume_with_above(&linked_above, docLines...)))
			i = end
			flush_above()

		case strings.HasPrefix(line, "signal"):
			blocks = append(blocks, make_block(tk.Signals, consume_with_above(&linked_above, []string{line}...)))
			flush_above()

		case strings.HasPrefix(line, "enum"):
			blocks = append(blocks, make_block(tk.Enum, consume_with_above(&linked_above, []string{line}...)))
			flush_above()

		case strings.HasPrefix(line, "const"):
			blocks = append(blocks, make_block(tk.Constants, consume_with_above(&linked_above, []string{line}...)))
			flush_above()

		case strings.HasPrefix(line, "@export"):
			end := findSegmentEnd(lines, i, "@export")
			fnLines := lines[i : end+1]
			blocks = append(blocks, make_block(tk.Export, consume_with_above(&linked_above, fnLines...)))
			i = end
			flush_above()

		case strings.HasPrefix(line, "@onready"):
			end := findSegmentEnd(lines, i, "@onready")
			fnLines := lines[i : end+1]
			blocks = append(blocks, make_block(tk.Ready, consume_with_above(&linked_above, fnLines...)))
			i = end
			flush_above()

		case strings.HasPrefix(line, "class "):
			// class body is indented
			end := findBlockEnd(lines, i)
			classLines := lines[i : end+1]
			blocks = append(blocks, make_block(tk.Class, consume_with_above(&linked_above, classLines...)))
			i = end
			flush_above()

		case strings.HasPrefix(line, "var"):
			end := findSegmentEnd(lines, i, "var")
			fnLines := lines[i : end+1]
			blocks = append(blocks, make_block(tk.LocalVar, consume_with_above(&linked_above, fnLines...)))
			i = end
			flush_above()

		case strings.HasPrefix(line, "func _init("):
			end := findBlockEnd(lines, i)
			initLines := lines[i : end+1]
			blocks = append(blocks, make_block(tk.Init, consume_with_above(&linked_above, initLines...)))
			i = end
			flush_above()

		case strings.HasPrefix(line, "func _ready("):
			end := findBlockEnd(lines, i)
			readyLines := lines[i : end+1]
			blocks = append(blocks, make_block(tk.Ready, consume_with_above(&linked_above, readyLines...)))
			i = end
			flush_above()

		case strings.HasPrefix(line, "func "):
			// generic function
			end := findBlockEnd(lines, i)
			fnLines := lines[i : end+1]
			blocks = append(blocks, make_block(tk.Function, consume_with_above(&linked_above, fnLines...)))
			i = end
			flush_above()

		// ---- “Above” LINES ----
		case strings.HasPrefix(line, "#"), line == "", strings.HasPrefix(line, indent):
			// comment, blank, or indent-only: hold for next block
			linked_above = append(linked_above, line)

		default:
			// nothing matched → Unknown standalone
			blocks = append(blocks, make_block(tk.Unknown, consume_with_above(&linked_above, []string{line}...)))
			flush_above()
		}
	}

	return blocks
}

// countIndent counts how many indent tabs at line start.
func countIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, indent))
}

func consume_with_above(linked_above *[]string, content ...string) []string {
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

func make_block(btype tk.BlockType, lines []string) tk.Block {
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
