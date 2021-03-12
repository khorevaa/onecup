package jobs

import (
	"errors"
)

type Action func(ctx Context) error

type HandlerType int

const (
	BeforeType HandlerType = iota
	DefaultType
	AfterType
	ErrorType
	AlwaysType
)

func (h HandlerType) String() string {
	switch h {
	case BeforeType:
		return "Before"
	case DefaultType:
		return "Default"
	case AfterType:
		return "After"
	case ErrorType:
		return "Error"
	case AlwaysType:
		return "Always"
	default:
		return "unknown handler type"
	}
}

type TaskStep struct {
	name    string
	handler HandlerType
	fn      Action
	status  CompletionStatus
}

func (s *TaskStep) Name() string {
	return s.name
}

func (s *TaskStep) run(ctx *jobContext) (err error) {

	defer func() {
		if rec := recover(); rec != nil {
			switch t := rec.(type) {
			case error:
				err = t
			case string:
				err = errors.New(t)
			default:
				panic(rec)
			}
		}
	}()

	if s.needSkip(ctx) {
		return nil
	}

	err = s.fn(ctx)

	if err != nil {
		ctx.err = err
		return err
	}

	return nil
}

func (s *TaskStep) needSkip(ctx Context) bool {
	if ctx.Fault() && !(s.handler == ErrorType || s.handler == AlwaysType) {
		s.status = Skip
		return true
	}

	return false
}
