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

	Use, err := uses.CreateUseWithParams(s.Uses, s.Params)
	if err != nil {
		s.State = Error
		return err
	}

	err = Use.Action(ctx, s.task.Infobase.Infobase, s.task.Outputs)

	if err != nil {
		s.State = Error
		return err
	}
	s.State = Success

	return nil
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
