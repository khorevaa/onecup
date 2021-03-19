package jobs

type TaskObject interface {
	Name() string
	Action() TaskAction
	Inputs() ValuesMap
	Outputs() ValuesMap
	//Check() CheckFunc // check setup on add task
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
	opts []TaskOption
}

func (b *taskObjectBuilder) Build() Task {
	options := TaskOptions{
		b.Inputs(),
		b.Outputs(),
		nil,
	}

	for _, opt := range b.opts {
		opt(&options)
	}

	return NewTask(b.Name(), b.Action(), WithOptions(options))

}
