package middleware

import "github.com/mikesimons/yaml-dsl/types"

// Middleware is an interface that all middleware needs to implement
type Middleware interface {
	Execute(action types.Action, vars map[string]interface{}, chain Chain) (*types.ActionResult, error)
}
