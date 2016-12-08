package parser

import (
	"container/list"

	"github.com/mikesimons/yaml-dsl/parser/middleware"
	"github.com/mikesimons/yaml-dsl/types"
)

type CommonAction struct {
	Name     string `mapstructure:"name"`
	Register string `mapstructure:"register"`
	When     string `mapstructure:"when"`
	Unless   string `mapstructure:"unless"`
}

type ActionList struct {
	list.List
	Dsl         *Dsl
	Middlewares *middleware.Chain
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
		action := current.Value.(types.Action)
		list.Middlewares.Next(action, vars)
		list.Middlewares.Reset()
	}

	return nil
}
