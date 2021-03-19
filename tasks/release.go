package tasks

import (
	"github.com/khorevaa/onecup/jobs"
)

var _ jobs.TaskObject = (*FileReleaseStep)(nil)

type FileReleaseStep struct {
	File string
	Hash string
}

func (j *FileReleaseStep) Action() jobs.TaskAction {
	return j.action
}

func (j *FileReleaseStep) Inputs() jobs.ValuesMap {
	return nil
}

func (j *FileReleaseStep) Outputs() jobs.ValuesMap {
	return map[string]string{
		"file": "file",
	}
}

func (j *FileReleaseStep) action(ctx jobs.Context) error {
	ctx.OutputValue("file", j.File)
	return nil
}

func (j *FileReleaseStep) Name() string {
	return "Getting release file"
}
