package tasks

import (
	"github.com/khorevaa/onecup/jobs"
)

var _ jobs.StepInterface = (*FileReleaseStep)(nil)

type FileReleaseStep struct {
	File string
	Hash string
}

func (j *FileReleaseStep) Handler() jobs.HandlerType {
	return jobs.DefaultType
}

func (j *FileReleaseStep) Action(ctx jobs.Context) error {
	ctx.StoreValue("release-file", j.File)
	return nil
}

func (j *FileReleaseStep) Name() string {
	return "Getting release file"
}
