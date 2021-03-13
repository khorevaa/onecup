package jobs

import (
	"errors"
	"time"
)

type Action func(ctx Context) error

func NewStep(name string, fn Action) Task {
	return &taskStep{
		name: name,
		fn:   fn,
	}
}

var _ Task = (*taskStep)(nil)

type taskStep struct {
	name string
	fn   Action

	status  CompletionStatus
	startAt time.Time
	endAt   time.Time
}

func (s *taskStep) Stats() Stats {
	return Stats{
		StartAt: s.startAt,
		EndAt:   s.endAt,
	}
}

func (s *taskStep) Fault() bool {
	return s.Status() == Error
}

func (s *taskStep) Success() bool {
	return s.Status() == Success
}

func (s *taskStep) Status() CompletionStatus {
	return s.Status()
}

func (s *taskStep) Name() string {
	return s.name
}

func (s *taskStep) Run(ctx Context) (output Values, err error) {

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

	err = s.fn(ctx)
	if err != nil {
		s.status = Error
		return
	}
	s.status = Success
	return
}
