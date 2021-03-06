package config

import "testing"

func TestSimulateRunJobConfig(t *testing.T) {

	tests := []struct {
		name       string
		yamlConfig string
		wantErr    bool
	}{
		{
			"full config",
			`
name: Пакетное обновление

env:
  cluster-admin: CLUSTER_ADMIN
params:
  uc: КодБлокировки
  lock-message: Сообщение блокировки
  v8-path: путь к платформе
concurrency: staging_environment
strategy:
  max-parallel: 2

infobase:
  auth:
    usr: "Администратор"
    pwd: "Пароль"
  items:
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
          lock-message: ${{ params.lock-message }}
        uses: user-sessions@v1
        outputs:
          unblock:
              type: Boolean
      -
        name: Резервная копия
        with:
          uc: ${{ params.uc }}
          path-template: ./backup/${{ target.path.ref }}_${{ now.Format "02-01-2006 15:04:05"}}.dt
        uses: backup@v1
        out:
          - backup-file:
              type: String
      -
        name: Получение релиза
        with:
          url:
          binary:
        uses: release-get@v1
        out:
          - release-file:
              type: String
        cache:
          path: ./cache/release
      -
        name: Обновление конфигурации
        with:
          uc: ${{ params.uc }}
          release: ${{ params.release-file }}
          load-config: false
          server: false
          dynamic: false
          warning-as-erros: false
          rollback-on-error: true
        uses: update-cf@v1
      -
        name: Выполнение обработчиков обновления в Предприятии
        with:
          uc: ${{ params.uc }}
        uses: enterprise-update@v1
      -
        name: Обновление расширений конфигурации
        uses: update-cfe@v1
        with:
          uc: ${{ params.uc }}
          cfe:
            - name: Расширение1
              version: 1.0.1
              update:
                path:
                  file: ./path-to-file
              hash: md5
            - name: Расширение2
              version: 1.0.2
              update:
                binary:
                  bin: base64
              hash: md5
        cache:
          path: ./cache/cfe

      -
        name: Восстановление бекапа
        if: ${{ failure }}
        uses: restore-ib@v1
        with:
          backup-file: ${{ params.backup-file }}
      -
        name: Разблокировка сеансов пользователей
        if: ${{ always && params.unblock }}
        uses: user-sessions@v1
        with:
          block: false
`,
			false,
		},
		{
			"list config",
			`
name: Пакетное обновление

infobase:
  auth:
    usr: "Администратор"
    pwd: "Пароль"
  items:
    - name: "База данных обмена"
      auth:
        usr: "Администратор"
        pwd: "Пароль"
      path:
        server:
          srv: "srv-app"
          ref: "База данных"

jobs:
`,
			false,
		},
		{
			"one infobase",
			`
name: Пакетное обновление

infobase:
  name: "База данных обмена"
  auth:
    usr: "Администратор"
    pwd: "Пароль"
  path:
    server:
      srv: "srv-app"
      ref: "База данных"

jobs:
`,
			false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfigFrom(tt.yamlConfig)

			if err != nil != tt.wantErr {
				t.Errorf("SimulateRunJobConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
