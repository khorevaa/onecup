package tasks

import (
	"bytes"
	"html/template"
	"os"
	"strings"
	"time"
)

type Auth struct {
	User     string `config:"usr" json:"usr" yaml:"usr"`
	Password string `config:"pwd" json:"pwd" yaml:"pwd"`
}

type Context struct {
	Auth        Auth              `config:"auth" json:"auth"`
	env         map[string]string `config:"env"`
	params      map[string]string `config:"params"`
	Concurrency string            `config:"concurrency"`
	Strategy    StrategyConfig    `config:"strategy"`
	secrets     map[string]string `config:"secrets"`
}

type ContextConfig struct {
	Auth        Auth              `config:"auth" json:"auth"`
	Env         map[string]string `config:"env"`
	Params      map[string]string `config:"params"`
	Concurrency string            `config:"concurrency"`
	Strategy    StrategyConfig    `config:"strategy"`
}

type StrategyConfig struct {
	MaxParallel int `config:"max-parallel" json:"max_parallel" yaml:"max-parallel"`
}

func NewContext(config ContextConfig) (*Context, error) {

	ctx := Context{

		Concurrency: config.Concurrency,
		Strategy:    config.Strategy,
	}

	ctx.env = buildEnv(config.Env)

	var err error

	ctx.params, err = buildParams(ctx, config.Params)
	if err != nil {
		return nil, err
	}
	// TODO Добавить secrets
	ctx.Auth = Auth{
		User:     ctx.MustExecuteTemplate(config.Auth.User),
		Password: ctx.MustExecuteTemplate(config.Auth.Password),
	}

	return &ctx, nil
}

func buildParams(ctx Context, paramsConfig map[string]string) (map[string]string, error) {

	params := make(map[string]string)

	for key, value := range paramsConfig {

		val := value
		if strings.Contains(value, "{{") {
			newval, err := ExecuteTemplate(ctx, value)
			if err != nil {
				return nil, err
			}
			val = newval
		}
		params[key] = val

	}
	return params, nil
}

func buildEnv(env map[string]string) map[string]string {

	envs := make(map[string]string)

	for key, value := range env {
		envs[key] = os.Getenv(value)
	}

	return envs
}

func (ctx *Context) ExecuteTemplate(tmpl string) (string, error) {

	if !strings.Contains(tmpl, "{{") {
		return tmpl, nil
	}

	funcMap := ctx.FuncMap()

	var strBuf bytes.Buffer
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("Context").Funcs(funcMap).Parse(tmpl))

	err := t.Execute(&strBuf, ctx)
	if err != nil {
		return "", err
	}

	return strBuf.String(), nil
}

func (ctx *Context) MustExecuteTemplate(tmpl string) string {

	result, err := ctx.ExecuteTemplate(tmpl)
	if err != nil {
		panic(err)
	}

	return result
}

func (ctx *Context) FuncMap() map[string]interface{} {
	return template.FuncMap{
		"env":     ctx.getEnv,
		"secrets": ctx.getEnv,
		"params":  ctx.getParams,
		"auth":    ctx.getAuth,
		"now":     time.Now,
	}
}

func (ctx *Context) getEnv() map[string]string {
	return ctx.env
}

func (ctx *Context) getParams() map[string]string {
	return ctx.params
}

func (ctx *Context) getAuth() Auth {
	return ctx.Auth
}

func ExecuteTemplate(ctx interface{}, tmpl string) (string, error) {

	var strBuf bytes.Buffer
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("Context").Parse(tmpl))

	err := t.Execute(&strBuf, ctx)
	if err != nil {
		return "", err
	}

	return strBuf.String(), nil
}
