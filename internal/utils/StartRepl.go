package utils

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/mdwiltfong/PokeDex/internal/pokeapiclient"
	"github.com/mdwiltfong/PokeDex/internal/types"
)

func StartRepl() {
	client := pokeapiclient.NewClient(50000, 5*time.Second)
	cfg := &types.Config{
		PREV_URL: nil,
		NEXT_URL: nil,
		Client:   client,
		Pokedex:  types.Pokedex{},
	}
	scanner := bufio.NewScanner(os.Stdin)
	cliMap := CliCommandMap()

	fmt.Print("PokeDex > ")

	for scanner.Scan() {
		input := scanner.Text()

		sanitizedInput := SanitizeInput(input)
		command, exists := cliMap[sanitizedInput[0]]
		if exists {
			var response types.CallbackResponse
			var err error
			if len(sanitizedInput) == 2 {
				response, err = command.Callback(cfg, sanitizedInput[1])
			} else {
				response, err = command.Callback(cfg, "")
			}
			if err != nil {
				fmt.Println(err.Error())
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
