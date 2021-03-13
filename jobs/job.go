package jobs

import (
	"context"
	"time"
)

type Job interface {
	Name() string
	Stats() Stats
	Fault() bool
	Success() bool
	Status() CompletionStatus
	Subscribe(subscribe *Subscribe)

	Run(ctx context.Context, input Values) (Values, error)
}

type Stats struct {
	StartAt    time.Time
	EndAt      time.Time
	TasksCount int
	TasksSkip  int
	TasksRun   int
	StepsCount int
	StepsSkip  int
	StepsRun   int
}

type job struct {
	name    string
	tasks   Steps
	startAt time.Time
	endAt   time.Time

	ranTasks int

	status CompletionStatus

	outputs     map[string]string
	inputs      map[string]string
	subscribers []*Subscribe
}

type Steps struct {
	Before []Task
	Steps  []Task
	After  []Task
	Error  []Task
	Always []Task

	err error
}

type RunStep func(step Task) error

func (s Steps) Len() int {

	return len(s.Before) +
		len(s.Steps) +
		len(s.After) +
		len(s.Error) +
		len(s.Always)
}

func (s Steps) Do(fn RunStep) (int, error) {
	var total, n int
	n, s.err = s.doSteps(fn, s.Before)
	total += n
	n, s.err = s.doSteps(fn, s.Steps)
	total += n
	n, s.err = s.doSteps(fn, s.After)
	total += n
	n, s.err = s.doError(fn, s.Error)
	total += n
	n, s.err = EachStep(s.Always).Do(fn)
	total += n

	return total, s.err

}

func (s Steps) doSteps(fn RunStep, steps []Task) (int, error) {

	if s.err != nil {
		return 0, s.err
	}

	return EachStep(steps).Do(fn)

}

func (s Steps) doError(fn RunStep, steps []Task) (int, error) {

	if s.err == nil {
		return 0, s.err
	}

	return EachStep(steps).Do(fn)

}

type EachStep []Task

func (list EachStep) Do(fn func(step Task) error) (int, error) {
	var n int
	for _, s := range list {
		n++
		err := fn(s)
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (j *job) EmitEvent(task string, step string, event string) {
	j.emmitEvent(Event{
		Type:  AllEvents,
		Job:   j.name,
		Task:  task,
		Step:  step,
		Event: event,
	})
}

func (j *job) EmitEventErr(task string, step string, event string, err error) error {
	j.emmitEvent(Event{
		Type:  ErrorEvents,
		Job:   j.name,
		Task:  task,
		Step:  step,
		Event: event,
		Err:   err,
	})

	return err
}

func (j *job) EmitTiming(task string, step string, event string, nanoseconds int64) {
	j.emmitEvent(Event{
		Type:   TimingEvents,
		Job:    j.name,
		Task:   task,
		Step:   step,
		Event:  event,
		Timing: nanoseconds,
	})
}

func (j *job) EmitComplete(task string, step string, status CompletionStatus, nanoseconds int64) {
	j.emmitEvent(Event{
		Type:   CompletionEvents,
		Job:    j.name,
		Task:   task,
		Step:   step,
		Status: status,
		Timing: nanoseconds,
	})
}
func (j *job) EmitCompleteErr(task string, step string, err error, nanoseconds int64) error {
	j.emmitEvent(Event{
		Type:   CompletionEvents,
		Job:    j.name,
		Task:   task,
		Step:   step,
		Err:    err,
		Status: Error,
		Timing: nanoseconds,
	})

	return err
}

func (j *job) emmitEvent(event Event) {

	for _, subscriber := range j.subscribers {
		subscriber.emmitEvent(event)
	}
}

func (j *job) Subscribe(subscribe *Subscribe) {

	j.subscribers = append(j.subscribers, subscribe)

}

func (j *job) Stats() Stats {

	count := j.tasks.Len()

	var stepsCount, stepsRun, stepsSkip int

	_, _ = j.tasks.Do(func(step Task) error {
		stat := step.Stats()
		stepsCount += stat.StepsCount
		stepsRun += stat.StepsRun
		stepsSkip += stat.StepsSkip
		return nil
	})

	return Stats{
		StartAt:    j.startAt,
		EndAt:      j.endAt,
		TasksCount: count,
		TasksSkip:  count - j.ranTasks,
		TasksRun:   j.ranTasks,
		StepsCount: stepsCount,
		StepsSkip:  stepsSkip,
		StepsRun:   stepsRun,
	}
}

func (j *job) Fault() bool {
	return j.status == Error
}

func (j *job) Success() bool {
	return j.status == Success
}

func (j *job) Status() CompletionStatus {
	return j.status
}

func (j *job) Run(ctx context.Context, input Values) (Values, error) {

	values := getValues(input, j.inputs)

	jobCtx := &jobContext{
		Context: ctx,
		job:     j,
		values:  values,
	}

	err := j.runTasks(jobCtx)
	if err != nil {
		j.status = Error
		return nil, err
	}
	output := getValues(jobCtx.values, j.outputs)

	return output, nil
}
func (j *job) runTasks(ctx Context) error {

	doTask := func(step Task) error {
		return j.runTask(ctx, step)
	}

	n, err := j.tasks.Do(doTask)

	j.ranTasks = n

	return err

}

func (j *job) runTask(ctx Context, step Task) error {

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

func (j *job) Name() string {
	return j.name
}
