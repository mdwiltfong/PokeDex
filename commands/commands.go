package commands

import (
	"fmt"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func() error
}

func CliCommandMap() map[string]CliCommand {

	return map[string]CliCommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    helpCommand,
		},
		"exit": {
			Name:        "exit",
			Description: "Exits the REPL",
			Callback:    exitCommand,
		},
	}

}
func helpCommand() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for key, value := range CliCommandMap() {
		fmt.Println(key + ": " + value.Description)
	}
	fmt.Println("")
	return nil
}
func exitCommand() error {
	fmt.Println("Okay! See you next time!")
	return nil
}
