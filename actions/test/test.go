package test

import (
	"fmt"
	"io"
	"os"

	"github.com/mikesimons/yaml-dsl/types"
)

// TestAction is a dummy action useful for testing & validation purpose
type TestAction struct {
	types.ActionCommon `mapstructure:",squash"`
	Test               string `mapstructure:"test"`
	Stdout             io.Writer
}

// Prototype returns a new instance of TestAction as a Handler
func Prototype() types.Handler {
	return &TestAction{}
}

// Execute dumps the action via Printf.
// It always returns a success result and no errors.
func (action *TestAction) Execute() (*types.ActionResult, error) {
	stdout := action.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	fmt.Fprintf(stdout, "%#v\n", action)
	return &types.ActionResult{Success: true}, nil
}
