package jobs

type HandlerType int

const (
	BeforeType HandlerType = iota
	DefaultType
	AfterType
	ErrorType
	AlwaysType
)

func (h HandlerType) String() string {
	switch h {
	case BeforeType:
		return "Before"
	case DefaultType:
		return "Default"
	case AfterType:
		return "After"
	case ErrorType:
		return "Error"
	case AlwaysType:
		return "Always"
	default:
		return "unknown handler type"
	}
}
