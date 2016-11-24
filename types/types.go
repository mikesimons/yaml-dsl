package types

type RawAction map[string]interface{}
type RawActionList []RawAction

type ScriptParser interface {
	Parse(script string) (interface{}, error)
	ParseList(script string) (interface{}, error)
	SetVars(vars map[string]interface{})
}
