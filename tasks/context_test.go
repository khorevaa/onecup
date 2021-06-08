package tasks

import (
	"github.com/khorevaa/onecup/internal/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadContext(t *testing.T) {

	cfg, err := common.NewConfigFrom(`
auth:
  usr: "Администратор"
  pwd: "Пароль"
env:
  cluster-admin: CLUSTER_ADMIN
params:
  uc: КодБлокировки
  lock-message: Сообщение блокировки
  v8-path: путь к платформе
concurrency: staging_environment
strategy:
  max-parallel: 2
`)
	require.NoError(t, err)

	ctx := ContextConfig{}

	err = cfg.Unpack(&ctx)
	require.NoError(t, err)

}
