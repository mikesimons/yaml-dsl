package task

import (
    "fmt"
    "github.com/mikesimons/yaml-dsl/dsl"
    "github.com/mitchellh/mapstructure"
    "github.com/pkg/errors"
)

type ShellAction struct {
    dsl.BaseAction `mapstructure:",squash"`
    Command        string `mapstructure:"shell"`
}

func Type(data map[string]interface{}, dsl *dsl.Dsl) (string, dsl.Action, error) {
    result := &ShellAction{BaseAction: dsl.NewAction()}
    err := mapstructure.WeakDecode(data, &result)

    if err != nil {
        return "", nil, errors.Wrap(err, fmt.Sprintf("shell action could not be parsed from %#v", data))
    }

    return "shell", result, nil
}

func (a *ShellAction) String() string {
    if a.Name != "" {
        return a.Name
    } else {
        return fmt.Sprintf("shell: %s", a.Command)
    }
}

func (a *ShellAction) Execute() error {
    return nil
}
