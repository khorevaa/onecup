package config

import (
	"github.com/khorevaa/onecup/internal/common"
	"strings"
)

// VersionNamespace storing at most one configuration section by name and sub-section.
type VersionNamespace struct {
	version string `config:"version,required"`
	config  common.ConfigFactory
}

// Unpack unpacks a configuration with at most one sub object. An sub object is
// ignored if it is disabled by setting `enabled: false`. If the configuration
// passed contains multiple active sub objects, Unpack will return an error.
func (ns *VersionNamespace) Unpack(cfg *common.Config) error {

	fields := cfg.GetFields()
	if len(fields) == 0 {
		return nil
	}

	configVersion, err := cfg.String("version", -1)
	if err != nil {
		// element is no configuration object -> continue so a namespace
		// configVersionFactory unpacked as a namespace can have other configuration
		// values as well
		return err
	}

	configVersion = strings.TrimPrefix(configVersion, "v")

	ns.version = configVersion

	config, err := newConfigFactory(configVersion, cfg)

	if err != nil {
		return err
	}

	ns.config = config

	return nil
}

// Name returns the configuration sections it's name if a section has been set.
func (ns *VersionNamespace) Name() string {
	return ns.version
}

// configVersionFactory return the sub-configuration section if a section has been set.
func (ns *VersionNamespace) Config() common.ConfigFactory {
	return ns.config
}

// IsSet returns true if a sub-configuration section has been set.
func (ns *VersionNamespace) IsSet() bool {
	return ns.config != nil
}
