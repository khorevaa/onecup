package jobs

import (
	"time"
)

var _ Task = (*task)(nil)

type TaskAction func(ctx Context) error

func NewTask(name string, action TaskAction, opts ...TaskOption) Task {

	t := &task{
		name:   name,
		action: action,
		check:  NotErrorCheck,
	}

	options := &TaskOptions{}

	for _, opt := range opts {
		opt(options)
	}

	t.applyOptions(options)

	return t
}

type task struct {
	job    *job
	name   string
	action TaskAction
	check  CheckFunc

	outputs map[string]string
	inputs  map[string]string

	status  CompletionStatus
	startAt time.Time
	endAt   time.Time
}

func (t *task) applyOptions(opt *TaskOptions) {

	t.inputs = opt.Inputs
	t.outputs = opt.Outputs
	t.check = opt.Check

}

func (t *task) Name() string {
	return t.name
}

func (t *task) Stats() Stats {

	return Stats{

		StartAt: t.startAt,
		EndAt:   t.endAt,
	}
}

func (t *task) Fault() bool {
	return t.status == Error
}

func (t *task) Success() bool {
	return t.status == Success
}

func (t *task) Status() CompletionStatus {
	return t.status
}

func (t *task) Skip(ctx Context, err error) bool {

	if !t.check(ctx, err) {
		t.status = Skip
		return true
	}

	return false
}

func (t *task) finishTask(err error) {
	t.endAt = time.Now()
}

func (t *task) Run(ctx Context) (output Values, err error) {

	t.startAt = time.Now()
	defer t.finishTask(err)

	values := getValues(ctx.Values(), t.inputs)

	taskCtx := newTaskContext(ctx, t)
	taskCtx.StoreValues(values)

	select {
	case <-ctx.Done():
		return output, ctx.Err()
	default:
	}

	err = t.action(taskCtx)
	if err != nil {
		t.status = Error
		return
	}

	output = getValues(taskCtx.Output(), t.outputs)
	t.status = Success

	return
}
