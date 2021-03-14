package jobs

type TaskInterface interface {
	Name() string
	Steps() []StepInterface
	Inputs() Inputs
	Outputs() Inputs
	Handler() HandlerType
}

type JobTaskBuilder interface {
	Handler() HandlerType
	Build(job *job) Task
}

type TaskBuilder interface {
	JobTaskBuilder

	NewStep(step ...StepInterface) TaskBuilder

	Outputs(outputs Inputs) TaskBuilder
	Inputs(inputs Inputs) TaskBuilder
	Before(step ...StepBuilder) TaskBuilder
	Steps(step ...StepBuilder) TaskBuilder
	After(step ...StepBuilder) TaskBuilder
	Error(step ...StepBuilder) TaskBuilder
	Always(step ...StepBuilder) TaskBuilder
}

func NewTaskBuilder(name string, h HandlerType, inputsOutputs ...Inputs) TaskBuilder {
	var inputs, outputs Inputs

	if len(inputsOutputs) == 1 {
		inputs = inputsOutputs[0]
	}
	if len(inputsOutputs) == 2 {
		outputs = inputsOutputs[1]
	}

	return &taskBuilder{
		name:    name,
		handler: h,
		inputs:  inputs,
		outputs: outputs,
		steps:   []StepBuilder{},
	}
}

func NewTaskBuilderI(i TaskInterface) JobTaskBuilder {

	t := NewTaskBuilder(i.Name(), i.Handler(), i.Inputs(), i.Outputs())

	t.NewStep(i.Steps()...)

	return t

}

var _ TaskBuilder = (*taskBuilder)(nil)

type taskBuilder struct {
	name            string
	inputs, outputs Inputs
	steps           []StepBuilder
	handler         HandlerType
}

func (b *taskBuilder) Handler() HandlerType {
	return b.handler
}

func (b *taskBuilder) Outputs(outputs Inputs) TaskBuilder {
	panic("implement me")
}

func (b *taskBuilder) Inputs(inputs Inputs) TaskBuilder {
	panic("implement me")
}

func (b *taskBuilder) NewStep(step ...StepInterface) TaskBuilder {
	panic("implement me")
}

func (b *taskBuilder) Before(step ...StepBuilder) TaskBuilder {
	panic("implement me")
}

func (b *taskBuilder) Steps(step ...StepBuilder) TaskBuilder {
	panic("implement me")
}

func (b *taskBuilder) After(step ...StepBuilder) TaskBuilder {
	panic("implement me")
}

func (b *taskBuilder) Error(step ...StepBuilder) TaskBuilder {
	panic("implement me")
}

func (b *taskBuilder) Always(step ...StepBuilder) TaskBuilder {
	panic("implement me")
}

func (b *taskBuilder) Build(job *job) Task {

	var steps Steps

	for _, step := range b.steps {

		switch step.Handler() {
		case BeforeType:
			steps.Before = append(steps.Before, step.Build())
		case DefaultType:
			steps.Steps = append(steps.Steps, step.Build())
		case AfterType:
			steps.After = append(steps.After, step.Build())
		case ErrorType:
			steps.Error = append(steps.Error, step.Build())
		case AlwaysType:
			steps.Always = append(steps.Always, step.Build())
		}
	}

	return NewTask(job, b.name, steps, b.inputs, b.outputs)

}
