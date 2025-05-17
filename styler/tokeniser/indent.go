package tokeniser

import (
	"strings"
	"unicode"
)

// Public API: ConvertSpaceIndentsToTabs takes lines of text and returns lines with space indents replaced by tabs.
func ConvertSpaceIndentsToTabs(lines []string) []string {
	indentSizes := getIndentSizes(lines)
	if len(indentSizes) == 0 {
		return lines // nothing to convert
	}

	indentUnit := gcdOfSlice(indentSizes)
	if indentUnit == 0 {
		return lines // avoid division by zero
	}

	var converted []string
	for _, line := range lines {
		converted = append(converted, convertLine(line, indentUnit))
	}
	return converted
}

// --- Internal helpers ---

// getIndentSizes scans lines to extract leading space counts (excluding empty and unindented lines)
func getIndentSizes(lines []string) []int {
	var sizes []int
	for _, line := range lines {
		spaceCount := countLeadingSpaces(line)
		if spaceCount > 0 {
			sizes = append(sizes, spaceCount)
		}
	}
	return sizes
}

// countLeadingSpaces returns the number of leading space characters
func countLeadingSpaces(s string) int {
	count := 0
	for _, r := range s {
		if r == ' ' {
			count++
		} else if r == '\t' {
			// Optional: normalize tabs to fixed space width (e.g., 4)
			count += 4
		} else {
			break
		}
	}
	return count
}

// convertLine replaces leading spaces with tabs according to indent unit
func convertLine(line string, indentUnit int) string {
	spaceCount := countLeadingSpaces(line)
	nTabs := spaceCount / indentUnit
	nSpaces := spaceCount % indentUnit

	// Strip the original leading spaces
	trimmed := strings.TrimLeftFunc(line, unicode.IsSpace)
	return strings.Repeat("\t", nTabs) + strings.Repeat(" ", nSpaces) + trimmed
}

// gcdOfSlice returns the GCD of a slice of integers
func gcdOfSlice(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	result := nums[0]
	for _, n := range nums[1:] {
		result = gcd(result, n)
	}
	return result
}

// gcd computes the greatest common divisor
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}
