package jobs

import (
	"fmt"
	"reflect"
	"testing"
)

type testStep struct {
	name    string
	handler HandlerType
}

func (t testStep) Name() string {
	return t.name
}

func (t testStep) Action(ctx Context) error {

	fmt.Printf("Do <%s> handler <%s> sort <%d>\n", t.name, t.handler, t.handler)

	v8input, _ := ctx.LoadValue("v8-task")
	fmt.Printf("input v8-task <%s>\n", v8input)
	fmt.Printf("Must input v8-opts <%s>\n", ctx.MustLoadValue("v8-opts"))

	return nil
}

func (t testStep) Handler() HandlerType {
	return t.handler
}

func TestRunner_Run(t *testing.T) {
	type fields struct {
		params Values
		job    Job
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
				params: Values{
					"v8-platform": "./path-to-v8",
				},
				job: NewJob("test", Inputs{
					"v8": "v8-platform",
				}).
					NewTask("task", DefaultType, Inputs{
						"v8-task": "v8",
					}, nil,
						testStep{
							name:    "step default",
							handler: DefaultType,
						},
						testStep{
							name:    "step before",
							handler: BeforeType,
						},
						testStep{
							name:    "step after",
							handler: AfterType,
						},
						testStep{
							name:    "step default 2",
							handler: DefaultType,
						},
						testStep{
							name:    "step always",
							handler: AlwaysType,
						}).
					Build(),
			},
			Stats{
				StepsCount: 1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SimulateRunner{
				params: tt.fields.params,
				jobs:   []Job{tt.fields.job},
			}
			if err := r.Run(r.params); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := r.jobs[0].Stats()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job() got = %v, want %v", got, tt.want)
			}
		})
	}
}
