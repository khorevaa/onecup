package jobs

type JobBuilder interface {
	Tasks(tasks ...TaskInterface) JobBuilder
	Group(name string, inputs, outputs Inputs, tasks ...TaskInterface) JobBuilder

	Before(tasks ...JobTaskBuilder) JobBuilder
	Task(tasks ...JobTaskBuilder) JobBuilder
	After(tasks ...JobTaskBuilder) JobBuilder
	Error(tasks ...JobTaskBuilder) JobBuilder
	Always(tasks ...JobTaskBuilder) JobBuilder

	Build() Job
}

var _ JobBuilder = (*jobBuilder)(nil)

type jobBuilder struct {
	name            string
	tasks           []JobTaskBuilder
	inputs, outputs Inputs
}

func (b *jobBuilder) Tasks(tasks ...TaskInterface) JobBuilder {
	for _, task := range tasks {
		b.tasks = append(b.tasks, NewTaskBuilder(task))
	}
	return b
}

func (b *jobBuilder) Group(name string, inputs, outputs Inputs, tasks ...TaskInterface) JobBuilder {

	task := NewGroupBuilder(name, DefaultType, inputs, outputs, tasks...)

	b.tasks = append(b.tasks, task)

	return b
}

func (b *jobBuilder) Before(tasks ...JobTaskBuilder) JobBuilder {
	return b.Task(tasks...)
}

func (b *jobBuilder) Task(tasks ...JobTaskBuilder) JobBuilder {
	b.tasks = append(b.tasks, tasks...)
	return b
}

func (b *jobBuilder) After(tasks ...JobTaskBuilder) JobBuilder {
	return b.Task(tasks...)
}

func (b *jobBuilder) Error(tasks ...JobTaskBuilder) JobBuilder {
	return b.Task(tasks...)
}

func (b *jobBuilder) Always(tasks ...JobTaskBuilder) JobBuilder {
	return b.Task(tasks...)
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
		tasks:   []JobTaskBuilder{},
	}
}

func (b *jobBuilder) Build() Job {

	j := &job{
		name:    b.name,
		inputs:  b.inputs,
		outputs: b.outputs,
	}

	steps := Steps{}

	for _, builder := range b.tasks {
		switch builder.Handler() {
		case BeforeType:
			steps.Before = append(steps.Before, builder.Build(j))
		case DefaultType:
			steps.Steps = append(steps.Steps, builder.Build(j))
		case AfterType:
			steps.After = append(steps.After, builder.Build(j))
		case ErrorType:
			steps.Error = append(steps.Error, builder.Build(j))
		case AlwaysType:
			steps.Always = append(steps.Always, builder.Build(j))
		}
	}

	return j
}
