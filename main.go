package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mdwiltfong/commands"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cliMap := commands.CliCommandMap()
	response := cliMap["help"].name
	fmt.Printf("cliMap[\"help\"]: %v\n", cliMap["help"])
	for scanner.Scan() {
	}
}
