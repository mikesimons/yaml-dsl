package withitems

import (
	"fmt"
	"reflect"

	"github.com/mikesimons/yaml-dsl/middleware"
	"github.com/mikesimons/yaml-dsl/parser"
	"github.com/mikesimons/yaml-dsl/types"
	//"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// Middleware implements an iterating mechanism for actions
type Middleware struct {
	Dsl     *parser.Dsl
	Key     string
	ItemKey string
}

// Execute runs the given action for a list of items.
// The list is taken from the `with_items` field of the action data.
// `with_items` may also be a script expression that returns a list.
// As the list is iterated, the current value will be set in the variable `item`.
//
// The key `with_items` may be changed by setting the `Key` field.
// The item var `item` may be changed by setting the `ItemKey` field.
func (wim *Middleware) Execute(action types.Action, vars map[string]interface{}, chain middleware.Chain) (*types.ActionResult, error) {

	// Key defaults
	middlewareKey := wim.Key
	if middlewareKey == "" {
		middlewareKey = "with_items"
	}

	itemKey := wim.ItemKey
	if itemKey == "" {
		itemKey = "item"
	}

	// Get the list
	var items []interface{}
	rawItems, hasItems := action.Data[middlewareKey]

	if !hasItems {
		return chain.Next(action, vars)
	}

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
		return nil, errors.New(fmt.Sprintf("invalid type for `%s` field", middlewareKey))
	}

	// Iterate the list
	var results []*types.ActionResult
	for _, v := range items {

		item, err := wim.Dsl.ScriptParser.ParseEmbedded(v.(string))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error parsing expression: %s", v.(string)))
		}
		wim.Dsl.ScriptParser.SetVar(itemKey, item)

		result, err := chain.Next(action, vars)

		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error executing action: %s", action))
		}
		results = append(results, result)
	}

	wim.Dsl.ScriptParser.SetVar(itemKey, "")

	return &types.ActionResult{Success: true}, nil
}
