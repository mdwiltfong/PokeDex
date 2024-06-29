package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/mdwiltfong/PokeDex/internal/pokeapiclient"
)

func StartRepl() {
	client := pokeapiclient.NewClient(50000, 100000)
	cfg := &Config{
		PREV_URL: nil,
		NEXT_URL: nil,
		Client:   client,
	}
	scanner := bufio.NewScanner(os.Stdin)
	cliMap := CliCommandMap()

	fmt.Print("PokeDex > ")

	for scanner.Scan() {
		input := scanner.Text()

		sanitizedInput := SanitizeInput(input)
		command, exists := cliMap[sanitizedInput]
		if exists {
			response, error := command.Callback(cfg)
			if error != nil {
				log.Fatalf(error.Error())
				return
			}
			response.Print()
		} else {
			fmt.Println("Hmm, this command doesn't exist. Try again")
		}
		if sanitizedInput == "exit" {
			return
		}
		fmt.Print("PokeDex > ")
	}
}
