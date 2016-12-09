package testparser

import (
	"fmt"
	"strings"
)

type TestScriptParser struct {
	ParseExpressionError error
	ParseEmbeddedError   error
	Vars                 map[string]interface{}
}

func New() *TestScriptParser {
	return &TestScriptParser{
		Vars: make(map[string]interface{}),
	}
}

func (parser *TestScriptParser) SetVar(key string, val interface{}) {
	parser.Vars[key] = val
}

func (parser *TestScriptParser) SetVars(vars map[string]interface{}) {
	for key, val := range vars {
		parser.SetVar(key, val)
	}
}

func (parser *TestScriptParser) ParseExpression(script string) (interface{}, error) {
	var result []interface{}
	for _, str := range strings.Split(script, ",") {
		result = append(result, interface{}(str))
	}
	return result, parser.ParseExpressionError
}

func (parser *TestScriptParser) ParseEmbedded(script string) (interface{}, error) {
	result := script
	for k, v := range parser.Vars {
		result = strings.Replace(result, fmt.Sprintf("{{%s}}", k), v.(string), -1)
	}
	return result, parser.ParseEmbeddedError
}
