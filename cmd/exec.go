package cmd

import (
	"errors"
	"github.com/khorevaa/onecup/config"
	"github.com/khorevaa/onecup/internal/common"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type execCommand struct {
	configFile string
	simulate   bool
}

func (c *execCommand) run(context *cli.Context) error {

	_, err := os.Stat(c.configFile)
	if err != nil {
		return err
	}

	ext := filepath.Ext(c.configFile)

	var cfg *common.Config

	bs, err := ioutil.ReadFile(c.configFile)
	if err != nil {
		return err
	}

	switch strings.ToLower(ext) {
	case "yaml", "yml":
		cfg, err = common.NewConfigWithYAML(bs, c.configFile)
		if err != nil {
			return err
		}
	case "json":
		cfg, err = common.NewConfigWithJSON(bs, c.configFile)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown config file format. Supposed `JSON` or `YAML")
	}

	jobConfig, err := config.NewConfig(cfg)
	if err != nil {
		return err
	}

	if c.simulate {
		err = config.SimulateRunJobConfig(jobConfig)
	} else {
		err = config.RunJobConfig(jobConfig)
	}

	return err
}

func (c *execCommand) Flags() []cli.Flag {

	return []cli.Flag{
		&cli.StringFlag{
			Destination: &c.configFile, Name: "config",
			Value: "", Usage: "job config file `FILE`"},
		&cli.BoolFlag{
			Destination: &c.simulate, Name: "simulate", Aliases: []string{"s"},
			Value: false, Usage: "simulate job work"},
	}
}

func (c *execCommand) Cmd() *cli.Command {

	cmd := &cli.Command{
		//Category:    "some_category",
		Name:        "exec",
		Usage:       "Execute any job from config",
		Description: "This command help execute any job from config file",
		Action:      c.run,
		Flags:       c.Flags(),
	}

	return cmd
}
