package mrubyparser

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

func New() types.ScriptParser {
	mrb := mruby.NewMrb()
	vs, _ := mrb.LoadString(`OpenStruct.new`)

	return &MrubyScriptParser{
		mrb:  mrb,
		vars: vs,
	}
}

func (msp *MrubyScriptParser) SetVar(key string, val interface{}) {
	valValue := reflect.ValueOf(val)
	switch valValue.Kind() {
	case reflect.String:
		msp.vars.Call("[]=", msp.mrb.StringValue(key), msp.mrb.StringValue(val.(string)))
	}
}

func (msp *MrubyScriptParser) SetVars(vars map[string]interface{}) {
	for key, val := range vars {
		msp.SetVar(key, val)
	}
}

func (msp *MrubyScriptParser) ParseExpression(script string) (interface{}, error) {
	return msp.eval(`Proc.new { |val, vars| vars.instance_eval(val) }`, script)
}

func (msp *MrubyScriptParser) ParseEmbedded(script string) (interface{}, error) {
	return msp.eval(`Proc.new { |val, vars| vars.instance_eval("\"#{val}\"") }`, script)
}

func (msp *MrubyScriptParser) eval(context string, eval string) (interface{}, error) {
	proc, err := msp.mrb.LoadString(context)
	if err != nil {
		panic(err.Error())
	}

	result, err := proc.Call("call", msp.mrb.StringValue(eval), msp.vars)
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
