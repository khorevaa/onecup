package tasks

import (
	"github.com/khorevaa/onecup/internal/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadConfig(t *testing.T) {

	cfg, err := common.NewConfigFrom(`
context:
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

targets:
  - name: "База данных обмена"
    auth:
      usr: "Администратор"
      pwd: "Пароль"
    path:
      server:
        srv: "srv-app"
        ref: "База данных"

jobs:
  update:
    steps:
      -
        name: Блокировка сеансов пользователей
        with:
          block: true
          uc-code: ${{ params.uc }}
          lock-message: "{{ params.lock-message }}"
        uses: user-sessions@v1
        out:
          - unblock:
              type: Boolean
`)
	require.NoError(t, err)

	ctx := Config{}

	err = cfg.Unpack(&ctx)
	require.NoError(t, err)

}
