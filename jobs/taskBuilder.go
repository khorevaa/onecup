package jobs

type TaskObject interface {
	Name() string
	Action() TaskAction
	Inputs() Inputs
	Outputs() Inputs
	Check() CheckFunc
}

type TaskBuilder interface {
	Build() Task
}

var _ TaskBuilder = (*taskBuilder)(nil)

type taskBuilder struct {
	name    string
	action  TaskAction
	options TaskOptions
}

func (b *taskBuilder) Build() Task {

	return NewTask(b.name, b.action, WithOptions(b.options))

}

var _ TaskBuilder = (*taskObjectBuilder)(nil)

type taskObjectBuilder struct {
	TaskObject
}

func (b *taskObjectBuilder) Build() Task {

	return NewTask(b.Name(), b.Action(),
		WithCheck(b.Check()),
		WithInput(b.Inputs()),
		WithOutput(b.Outputs()))

}
