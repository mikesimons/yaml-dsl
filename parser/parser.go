package parser

import (
	"fmt"
	"reflect"

	"github.com/mikesimons/yaml-dsl/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/mikesimons/yaml-dsl/parser/middleware"
)

type Dsl struct {
	Handlers     map[string]types.HandlerPrototypeFunc
	ScriptParser types.ScriptParser
	Middleware   []types.Middleware
}

func New() *Dsl {
	return &Dsl{
		Handlers: make(map[string]types.HandlerPrototypeFunc),
	}
}

func (dsl *Dsl) Parse(unparsedList *types.UnparsedActionList) (*ActionList, error) {
	out := &ActionList{
		Dsl: dsl,
		Middlewares: &middleware.Chain{
			Middleware: dsl.Middleware,
			DecodeFunc: dsl.Decode,
		},
	}

	for _, unparsedAction := range *unparsedList {
		var action types.Action

		for key := range unparsedAction {
			if protoFunc, ok := dsl.Handlers[key]; ok {
				action = types.Action{
					Handler: protoFunc,
					Data:    unparsedAction,
				}
				break
			}
		}

		if action.Handler == nil {
			fmt.Printf("Warning: Could not match action to a handler - %#v\n", unparsedAction)
		}

		if action.Handler != nil {
			out.PushBack(action)
		}
	}
	return out, nil
}

func (dsl *Dsl) Decode(raw types.UnparsedAction, fn func(*mapstructure.DecoderConfig)) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		DecodeHook:       dsl.scriptDecodeHook(),
	}

	fn(config)

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(raw)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("shell action could not be parsed from %#v", raw))
	}

	return nil
}

func (dsl *Dsl) scriptDecodeHook() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {

		if f.Kind() != reflect.String {
			return data, nil
		}

		return dsl.ScriptParser.ParseEmbedded(data.(string))
	}
}
