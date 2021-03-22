package v1

import (
	"github.com/khorevaa/onecup/internal/common"
	"github.com/khorevaa/onecup/jobs"
	"github.com/stretchr/testify/require"
	v8 "github.com/v8platform/api"
	"testing"
)

type TestConfig struct {
	Name     string
	Infobase *v8.Infobase
	Options  []interface{}
	Jobs     []jobs.Job
}

func (c *TestConfig) Build(name string, ib *v8.Infobase, list jobs.List) error {

	c.Name = name
	c.Infobase = ib
	c.Jobs = list

	return nil
}

func TestLoadConfig(t *testing.T) {

	cfg, err := common.NewConfigFrom(`
version: Config.0
name: test_version1
infobase:
  user: admin
  path: 
	server:
	  serv: app1
	  ref: testib

configuration:
  release:
	file:
	  path: ./update.cfu
  load-rawConfig: true
`)
	require.NoError(t, err)
	config, err := New(cfg)
	require.NoError(t, err)
	var runConfig TestConfig
	err = config.Build(&runConfig)
	require.NoError(t, err)

	require.Equal(t, runConfig.Name, "test_version1")
	require.Equal(t, len(runConfig.Jobs), 1)
	require.Equal(t, runConfig.Name, "test_version1")
	//require.Equal(t, rawConfig.Flows.Index, "logs-packet.flow-default")
	//require.Len(t, rawConfig.ProtocolsList, 1)

	//var protocol map[string]interface{}
	//require.NoError(t, rawConfig.ProtocolsList[0].Unpack(&protocol))
	//require.Len(t, protocol["processors"].([]interface{}), 3)
	//require.Equal(t, rawConfig.Interfaces.Device, "en1")
	//require.Len(t, rawConfig.Procs.Monitored, 2)

}

func TestLoadConfigJobs(t *testing.T) {

	cfg, err := common.NewConfigFrom(`
version: 1.0
jobs:
  job1: 
    infobase:
      user: admin
      path: 
       server:
         serv: app1
         ref: testib

    configuration:
      release:
        file:
          path: ./update.cfu
      load-config: true
`)
	require.NoError(t, err)
	config, err := New(cfg)
	require.NoError(t, err)

	var runConfig TestConfig
	err = config.Build(&runConfig)
	require.NoError(t, err)

	require.Equal(t, runConfig.Name, "test_version1")
	require.Equal(t, len(runConfig.Jobs), 1)
	require.Equal(t, runConfig.Name, "test_version1")

}

func TestLoadConfigEmpty(t *testing.T) {

	cfg, err := common.NewConfigFrom(`
version: Config.0
name: test_version1
infobase:
  user: admin
  path: 
	server:
	  serv: app1
	  ref: testib

`)
	require.NoError(t, err)
	config, err := New(cfg)
	require.NoError(t, err)

	require.Equal(t, config, "test_version1")
	//require.Equal(t, rawConfig.Flows.Index, "logs-packet.flow-default")
	//require.Len(t, rawConfig.ProtocolsList, 1)

	//var protocol map[string]interface{}
	//require.NoError(t, rawConfig.ProtocolsList[0].Unpack(&protocol))
	//require.Len(t, protocol["processors"].([]interface{}), 3)
	//require.Equal(t, rawConfig.Interfaces.Device, "en1")
	//require.Len(t, rawConfig.Procs.Monitored, 2)

}
