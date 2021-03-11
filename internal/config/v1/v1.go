package v1

import (
	"github.com/khorevaa/onecup/internal/common"
	"github.com/khorevaa/onecup/jobs"
)

type Config struct {
	Name       string            `rawConfig:"name,required" json:"name"`
	Uuid       string            `rawConfig:"id" json:"id"`
	Infobase   *InfobaseConfig   `rawConfig:"infobase,replace,required" json:"infobase"`
	Update     *UpdateConfig     `rawConfig:"update" json:"update"`
	Enterprise *EnterpriseConfig `rawConfig:"enterprise" json:"enterprise"`
	Extension  *ExtensionConfig  `rawConfig:"extension" json:"extension"`
	Backup     *BackupConfig     `rawConfig:"backup" json:"backup"`
	Sessions   *SessionsConfig   `rawConfig:"sessions" json:"sessions"`
}

func (v Config) Build(builder common.Builder) error {

	infobase, err := unpackInfobase(v.Infobase)
	if err != nil {
		return err
	}

	jobList := &jobs.List{}

	if err := v.addBlockSessionJob(jobList); err != nil {
		return err
	}
	if err := v.addBackupJob(jobList); err != nil {
		return err
	}
	if err := v.addUpdateJob(jobList); err != nil {
		return err
	}
	if err := v.addExtensionJob(jobList); err != nil {
		return err
	}
	if err := v.addEnterpriseJob(jobList); err != nil {
		return err
	}
	if err := v.addRestoreJob(jobList); err != nil {
		return err
	}
	if err := v.addUnblockSessionJob(jobList); err != nil {
		return err
	}

	return builder.Build(v.Name, infobase, *jobList)
}

func (v Config) addUpdateJob(list *jobs.List) error {

	if v.Update == nil {
		return nil
	}

	job, err := v.Update.Job()
	if err != nil {
		return err
	}

	list.Add(job)

	return nil

}

func (v Config) addBackupJob(list *jobs.List) error {

	if v.Backup == nil {
		return nil
	}

	job, err := v.Backup.Unpack()
	if err != nil {
		return err
	}

	list.Add(job)

	return nil
}

func (v Config) addBlockSessionJob(list *jobs.List) error {
	if v.Sessions == nil {
		return nil
	}
	return nil
}

func (v Config) addExtensionJob(list *jobs.List) error {
	if v.Extension == nil {
		return nil
	}
	return nil
}

func (v Config) addEnterpriseJob(list *jobs.List) error {
	if v.Enterprise == nil {
		return nil
	}
	return nil
}

func (v Config) addUnblockSessionJob(list *jobs.List) error {
	if v.Sessions == nil {
		return nil
	}
	return nil
}

func (v Config) addRestoreJob(list *jobs.List) error {
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