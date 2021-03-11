package tasks

import (
	"github.com/khorevaa/onecup/jobs"
	v8 "github.com/v8platform/api"
	"github.com/v8platform/designer"
	"os"
)

type Update struct {
	File             string
	LoadConfig       bool `config:"load-config" json:"load-config"`
	Server           bool `config:"on-server" json:"on-server"`
	Dynamic          bool `config:"dynamic" json:"dynamic"`
	WarningsAsErrors bool `config:"warnings-as-errors" json:"warnings-as-errors"`
}

func (j *Update) Action(ctx *jobs.Context) error {

	if len(j.File) == 0 {
		if val, ok := ctx.Value("release-file"); ok {
			j.File = val.(string)
		}
	}

	_, err := os.Stat(j.File)
	if err != nil {
		return err
	}

	var command v8.Command

	updateDBConfig := designer.UpdateDBCfgOptions{
		Server:           j.Server,
		Dynamic:          j.Dynamic,
		WarningsAsErrors: j.WarningsAsErrors,
	}

	if j.LoadConfig {
		command = v8.LoadCfg(j.File, updateDBConfig)
	} else {
		command = v8.UpdateCfg(j.File, false, updateDBConfig)
	}

	return v8.Run(ctx.Infobase(), command, ctx.Options()...)
}

func (j *Update) Params() jobs.Params {
	return map[string]interface{}{}
}

func (j *Update) Name() string {
	name := "Update configuration"
	if j.LoadConfig {
		name = "load configuration"
	}

	return name
}
