package workflow

import (
	"errors"
	"github.com/khorevaa/onecup/config"
	v8 "github.com/v8platform/api"
)

type Infobase struct {
	v8.Infobase
	Name string
	ID   string
}

type FileInfobaseConfig struct {
	File string `config:"file,required" json:"file"`
}

type ServerInfobaseConfig struct {
	Serv string `config:"srv,required" json:"srv"`
	Ref  string `config:"ref,required" json:"ref"`
}

func NewInfobase(ctx interface{}, globalAuth Auth, config config.InfobaseConfig) (Infobase, error) {

	ib := Infobase{
		Name: config.Name,
	}

	if len(config.ID) == 0 {
		// ib.ID = uuid.NewV4().String()
	}
	auth := buildAuth(ctx, config.Auth)

	if len(auth.User) == 0 {
		auth = globalAuth
	}

	ib.Infobase = v8.Infobase{
		User:     auth.User,
		Password: auth.Password,
	}

	switch config.Path.Name() {
	case "file":
		var c FileInfobaseConfig
		if err := config.Path.Config().Unpack(&c); err != nil {
			return ib, err
		}

		ib.Infobase.Connect = v8.FilePath{File: c.File}

	case "server":
		var c ServerInfobaseConfig
		if err := config.Path.Config().Unpack(&c); err != nil {
			return ib, err
		}

		ib.Infobase.Connect = v8.ServerPath{Server: c.Serv, Ref: c.Ref}

	default:
		return ib, errors.New("error connection infobase string")
	}

	return ib, nil
}
