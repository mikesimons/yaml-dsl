package mrubyparser

import (
	"encoding/json"
	"fmt"
	"github.com/mikesimons/yaml-dsl/types"
	"github.com/mitchellh/go-mruby"
)

type MrubyScriptParser struct {
	mrb  *mruby.Mrb
	vars *mruby.MrbValue
}

func New() types.ScriptParser {
	mrb := mruby.NewMrb()
	vs, err := mrb.LoadString(`
		o = OpenStruct.new
		o.instance_eval do
		  def set_from_json k, v
			self[k] = YAML.load(v)
		  end
		end
		o
	`)

	if err != nil {
		panic(fmt.Sprintf("Error initializing var store: %s", err))
	}

	return &MrubyScriptParser{
		mrb:  mrb,
		vars: vs,
	}
}

func (msp *MrubyScriptParser) Vars() {
	msp.eval(`Proc.new { |val, vars| puts vars }`, `puts vars`)
}

func (msp *MrubyScriptParser) SetVar(key string, val interface{}) {
	data, err := json.Marshal(val)
	if err != nil {
		panic(fmt.Sprintf("Oh noes, JSON marshal error: %s", err))
	}
	msp.vars.Call("set_from_json", msp.mrb.StringValue(key), msp.mrb.StringValue(string(data)))
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
