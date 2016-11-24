package dsl

import (
	"container/list"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type CommonAction struct {
	Name     string `mapstructure:"name"`
	Register string `mapstructure:"register"`
	When     string `mapstructure:"when"`
	Unless   string `mapstructure:"unless"`
}

type ActionList struct {
	list.List
	Dsl *Dsl
}

func (list *ActionList) InsertListAfter(otherList *ActionList, element *list.Element) {
	prev := element
	for current := otherList.Front(); current != nil; current = current.Next() {
		prev = list.InsertAfter(current.Value, prev)
	}
}

func (list *ActionList) Execute() error {
	vars := make(map[string]interface{})

	for current := list.Front(); current != nil; current = current.Next() {
		action := current.Value.(Action)
		fmt.Printf("%#v\n", action)

		items := list.itemsForAction(action)

		for _, v := range items {
			item, _ := list.Dsl.ScriptParser.Parse(v.(string))
			vars["item"] = item
			list.Dsl.ScriptParser.SetVars(vars)

			err := action.Handler(action.Data, list.Dsl)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error occured executing action: %s", action))
			}
		}
	}

	return nil
}

func (list *ActionList) itemsForAction(action Action) []interface{} {
	var items []interface{}
	if rawItems, hasItems := action.Data["with_items"]; hasItems {
		value := reflect.ValueOf(rawItems)
		switch value.Kind() {
		case reflect.String:
			result, _ := list.Dsl.ScriptParser.ParseList(rawItems.(string))
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

	return items
}
