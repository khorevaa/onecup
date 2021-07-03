package sessionControl

import (
	"github.com/khorevaa/onecup/internal/common"
	"github.com/khorevaa/onecup/workflow/context"
	"time"
)

type SessionControl struct {
	Block         bool          `json:"block,omitempty"`
	WaitTime      time.Duration `json:"wait-time,omitempty"`
	BlockDuration time.Duration `json:"block-duration,omitempty"`
	UnlockCode    string        `json:"uc-code,omitempty"`
	BlockMessage  string        `json:"lock-message,omitempty"`
	SessionFilter []string      `json:"session-filter,omitempty"`
}

func (s SessionControl) Action(ctx context.Context) (outputs context.Outputs, err error) {
	panic("implement me")
}

var defaultControl = SessionControl{
	Block:         false,
	WaitTime:      time.Minute * 5,
	BlockDuration: time.Hour * 3,
	UnlockCode:    "КодРазрешения",
	BlockMessage:  "Установлена блокировка сеансов пользователя. Код доступа: <КодРазрешения>",
	SessionFilter: []string{
		"Configurator",
	},
}

func New(cfg *common.Config) (context.Use, error) {

	control := defaultControl

	err := cfg.Unpack(&control)
	if err != nil {
		return nil, err
	}

	return &control, nil
}
