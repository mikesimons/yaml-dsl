package types

type RawActionList []map[string]interface{}
type ScriptParser interface {
    Parse(script string) interface{}
}
