package tasks

import (
	"context"
	v8 "github.com/v8platform/api"
)

func Run(ctx context.Context, ib v8.ConnectionString, what v8.Command, opts ...interface{}) error {
	process, err := v8.Background(ctx, ib, what, opts...)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-process.Wait():
		return err
	}

}
