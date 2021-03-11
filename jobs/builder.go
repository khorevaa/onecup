package jobs

type JobBuilder struct {
	name    string
	onError []Step
	after   []Step
	steps   []Step
	before  []Step
	always  []Step
}

func NewJob(name string) *JobBuilder {
	return &JobBuilder{
		name: name,
	}
}

func (b *JobBuilder) addStep(s Step, in *[]Step) {

	*in = append(*in, s)

}

func (b *JobBuilder) Step(name string, action Action, params ...Params) *JobBuilder {

	step := Step{
		Name:   name,
		On:     DefaultType,
		Action: action,
	}

	if len(params) > 0 {
		step.params = params[0]
	}

	b.addStep(step, &b.steps)

	return b
}

type Task interface {
	Name() string
	Action(ctx *Context) error
	Params() Params
}

func (b *JobBuilder) Task(t Task) *JobBuilder {

	step := Step{
		Name:   t.Name(),
		On:     DefaultType,
		Action: t.Action,
		params: t.Params(),
	}
	b.addStep(step, &b.steps)

	return b
}

func (b *JobBuilder) OnError(name string, action Action, params ...Params) *JobBuilder {

	step := Step{
		Name:   name,
		On:     ErrorType,
		Action: action,
	}

	if len(params) > 0 {
		step.params = params[0]
	}

	b.addStep(step, &b.onError)
	return b
}

func (b *JobBuilder) After(name string, action Action, params ...Params) *JobBuilder {

	step := Step{
		Name:   name,
		On:     AfterType,
		Action: action,
	}

	if len(params) > 0 {
		step.params = params[0]
	}

	b.addStep(step, &b.after)
	return b
}

func (b *JobBuilder) Before(name string, action Action, params ...Params) *JobBuilder {

	step := Step{
		Name:   name,
		On:     BeforeType,
		Action: action,
	}

	if len(params) > 0 {
		step.params = params[0]
	}

	b.addStep(step, &b.before)
	return b
}

func (b *JobBuilder) Always(name string, action Action, params ...Params) *JobBuilder {

	step := Step{
		Name:   name,
		On:     AlwaysType,
		Action: action,
	}

	if len(params) > 0 {
		step.params = params[0]
	}

	b.addStep(step, &b.always)
	return b
}

func (b *JobBuilder) Build() Job {

	j := &job{
		name:  b.name,
		steps: []Step{},
	}

	j.steps = append(j.steps, b.before...)
	j.steps = append(j.steps, b.steps...)
	j.steps = append(j.steps, b.after...)
	j.steps = append(j.steps, b.onError...)
	j.steps = append(j.steps, b.always...)

	return j
}
