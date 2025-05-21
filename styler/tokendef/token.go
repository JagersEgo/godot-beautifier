package tokendef

type BlockType int8

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

var Prefixes = []string{
	"@tool",
	"class_name",
	"extends",
	"\"\"\"",
	"signal",
	"enum",
	"const",
	"@export",
	"@onready",
	"class ",
	"var",
	"func ",
	"#", // for comment lines
}

type Block struct {
	Type    BlockType
	Content []string
}

func BlockTypeToString(bt BlockType) string {
	switch bt {
	case Tool:
		return "Tool"
	case ClassName:
		return "ClassName"
	case Extend:
		return "Extend"
	case DocString:
		return "DocString"
	case Signals:
		return "Signals"
	case Enum:
		return "Enum"
	case Constants:
		return "Constants"
	case Export:
		return "Export"
	case Onready:
		return "Onready"
	case Class:
		return "Class"
	case LocalVar:
		return "LocalVar"
	case Init:
		return "Init"
	case Ready:
		return "Ready"
	case Function:
		return "Function"
	case Unknown:
		return "Unknown"
	default:
		return "Invalid"
	}
}
