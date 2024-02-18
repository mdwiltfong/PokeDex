package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mdwiltfong/PokeDex/internal/utils"
)

type Config struct {
	PREV_URL *string
	NEXT_URL *string
}

func StartRepl(cfg *Config) {

	scanner := bufio.NewScanner(os.Stdin)
	cliMap := utils.CliCommandMap()
	fmt.Print("PokeDex > ")
	config := &Config{
		PREV_URL: nil,
		NEXT_URL: nil,
	}
	for scanner.Scan() {
		input := scanner.Text()
		sanitizedInput := utils.SanitizeInput(input)
		command, exists := cliMap[sanitizedInput]

		if exists {
			command.Callback(config)
		} else {
			fmt.Println("Hmm, this command doesn't exist. Try again")
		}
		if sanitizedInput == "exit" {
			return
		}
		fmt.Print("PokeDex > ")
	}
}
