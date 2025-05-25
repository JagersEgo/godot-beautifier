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
	Class
	LocalVar

	Init
	Ready
	Function
	Unknown
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
	BlockTypeToString() string
}
