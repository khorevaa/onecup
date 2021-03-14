package jobs

type TaskOption func(o *TaskOptions)

type TaskOptions struct {
	Inputs  Inputs
	Outputs Inputs
	Check   CheckFunc
}

type ApplyOption interface {
	apply(opt *TaskOptions)
}

func WithInput(inputs Inputs) TaskOption {
	return func(o *TaskOptions) {
		o.Inputs = inputs
	}
}

func WithOutput(outputs Inputs) TaskOption {
	return func(o *TaskOptions) {
		o.Outputs = outputs
	}
}

type CheckFunc func(ctx Context, err error) bool

var (
	NotErrorCheck = func(ctx Context, err error) bool {
		return err == nil
	}

	OnError = func() TaskOption {
		return func(ctx Context, err error) bool {
			return err != nil
		}
	}

	Always = func() TaskOption {
		return func(o *TaskOptions) {
			o.Check = func(_ Context, _ error) bool {
				return true
			}
		}
	}
)

func CheckAll(checks ...CheckFunc) CheckFunc {
	return func(ctx Context, err error) bool {
		for _, check := range checks {
			if ok := check(ctx, err); !ok {
				return false
			}
		}
		return true
	}
}

func CheckOneOf(checks ...CheckFunc) CheckFunc {
	return func(ctx Context, err error) bool {
		for _, check := range checks {
			if ok := check(ctx, err); ok {
				return true
			}
		}
		return false
	}
}

func WithOptions(opts TaskOptions) TaskOption {
	return func(o *TaskOptions) {
		o = &opts
	}
}

func WithCheck(check CheckFunc) TaskOption {
	return func(o *TaskOptions) {
		o.Check = check
	}
}

type CheckObject interface {
	Check(ctx Context, err error) bool
}

func WithCheckObj(obj CheckObject) TaskOption {
	return func(o *TaskOptions) {
		o.Check = obj.Check
	}
}
