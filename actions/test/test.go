package test

import (
	"fmt"
	"io"
	"os"

	"github.com/mikesimons/yaml-dsl/parser"
	"github.com/mikesimons/yaml-dsl/types"
)

type TestAction struct {
	parser.CommonAction `mapstructure:",squash"`
	Test                string `mapstructure:"test"`
	Stdout              io.Writer
}

func Prototype() types.Handler {
	return &TestAction{}
}

func (action *TestAction) Execute() (*types.ActionResult, error) {
	stdout := action.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	fmt.Fprintf(stdout, "%#v\n", action)
	return &types.ActionResult{Success: true}, nil
}
