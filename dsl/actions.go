package dsl

import (
    "container/list"
    "fmt"
    "github.com/pkg/errors"
)

type Action interface{}

type ListModifierAction interface {
    ModifyList(list *ActionList, current *list.Element) error
}

type ExecutableAction interface {
    Execute() error
}

type ActionList struct {
    list.List
}

type BaseAction struct {
    _state    *Dsl
    Name      string   `mapstructure:"name"`
    WithItems []string `mapstructure:"with_items"`
    Register  string   `mapstructure:"register"`
    When      string   `mapstructure:"when"`
    Unless    string   `mapstructure:"unless"`
}

func (a *BaseAction) State() *Dsl {
    return a._state
}

func (list *ActionList) InsertListAfter(otherList *ActionList, element *list.Element) {
    prev := element
    for current := otherList.Front(); current != nil; current = current.Next() {
        prev = list.InsertAfter(current.Value, prev)
    }
}

func (list *ActionList) Execute() error {
    for current := list.Front(); current != nil; current = current.Next() {
        action := current.Value.(Action)
        //fmt.Printf("%s\n", action)
        fmt.Printf("%#v\n", action)

        if modifier, ok := current.Value.(ListModifierAction); ok {
            err := modifier.ModifyList(list, current)
            if err != nil {
                return errors.Wrap(err, fmt.Sprintf("error occured executing modifier: %s", modifier))
            }
        }

        if executable, ok := current.Value.(ExecutableAction); ok {
            err := executable.Execute()
            if err != nil {
                return errors.Wrap(err, fmt.Sprintf("error occured executing action: %s", executable))
            }
        }
    }

    return nil
}
