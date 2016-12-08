package types

type UnparsedAction map[string]interface{}
type UnparsedActionList []UnparsedAction

type ScriptParser interface {
	ParseEmbedded(script string) (interface{}, error)
	ParseExpression(script string) (interface{}, error)
	SetVars(vars map[string]interface{})
	SetVar(key string, val interface{})
}

type ActionResult struct {
	Success bool
	Result  map[string]interface{}
}

type HandlerPrototypeFunc func() Handler
type Handler interface {
	Execute() (*ActionResult, error)
}

type Middleware interface {
	Execute(action Action, vars map[string]interface{}) (*ActionResult, error)
}

type Action struct {
	Handler HandlerPrototypeFunc
	Data    UnparsedAction
}
