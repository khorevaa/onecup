package jobs

type GroupBuilder interface {
	JobTaskBuilder

	NewTask(steps ...TaskInterface) GroupBuilder
	Task(step ...JobTaskBuilder) GroupBuilder
}

func newGroupBuilder(name string, h HandlerType, inputsOutputs ...Inputs) GroupBuilder {
	var inputs, outputs Inputs

	if len(inputsOutputs) == 1 {
		inputs = inputsOutputs[0]
	}
	if len(inputsOutputs) == 2 {
		outputs = inputsOutputs[1]
	}

	return &groupBuilder{
		name:    name,
		handler: h,
		inputs:  inputs,
		outputs: outputs,
		tasks:   []JobTaskBuilder{},
	}
}

func NewGroupBuilder(name string, h HandlerType, inputs, output Inputs, tasks ...TaskInterface) GroupBuilder {

	t := newGroupBuilder(name, h, inputs, output)

	t.NewTask(tasks...)

	return t

}

var _ GroupBuilder = (*groupBuilder)(nil)

type groupBuilder struct {
	name            string
	inputs, outputs Inputs
	tasks           []JobTaskBuilder
	handler         HandlerType
}

func (g groupBuilder) Handler() HandlerType {
	return g.handler
}

func (g groupBuilder) Build(job *job) Task {
	panic("implement me")
}

func (g groupBuilder) NewTask(steps ...TaskInterface) GroupBuilder {
	panic("implement me")
}

func (g groupBuilder) Task(step ...JobTaskBuilder) GroupBuilder {
	panic("implement me")
}
