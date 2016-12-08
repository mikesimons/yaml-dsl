package middleware

import (
	"github.com/mikesimons/yaml-dsl/types"
	"github.com/mitchellh/mapstructure"
)

type Chain struct {
	position    int
	middlewares []types.Middleware
	DecodeFunc  func(raw types.UnparsedAction, fn func(*mapstructure.DecoderConfig)) error
}

func (chain *Chain) Reset() {
	chain.position = 0
}

func (chain *Chain) Add(fn types.Middleware) {
	chain.middlewares = append(chain.middlewares, fn)
}

func (chain *Chain) Next(action types.Action, vars map[string]interface{}) (*types.ActionResult, error) {
	defer func() {
		chain.position += 1
	}()

	if chain.position == len(chain.middlewares) {
		instance := action.Handler()
		chain.DecodeFunc(action.Data, func(config *mapstructure.DecoderConfig) {
			config.Result = instance
		})
		return instance.Execute()
	}

	if chain.position > len(chain.middlewares) {
		panic("Middleware chain beyond bounds")
	}

	return chain.middlewares[chain.position].Execute(action, vars)
}
