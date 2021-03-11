package tasks

import (
	"github.com/khorevaa/onecup/jobs"
	"log"
)

type FileReleaseJob struct {
	File string
	Hash string
}

func (j *FileReleaseJob) Action(ctx *jobs.Context) error {
	log.Println("TODO Сделать получение файла релиза")

	ctx.Out("release-file", j.File)
	return nil
}

func (j *FileReleaseJob) Params() jobs.Params {
	return map[string]interface{}{}
}

func (j *FileReleaseJob) Name() string {
	return "Getting release file"
}
