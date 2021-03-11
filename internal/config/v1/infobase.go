package v1

import (
	"errors"
	"github.com/khorevaa/onecup/internal/common"
	v8 "github.com/v8platform/api"
)

type InfobaseConfig struct {
	User          string
	Password      string
	ConnectString common.ConfigNamespace `rawConfig:"path,replace,required" json:"path"`
}

type FileInfobaseConfig struct {
	File string `rawConfig:"file,required" json:"file"`
}

type ServerInfobaseConfig struct {
	Serv string `rawConfig:"serv,required" json:"srv"`
	Ref  string `rawConfig:"ref,required" json:"ref"`
}

func unpackInfobase(ibConfig *InfobaseConfig) (*v8.Infobase, error) {

	if ibConfig == nil {
		return nil, nil
	}

	var ib v8.Infobase

	switch ibConfig.ConnectString.Name() {
	case "file":
		var c FileInfobaseConfig
		if err := ibConfig.ConnectString.Config().Unpack(&c); err != nil {
			return nil, err
		}

		ib = *v8.NewFileInfobase(c.File)

	case "server":
		var c ServerInfobaseConfig
		if err := ibConfig.ConnectString.Config().Unpack(&c); err != nil {
			return nil, err
		}

		ib = v8.NewServerIB(c.Serv, c.Ref)

	default:
		return nil, errors.New("error connection infobase string")
	}
	ib.User = ibConfig.User
	ib.Password = ibConfig.Password
	return &ib, nil

}
