package v1

import (
	"errors"
	"github.com/khorevaa/onecup/internal/common"
	"github.com/khorevaa/onecup/jobs"
	"github.com/khorevaa/onecup/tasks"
)

type UpdateConfig struct {
	Release          common.ConfigNamespace `config:"release,replace,required" json:"release"`
	LoadConfig       bool                   `config:"load-config" json:"load-config"`
	Server           bool                   `config:"server" json:"server"`
	Dynamic          bool                   `config:"dynamic" json:"dynamic"`
	WarningsAsErrors bool                   `config:"warnings-as-errors" json:"warnings-as-errors"`
	RollbackOnError  bool                   `config:"rollback-on-error" json:"rollback-on-error"`
	Hooks            common.ConfigNamespace `config:"hooks,replace" json:"hooks"`
}

type FileReleaseConfig struct {
	File string `config:"path,required" json:"path"`
	Hash string `config:"hash" json:"hash"`
}

func (c *UpdateConfig) Task() (*jobs.TaskBuilder, error) {

	releaseStep, err := c.releaseStep()
	if err != nil {
		return nil, err
	}

	task := jobs.NewTask("update").Steps(
		releaseStep,
		&tasks.Update{
			LoadConfig:       c.LoadConfig,
			Server:           c.Server,
			Dynamic:          c.Dynamic,
			WarningsAsErrors: c.WarningsAsErrors,
		})

	if c.RollbackOnError {
		task.Steps(&tasks.RollbackUpdate{
			HandlerType: jobs.ErrorType,
		})
	}

	return task, nil
}

func (c *UpdateConfig) releaseStep() (jobs.Step, error) {

	switch c.Release.Name() {

	case "file":

		var fileConfig FileReleaseConfig
		if err := c.Release.Config().Unpack(&fileConfig); err != nil {
			return nil, err
		}

		return &tasks.FileReleaseStep{
			File: fileConfig.File,
			Hash: fileConfig.Hash,
		}, nil

	case "binary":
		panic("not implement")
	default:
		return nil, errors.New("required release config")
	}
}
