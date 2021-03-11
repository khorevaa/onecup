package jobs

import (
	"log"
	"reflect"
	"testing"
)

func TestRunner_Run(t *testing.T) {
	type fields struct {
		params Params
		jobs   []Job
	}

	action := func(ctx *Context) error {
		log.Println(ctx.Step().Name)
		return nil
	}

	tests := []struct {
		name    string
		fields  fields
		want    Stats
		wantErr bool
	}{
		{
			"one job",
			fields{
				params: map[string]interface{}{},
				jobs: []Job{
					NewJob("test").
						Step("step1", action).
						Build(),
				},
			},
			Stats{
				StepsCount: 1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Runner{
				params: tt.fields.params,
				jobs:   tt.fields.jobs,
			}
			if err := r.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := r.jobs[0].Stats()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job() got = %v, want %v", got, tt.want)
			}
		})
	}
}
