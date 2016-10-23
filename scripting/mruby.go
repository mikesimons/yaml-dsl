package dsl

import (
    "github.com/mikesimons/yaml-dsl/types"
    "github.com/mitchellh/go-mruby"
)

type MrubyScriptParser struct {
    mrb mruby.Mrb
}

func NewMrubyScriptParser() types.ScriptParser {
    return &MrubyScriptParser{
        mrb: mruby.NewMrb(),
    }.(types.ScriptParser)
}

func (msp *MrubyScriptParser) Parse(script string) interface{} {
    result, err := msp.mrb.LoadString(script)
    if err != nil {
        panic(err.Error())
    }

    fmt.Printf("Result: %s\n", result.String())

    return nil
}
