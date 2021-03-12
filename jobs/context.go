package jobs

import "fmt"

func newCtx(parent Context) *jobContext {
	return &jobContext{
		job:      parent.Job(),
		values:   make(Values),
		parent:   parent,
		simulate: parent.Simulate(),
	}
}

type Context interface {
	LoadValue(name string) (interface{}, bool)
	MustLoadValue(name string) interface{}

	StoreValue(name string, value interface{})
	OutputValue(name string, value interface{})

	StoreValues(values Values)
	OutputValues(values Values)

	Job() Job
	Err() error
	Fault() bool

	Simulate() bool
}

type jobContext struct {
	job     Job
	values  Values
	outputs Values
	err     error
	parent  Context

	simulate bool
}

func (c *jobContext) Simulate() bool {
	return c.simulate
}

func (c *jobContext) Job() Job {

	return c.job
}

func (c *jobContext) LoadValue(name string) (interface{}, bool) {

	value, ok := c.values[name]

	if c.parent != nil && !ok {
		return c.parent.LoadValue(name)
	}

	return value, ok
}
func (c *jobContext) MustLoadValue(name string) interface{} {

	value, ok := c.values[name]

	if c.parent != nil && !ok {
		return c.parent.MustLoadValue(name)
	}

	if !ok {
		panic(fmt.Sprintf("context: must have value for key <%s>", name))
	}

	return value
}

func (c *jobContext) StoreValue(name string, value interface{}) {
	c.values[name] = value
}

func (c *jobContext) StoreValues(values Values) {
	for name, value := range values {
		c.StoreValue(name, value)
	}
}

func (c *jobContext) OutputValue(name string, value interface{}) {
	c.outputs[name] = value
}

func (c *jobContext) OutputValues(values Values) {
	for name, value := range values {
		c.OutputValue(name, value)
	}
}

func (c *jobContext) Err() error {

	if c.err != nil {
		return c.err
	}

	var err error

	if c.parent != nil {
		err = c.parent.Err()
	}

	return err
}

func (c *jobContext) Fault() bool {
	return c.Err() != nil
}
