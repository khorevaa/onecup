package jobs

type StepInterface interface {
	Name() string
	Action(ctx Context) error
	Handler() HandlerType
}

type StepBuilder interface {
	Handler() HandlerType
	Build() Task
}

func Step(step StepInterface) StepBuilder {
	return stepBuilder{
		step,
	}
}

func NewStepBuilder(name string, fn func(ctx Context) error, h HandlerType) StepBuilder {

	return stepBuilder{customStepBuilder{
		name:    name,
		fn:      fn,
		handler: h,
	}}
}

type stepBuilder struct {
	StepInterface
}

var _ StepInterface = (*customStepBuilder)(nil)

type customStepBuilder struct {
	name    string
	fn      func(ctx Context) error
	handler HandlerType
}

func (c customStepBuilder) Name() string {
	return c.name
}

func (c customStepBuilder) Action(ctx Context) error {
	return c.fn(ctx)
}

func (c customStepBuilder) Handler() HandlerType {
	return c.handler
}

func (b stepBuilder) Build() Task {
	return NewStep(
		b.Name(),
		b.Action,
	)
}
