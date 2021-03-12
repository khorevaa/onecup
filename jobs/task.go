package jobs

import (
	"sort"
	"time"
)

type task struct {
	job     *job
	name    string
	steps   []TaskStep
	startAt time.Time
	endAt   time.Time
	skip    bool
	status  CompletionStatus

	handler HandlerType

	outputs map[string]string
	inputs  map[string]string
}

func (t *task) Stats() Stats {

	return Stats{
		StartAt:    t.startAt,
		EndAt:      t.endAt,
		StepsCount: len(t.steps),
		StepsSkip:  0,
		StepsRun:   0,
	}
}

func (t *task) Fault() bool {
	return t.status == Error
}

func (t *task) Skiped() bool {
	return t.status == Skip
}

func (t *task) Success() bool {
	return t.status == Success
}

func (t *task) Status() CompletionStatus {
	return t.status
}

func (t *task) eventComplete(err error) {
	t.endAt = time.Now()
	if err != nil {
		_ = t.job.EmitCompleteErr(t.name, "", err, t.endAt.Sub(t.startAt).Nanoseconds())
	} else {
		t.job.EmitComplete(t.name, "", t.status, t.endAt.Sub(t.startAt).Nanoseconds())
	}
}

func (t *task) run(ctx *jobContext, input Values) (v Values, err error) {
	t.startAt = time.Now()
	defer t.eventComplete(err)

	if t.needSkip(ctx) {
		return Values{}, nil
	}
	values := getValues(input, t.inputs)

	taskCtx := withTask(ctx, t)
	taskCtx.StoreValues(values)

	t.runSteps(taskCtx)

	output := getValues(taskCtx.outputs, t.outputs)
	err = taskCtx.Err()

	if err != nil {
		t.status = Error
		return nil, err
	}

	t.status = Success
	return output, nil
}

func (t *task) runSteps(ctx *jobContext) {
	t.sortSteps()

	for _, step := range t.steps {
		startAt := time.Now()
		err := step.run(ctx)

		if err != nil {
			ctx.err = err
		}

		if err != nil {
			_ = t.job.EmitCompleteErr(t.name, step.name, err, time.Since(startAt).Nanoseconds())
		} else {
			t.job.EmitComplete(t.name, step.name, step.status, time.Since(startAt).Nanoseconds())
		}

	}

}

func (t *task) needSkip(ctx Context) bool {

	if ctx.Fault() && !(t.handler == ErrorType || t.handler == AlwaysType) {
		t.status = Skip
		return true
	}
	return false

}

func (t *task) Name() string {
	return t.name
}

func (t *task) sortSteps() {

	sort.Slice(t.steps, func(i, k int) bool {
		return t.steps[i].handler < t.steps[k].handler
	})

}
