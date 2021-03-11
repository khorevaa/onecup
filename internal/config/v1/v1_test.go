package v1

import (
	"github.com/khorevaa/onecup/internal/common"
	"github.com/stretchr/testify/require"
	"testing"
)

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

	require.Equal(t, config, "test_version1")
	//require.Equal(t, rawConfig.Flows.Index, "logs-packet.flow-default")
	//require.Len(t, rawConfig.ProtocolsList, 1)

	//var protocol map[string]interface{}
	//require.NoError(t, rawConfig.ProtocolsList[0].Unpack(&protocol))
	//require.Len(t, protocol["processors"].([]interface{}), 3)
	//require.Equal(t, rawConfig.Interfaces.Device, "en1")
	//require.Len(t, rawConfig.Procs.Monitored, 2)

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
