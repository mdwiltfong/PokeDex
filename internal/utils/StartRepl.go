package utils

import (
	"bufio"
	"fmt"
	"log"
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
		command, exists := cliMap[sanitizedInput[0]]
		if exists {
			response, error := command.Callback(cfg, client, sanitizedInput[1])
			if error != nil {
				log.Fatalf(error.Error())
				return
			}
			response.Print()
		} else {
			fmt.Println("Hmm, this command doesn't exist. Try again")
		}
		if sanitizedInput[0] == "exit" {
			return
		}
		fmt.Print("PokeDex > ")
	}
}
