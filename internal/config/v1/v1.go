package v1

import (
	"github.com/khorevaa/onecup/internal/common"
	"github.com/khorevaa/onecup/jobs"
)

type Config struct {
	Name       string            `config:"name,required" json:"name"`
	Uuid       string            `config:"id" json:"id"`
	Infobase   *InfobaseConfig   `config:"infobase,replace,required" json:"infobase"`
	Update     *UpdateConfig     `config:"update" json:"update"`
	Enterprise *EnterpriseConfig `config:"enterprise" json:"enterprise"`
	Extension  *ExtensionConfig  `config:"extension" json:"extension"`
	Backup     *BackupConfig     `config:"backup" json:"backup"`
	Sessions   *SessionsConfig   `config:"sessions" json:"sessions"`
}

func init() {
	common.RegisterConfigVersion("1.0", New)
}

func (v Config) Build(builder common.Builder) error {

	infobase, err := unpackInfobase(v.Infobase)
	if err != nil {
		return err
	}

	job := jobs.NewJobBuilder(v.Uuid, jobs.Inputs{
		"infobase": "infobase",
		"options":  "options",
	})

	if err := v.addBlockSessionJob(job); err != nil {
		return err
	}
	if err := v.addBackupJob(job); err != nil {
		return err
	}
	if err := v.addUpdateJob(job); err != nil {
		return err
	}
	if err := v.addExtensionJob(job); err != nil {
		return err
	}
	if err := v.addEnterpriseJob(job); err != nil {
		return err
	}
	if err := v.addRestoreJob(job); err != nil {
		return err
	}
	if err := v.addUnblockSessionJob(job); err != nil {
		return err
	}

	return builder.Build(v.Name, infobase, job.Build())
}

func (v Config) addUpdateJob(job *jobs.jobBuilder) error {

	if v.Update == nil {
		return nil
	}

	task, err := v.Update.Task()

	if err != nil {
		return err
	}

	job.Task(task)

	return nil

}

func (v Config) addBackupJob(job *jobs.jobBuilder) error {

	if v.Backup == nil {
		return nil
	}
	task, err := v.Backup.Task()
	if err != nil {
		return err
	}

	job.Task(task)
	return nil
}

func (v Config) addBlockSessionJob(job *jobs.jobBuilder) error {
	if v.Sessions == nil {
		return nil
	}
	return nil
}

func (v Config) addExtensionJob(job *jobs.jobBuilder) error {
	if v.Extension == nil {
		return nil
	}
	return nil
}

func (v Config) addEnterpriseJob(job *jobs.jobBuilder) error {
	if v.Enterprise == nil {
		return nil
	}
	return nil
}

func (v Config) addUnblockSessionJob(job *jobs.jobBuilder) error {
	if v.Sessions == nil {
		return nil
	}
	return nil
}

func (v Config) addRestoreJob(job *jobs.jobBuilder) error {
	if v.Backup == nil {
		return nil
	}
	return nil
}

func New(cfg *common.Config) (common.ConfigFactory, error) {

	config := Config{}
	if cfg != nil {
		if err := cfg.Unpack(&config); err != nil {
			return nil, err
		}
	}

	return config, nil
}
