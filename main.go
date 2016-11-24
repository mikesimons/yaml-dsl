package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	shellAction "github.com/mikesimons/yaml-dsl/actions/shell"
	testAction "github.com/mikesimons/yaml-dsl/actions/test"
	"github.com/mikesimons/yaml-dsl/dsl"
	scripting "github.com/mikesimons/yaml-dsl/scripting"
	"github.com/mikesimons/yaml-dsl/types"
	"gopkg.in/yaml.v2"
)

/**
This is an interpreter / execution machine for a YAML DSL in the style of Ansible.

The design (when scripting is implemented) allows for highly dynamic features such as:

```yaml
- shell: ls /includes/*.yml
  register: some_var

- include: {{ item }}
  with_items: {{ some_var.stdout.split("\n") }}
```

Dependencies will be implemented as includes? How to name them?
**/

type Config map[string]Task

type Task struct {
	RawActions types.RawActionList `yaml:"actions"`
}

func main() {
	config := make(Config)

	f, _ := os.Open("config.yml")
	bytes, _ := ioutil.ReadAll(f)
	yaml.Unmarshal(bytes, &config)

	engine := dsl.New()
	engine.ScriptParser = scripting.NewMrubyScriptParser()
	engine.Handlers["shell"] = shellAction.Execute
	engine.Handlers["test"] = testAction.Execute

	tasks := make(map[string]*dsl.ActionList)
	for name, task := range config {
		actions, err := engine.ProcessRawActions(&task.RawActions)
		if err != nil {
			log.Fatalf("%#v", err)
		}
		tasks[name] = actions
	}

	tasks["install-docker"].Execute()
}

func toString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
