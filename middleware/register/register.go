package register

import (
	"github.com/mikesimons/yaml-dsl/middleware"
	"github.com/mikesimons/yaml-dsl/parser"
	"github.com/mikesimons/yaml-dsl/types"
)

// Middleware implements variable registration for action
type Middleware struct {
	Dsl *parser.Dsl
}

func (wim *Middleware) Execute(action types.Action, vars map[string]interface{}, chain middleware.Chain) (*types.ActionResult, error) {
	result, err := chain.Next(action, vars)

	if action.Data["register"] != nil {
		wim.Dsl.ScriptParser.SetVar(action.Data["register"].(string), result)
	}

	return result, err
}
