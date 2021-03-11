package tasks

import (
	"github.com/khorevaa/onecup/jobs"
	v8 "github.com/v8platform/api"
	"os"
	"path/filepath"
)

type DumpInfobase struct {
	FileTemplate string
	Dir          string
}

func (c *DumpInfobase) Name() string {
	return "Dump infobase data"
}

func (c *DumpInfobase) Params() jobs.Params {

	return map[string]interface{}{
		"file-template": c.FileTemplate,
		"dir":           c.Dir,
	}
}

func (c *DumpInfobase) Action(ctx *jobs.Context) error {

	backupFileName := filepath.Join(c.Dir, c.FileTemplate)
	err := v8.Run(ctx.Infobase(), v8.DumpIB(backupFileName), ctx.Options()...)
	ctx.Out("backup-file", backupFileName)

	return err
}

type RestoreInfobase struct {
	File string
}

func (j *RestoreInfobase) Action(ctx *jobs.Context) error {

	if len(j.File) == 0 {
		if val, ok := ctx.Value("backup-file"); ok {
			j.File = val.(string)
		}
	}

	_, err := os.Stat(j.File)
	if err != nil {
		return err
	}

	err = v8.Run(ctx.Infobase(), v8.RestoreIB(j.File), ctx.Options()...)

	return err

}
