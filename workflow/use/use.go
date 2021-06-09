package use

import (
	"context"
	v8 "github.com/v8platform/api"
)

type Use interface {
	Action(ctx context.Context, infobase v8.Infobase, outputs map[string]interface{}) error
}
