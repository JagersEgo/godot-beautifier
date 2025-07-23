package tokeniser

import (
	"reflect"
	"testing"
)

func TestConvertSpaceIndentsToTabs(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name: "Simple 4-space indent",
			input: []string{
				"    func main() {",
				"        fmt.Println(\"hello\")",
				"    }",
			},
			expected: []string{
				"\tfunc main() {",
				"\t\tfmt.Println(\"hello\")",
				"\t}",
			},
		},
		{
			name: "Mixed indentation widths",
			input: []string{
				"  level1",
				"    level2",
				"      level3",
			},
			expected: []string{
				"\tlevel1",
				"\t\tlevel2",
				"\t\t\tlevel3",
			},
		},
		{
			name: "No indentation",
			input: []string{
				"func main() {}",
			},
			expected: []string{
				"func main() {}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ConvertSpaceIndentsToTabs(tt.input)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("ConvertSpaceIndentsToTabs failed.\nInput:\n%v\nExpected:\n%v\nGot:\n%v",
					tt.input, tt.expected, actual)
			}
		})
	}
}
