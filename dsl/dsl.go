package dsl

import (
	"fmt"
	"reflect"

	"github.com/mikesimons/yaml-dsl/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	//"fmt"
)

type HandlerFunc func(raw types.RawAction, dsl *Dsl) error

type Dsl struct {
	Handlers     map[string]HandlerFunc
	ScriptParser types.ScriptParser
}

type Action struct {
	Handler HandlerFunc
	Data    types.RawAction
}

func New() *Dsl {
	return &Dsl{
		Handlers: make(map[string]HandlerFunc),
	}
}

func (dsl *Dsl) ProcessRawActions(raw *types.RawActionList) (*ActionList, error) {
	out := &ActionList{
		Dsl: dsl,
	}

	for _, rawAction := range *raw {
		var action Action

		for key := range rawAction {
			if handlerFn, ok := dsl.Handlers[key]; ok {
				action = Action{
					Handler: handlerFn,
					Data:    rawAction,
				}
				break
			}
		}

		if action.Handler == nil {
			fmt.Printf("Warning: Could not match action to a handler - %#v\n", rawAction)
		}

		if action.Handler != nil {
			out.PushBack(action)
		}
	}
	return out, nil
}

func (dsl *Dsl) Decode(raw types.RawAction, fn func(*mapstructure.DecoderConfig)) error {
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

		return dsl.ScriptParser.Parse(data.(string))
	}
}
