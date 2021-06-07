package tasks

import (
	"errors"
	"fmt"
	"github.com/khorevaa/onecup/internal/common"
	v8 "github.com/v8platform/api"
	"github.com/v8platform/runner"
)

var contexts = make(map[string]ContextFactory)

type ContextFactory func(config *common.Config) (JobContext, error)

type JobContext interface {
}

type context struct {
	infobase *v8.Infobase
	options  []runner.Option
}

type singleContextConfig struct {
	Id            string                 `json:"id"`
	User          string                 `json:"user"`
	Password      string                 `json:"password"`
	ConnectString common.ConfigNamespace `config:"path,replace,required" json:"path"`
}

type matrixContextConfig struct {
	Matrix []singleContextConfig `config:",inline,required" json:"matrix"`
}

type FileInfobaseConfig struct {
	File string `config:"file,required" json:"file"`
}

type ServerInfobaseConfig struct {
	Serv string `config:"serv,required" json:"srv"`
	Ref  string `config:"ref,required" json:"ref"`
}

func NewSingle(cfg *common.Config) (JobContext, error) {

	config := singleContextConfig{}

	if cfg != nil {
		if err := cfg.Unpack(&config); err != nil {
			return nil, err
		}
	}

	var ib v8.Infobase

	switch config.ConnectString.Name() {
	case "file":
		var c FileInfobaseConfig
		if err := config.ConnectString.Config().Unpack(&c); err != nil {
			return nil, err
		}

		ib = *v8.NewFileInfobase(c.File)

	case "server":
		var c ServerInfobaseConfig
		if err := config.ConnectString.Config().Unpack(&c); err != nil {
			return nil, err
		}

		ib = v8.NewServerIB(c.Serv, c.Ref)

	default:
		return nil, errors.New("error connection infobase string")
	}
	return &ib, nil

}

func init() {
	RegisterContextType("single", NewSingle)
	RegisterContextType("matrix", NewMatrix)
}

func NewMatrix(config *common.Config) (JobContext, error) {
	return nil, nil
}

func RegisterContextType(name string, f ContextFactory) {
	if contexts[name] != nil {
		panic(fmt.Errorf("context type '%v' exists already", name))
	}
	contexts[name] = f
}

func CreateContext(contextType string, config *common.Config) (JobContext, error) {

	return nil, nil

}

func NewContext(name string, config *common.Config) (JobContext, error) {
	factory := contexts[name]
	if factory == nil {
		return nil, fmt.Errorf("context type %v undefined", name)
	}
	return factory(config)
}
