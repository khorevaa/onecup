package jobs

type JobBuilder interface {
	NewTask(name string, action TaskAction, opts ...TaskOption) JobBuilder
	AddTask(task TaskBuilder) JobBuilder
	Task(task TaskObject, opts ...TaskOption) JobBuilder

	Build() Job
}

var _ JobBuilder = (*jobBuilder)(nil)

type jobBuilder struct {
	name            string
	tasks           []TaskBuilder
	inputs, outputs ValuesMap
}

func (b *jobBuilder) AddTask(task TaskBuilder) JobBuilder {

	b.tasks = append(b.tasks, task)

	return b
}
func (b *jobBuilder) Task(task TaskObject, opts ...TaskOption) JobBuilder {

	b.tasks = append(b.tasks, &taskObjectBuilder{task, opts})
	return b
}

func (b *jobBuilder) NewTask(name string, action TaskAction, opts ...TaskOption) JobBuilder {
	t := &taskBuilder{
		name:   name,
		action: action,
	}
	for _, opt := range opts {
		opt(&t.options)
	}

	b.tasks = append(b.tasks, t)

	return b
}

func NewJobBuilder(name string, inputsOutputs ...ValuesMap) JobBuilder {
	var inputs, outputs ValuesMap

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
