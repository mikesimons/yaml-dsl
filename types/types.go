package types

// UnparsedAction is a type alias for an individual raw action that needs to be parsed
type UnparsedAction map[string]interface{}

// UnparsedActionList is a type alias for a collection of raw actions needing to be parsed
type UnparsedActionList []UnparsedAction

// ScriptParser is a simple interface that script engines need to support to enable embedded scripting
type ScriptParser interface {
	ParseEmbedded(script string) (interface{}, error)
	ParseExpression(script string) (interface{}, error)
	SetVars(vars map[string]interface{})
	SetVar(key string, val interface{})
}

// ActionResult is the common return value type for actions
type ActionResult struct {
	Success bool
	Result  map[string]interface{}
}

// HandlerPrototypeFunc is a type alias for the a function that should return new action instances
type HandlerPrototypeFunc func() Handler

// Handler is an interface that all actions should implement
type Handler interface {
	Execute() (*ActionResult, error)
}

// Action is the type used to pair unparsed action data with a handler for execution
type Action struct {
	Handler HandlerPrototypeFunc
	Data    UnparsedAction
}

// ActionCommon is a type that all actions should embed in order to recieve core fields
type ActionCommon struct {
	Name string `mapstructure:"name"`
}
