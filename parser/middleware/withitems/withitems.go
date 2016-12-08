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

func New(dsl *parser.Dsl) *Middleware {
	return &Middleware{
		Dsl:     dsl,
		Key:     "with_items",
		ItemKey: "item",
	}
}

func (wim *Middleware) Execute(action types.Action, vars map[string]interface{}) (*types.ActionResult, error) {
	var items []interface{}
	if rawItems, hasItems := action.Data["with_items"]; hasItems {
		value := reflect.ValueOf(rawItems)
		switch value.Kind() {
		case reflect.String:
			result, _ := wim.Dsl.ScriptParser.ParseExpression(rawItems.(string))
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

	for _, v := range items {
		item, _ := wim.Dsl.ScriptParser.ParseEmbedded(v.(string))
		vars["item"] = item
		wim.Dsl.ScriptParser.SetVars(vars)

		instance := action.Handler()
		wim.Dsl.Decode(action.Data, func(config *mapstructure.DecoderConfig) {
			config.Result = instance
		})

		_, err := instance.Execute()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error occured executing action: %s", action))
		}
	}

	wim.Dsl.ScriptParser.SetVar("item", nil)

	return &types.ActionResult{}, nil
}
