package tasks

import (
	"errors"
	"github.com/khorevaa/onecup/internal/common"
	v8 "github.com/v8platform/api"
	"github.com/v8platform/runner"
)

type Target struct {
	ID       string
	Name     string
	Auth     Auth
	Infobase v8.Infobase
	Options  []runner.Option
}

type TargetConfig struct {
	ID   string                 `json:"id,omitempty"`
	Name string                 `json:"name,omitempty"`
	Auth Auth                   `json:"auth,omitempty"`
	Path common.ConfigNamespace `config:"path,required" json:"path"`
}

type FileInfobaseConfig struct {
	File string `config:"file,required" json:"file"`
}

type ServerInfobaseConfig struct {
	Serv string `config:"srv,required" json:"srv"`
	Ref  string `config:"ref,required" json:"ref"`
}

func NewTarget(ctx Context, config TargetConfig) (*Target, error) {

	target := Target{
		Name: config.Name,
	}

	if len(config.ID) == 0 {
		// target.ID = uuid.NewV4().String()
	}

	auth := getTargetAuth(ctx, config.Auth)
	target.Auth = auth
	target.Infobase = v8.Infobase{
		User:     auth.User,
		Password: auth.Password,
	}

	switch config.Path.Name() {
	case "file":
		var c FileInfobaseConfig
		if err := config.Path.Config().Unpack(&c); err != nil {
			return nil, err
		}

		target.Infobase.Connect = v8.FilePath{File: c.File}

	case "server":
		var c ServerInfobaseConfig
		if err := config.Path.Config().Unpack(&c); err != nil {
			return nil, err
		}

		target.Infobase.Connect = v8.ServerPath{Server: c.Serv, Ref: c.Ref}

	default:
		return nil, errors.New("error connection infobase string")
	}

	return &target, nil
}

func getTargetAuth(ctx Context, authConfig Auth) Auth {

	if len(authConfig.User) > 0 {

		return Auth{
			User:     ctx.MustExecuteTemplate(authConfig.User),
			Password: ctx.MustExecuteTemplate(authConfig.Password),
		}
	}

	return ctx.Auth
}
