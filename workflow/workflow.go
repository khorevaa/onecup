package workflow

import (
	"context"
	"fmt"
	"github.com/khorevaa/onecup/config"
	uuid "github.com/satori/go.uuid"
	"html/template"
	"os"
	"strings"
)

type Task struct {
	ID        string
	Name      string
	Group     string
	Infobase  Infobase
	Outputs   Values
	Steps     []Step
	State     TaskState
	Condition Condition
	Needs     []string

	jobConfig config.JobConfig
	workflow  *Workflow
}

func (t *Task) FuncMap() map[string]interface{} {
	return template.FuncMap{
		"env":     t.workflow.getEnv,
		"secrets": t.workflow.getEnv,
		"params":  t.workflow.getParams,
		"output":  t.getOutputs,
	}
}

func (t *Task) getOutputs() Values {
	return t.Outputs
}

type TaskState int

const (
	Pending TaskState = iota
	Running
	Skip
	Success
	Error
)

func (h TaskState) String() string {
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
		return "unknown"
	}
}

type Workflow struct {
	Name   string
	Env    Values
	Params Values

	MaxParallel int

	InfobaseList []Infobase
	Tasks        []*Task
}

func NewWorkflow(cfg config.Config) (*Workflow, error) {

	workflow := &Workflow{
		Name:        cfg.Name,
		MaxParallel: cfg.Strategy.MaxParallel,
	}

	workflow.Env = buildEnv(cfg.Env)
	workflow.Params = buildParams(workflow, cfg.Params)
	workflow.InfobaseList = buildInfobaseList(workflow, cfg.InfobaseList)
	workflow.buildTasks(cfg.Jobs)

	return workflow, nil
}

func (w *Workflow) buildTasks(jobsConfig map[string]config.JobConfig) {

	for _, infobase := range w.InfobaseList {

		group := uuid.NewV4().String()

		for jobKey, jobConfig := range jobsConfig {

			newTask := &Task{
				ID:        jobKey,
				Name:      fmt.Sprintf("%s (%s)", jobKey, infobase.Name),
				Group:     group,
				Outputs:   make(Values),
				Condition: Condition(jobConfig.If),
				State:     Pending,
				Infobase:  infobase,

				workflow:  w,
				jobConfig: jobConfig,
			}

			buildSteps(newTask)

			w.Tasks = append(w.Tasks, newTask)

		}
	}

}

func buildSteps(task *Task) {

	stepsConfig := task.jobConfig.Steps

	for _, stepConfig := range stepsConfig {

		step := Step{
			ID:        generateIDFromString(stepConfig.Name),
			Name:      stepConfig.Name,
			Condition: Condition(stepConfig.If),
			State:     Pending,
			task:      task,
			Uses:      stepConfig.Uses,
			Params:    stepConfig.With,
			// Outputs:   stepConfig.Out,
			// Cache: stepConfig.Cache,
		}

		task.Steps = append(task.Steps, step)
	}

}

func generateIDFromString(name string) string {

	val := strings.ReplaceAll(name, " ", "_")
	return val

}

func (w *Workflow) Run(ctx context.Context) error {

	return nil
}

func (w *Workflow) FuncMap() map[string]interface{} {
	return template.FuncMap{
		"env":     w.getEnv,
		"secrets": w.getEnv,
		"params":  w.getParams,
	}
}

func (w *Workflow) getEnv() Values {
	return w.Env
}

func (w *Workflow) getParams() Values {
	return w.Params
}

type Auth struct {
	User     string
	Password string
}

func buildInfobaseList(ctx interface{}, listConfig config.InfobaseListConfig) []Infobase {

	var list []Infobase
	var globalAuth Auth

	globalAuth = buildAuth(ctx, listConfig.Auth)

	if len(listConfig.Items) > 0 {

		for _, item := range listConfig.Items {
			ib, err := NewInfobase(ctx, globalAuth, item)
			if err != nil {
				panic(err)
			}
			list = append(list, ib)
		}

	} else {
		ib, err := NewInfobase(ctx, globalAuth, listConfig.Infobase)
		if err != nil {
			panic(err)
		}
		list = append(list, ib)

	}

	return list
}

func buildAuth(ctx interface{}, authConfig config.AuthConfig) Auth {

	if len(authConfig.User) > 0 {

		return Auth{
			User:     authConfig.User.MustExecute(ctx),
			Password: authConfig.Password.MustExecute(ctx),
		}
	}

	return Auth{}
}

func buildParams(ctx interface{}, paramsConfig map[string]config.TemplateValue) Values {

	params := make(Values)

	for key, value := range paramsConfig {
		params[key] = value.MustExecute(ctx)

	}
	return params
}

func buildEnv(env map[string]string) Values {

	envs := make(Values)

	for key, value := range env {
		envs[key] = os.Getenv(value)
	}

	return envs
}
