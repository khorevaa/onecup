package jobs

import (
	"github.com/hashicorp/go-multierror"
	"time"
)

func NewGroup(name string, steps []Task, opts ...TaskOption) Task {

	t := &groupTask{
		task: task{
			name: name,
		},
		steps: steps,
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
	steps []Task
}

func (t *groupTask) Stats() Stats {
	return Stats{
		StartAt: t.startAt,
		EndAt:   t.endAt,
	}
}

func (t *groupTask) Run(ctx Context) (output Values, err error) {
	t.startAt = time.Now()
	defer t.finishTask(err)

	if len(t.steps) == 0 {
		output = getValues(ctx.Values(), t.outputs)
		return
	}

	values := getValues(ctx.Values(), t.inputs)
	taskCtx := newTaskContext(ctx, t)
	taskCtx.StoreValues(values)

	select {
	case <-ctx.Done():
		return output, ctx.Err()
	default:
	}

	var groupErr error
	for _, step := range t.steps {
		if !step.Check(taskCtx, groupErr) {
			continue
		}

		err := t.runTask(taskCtx, step)

		if err != nil {
			groupErr = multierror.Append(groupErr, err)
		}
	}

	if groupErr != nil {
		t.status = Error
		return nil, groupErr
	}

	output = getValues(taskCtx.Output(), t.outputs)
	t.status = Success

	return
}

func (t *groupTask) runTask(ctx Context, step Task) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	output, err := step.Run(ctx)
	if err != nil {
		return err
	}
	ctx.StoreValues(output)

	return nil
}
