package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	PREV_URL *string
	NEXT_URL *string
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	cliMap := CliCommandMap()
	fmt.Print("PokeDex > ")
	for scanner.Scan() {
		input := scanner.Text()
		sanitizedInput := sanitizeInput(input)
		command, exists := cliMap[sanitizedInput]
		config := &Config{
			PREV_URL: nil,
			NEXT_URL: nil,
		}

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

func sanitizeInput(input string) string {
	output := strings.TrimSpace(input)
	output = strings.ToLower(input)
	return output
}

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*Config) error
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
		"map": {
			Name:        "map",
			Description: "Sends a get request of maps in the pokemon game",
			Callback:    Map,
		},
	}

}
func helpCommand(*Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for key, value := range CliCommandMap() {
		fmt.Println(key + ": " + value.Description)
	}
	fmt.Println("")
	return nil
}
func exitCommand(*Config) error {
	fmt.Println("Okay! See you next time!")
	return nil
}

type GetLocationsResponse struct {
	Count    int
	Next     string
	Previous string
	Results  []struct {
		Name string
		URL  string
	}
}

func Map(config *Config) error {
	url := "https://pokeapi.co/api/v2/location/"
	if config.NEXT_URL != nil {
		fmt.Println("Next URL is not nil")
		url = *config.NEXT_URL
	}
	response, err := http.Get(url)

	if err != nil {
		errors.New("There was an issue with the API request")
	}
	body, _ := io.ReadAll(response.Body)
	if response.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", response.StatusCode, body)
	}
	responseBytes := []byte(body)
	var locations GetLocationsResponse
	marshalingError := json.Unmarshal(responseBytes, &locations)
	config.NEXT_URL = &locations.Next
	fmt.Println("Next URL: ", *config.NEXT_URL)
	if marshalingError != nil {
		log.Fatalf("Failed to unmarshal response: %s\n", marshalingError)
	}
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func Mapb(config *Config) error {

	url := *config.PREV_URL

	response, err := http.Get(url)
	if err != nil {
		errors.New("There was an issue with the API request")
	}
	body, _ := io.ReadAll(response.Body)
	if response.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", response.StatusCode, body)
	}
	responseBytes := []byte(body)
	var locations GetLocationsResponse
	marshalingError := json.Unmarshal(responseBytes, &locations)
	if marshalingError != nil {
		log.Fatalf("Failed to unmarshal response: %s\n", marshalingError)
	}
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}
