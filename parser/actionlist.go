package parser

import (
	"container/list"

	"github.com/mikesimons/yaml-dsl/middleware"
	"github.com/mikesimons/yaml-dsl/types"
)

type ActionList struct {
	list.List
	Dsl        *Dsl
	Middleware middleware.Chain
}

func (list *ActionList) InsertListAfter(otherList *ActionList, element *list.Element) {
	prev := element
	for current := otherList.Front(); current != nil; current = current.Next() {
		prev = list.InsertAfter(current.Value, prev)
	}
}

func (list *ActionList) Execute() error {
	vars := make(map[string]interface{})

	var results []*types.ActionResult
	for current := list.Front(); current != nil; current = current.Next() {
		action := current.Value.(types.Action)
		result, _ := list.Middleware.Next(action, vars)
		results = append(results, result)
		list.Middleware.Reset()
	}

	return nil
}
