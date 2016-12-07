package dsl

import (
	"reflect"

	"github.com/mikesimons/yaml-dsl/types"
	"github.com/mitchellh/go-mruby"
	//"fmt"
)

type MrubyScriptParser struct {
	mrb  *mruby.Mrb
	vars *mruby.MrbValue
}

func NewMrubyScriptParser() types.ScriptParser {
	mrb := mruby.NewMrb()
	vs, _ := mrb.LoadString(`OpenStruct.new`)

	return &MrubyScriptParser{
		mrb:  mrb,
		vars: vs,
	}
}

func (msp *MrubyScriptParser) SetVars(vars map[string]interface{}) {

	for key, val := range vars {
		valValue := reflect.ValueOf(val)
		switch valValue.Kind() {
		case reflect.String:
			msp.vars.Call("[]=", msp.mrb.StringValue(key), msp.mrb.StringValue(val.(string)))
		}
	}
}

func (msp *MrubyScriptParser) ParseList(script string) (interface{}, error) {
	proc, err := msp.mrb.LoadString(`Proc.new { |val, vars| vars.instance_eval(val) }`)
	if err != nil {
		panic(err.Error())
	}

	result, err := proc.Call("call", msp.mrb.StringValue(script), msp.vars)
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
	proc, err := msp.mrb.LoadString(`Proc.new { |val, vars| vars.instance_eval("\"#{val}\"") }`)

	if err != nil {
		panic(err.Error())
	}

	result, err := proc.Call("call", msp.mrb.StringValue(script), msp.vars)
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
