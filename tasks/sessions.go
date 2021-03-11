package tasks

import (
	"time"
)

type BlockSessionsRas struct {
	RasClientConfig string
	WaitTime        time.Duration
	BlockDuration   time.Duration
	UnlockCode      string
	BlockMessage    string
	SessionFilter   string
}

func (j *BlockSessionsRas) Name() string {
	return "Block sessions"
}

type UnblockSessionsRas struct {
	File           string
	RestoreOnError bool
	doRestore      bool
	backupFileName string
	Dir            string
}

func (j *UnblockSessionsRas) Name() string {
	return "Sessions"
}
