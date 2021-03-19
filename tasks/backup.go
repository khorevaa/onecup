package tasks

import (
	"github.com/khorevaa/onecup/jobs"
	v8 "github.com/v8platform/api"
	"os"
	"path/filepath"
)

var _ jobs.TaskObject = (*DumpInfobase)(nil)

type DumpInfobase struct {
	FileTemplate string
	Dir          string
}

func (c *DumpInfobase) Action() jobs.TaskAction {
	return c.action
}

func (c *DumpInfobase) Inputs() jobs.ValuesMap {
	return nil
}

func (c *DumpInfobase) Outputs() jobs.ValuesMap {
	return map[string]string{
		"backup-file": "backup",
	}
}

func (c *DumpInfobase) Name() string {
	return "Dump infobase data"
}

func (c *DumpInfobase) action(ctx jobs.Context) error {

	ib := jobs.InfobaseFromCtx(ctx)
	opts := jobs.OptionsFromCtx(ctx)

	backupFileName := filepath.Join(c.Dir, c.FileTemplate)
	err := v8.Run(ib, v8.DumpIB(backupFileName), opts...)
	ctx.StoreValue("backup", backupFileName)

	return err
}

var _ jobs.TaskObject = (*RestoreInfobase)(nil)

type RestoreInfobase struct {
	File string
}

func (c *RestoreInfobase) Action() jobs.TaskAction {
	return c.action
}

func (c *RestoreInfobase) Inputs() jobs.ValuesMap {
	return map[string]string{
		"backup-file": "backup",
	}
}

func (c *RestoreInfobase) Outputs() jobs.ValuesMap {
	return nil
}

func (c *RestoreInfobase) Name() string {
	return "Restore infobase data"
}

func (j *RestoreInfobase) action(ctx jobs.Context) error {

	if len(j.File) == 0 {
		j.File = ctx.MustLoadValue("file").(string)
	}

	_, err := os.Stat(j.File)
	if err != nil {
		return err
	}

	ib := jobs.InfobaseFromCtx(ctx)
	opts := jobs.OptionsFromCtx(ctx)

	err = v8.Run(ib, v8.RestoreIB(j.File), opts...)

	return err

}
