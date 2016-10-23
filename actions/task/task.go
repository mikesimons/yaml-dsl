package task

import (
    "container/list"
    "fmt"
    "github.com/mikesimons/yaml-dsl/dsl"
    "github.com/mitchellh/mapstructure"
    "github.com/pkg/errors"
)

type TaskAction struct {
    dsl.BaseAction `mapstructure:",squash"`
    Task           string `mapstructure:"task"`
}

func Type(data map[string]interface{}, dsl *dsl.Dsl) (string, dsl.Action, error) {
    result := &TaskAction{BaseAction: dsl.NewAction()}
    err := mapstructure.WeakDecode(data, result)
    if err != nil {
        return "", nil, errors.Wrap(err, fmt.Sprintf("task action could not be parsed from %#v", data))
    }

    return "task", result, nil
}

func (a *TaskAction) String() string {
    if a.Name != "" {
        return a.Name
    } else {
        return fmt.Sprintf("task: %s", a.Task)
    }
}

func (a *TaskAction) ModifyList(list *dsl.ActionList, current *list.Element) error {
    other := a.State().GetTask(a.Task)
    if other == nil {
        return errors.New(fmt.Sprintf("Invalid task: %s", a.Task))
    }

    list.InsertListAfter(other, current)

    return nil
}
