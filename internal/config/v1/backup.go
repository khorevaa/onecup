package v1

import (
	"github.com/khorevaa/onecup/internal/common"
	"github.com/khorevaa/onecup/jobs"
	"github.com/khorevaa/onecup/tasks"
)

type BackupConfig struct {
	config common.ConfigNamespace `config:",inline,required"`
}

type FileBackupConfig struct {
	Dir          string `config:"dir,required" json:"dir"`
	FileTemplate string `config:"file-template" json:"file-template"`
}

type SqlBackupConfig struct {
}

func (c *BackupConfig) Task() (task *jobs.TaskBuilder, err error) {

	switch c.config.Name() {

	case "file":

		var fileConfig FileBackupConfig
		if err := c.config.Config().Unpack(&fileConfig); err != nil {
			return nil, err
		}

		b := jobs.NewTask("Backup (file)", jobs.Inputs{}, jobs.Inputs{
			"backup-file": "backup",
		})

		b.Steps(&tasks.DumpInfobase{
			FileTemplate: fileConfig.FileTemplate,
			Dir:          fileConfig.Dir,
		})

		return b, err

	case "sql":
		panic("not implement")
	default:
		return
	}

	return
}
