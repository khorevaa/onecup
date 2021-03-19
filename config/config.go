package config

import (
	"context"
	"github.com/khorevaa/onecup/internal/common"
	_ "github.com/khorevaa/onecup/internal/config/v1"
	"github.com/khorevaa/onecup/jobs"
	v8 "github.com/v8platform/api"
	"time"
)

type Config struct {
	Name     string
	Infobase *v8.Infobase
	Options  []interface{}
	Jobs     []jobs.Job
}

func (c *Config) Build(name string, ib *v8.Infobase, list jobs.List) error {

	c.Name = name
	c.Infobase = ib
	c.Jobs = list

	return nil
}

func RunJobConfig(c *Config) error {
	runner := jobs.NewJobsRunner(c.Jobs...)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return runner.Run(ctx, jobs.Values{
		"infobase": c.Infobase,
		"options":  c.Options,
	})
}

func NewConfig(cfg *common.Config) (*Config, error) {

	config := &Config{}

	err := common.Unpack(cfg, config)
	if err != nil {
		return nil, err
	}
	return config, nil

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
