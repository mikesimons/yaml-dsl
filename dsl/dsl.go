package dsl

import (
    "github.com/mikesimons/yaml-dsl/types"
    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"
    "io"
    "io/ioutil"
    "os"
)

type factoryFn func(data map[string]interface{}, dsl *Dsl) (string, Action, error)

type Dsl struct {
    factories    map[string]factoryFn
    tasks        map[string]*ActionList
    scriptParser *types.ScriptParser
}

func New() *Dsl {
    return &Dsl{
        factories: make(map[string]factoryFn),
        tasks:     make(map[string]*ActionList),
    }
}

func (dsl *Dsl) AddActionType(f factoryFn) {
    key, _, _ := f(make(map[string]interface{}), dsl)
    dsl.factories[key] = f
}

func (dsl *Dsl) AddTask(name string, actions *ActionList) {
    dsl.tasks[name] = actions
}

func (dsl *Dsl) GetTask(name string) *ActionList {
    return dsl.tasks[name]
}

func (dsl *Dsl) LoadActionsFromFile(file string) (*ActionList, error) {
    reader, err := os.Open(file)
    if err != nil {
        return nil, errors.Wrap(err, "Unable to open file")
    }

    return dsl.LoadActionsFromReader(reader)
}

func (dsl *Dsl) LoadActionsFromReader(in io.Reader) (*ActionList, error) {
    data, err := ioutil.ReadAll(in)

    if err != nil {
        return nil, errors.Wrap(err, "Unable to read file")
    }

    raw := make(types.RawActionList, 0)
    err = yaml.Unmarshal(data, &raw)

    if err != nil {
        return nil, errors.Wrap(err, "Unable to parse YAML")
    }

    return dsl.ProcessRawActions(&raw)
}

func (dsl *Dsl) ProcessRawActions(raw *types.RawActionList) (*ActionList, error) {
    out := &ActionList{}
    for _, rawAction := range *raw {
        var action Action
        var err error

        for key := range rawAction {
            if factoryFn, ok := dsl.factories[key]; ok {
                _, action, err = factoryFn(rawAction, dsl)
            }

            if err != nil {
                return nil, errors.Wrap(err, "failed to process raw actions")
            }
        }

        if action != nil {
            out.PushBack(action)
        }
    }
    return out, nil
}

func (dsl *Dsl) ScriptParser() *types.ScriptParser {
    return dsl.scriptParser
}

func (dsl *Dsl) NewAction() BaseAction {
    return BaseAction{
        _state: dsl,
    }
}
