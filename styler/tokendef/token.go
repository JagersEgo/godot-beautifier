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

type VariableType int8 (
const (
	Local VariableType = iota
	Export
	Onready
)

const ()

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
	"static var",
	"func ",
	"static func ",
	"#",
}

type Block interface {
	GetType() BlockType
	GetContent() []string
}

type GenericBlock struct {
	Type    BlockType
	Content []string
}

func (g GenericBlock) GetType() BlockType {
	return g.Type
}

func (g GenericBlock) GetContent() []string {
	return g.Content
}


type VariableBlock struct {
	Content []string
	VariableType VariableBlocKType
}

func (v VariableBlock) GetType() BlockType {
	return LocalVar
}

func (v VariableBlock) GetContent() []string {
	return v.Content
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
