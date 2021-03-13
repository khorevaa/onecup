package jobs

type Task interface {
	Name() string
	Stats() Stats
	Fault() bool
	Success() bool
	Status() CompletionStatus
	Run(ctx Context) (output Values, err error)
}
