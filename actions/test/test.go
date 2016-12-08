package test

import (
	"fmt"

	"github.com/mikesimons/yaml-dsl/parser"
	"github.com/mikesimons/yaml-dsl/types"
)

type TestAction struct {
	parser.CommonAction `mapstructure:",squash"`
	Test                string `mapstructure:"test"`
}

func Prototype() types.Handler {
	return &TestAction{}
}

func (action *TestAction) Execute() (*types.ActionResult, error) {
	fmt.Printf("%#v\n", action)
	return &types.ActionResult{}, nil
}
