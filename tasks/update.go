package tasks

import (
	"github.com/khorevaa/onecup/jobs"
	v8 "github.com/v8platform/api"
	"github.com/v8platform/designer"
)

var _ jobs.TaskObject = (*Update)(nil)

type Update struct {
	File             string
	LoadConfig       bool `json:"load-config"`
	Server           bool `json:"on-server"`
	Dynamic          bool `json:"dynamic"`
	WarningsAsErrors bool `json:"warnings-as-errors"`
}

func (j *Update) Action() jobs.TaskAction {
	return j.action
}

func (j *Update) Inputs() jobs.ValuesMap {
	return map[string]string{
		"file": "file",
	}
}

func (j *Update) Outputs() jobs.ValuesMap {
	return nil
}

func (j *Update) action(ctx jobs.Context) error {

	if len(j.File) == 0 {
		j.File = ctx.MustLoadValue("file").(string)
	}
	//_, err := os.Stat(j.File)
	//if err != nil {
	//	return err
	//}

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

	process, err := v8.Background(ctx, ib, command, opts...)
	if err != nil {
		return err
	}
	<-process.Ready()
	return <-process.Wait()
}

func (j *Update) Name() string {
	name := "Update configuration"
	if j.LoadConfig {
		name = "load configuration"
	}
	return name
}

var _ jobs.TaskObject = (*RollbackUpdate)(nil)

type RollbackUpdate struct{}

func (j *RollbackUpdate) Action() jobs.TaskAction {
	return j.action
}

func (j *RollbackUpdate) Inputs() jobs.ValuesMap {
	return nil
}

func (j *RollbackUpdate) Outputs() jobs.ValuesMap {
	return nil
}

func (j *RollbackUpdate) action(ctx jobs.Context) error {

	ib := jobs.InfobaseFromCtx(ctx)
	opts := jobs.OptionsFromCtx(ctx)

	process, err := v8.Background(ctx, ib, v8.RollbackCfg(), opts...)
	if err != nil {
		return err
	}
	return <-process.Wait()
}

func (j *RollbackUpdate) Name() string {
	name := "Rollback configuration"
	return name
}
