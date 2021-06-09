package uses

import (
	"fmt"
	"github.com/khorevaa/onecup/internal/common"
	"github.com/khorevaa/onecup/uses/sessionControl"
	"github.com/khorevaa/onecup/workflow/use"
	"strings"
)

var registeredUses = map[string]UseFactory{}

type UseFactory func(cfg *common.Config) (use.Use, error)

func init() {
	RegisterUseType("session-control", sessionControl.New)
}

func getUseTypeVersion(useTypeString string) (string, string) {

	typeVersion := strings.Split(useTypeString, "@")

	useTyoe := typeVersion[0]
	useVersion := "latest"
	if len(typeVersion) == 2 {
		useVersion = typeVersion[1]
	}

	return useTyoe, useVersion
}

func CreateUseWithParams(useTypeString string, params map[string]interface{}) (use.Use, error) {

	config := common.MustNewConfigFrom(params)

	return CreateUse(useTypeString, config)

}

func CreateUse(useTypeString string, config *common.Config) (use.Use, error) {

	useType, _ := getUseTypeVersion(useTypeString)

	use, err := NewUse(useType, config)
	if err != nil {
		return nil, err
	}

	return use, nil
}

func RegisterUseType(name string, f UseFactory) {
	if registeredUses[name] != nil {
		panic(fmt.Errorf("use type '%v' exists already", name))
	}
	registeredUses[name] = f
}

func NewUse(name string, config *common.Config) (use.Use, error) {
	factory := registeredUses[name]
	if factory == nil {
		return nil, fmt.Errorf("use type %v undefined", name)
	}
	return factory(config)
}
