package jobs

type GroupBuilder interface {
	TaskBuilder

	NewTask(name string, action TaskAction, opts ...TaskOption) GroupBuilder
	AddTask(task TaskBuilder) GroupBuilder
	Task(task TaskObject, opts ...TaskOption) GroupBuilder
}

func Group(name string, opts ...TaskOption) GroupBuilder {

	b := &groupBuilder{
		name:    name,
		options: opts,
		tasks:   []TaskBuilder{},
	}

	return b

}

var _ GroupBuilder = (*groupBuilder)(nil)

type groupBuilder struct {
	name    string
	options []TaskOption
	tasks   []TaskBuilder
}

func (g *groupBuilder) Build() Task {

	var steps []Task

	for _, builder := range g.tasks {
		steps = append(steps, builder.Build())
	}
	return NewGroup(g.name, steps, g.options...)

}

func (b *groupBuilder) AddTask(task TaskBuilder) GroupBuilder {

	b.tasks = append(b.tasks, task)

	return b
}

func (b *groupBuilder) Task(task TaskObject, opts ...TaskOption) GroupBuilder {

	b.tasks = append(b.tasks, &taskObjectBuilder{task, opts})
	return b
}

func (b *groupBuilder) NewTask(name string, action TaskAction, opts ...TaskOption) GroupBuilder {
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
