package jobs

type JobStatus int

const (
	InActiveStatus JobStatus = iota
	RunningStatus
	SuccessStatus
	SkipStatus
	FaultStatus
)
