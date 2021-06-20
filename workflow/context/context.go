package context

import v8 "github.com/v8platform/api"

type Context interface {
	Infobase() v8.Infobase
}
