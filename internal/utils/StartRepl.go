package utils

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mdwiltfong/PokeDex/internal/pokeapiclient"
)

func StartRepl(cfg *Config, client *pokeapiclient.Client) {

	scanner := bufio.NewScanner(os.Stdin)
	cliMap := CliCommandMap()

	fmt.Print("PokeDex > ")

	for scanner.Scan() {
		input := scanner.Text()

		sanitizedInput := SanitizeInput(input)
		command, exists := cliMap[sanitizedInput]
		if exists {
			command.Callback(cfg, client)
		} else {
			fmt.Println("Hmm, this command doesn't exist. Try again")
		}
		if sanitizedInput == "exit" {
			return
		}
		fmt.Print("PokeDex > ")
	}
}
