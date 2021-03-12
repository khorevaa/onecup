package jobs

import (
	"sort"
	"time"
)

type Job interface {
	Name() string
	Stats() Stats
	Fault() bool
	Skiped() bool
	Success() bool
	Status() CompletionStatus
	Error() error
	Subscribe(subscribe *Subscribe)

	run(input Values) (Values, error)
	simulate(input Values) (Values, error)
}

type Stats struct {
	StartAt    time.Time
	EndAt      time.Time
	StepsCount int
	StepsSkip  int
	StepsRun   int
}

type job struct {
	name    string
	tasks   []task
	startAt time.Time
	endAt   time.Time
	skip    bool
	err     error
	status  CompletionStatus

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

	return Stats{
		StartAt:    j.startAt,
		EndAt:      j.endAt,
		StepsCount: len(j.tasks),
		StepsSkip:  0,
		StepsRun:   0,
	}
}

func (j *job) Fault() bool {
	return j.err != nil
}

func (j *job) Skiped() bool {
	panic("implement me")
}

func (j *job) Success() bool {
	return j.status == Success
}

func (j *job) Status() CompletionStatus {
	return j.status
}

func (j *job) Error() error {
	return j.err
}

func (j *job) run(input Values) (Values, error) {

	values := getValues(input, j.inputs)

	ctx := &jobContext{
		job:    j,
		values: values,
	}

	j.runTasks(ctx)

	output := getValues(ctx.values, j.outputs)

	return output, nil
}
func (j *job) runTasks(ctx *jobContext) {

	j.sortTasks()

	for _, t := range j.tasks {
		out, err := t.run(ctx, ctx.values)
		if err != nil {
			ctx.err = err
		}
		ctx.StoreValues(out)
	}

}

func getValues(values Values, keyMap map[string]string) Values {

	val := make(Values)

	for k1, k2 := range keyMap {
		val[k1] = values[k2]
	}

	return val
}

func (j *job) simulate(input Values) (Values, error) {

	values := getValues(input, j.inputs)

	ctx := &jobContext{
		job:      j,
		values:   values,
		simulate: true,
	}

	j.runTasks(ctx)

	output := getValues(ctx.values, j.outputs)

	return output, ctx.Err()

}

func (j *job) Name() string {
	return j.name
}

func (j *job) sortTasks() {

	sort.Slice(j.tasks, func(i, k int) bool {
		return j.tasks[i].handler < j.tasks[k].handler
	})

}

type Params map[string]interface{}
type Inputs map[string]string
type Values map[string]interface{}

func copyValues(dest, src Values) {
	if src == nil {
		return
	}

	for key, value := range src {
		dest[key] = value
	}
}
