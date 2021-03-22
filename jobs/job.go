package jobs

import (
	"context"
	"github.com/hashicorp/go-multierror"
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
	tasks   []Task
	startAt time.Time
	endAt   time.Time

	ranTasks int

	values map[string]interface{}

	status CompletionStatus

	outputs     map[string]string
	inputs      map[string]string
	subscribers []*Subscribe
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

	count := len(j.tasks)

	return Stats{
		StartAt:    j.startAt,
		EndAt:      j.endAt,
		TasksCount: count,
		TasksSkip:  count - j.ranTasks,
		TasksRun:   j.ranTasks,
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

	var groupErr error
	for _, step := range j.tasks {
		if !step.Check(ctx, groupErr) {
			continue
		}
		err := j.runTask(ctx, step)
		if err != nil {
			groupErr = multierror.Append(groupErr, err)
		}
		j.ranTasks++
	}

	return groupErr

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
