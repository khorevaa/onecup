package workflow

import (
	"github.com/khorevaa/onecup/internal/common"
	"github.com/khorevaa/onecup/uses"
	"github.com/khorevaa/onecup/workflow/context"
)

type Step struct {
	ID        string
	Name      string
	Params    map[string]interface{}
	Outputs   map[string]OutputConfig
	Uses      string
	Condition Condition
	Cache     *CacheData
	State     JobState

	job *Job
}

func (s *Step) Run(ctx context.Context) error {

	if s.Condition.False(s.job) {
		s.State = Skip
		return nil
	}

	params, err := buildStepParams(s.job, s.Params)

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

	_, err = Use.Action(ctx)

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

			val[key], err = common.TemplateValue(typed).Execute(ctx)

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

type Condition common.TemplateValue

func (c Condition) False(ctx interface{}) bool {

	return false
}
