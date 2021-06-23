package context

type Use interface {
	Action(ctx Context) (outputs Outputs, err error)
}

type Outputs map[string]interface{}
