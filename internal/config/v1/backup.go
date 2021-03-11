package v1

import (
	"github.com/khorevaa/onecup/internal/common"
	"github.com/khorevaa/onecup/jobs"
	"github.com/khorevaa/onecup/tasks"
)

type BackupConfig struct {
	config common.ConfigNamespace `rawConfig:",inline,required"`
}

type FileBackupConfig struct {
	Dir          string `rawConfig:"dir,required" json:"dir"`
	FileTemplate string `rawConfig:"file-template" json:"file-template"`
}

type SqlBackupConfig struct {
}

func (c *BackupConfig) Unpack() (job jobs.Job, err error) {

	switch c.config.Name() {

	case "file":

		var fileConfig FileBackupConfig
		if err := c.config.Config().Unpack(&fileConfig); err != nil {
			return nil, err
		}

		b := jobs.NewJob("Backup (file)")

		b.Task(&tasks.DumpInfobase{
			FileTemplate: fileConfig.FileTemplate,
			Dir:          fileConfig.Dir,
		})

	case "sql":
		panic("not implement")
	default:
		return
	}

	return
}
