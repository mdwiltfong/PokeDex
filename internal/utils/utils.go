package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Config struct {
	PREV_URL *string
	NEXT_URL *string
}

func SanitizeInput(input string) string {
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
		"mapb": {
			Name:        "mapb",
			Description: "Sends a get request of maps in the pokemon game",
			Callback:    Mapb,
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
	config.PREV_URL = &locations.Previous
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

	if config.PREV_URL == nil || *config.PREV_URL == "" {
		fmt.Println("There are no previous pages")
		return nil
	}
	url := *config.PREV_URL
	fmt.Println("Previous URL: ", *config.PREV_URL)
	fmt.Println("Next URL: ", *config.NEXT_URL)
	response, err := http.Get(url)
	if err != nil {
		return errors.New("there was an issue with the API request")
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
	config.NEXT_URL = &locations.Next
	config.PREV_URL = &locations.Previous
	return nil
}
