package jobs

import v8 "github.com/v8platform/api"

type Context struct {
	job         *job
	currentStep *Step
	params      Input
	outputs     Output
	err         error
}

func (c *Context) Job() *job {
	return c.job
}

func (c *Context) Step() *Step {
	return c.currentStep
}

func (c *Context) Error() error {
	return c.err
}

func (c *Context) Fault() bool {
	return c.err != nil
}

func (c *Context) Out(name string, value interface{}) {

	if c.outputs == nil {
		c.outputs = make(map[string]interface{})
	}
	c.outputs[name] = value

	if c.params == nil {
		c.params = make(map[string]interface{})
	}
	c.params[name] = value
}

func (c *Context) Param(name string) (interface{}, bool) {

	if c.currentStep == nil {
		return nil, false
	}

	value, ok := c.currentStep.params[name]
	return value, ok
}

func (c *Context) Infobase() *v8.Infobase {

	if ib, ok := c.Value("infobase"); ok {
		return ib.(*v8.Infobase)
	}
	return nil

}

func (c *Context) Options() []interface{} {

	if val, ok := c.Value("options"); ok {
		return val.([]interface{})
	}
	return nil

}

func (c *Context) Value(name string) (interface{}, bool) {

	if c.params == nil {
		return nil, false
	}
	value, ok := c.params[name]
	return value, ok

}
