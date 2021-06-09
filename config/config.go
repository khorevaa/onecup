package config

import (
	"github.com/khorevaa/onecup/internal/common"
	_ "github.com/khorevaa/onecup/internal/config/v1"
)

type Config struct {
	Name        string
	Env         map[string]string        `config:"env"`
	Params      map[string]TemplateValue `config:"params"`
	Concurrency string                   `config:"concurrency"`
	Strategy    StrategyConfig           `config:"strategy"`
	// Infobase    InfobaseConfig       `config:"infobase"`
	InfobaseList InfobaseListConfig   `config:"infobase"`
	Jobs         map[string]JobConfig `config:"jobs"`
}

type JobConfig struct {
	Steps []StepConfig  `config:"steps,required" json:"steps"`
	Need  []string      `config:"need" json:"need"`
	If    TemplateValue `config:"if" json:"if"`
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
	If    TemplateValue          `config:"if" json:"if"`
	Out   map[string]ParamConfig `json:"out"`
	Cache *CacheConfig           `config:"cache" json:"path"`
}

type AuthConfig struct {
	User     TemplateValue `config:"usr" json:"usr" yaml:"usr"`
	Password TemplateValue `config:"pwd" json:"pwd" yaml:"pwd"`
}

type StrategyConfig struct {
	MaxParallel int `config:"max-parallel" json:"max_parallel" yaml:"max-parallel"`
}

type InfobaseListConfig struct {
	Auth     AuthConfig       `config:"auth" json:"auth"`
	Infobase InfobaseConfig   `config:",inline" json:"infobase"`
	Items    []InfobaseConfig `config:"items" json:"items"`
}

type InfobaseConfig struct {
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
