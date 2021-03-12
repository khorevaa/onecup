package jobs

type EventType int

const (
	AllEvents = iota
	ErrorEvents
	TimingEvents
	CompletionEvents
)

func (h EventType) String() string {
	switch h {
	case AllEvents:
		return "Event"
	case ErrorEvents:
		return "Error"
	case TimingEvents:
		return "Timing"
	case CompletionEvents:
		return "Complete"
	default:
		return "Unknown"
	}
}

type Event struct {
	Type   EventType        `json:"type"`
	Job    string           `json:"job,omitempty"`
	Task   string           `json:"task,omitempty"`
	Step   string           `json:"step,omitempty"`
	Event  string           `json:"event,omitempty"`
	Err    error            `json:"err"`
	Status CompletionStatus `json:"status"`
	Timing int64            `json:"timing"`
}
