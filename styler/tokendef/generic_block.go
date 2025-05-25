package tokendef

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

func (g GenericBlock) BlockTypeToString() string {
	switch g.Type {
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
