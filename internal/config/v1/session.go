package v1

import (
	"github.com/khorevaa/onecup/internal/common"
	"time"
)

type SessionsConfig struct {
	Blocker       common.ConfigNamespace `config:"use,replace,required" json:"use"`
	WaitTime      time.Duration
	BlockDuration time.Duration
	UnlockCode    string
	BlockMessage  string
	SessionFilter string
}
