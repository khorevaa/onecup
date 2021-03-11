package common

import (
	"fmt"
)

var (
	configVersions = map[string]VersionFactory{}
)

type VersionFactory func(config *Config) (ConfigFactory, error)

func RegisterConfigVersion(version string, f VersionFactory) {
	configVersions[version] = f
}

type configInput struct {
	VersionNamespace `config:",inline,replace"`
}

func Unpack(cfg *Config, builder Builder) error {

	var input configInput

	if err := cfg.Unpack(&input); err != nil {
		return err
	}

	factory := input.Config()

	err := factory.Build(builder)
	if err != nil {
		return err
	}

	return nil

}

func newConfigFactory(version string, cfg *Config) (ConfigFactory, error) {
	factory := configVersions[version]
	if factory == nil {
		return nil, fmt.Errorf("config version %v undefined", version)
	}
	return factory(cfg)

}
