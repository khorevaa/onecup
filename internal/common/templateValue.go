package common

import (
	"bytes"
	"html/template"
	"strings"
	"time"
)

type TemplateValue string

type FuncMapInterface interface {
	FuncMap() template.FuncMap
}

func (val TemplateValue) MustExecute(context interface{}) string {
	value, err := val.Execute(context)
	if err != nil {
		panic(err)
	}
	return value
}

func (val TemplateValue) Execute(context interface{}) (string, error) {

	tmpl := string(val)

	if !strings.Contains(tmpl, "{{") {
		return tmpl, nil
	}

	tmpl = strings.ReplaceAll(tmpl, "$", "")

	funcMap := defaultFuncMap()

	switch typed := context.(type) {
	case template.FuncMap:
		for key, value := range typed {
			funcMap[key] = value
		}
	case FuncMapInterface:
		for key, value := range typed.FuncMap() {
			funcMap[key] = value
		}
	}

	var strBuf bytes.Buffer
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("Context").Funcs(funcMap).Parse(tmpl))

	err := t.Execute(&strBuf, context)
	if err != nil {
		return "", err
	}

	return strBuf.String(), nil

}

func defaultFuncMap() map[string]interface{} {
	return template.FuncMap{
		"now": time.Now,
	}
}
