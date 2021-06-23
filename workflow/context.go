package workflow

import (
	"context"
	"github.com/khorevaa/logos"
	v8 "github.com/v8platform/api"
)

type jobContext struct {
	context.Context

	infobase v8.Infobase
	options  []interface{}
	logger   logos.SugaredLogger
}

func (j jobContext) Infobase() v8.Infobase {
	return j.infobase
}

func (j jobContext) Options() []interface{} {
	return j.options
}

func (j jobContext) Logger() logos.SugaredLogger {
	return j.logger
}
