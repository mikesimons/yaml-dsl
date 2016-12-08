# yaml-dsl

This is a proof of concept project that implements parsing and execution of an Ansible-like imperative DSL in golang.

Only the action list parsing is implemented to enable consumers to embed it within a larger structure.

Custom actions and middlewares can be implemented on a per project basis although some potentially useful ones are available in the `actions` and `parser/middleware` directories respectively.

## Usage
There are three steps to using it...

### Initialisation
Initialisation is where you configure the parser with script engine, middleware and registered handlers.

```go
p := parser.New()
p.ScriptParser = scripting.NewMrubyScriptParser()
p.Middleware = []types.Middleware{
  &withitems.Middleware{Dsl: parser},
}
p.Handlers["test"] = testAction.Prototype
```

### Parse
Parsing is where we transform the input data to an executable object. The result is an `ActionList`.

```go
file, _ := os.Open(os.Args[1])
data, _ := ioutil.ReadAll(file)
var raw *types.UnparsedActionList
yaml.Unmarshal(data, &raw)

actionList, _ := p.Parse(raw)
```

### Execute
Finally, you can execute the action list. This is where each step gets executed.

``` go
actionList.Execute()
```

It should be noted that with the possible exception of middleware fields, all script execution takes place just before the action is executed.
The action itself does not need to (and at present can't) handle script execution.

This allows* actions to reference variables that come from another action earlier in the action list. (*variable registration middleware currently not implemented)
  
## Actions

### test
The test action simply dumps itself through `fmt.Printf`. This is useful for testing middleware and generally confirming that you have everything initialized correctly.
```
- name: "Example test action"
  test: "This is a test"
```

### shell
The shell action will execute a shell command locally. This may be expanded to be able to execute remotely at some point in the future.
```
- name: "Example shell action"
  shell: "ls -l /"
```

## Middleware
### withitems
The withitems middleware behaves like Ansible `with_items`. It allows you to loop over a list of data and execute the action for each.
Access to the current item is available through the `item` variable.
`with_items` may either be a yaml list or a script expression in the form of a string.

```
- name: "Example with_items middleware using yaml list"
  test: "Test item #{item}"
  with_items:
    - item1
    - item2
```

```
- name: "Example with_items middleware using scripting expression"
  test: "Test item #{item}"
  with_items: >
    "item1,item2,item3".split(",")
```
