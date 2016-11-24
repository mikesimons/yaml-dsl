package test

import (
	"fmt"

	"github.com/mikesimons/yaml-dsl/dsl"
	"github.com/mikesimons/yaml-dsl/types"
	"github.com/mitchellh/mapstructure"
)

type TestAction struct {
	dsl.CommonAction `mapstructure:",squash"`
	Test             string `mapstructure:"test"`
}

func Execute(raw types.RawAction, dsl *dsl.Dsl) error {
	action := &TestAction{}
	dsl.Decode(raw, func(config *mapstructure.DecoderConfig) {
		config.Result = &action
	})

	fmt.Printf("%#v\n", action)

	return nil
}
