package jobs

import (
	"time"
)

var _ Task = (*task)(nil)

func NewTask(job2 *job, name string, steps Steps, inputsOutputs ...Inputs) Task {

	var inputs, outputs Inputs

	if len(inputsOutputs) == 1 {
		inputs = inputsOutputs[0]
	}
	if len(inputsOutputs) == 2 {
		outputs = inputsOutputs[1]
	}

	return &task{
		job:     job2,
		name:    name,
		inputs:  inputs,
		outputs: outputs,
		steps:   steps,
	}

}

type task struct {
	job     *job
	name    string
	steps   Steps
	startAt time.Time
	endAt   time.Time

	ranSteps int

	status CompletionStatus

	outputs map[string]string
	inputs  map[string]string
}

func (t *task) Stats() Stats {

	count := t.steps.Len()

	return Stats{

		StartAt:    t.startAt,
		EndAt:      t.endAt,
		StepsCount: count,
		StepsSkip:  count - t.ranSteps,
		StepsRun:   t.ranSteps,
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

func (t *task) eventComplete(err error) {
	t.endAt = time.Now()
	if err != nil {
		_ = t.job.EmitCompleteErr(t.name, "", err, t.endAt.Sub(t.startAt).Nanoseconds())
	} else {
		t.job.EmitComplete(t.name, "", t.status, t.endAt.Sub(t.startAt).Nanoseconds())
	}
}

func (t *task) Run(ctx Context) (output Values, err error) {

	t.startAt = time.Now()
	defer t.eventComplete(err)

	values := getValues(ctx.Values(), t.inputs)

	taskCtx := newTaskContext(ctx, t)
	taskCtx.StoreValues(values)

	err = t.runSteps(taskCtx)
	if err != nil {
		t.status = Error
		return
	}

	output = getValues(taskCtx.Output(), t.outputs)
	t.status = Success

	return
}

func (t *task) runSteps(ctx Context) error {

	runStep := func(step Task) error {
		return t.runStep(ctx, step)
	}

	n, err := t.steps.Do(runStep)
	t.ranSteps = n

	return err

}

func (t *task) runStep(ctx Context, step Task) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	startAt := time.Now()

	output, err := step.Run(ctx)

	if err != nil {
		return t.job.EmitCompleteErr(t.name, step.Name(), err, time.Since(startAt).Nanoseconds())
	}

	t.job.EmitComplete(t.name, step.Name(), step.Status(), time.Since(startAt).Nanoseconds())
	ctx.StoreValues(output)

	return nil
}

func (t *task) Name() string {
	return t.name
}
