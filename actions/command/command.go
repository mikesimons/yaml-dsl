package command

import (
	"bytes"
	"os/exec"

	"github.com/mikesimons/yaml-dsl/types"
	"github.com/pkg/errors"
)

type CommandAction struct {
	types.ActionCommon `mapstructure:",squash"`
	Command            string            `mapstructure:"command"`
	Env                map[string]string `mapstructure:"env"`
	WorkDir            string            `mapstructure:"workdir"`
}

func Prototype() types.Handler {
	return &CommandAction{}
}

func (action *CommandAction) Execute() (*types.ActionResult, error) {
	var outbuf bytes.Buffer
	var errbuf bytes.Buffer

	executable, args := parseCmd(action.Command)
	cmd := exec.Command(executable, args...)

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	if err := cmd.Start(); err != nil {
		// Couldn't even start the command
		return &types.ActionResult{
			Success: false,
			Action:  action,
		}, errors.Wrap(err, "failed to run command")
	}

	if err := cmd.Wait(); err != nil {
		// Started but command failed in some way
		return &types.ActionResult{
			Success: false,
			Action:  action,
			Result: map[string]interface{}{
				"stdout": outbuf.String(),
				"stderr": errbuf.String(),
			},
		}, errors.Wrap(err, "command exited with error")
	}

	return &types.ActionResult{
		Success: true,
		Action:  action,
		Result: map[string]interface{}{
			"stdout": outbuf.String(),
			"stderr": errbuf.String(),
		},
	}, nil
}

func parseCmd(input string) (string, []string) {
	inDoubleQuotes := false
	inSingleQuotes := false
	escapeMode := false

	var current []rune
	var args []string

	for _, ch := range input {
		if ch == '\\' {
			escapeMode = true
			continue
		}

		if escapeMode {
			current = append(current, ch)
			escapeMode = false
			continue
		}

		if ch == '"' {
			inDoubleQuotes = !inDoubleQuotes
			continue
		}

		if ch == '\'' {
			inSingleQuotes = !inSingleQuotes
			continue
		}

		if !inSingleQuotes && !inDoubleQuotes && (ch == ' ' || ch == '\t' || ch == '\n') {
			if len(current) > 0 {
				args = append(args, string(current))
				current = make([]rune, 0)
			}
			continue
		}

		current = append(current, ch)
	}

	if len(current) > 0 {
		args = append(args, string(current))
	}

	if len(args) > 0 {
		return args[0], args[1:]
	}

	return "", args
}
