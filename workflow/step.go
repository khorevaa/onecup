package workflow

import (
	"context"
	"github.com/khorevaa/onecup/config"
	"github.com/khorevaa/onecup/uses"
)

type Step struct {
	ID        string
	Name      string
	Params    map[string]interface{}
	Outputs   map[string]OutputConfig
	Uses      string
	Condition Condition
	Cache     *CacheData
	State     TaskState

	task *Task
}

func (s *Step) Run(ctx context.Context) error {

	if s.Condition.False(s.task) {
		s.State = Skip
		return nil
	}

	params, err := buildStepParams(s.task, s.Params)

	if err != nil {
		s.State = Error
		return err
	}

	Use, err := uses.CreateUseWithParams(s.Uses, params)
	if err != nil {
		s.State = Error
		return err
	}

	s.State = Running

	err = Use.Action(ctx, s.task.Infobase.Infobase, s.task.Outputs)

	if err != nil {
		s.State = Error
		return err
	}
	s.State = Success

	return nil
}

func buildStepParams(ctx interface{}, params map[string]interface{}) (Values, error) {

	val := make(Values)
	var err error
	for key, value := range params {

		switch typed := value.(type) {
		case string:

			val[key], err = config.TemplateValue(typed).Execute(ctx)

			if err != nil {
				return nil, err
			}

		default:
			val[key] = value
		}

	}

	return val, nil
}

type OutputConfig struct {
	Type string `json:"type" yaml:"type"`
}

type CacheData struct {
	Path string
}

type Condition config.TemplateValue

func (c Condition) False(ctx interface{}) bool {

	return false
}
