package jobs

type JobBuilder interface {
	NewTask(name string, action TaskAction, opts ...TaskOption) JobBuilder
	Tasks(tasks ...TaskObject) JobBuilder
	Task(tasks ...TaskBuilder) JobBuilder
	Group(task TaskBuilder) JobBuilder

	Build() Job
}

var _ JobBuilder = (*jobBuilder)(nil)

type jobBuilder struct {
	name            string
	tasks           []TaskBuilder
	inputs, outputs Inputs
}

func (b *jobBuilder) Group(task TaskBuilder) JobBuilder {
	return b.Task(task)
}

func (b *jobBuilder) NewTask(name string, action TaskAction, opts ...TaskOption) JobBuilder {
	t := taskBuilder{
		name:   name,
		action: action,
	}
	for _, opt := range opts {
		opt(&t.options)
	}

	b.tasks = append(b.tasks, t)

	return b
}

func (b *jobBuilder) Task(tasks ...TaskBuilder) JobBuilder {
	b.tasks = append(b.tasks, tasks...)
}

func (b *jobBuilder) Tasks(tasks ...TaskObject) JobBuilder {
	for _, task := range tasks {
		b.tasks = append(b.tasks, &taskObjectBuilder{task})
	}
	return b
}

func NewJobBuilder(name string, inputsOutputs ...Inputs) JobBuilder {
	var inputs, outputs Inputs

	if len(inputsOutputs) == 1 {
		inputs = inputsOutputs[0]
	}
	if len(inputsOutputs) == 2 {
		outputs = inputsOutputs[1]
	}

	return &jobBuilder{
		name:    name,
		inputs:  inputs,
		outputs: outputs,
		tasks:   []TaskBuilder{},
	}
}

func (b *jobBuilder) Build() Job {

	j := &job{
		name:    b.name,
		inputs:  b.inputs,
		outputs: b.outputs,
	}

	for _, builder := range b.tasks {
		j.tasks = append(j.tasks, builder.Build())
	}

	return j
}
