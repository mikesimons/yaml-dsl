package withitems

import (
	"fmt"
	"reflect"

	"github.com/mikesimons/yaml-dsl/parser"
	"github.com/mikesimons/yaml-dsl/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type Middleware struct {
	Dsl     *parser.Dsl
	Key     string
	ItemKey string
}

func (wim *Middleware) Execute(action types.Action, vars map[string]interface{}) (*types.ActionResult, error) {

	middlewareKey := wim.Key
	if middlewareKey == "" {
		middlewareKey = "with_items"
	}

	itemKey := wim.ItemKey
	if itemKey == "" {
		itemKey = "item"
	}

	var items []interface{}
	if rawItems, hasItems := action.Data[middlewareKey]; hasItems {
		value := reflect.ValueOf(rawItems)
		switch value.Kind() {
		case reflect.String:
			result, err := wim.Dsl.ScriptParser.ParseExpression(rawItems.(string))
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("error parsing expression: %s", rawItems.(string)))
			}
			items = result.([]interface{})
		case reflect.Slice:
			items = rawItems.([]interface{})
		default:
			// Should raise an error here
		}
	}

	if len(items) == 0 {
		items = append(items, "")
	}

	var results []*types.ActionResult
	for _, v := range items {
		instance := action.Handler()

		item, err := wim.Dsl.ScriptParser.ParseEmbedded(v.(string))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error parsing expression: %s", v.(string)))
		}
		wim.Dsl.ScriptParser.SetVar(itemKey, item)

		wim.Dsl.Decode(action.Data, func(config *mapstructure.DecoderConfig) {
			config.Result = instance
		})

		result, err := instance.Execute()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error executing action: %s", action))
		}
		results = append(results, result)
	}

	wim.Dsl.ScriptParser.SetVar(itemKey, "")

	return &types.ActionResult{Success: true}, nil
}
