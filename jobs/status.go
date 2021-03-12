package jobs

type CompletionStatus int

const (
	Inactive CompletionStatus = iota
	Running
	Success
	Skip
	Error
)

func (h CompletionStatus) String() string {
	switch h {
	case Inactive:
		return "Inactive"
	case Running:
		return "Running"
	case Success:
		return "Success"
	case Skip:
		return "Skip"
	case Error:
		return "Error"
	default:
		return "unknown"
	}
}
