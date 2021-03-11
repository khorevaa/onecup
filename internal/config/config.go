package config

import (
	"fmt"
	"github.com/khorevaa/onecup/internal/common"
	v1 "github.com/khorevaa/onecup/internal/config/v1"
	"github.com/khorevaa/onecup/jobs"
	"github.com/v8platform/api"
)

var (
	configVersions = map[string]VersionFactory{}
)

type VersionFactory func(config *common.Config) (common.ConfigFactory, error)

func init() {
	RegisterConfigVersion("1.0", v1.New)
}

func RegisterConfigVersion(version string, f VersionFactory) {
	configVersions[version] = f
}

type Config struct {
	Name     string
	Infobase *v8.Infobase
	Jobs     []jobs.Job
}

func (c *Config) Build(name string, ib *v8.Infobase, list jobs.List) error {

	c.Name = name
	c.Infobase = ib
	c.Jobs = list

	return nil
}

type configInput struct {
	VersionNamespace `config:",inline,replace"`
}

func NewConfig(cfg *common.Config) (*Config, error) {

	var input configInput

	if err := cfg.Unpack(&input); err != nil {
		return nil, err
	}

	factory := input.Config()
	config := &Config{}

	err := factory.Build(config)
	if err != nil {
		return nil, err
	}

	return config, nil

}

func newConfigFactory(version string, cfg *common.Config) (common.ConfigFactory, error) {
	factory := configVersions[version]
	if factory == nil {
		return nil, fmt.Errorf("config version %v undefined", version)
	}
	return factory(cfg)

}
