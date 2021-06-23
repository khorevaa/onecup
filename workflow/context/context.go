package context

import (
	"context"
	"github.com/khorevaa/logos"
	v8 "github.com/v8platform/api"
)

type Context interface {
	context.Context
	Infobase() v8.Infobase
	Options() []interface{}
	Logger() logos.SugaredLogger
}
