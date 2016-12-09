package middleware

import (
	"github.com/mikesimons/yaml-dsl/types"
	"github.com/mitchellh/mapstructure"
)

// Chain implements a simple middleware iterator
type Chain struct {
	position   int
	Middleware []Middleware
	DecodeFunc func(raw types.UnparsedAction, fn func(*mapstructure.DecoderConfig)) error
}

// Reset resets the chain position back to the start
func (chain *Chain) Reset() {
	chain.position = 0
}

// Next executes the next middleware in the chain or the action itself if we are at the end of the chain
// Calling `Next` once the action has been executed will result in a panic
func (chain *Chain) Next(action types.Action, vars map[string]interface{}) (*types.ActionResult, error) {

	if chain.position == len(chain.Middleware) {
		instance := action.Handler()
		chain.DecodeFunc(action.Data, func(config *mapstructure.DecoderConfig) {
			config.Result = instance
		})
		return instance.Execute()
	}

	if chain.position > len(chain.Middleware) {
		panic("Middleware chain beyond bounds")
	}

	position := chain.position
	chain.position++
	return chain.Middleware[position].Execute(action, vars, *chain)
}
