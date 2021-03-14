package jobs

func NewGroup(name string, tasks []Task, opts ...TaskOption) Task {

	t := &groupTask{
		task: task{
			name:  name,
			check: NotErrorCheck,
		},
		tasks: tasks,
	}

	options := &TaskOptions{}

	for _, opt := range opts {
		opt(options)
	}

	t.applyOptions(options)

	return t

}

var _ Task = (*groupTask)(nil)

type groupTask struct {
	task
	tasks []Task
}

func (g groupTask) Stats() Stats {
	panic("implement me")
}

func (g groupTask) Run(ctx Context) (output Values, err error) {

}
