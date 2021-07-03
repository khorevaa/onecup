package workflow

import (
	"context"
	"fmt"
	"github.com/khorevaa/onecup/config"
	"github.com/khorevaa/onecup/internal/common"
	uuid "github.com/satori/go.uuid"
	"html/template"
	"os"
	"strings"
	"sync"
)

type Workflow struct {
	Name   string
	Env    Values
	Params Values

	MaxParallel int

	InfobaseList []Infobase
	Tasks        []*Job

	CurrentTaskIdx int
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

			newTask := &Job{
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

func buildSteps(job *Job) {

	stepsConfig := job.jobConfig.Steps

	for _, stepConfig := range stepsConfig {

		step := Step{
			ID:        generateIDFromString(stepConfig.Name),
			Name:      stepConfig.Name,
			Condition: Condition(stepConfig.If),
			State:     Pending,
			job:       job,
			Uses:      stepConfig.Uses,
			Params:    stepConfig.With,
			// Outputs:   stepConfig.Out,
			// Cache: stepConfig.Cache,
		}

		job.Steps = append(job.Steps, step)
	}

}

func generateIDFromString(name string) string {

	val := strings.ReplaceAll(name, " ", "_")
	return val

}

func (w *Workflow) Run(ctx context.Context) error {

	limit := make(chan struct{}, w.MaxParallel)
	wg := sync.WaitGroup{}
	for _, task := range w.Tasks {
		wg.Add(1)
		go func(t *Job) {

			limit <- struct{}{}
			t.Run(ctx)
			wg.Done()
			<-limit

		}(task)
	}

	wg.Wait()

	return nil
}

func (w *Workflow) doNextTask(ctx context.Context) bool {

	if len(w.Tasks) > 0 && len(w.Tasks) < w.CurrentTaskIdx {
		return false
	}

	task := w.Tasks[w.CurrentTaskIdx]

	task.Run(ctx)

	w.CurrentTaskIdx++

	return true

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
			User:     common.TemplateValue(authConfig.User).MustExecute(ctx),
			Password: common.TemplateValue(authConfig.Password).MustExecute(ctx),
		}
	}

	return Auth{}
}

func buildParams(ctx interface{}, paramsConfig map[string]string) Values {

	params := make(Values)

	for key, value := range paramsConfig {
		params[key] = common.TemplateValue(value).MustExecute(ctx)

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
