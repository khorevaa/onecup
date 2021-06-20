package workflow

import (
	"context"
	"github.com/khorevaa/onecup/config"
	"html/template"
)

type jobError struct {
	step Step
	err  error
}

func (j jobError) Error() string {
	return j.err.Error()
}

type JobState int

const (
	Pending JobState = iota
	Running
	Skip
	Success
	Error
)

func (h JobState) String() string {
	switch h {
	case Pending:
		return "Pending"
	case Running:
		return "Running"
	case Success:
		return "Success"
	case Skip:
		return "Skip"
	case Error:
		return "Error"
	default:
		return "unknown state"
	}
}

type Job struct {
	ID        string
	Name      string
	Group     string
	Infobase  Infobase
	Outputs   Values
	Steps     []Step
	State     JobState
	Condition Condition
	Needs     []string
	err       jobError

	jobConfig config.JobConfig
	workflow  *Workflow
}

func (j *Job) Err() error {
	return j.err
}

func (t *Job) FuncMap() map[string]interface{} {
	return template.FuncMap{
		"env":     t.workflow.getEnv,
		"secrets": t.workflow.getEnv,
		"params":  t.workflow.getParams,
		"output":  t.getOutputs,
		"failure": t.failure,
		"always":  t.always,
	}
}

func (t *Job) Run(ctx context.Context) {

	t.State = Running

	for _, step := range t.Steps {

		err := step.Run(ctx)

		// TODO add multi err
		if err != nil {
			t.State = Error
			t.err = jobError{
				step: step,
				err:  err,
			}
		}
	}

}

func (t *Job) failure() bool {
	return t.State == Error
}

func (t *Job) always() bool {
	return true
}

func (t *Job) getOutputs() Values {
	return t.Outputs
}
