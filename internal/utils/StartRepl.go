package utils

import (
	"bufio"
	"fmt"
	"os"
)

func StartRepl(cfg *Config) {

	scanner := bufio.NewScanner(os.Stdin)
	cliMap := CliCommandMap()

	fmt.Print("PokeDex > ")

	for scanner.Scan() {
		input := scanner.Text()

		sanitizedInput := SanitizeInput(input)
		fmt.Println(sanitizedInput)
		command, exists := cliMap[sanitizedInput]
		fmt.Println(command, exists)
		if exists {
			command.Callback(cfg)
		} else {
			fmt.Println("Hmm, this command doesn't exist. Try again")
		}
		if sanitizedInput == "exit" {
			return
		}
		fmt.Print("PokeDex > ")
	}
}
