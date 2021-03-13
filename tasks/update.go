package tasks

import (
	"github.com/khorevaa/onecup/jobs"
	v8 "github.com/v8platform/api"
	"github.com/v8platform/designer"
	"os"
)

var _ jobs.StepInterface = (*Update)(nil)

type Update struct {
	File             string
	LoadConfig       bool `json:"load-config"`
	Server           bool `json:"on-server"`
	Dynamic          bool `json:"dynamic"`
	WarningsAsErrors bool `json:"warnings-as-errors"`
}

func (j *Update) Handler() jobs.HandlerType {
	return jobs.DefaultType
}

func (j *Update) Action(ctx jobs.Context) error {

	if len(j.File) == 0 {
		j.File = ctx.MustLoadValue("file").(string)
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
	ib := jobs.InfobaseFromCtx(ctx)
	opts := jobs.OptionsFromCtx(ctx)

	return v8.Run(ib, command, opts...)
}

func (j *Update) Name() string {
	name := "Update configuration"
	if j.LoadConfig {
		name = "load configuration"
	}
	return name
}

var _ jobs.StepInterface = (*Update)(nil)

type RollbackUpdate struct {
	HandlerType jobs.HandlerType
}

func (j *RollbackUpdate) Handler() jobs.HandlerType {
	return j.HandlerType
}

func (j *RollbackUpdate) Action(ctx jobs.Context) error {

	ib := jobs.InfobaseFromCtx(ctx)
	opts := jobs.OptionsFromCtx(ctx)

	return v8.Run(ib, v8.RollbackCfg(), opts...)
}

func (j *RollbackUpdate) Name() string {
	name := "Rollback configuration"
	return name
}
