package main

import (
	"bufio"
	"os"
	"fmt"
	"github.com/mdwiltfong/commands"
	"strings"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	cliMap := commands.CliCommandMap()
	fmt.Print("PokeDex > ")
	for scanner.Scan() {
		input := scanner.Text()
		sanitizedInput:=sanitizeInput(input)
		command, exists := cliMap[sanitizedInput]
		if exists {
			command.Callback()
		}else{
			fmt.Println("Hmm, this command doesn't exist. Try again")
		}
		if sanitizedInput == "exit"{
			return 
		}
	fmt.Print("PokeDex > ")
	}
}

func sanitizeInput(input string) string  {
	output:=strings.TrimSpace(input)
	output=strings.ToLower(input)
	return output
}