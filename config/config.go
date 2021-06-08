package config

import (
	"github.com/khorevaa/onecup/internal/common"
	_ "github.com/khorevaa/onecup/internal/config/v1"
)

type Config struct {
	Context ContextConfig `yaml:"context"`
	Targets []TargetConfig
	Jobs    map[string]JobConfig
}

type JobConfig struct {
	Steps []StepConfig `config:"steps,required" json:"steps"`
	Need  []string     `config:"need" json:"need"`
	If    string       `config:"if" json:"if"`
}

type ParamConfig struct {
	Type string `json:"type" yaml:"type"`
}

type CacheConfig struct {
	Path string
}

type StepConfig struct {
	ID    string                 `json:"id,omitempty"`
	Name  string                 `config:"name,required" json:"name,omitempty"`
	With  map[string]interface{} `json:"with,omitempty"`
	Uses  string                 `config:"uses,required" json:"uses"`
	If    string                 `config:"if" json:"if"`
	Out   map[string]ParamConfig `json:"out"`
	Cache *CacheConfig           `config:"cache" json:"path"`
}

type AuthConfig struct {
	User     string `config:"usr" json:"usr" yaml:"usr"`
	Password string `config:"pwd" json:"pwd" yaml:"pwd"`
}

type ContextConfig struct {
	Auth        AuthConfig        `config:"auth" json:"auth"`
	Env         map[string]string `config:"env"`
	Params      map[string]string `config:"params"`
	Concurrency string            `config:"concurrency"`
	Strategy    StrategyConfig    `config:"strategy"`
}

type StrategyConfig struct {
	MaxParallel int `config:"max-parallel" json:"max_parallel" yaml:"max-parallel"`
}

type TargetConfig struct {
	ID   string                 `json:"id,omitempty"`
	Name string                 `json:"name,omitempty"`
	Auth AuthConfig             `json:"auth,omitempty"`
	Path common.ConfigNamespace `config:"path,required" json:"path"`
}

type FileInfobaseConfig struct {
	File string `config:"file,required" json:"file"`
}

type ServerInfobaseConfig struct {
	Serv string `config:"srv,required" json:"srv"`
	Ref  string `config:"ref,required" json:"ref"`
}

func NewConfig(cfg *common.Config) (*Config, error) {

	config := Config{}

	err := cfg.Unpack(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil

}

func NewConfigFrom(from interface{}) (*Config, error) {

	cfg, err := common.NewConfigFrom(from)
	if err != nil {
		return nil, err
	}

	return NewConfig(cfg)

}

func NewConfigFromJSON(in []byte, source string) (*Config, error) {

	cfg, err := common.NewConfigWithJSON(in, source)
	if err != nil {
		return nil, err
	}

	return NewConfig(cfg)

}
