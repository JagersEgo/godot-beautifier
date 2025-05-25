package tokendef

type VariableType int8

const (
	Static VariableType = iota
	Export
	Onready
	Local
)

var VariablePrefixes = []string{
	"@export",
	"@onready",
	"var",
	"static",
}

type VariableBlock struct {
	Content      []string
	VariableType VariableType
}

func (v VariableBlock) GetType() BlockType {
	return LocalVar
}

func (v VariableBlock) GetContent() []string {
	return v.Content
}

func (v VariableBlock) BlockTypeToString() string {
	switch v.VariableType {
	case Local:
		return "Variable"
	case Export:
		return "ExportedVariable"
	case Onready:
		return "OnReadyVariable"
	case Static:
		return "StaticVar"
	default:
		panic("out of range- not in the enum")
	}
}
