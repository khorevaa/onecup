package common

import (
	"github.com/khorevaa/onecup/jobs"
	v8 "github.com/v8platform/api"
)

type Builder interface {
	Build(name string, ib *v8.Infobase, list jobs.List) error
}

type ConfigFactory interface {
	Build(factory Builder) error
}
