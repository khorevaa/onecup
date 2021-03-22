package v1

import (
	"github.com/khorevaa/onecup/internal/common"
	"github.com/khorevaa/onecup/jobs"
	uuid "github.com/satori/go.uuid"
)

type JobConfig struct {
	Name       string            `config:"name,required" json:"name"`
	Uuid       string            `config:"id" json:"id"`
	Infobase   *InfobaseConfig   `config:"infobase,replace,required" json:"infobase"`
	Update     *UpdateConfig     `config:"update" json:"update"`
	Enterprise *EnterpriseConfig `config:"enterprise" json:"enterprise"`
	Extension  *ExtensionConfig  `config:"extension" json:"extension"`
	Backup     *BackupConfig     `config:"backup" json:"backup"`
	Sessions   *SessionsConfig   `config:"sessions" json:"sessions"`
}

type ContextConfig struct {
	Strategy common.ConfigNamespace `config:"strategy" json:"strategy"`
}

type MatrixConfig struct {
	Infobases map[string]*common.Config `config:"infobases" json:"infobases"`
}

type Config struct {
	Context   *ContextConfig `config:"context" json:"context"`
	JobConfig `config:",inline" json:",inline"`
	Jobs      map[string]*common.Config `config:"jobs,replace" json:"jobs"`
}

func init() {
	common.RegisterConfigVersion("1.0", New)
}

func (v Config) Build(builder common.Builder) error {

	if len(v.Jobs) == 0 {

		return v.JobConfig.Build(builder)

	}

	job := jobs.NewJobBuilder(v.Uuid, nil, jobs.ValuesMap{
		"infobase": "infobase",
		"options":  "options",
	})

	return builder.Build(v.Name, nil, jobs.List{
		job.Build(),
	})
}

func newJobBuilder(name string, cfg *common.Config) (jobs.JobBuilder, error) {

	var jobConfig JobConfig

	if err := cfg.Unpack(&jobConfig); err != nil {
		return nil, err
	}

	jobConfig.Name = name

	if len(jobConfig.Uuid) == 0 {
		jobConfig.Uuid = uuid.NewV1().String()
	}

	return jobConfig.JobBuilder()

}

func (v JobConfig) JobBuilder() (jobs.JobBuilder, error) {

	infobase, err := unpackInfobase(v.Infobase)
	if err != nil {
		return nil, err
	}

	job := jobs.NewJobBuilder(v.Uuid, jobs.Values{
		"infobase": infobase,
		"options":  nil,
	}, jobs.ValuesMap{
		"infobase": "infobase",
		"options":  "options",
	})

	if err := v.addBlockSessionJob(job); err != nil {
		return nil, err
	}
	if err := v.addBackupJob(job); err != nil {
		return nil, err
	}
	if err := v.addUpdateJob(job); err != nil {
		return nil, err
	}
	if err := v.addExtensionJob(job); err != nil {
		return nil, err
	}
	if err := v.addEnterpriseJob(job); err != nil {
		return nil, err
	}
	if err := v.addRestoreJob(job); err != nil {
		return nil, err
	}
	if err := v.addUnblockSessionJob(job); err != nil {
		return nil, err
	}

	return job, nil

}

func (v JobConfig) Build(builder common.Builder) error {

	infobase, err := unpackInfobase(v.Infobase)
	if err != nil {
		return err
	}

	job := jobs.NewJobBuilder(v.Uuid, nil, jobs.ValuesMap{
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

	return builder.Build(v.Name, infobase, jobs.List{
		job.Build(),
	})
}

func (v JobConfig) addUpdateJob(job jobs.JobBuilder) error {

	if v.Update == nil {
		return nil
	}

	err := v.Update.Task(job)

	return err

}

func (v JobConfig) addBackupJob(job jobs.JobBuilder) error {

	if v.Backup == nil {
		return nil
	}
	err := v.Backup.Task(job)
	return err

}

func (v JobConfig) addBlockSessionJob(job jobs.JobBuilder) error {
	if v.Sessions == nil {
		return nil
	}
	return nil
}

func (v JobConfig) addExtensionJob(job jobs.JobBuilder) error {
	if v.Extension == nil {
		return nil
	}
	return nil
}

func (v JobConfig) addEnterpriseJob(job jobs.JobBuilder) error {
	if v.Enterprise == nil {
		return nil
	}
	return nil
}

func (v JobConfig) addUnblockSessionJob(job jobs.JobBuilder) error {
	if v.Sessions == nil {
		return nil
	}
	return nil
}

func (v JobConfig) addRestoreJob(job jobs.JobBuilder) error {
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
