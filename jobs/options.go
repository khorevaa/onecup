package jobs

type Option func(o *Options)

type Options struct {
	Inputs  Inputs
	Outputs Inputs
	Check   Check
}

type ApplyOption interface {
	apply(opt Options)
}

func WithInput(inputs Inputs) Option {
	return func(o *Options) {
		o.Inputs = inputs
	}
}

func WithOutput(outputs Inputs) Option {
	return func(o *Options) {
		o.Outputs = outputs
	}
}

type Check func(state CompletionStatus, ctx Context) bool

var (
	NotError = func(state CompletionStatus, ctx Context) bool {
		return state != Error
	}

	OnError = func(state CompletionStatus, ctx Context) bool {
		return state == Error
	}

	Always = func(_ CompletionStatus, _ Context) bool {
		return true
	}
)

func WithCheck(check Check) Option {
	return func(o *Options) {
		o.Check = check
	}
}

type CheckObject interface {
	Check(state CompletionStatus, ctx Context) bool
}

func WithCheckObj(obj CheckObject) Option {
	return func(o *Options) {
		o.Check = obj.Check
	}
}
