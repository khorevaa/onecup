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
	Status() JobStatus
	Error() error

	run(input Input) (Output, error)
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
	steps   []Step
	startAt time.Time
	endAt   time.Time
	skip    bool
	err     error
	status  JobStatus
}

func (j *job) Stats() Stats {

	return Stats{
		StartAt:    j.startAt,
		EndAt:      j.endAt,
		StepsCount: len(j.steps),
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
	return j.status == SuccessStatus
}

func (j *job) Status() JobStatus {
	return j.status
}

func (j *job) Error() error {
	return j.err
}

func (j *job) run(input Input) (Output, error) {

	ctx := &Context{
		job:     j,
		inputs:  input,
		outputs: make(Output),
	}

	for _, step := range j.steps {
		step.Run(ctx)
	}

	ctx.currentStep = nil

	return ctx.outputs, nil
}

func (j *job) Name() string {
	return j.name
}

func (j *job) sortSteps() {

	sort.Slice(j.steps, func(i, k int) bool {
		return j.steps[i].On > j.steps[k].On
	})

}

type Params map[string]interface{}
type Input map[string]interface{}
type Output map[string]interface{}
