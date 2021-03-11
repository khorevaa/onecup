package v1

import (
	"github.com/khorevaa/onecup/internal/common"
	"github.com/khorevaa/onecup/jobs"
	"reflect"
	"testing"
)

func TestUpdateConfig_Job(t *testing.T) {

	tests := []struct {
		name      string
		rawConfig string
		want      jobs.Job
		wantErr   bool
	}{
		{"simple",
			`
release:
  file:
    path: ./update.cfu
load-config: true
`,
			jobs.Group("update", &jobs.FileReleaseJob{
				File: "./update.cfu",
			}, &jobs.UpdateJob{
				LoadConfig: true,
			}),
			false,
		},
		{"all",
			`
release:
  file:
    path: ./update.cfu
load-config: true
server: true
dynamic: true
warnings-as-errors: true
rollback-on-error: true
`,
			jobs.Group("update", &jobs.FileReleaseJob{
				File: "./update.cfu",
			}, &jobs.UpdateJob{
				LoadConfig:       true,
				Server:           true,
				Dynamic:          true,
				WarningsAsErrors: true,
				RollbackConfig:   true,
			}),
			false,
		},
		{"error",
			`
  load-config: true
`, nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := common.NewConfigFrom(tt.rawConfig)
			if err != nil {
				t.Error(err)
			}

			c := &UpdateConfig{}
			if err != nil {
				t.Error(err)
			}

			err = cfg.Unpack(c)
			if err != nil {
				t.Error(err)
			}

			got, err := c.Job()
			if (err != nil) != tt.wantErr {
				t.Errorf("Job() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job() got = %v, want %v", got, tt.want)
			}
		})
	}
}
