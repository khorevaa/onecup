package jobs

type Action func(ctx *Context) error

type StepType int

const (
	BeforeType StepType = iota
	DefaultType
	AfterType
	ErrorType
	AlwaysType
)

type Step struct {
	Name   string
	On     StepType
	Action Action

	params Params

	Skipped        bool
	SkippedMessage string
}

func (s *Step) Run(ctx *Context) {

	ctx.currentStep = s

	if ctx.Fault() && !(s.On == ErrorType || s.On == AlwaysType) {
		s.Skip()
		return
	}

	err := s.Action(ctx)
	if err != nil {
		ctx.err = err
	}
}

func (s *Step) Skip() {
	s.Skipped = true
}

func (s *Step) SkipMsg(msg string) {

	s.SkippedMessage = msg
	s.Skipped = true

}
