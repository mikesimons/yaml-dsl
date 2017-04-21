package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	commandAction "github.com/mikesimons/yaml-dsl/actions/command"
	shellAction "github.com/mikesimons/yaml-dsl/actions/shell"
	testAction "github.com/mikesimons/yaml-dsl/actions/test"
	"github.com/mikesimons/yaml-dsl/middleware"
	"github.com/mikesimons/yaml-dsl/middleware/register"
	"github.com/mikesimons/yaml-dsl/middleware/withitems"
	parserpkg "github.com/mikesimons/yaml-dsl/parser"
	"github.com/mikesimons/yaml-dsl/scripting/mrubyparser"
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
	UnparsedActions types.UnparsedActionList `yaml:"actions"`
}

func main() {
	config := make(Config)

	f, _ := os.Open("config.yml")
	bytes, _ := ioutil.ReadAll(f)
	yaml.Unmarshal(bytes, &config)

	tasks := make(map[string]*parserpkg.ActionList)

	parser := parserpkg.New()
	parser.ScriptParser = mrubyparser.New()
	parser.Middleware = []middleware.Middleware{
		&withitems.Middleware{Dsl: parser},
		&register.Middleware{Dsl: parser},
	}
	parser.Handlers["command"] = commandAction.Prototype
	parser.Handlers["shell"] = shellAction.Prototype
	parser.Handlers["test"] = testAction.Prototype

	for name, task := range config {
		actions, err := parser.Parse(&task.UnparsedActions)
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
