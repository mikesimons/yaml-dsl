package shell

import (
	"fmt"

	"github.com/mikesimons/yaml-dsl/parser"
	"github.com/mikesimons/yaml-dsl/types"
)

type ShellAction struct {
	parser.CommonAction `mapstructure:",squash"`
	Command             string `mapstructure:"shell"`
}

func Prototype() types.Handler {
	return &ShellAction{}
}

func (action *ShellAction) Execute() (*types.ActionResult, error) {
	fmt.Printf("%#v\n", action)
	return &types.ActionResult{}, nil
}
