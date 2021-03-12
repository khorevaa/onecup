package jobs

type Task interface {
	Name() string
	Steps() []Step
	Inputs() Inputs
	Outputs() Inputs
	Handler() HandlerType
}

type Step interface {
	Name() string
	Action(ctx Context) error
	Handler() HandlerType
}

type TaskBuilder struct {
	name            string
	handler         HandlerType
	inputs, outputs Inputs

	steps []TaskStep
}

type JobBuilder struct {
	name            string
	tasks           []TaskBuilder
	inputs, outputs Inputs
}

func NewJob(name string, inputsOutputs ...Inputs) *JobBuilder {
	var inputs, outputs Inputs

	if len(inputsOutputs) == 1 {
		inputs = inputsOutputs[0]
	}
	if len(inputsOutputs) == 2 {
		outputs = inputsOutputs[1]
	}

	return &JobBuilder{
		name:    name,
		inputs:  inputs,
		outputs: outputs,
		tasks:   []TaskBuilder{},
	}
}

func NewTask(name string, inputsOutputs ...Inputs) *TaskBuilder {
	var inputs, outputs Inputs

	if len(inputsOutputs) == 1 {
		inputs = inputsOutputs[0]
	}
	if len(inputsOutputs) == 2 {
		outputs = inputsOutputs[1]
	}

	return &TaskBuilder{
		name:    name,
		inputs:  inputs,
		outputs: outputs,
		steps:   []TaskStep{},
	}
}

func (b *JobBuilder) Task(task TaskBuilder) *JobBuilder {

	b.tasks = append(b.tasks, task)

	return b
}

func (b *JobBuilder) NewTask(name string, handler HandlerType, inputs, outputs Inputs, steps ...Step) *JobBuilder {

	t := TaskBuilder{
		name:    name,
		handler: handler,
		inputs:  inputs,
		outputs: outputs,
	}

	t.Steps(steps...)

	b.tasks = append(b.tasks, t)

	return b
}

func (b *JobBuilder) NewTasks(tasks ...Task) *JobBuilder {

	for _, task := range tasks {

		b.NewTask(task.Name(), task.Handler(), task.Inputs(), task.Outputs(), task.Steps()...)
	}

	return b
}

func (b *JobBuilder) Build() Job {

	t := &job{
		name:    b.name,
		inputs:  b.inputs,
		outputs: b.outputs,
	}

	for _, builder := range b.tasks {
		t.tasks = append(t.tasks, builder.build(t))
	}

	return t
}

func (b *TaskBuilder) addStep(s TaskStep) {

	b.steps = append(b.steps, s)

}

func (b *TaskBuilder) Step(name string, action Action) *TaskBuilder {

	step := TaskStep{
		name:    name,
		handler: DefaultType,
		fn:      action,
	}

	b.addStep(step)

	return b
}

func (b *TaskBuilder) Steps(steps ...Step) *TaskBuilder {

	for _, step := range steps {

		step := TaskStep{
			name:    step.Name(),
			handler: step.Handler(),
			fn:      step.Action,
		}

		b.addStep(step)
	}

	return b
}

func (b *TaskBuilder) OnError(name string, action Action) *TaskBuilder {

	step := TaskStep{
		name:    name,
		handler: ErrorType,
		fn:      action,
	}

	b.addStep(step)
	return b
}

func (b *TaskBuilder) After(name string, action Action) *TaskBuilder {

	step := TaskStep{
		name:    name,
		handler: AfterType,
		fn:      action,
	}

	b.addStep(step)
	return b
}

func (b *TaskBuilder) Before(name string, action Action) *TaskBuilder {

	step := TaskStep{
		name:    name,
		handler: BeforeType,
		fn:      action,
	}

	b.addStep(step)
	return b
}

func (b *TaskBuilder) Always(name string, action Action) *TaskBuilder {
	step := TaskStep{
		name:    name,
		handler: AlwaysType,
		fn:      action,
	}

	b.addStep(step)
	return b
}

func (b *TaskBuilder) build(j *job) task {

	t := task{
		job:     j,
		name:    b.name,
		steps:   b.steps,
		handler: b.handler,
		inputs:  b.inputs,
		outputs: b.outputs,
	}

	return t
}
