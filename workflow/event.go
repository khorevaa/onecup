package workflow

type Event struct {
	Job      string `json:"job,omitempty"`
	Task     string `json:"task,omitempty"`
	Step     string `json:"step,omitempty"`
	Event    string `json:"event,omitempty"`
	Err      error  `json:"err"`
	JobState `json:"status"`
	Timing   int64 `json:"timing"`
}
