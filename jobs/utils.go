package jobs

import (
	v8 "github.com/v8platform/api"
)

const (
	infobaseKey = "_infobase_"
	optionsKey  = "_options_"

	stepKey   = "_step_"
	taskKey   = "_task_"
	outputKey = "_output_"
)

func OutputFromCtx(ctx Context) Values {

	if val, ok := ctx.LoadValue(outputKey); ok {
		return val.(Values)
	}
	return nil

}

func InfobaseFromCtx(ctx Context) *v8.Infobase {

	if ib, ok := ctx.LoadValue(infobaseKey); ok {
		return ib.(*v8.Infobase)
	}
	return nil

}
func OptionsFromCtx(ctx Context) []interface{} {

	if val, ok := ctx.LoadValue(optionsKey); ok {
		return val.([]interface{})
	}
	return nil

}

func WithValues(parent Context, values Values) Context {

	ctx := newCtx(parent)
	ctx.StoreValues(values)
	return ctx

}

func withValue(parent Context, key string, value interface{}) *jobContext {

	ctx := newCtx(parent)
	ctx.StoreValue(key, value)
	return ctx

}

func WithInfobase(parent Context, infobase *v8.Infobase) Context {
	return withValue(parent, infobaseKey, infobase)
}

func newTaskContext(parent Context, t *task) *jobContext {
	return withValue(parent, taskKey, t)
}

func getValues(values Values, keyMap map[string]string) Values {

	val := make(Values)

	for k1, k2 := range keyMap {
		val[k1] = values[k2]
	}

	return val
}

func copyValues(dest, src Values) {
	if src == nil {
		return
	}

	for key, value := range src {
		dest[key] = value
	}
}
