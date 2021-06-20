package workflow

import (
	"context"
	"github.com/khorevaa/onecup/internal/common"
	"github.com/khorevaa/onecup/uses"
	context2 "github.com/khorevaa/onecup/workflow/context"
	"testing"
)

func TestStep_Run(t *testing.T) {
	type fields struct {
		ID        string
		Name      string
		Params    map[string]interface{}
		Outputs   map[string]OutputConfig
		Uses      string
		Condition Condition
		Cache     *CacheData
		State     JobState
		job       *Job
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Step{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				Params:    tt.fields.Params,
				Outputs:   tt.fields.Outputs,
				Uses:      tt.fields.Uses,
				Condition: tt.fields.Condition,
				Cache:     tt.fields.Cache,
				State:     tt.fields.State,
				job:       tt.fields.job,
			}
			if err := s.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func init() {
	uses.RegisterUseType("testUses", NewTestUses)
}

func NewTestUses(cfg *common.Config) (context2.Use, error) {

}

type testUses struct {
}
