package commands

import (
	"errors"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func CliCommandMap() map[string]cliCommand {

	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    helpCommand,
		},
	}

}
func helpCommand() error {
	return errors.New("THere is no help")
}
