package shell

import (
	"fmt"

	"github.com/mikesimons/yaml-dsl/dsl"
	"github.com/mikesimons/yaml-dsl/types"
	"github.com/mitchellh/mapstructure"
)

type ShellAction struct {
	dsl.CommonAction `mapstructure:",squash"`
	Command          string `mapstructure:"shell"`
}

func Execute(raw types.RawAction, dsl *dsl.Dsl) error {
	action := &ShellAction{}
	dsl.Decode(raw, func(config *mapstructure.DecoderConfig) {
		config.Result = &action
	})

	fmt.Printf("%#v\n", action)

	return nil
}
