package workflow

import (
	"context"
	"github.com/khorevaa/onecup/config"
)

type Context struct {
	parent context.Context
}

func NewContext(ctx context.Context, config config.ContextConfig) {

}
