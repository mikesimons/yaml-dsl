package dsl

import (
	"github.com/mikesimons/yaml-dsl/types"
	"github.com/mitchellh/go-mruby"
	//"fmt"
)

type MrubyScriptParser struct {
	mrb  *mruby.Mrb
	vars map[string]interface{}
}

func NewMrubyScriptParser() types.ScriptParser {
	return &MrubyScriptParser{
		mrb:  mruby.NewMrb(),
		vars: make(map[string]interface{}),
	}
}

func (msp *MrubyScriptParser) SetVars(vars map[string]interface{}) {
	msp.vars = vars
}

func (msp *MrubyScriptParser) ParseList(script string) (interface{}, error) {
	proc, err := msp.mrb.LoadString(`Proc.new { |val| eval(val) }`)
	if err != nil {
		panic(err.Error())
	}

	result, err := proc.Call("call", msp.mrb.StringValue(script))
	if err != nil {
		panic(err.Error())
	}

	var ret interface{}
	err = mruby.Decode(&ret, result)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (msp *MrubyScriptParser) Parse(script string) (interface{}, error) {
	proc, err := msp.mrb.LoadString(`Proc.new { |val| eval("\"#{val}\"") }`)

	if err != nil {
		panic(err.Error())
	}

	result, err := proc.Call("call", msp.mrb.StringValue(script))
	if err != nil {
		panic(err.Error())
	}

	var ret interface{}
	err = mruby.Decode(&ret, result)

	if err != nil {
		panic(err.Error())
		return nil, err
	}

	return ret, nil
}
