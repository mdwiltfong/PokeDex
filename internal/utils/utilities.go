package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/mdwiltfong/PokeDex/internal/pokeapiclient"
)

type Config struct {
	PREV_URL *string
	NEXT_URL *string
}

func SanitizeInput(input string) string {
	output := strings.TrimSpace(input)
	return strings.ToLower(output)
}

type CallbackResponse interface {
	Response() string
	Print()
}

type HelpCommandResponse struct {
	CliCommandMapType
}

func (h HelpCommandResponse) Response() string {
	return ""
}
func (h HelpCommandResponse) Print() {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for key, value := range h.CliCommandMapType {
		fmt.Println(key + ": " + value.Description)
	}
	fmt.Println("")
}

type ExitCommandResponse struct {
	Message string
}

func (h ExitCommandResponse) Response() string {
	return h.Message
}

func (h ExitCommandResponse) Print() {
	fmt.Println("Okay! See you next time!")
}

type MapCommandResponse struct {
	Locations []Location
}

func (h MapCommandResponse) Response() string {
	return ""
}
func (h MapCommandResponse) Print() {
	for loc := range h.Locations {
		fmt.Println(loc)
	}
}

type CallbackFunction func(*Config, *pokeapiclient.Client) (CallbackResponse, error)

type CliCommand struct {
	Name        string
	Description string
	Callback    CallbackFunction
}

type CliCommandMapType map[string]CliCommand

func CliCommandMap() CliCommandMapType {

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
func helpCommand(*Config, *pokeapiclient.Client) (CallbackResponse, error) {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for key, value := range CliCommandMap() {
		fmt.Println(key + ": " + value.Description)
	}
	fmt.Println("")
	return HelpCommandResponse{CliCommandMap()}, nil
}
func exitCommand(*Config, *pokeapiclient.Client) (CallbackResponse, error) {
	return ExitCommandResponse{"Okay! See you next time!"}, nil
}

type Location struct {
	Name string
	URL  string
}
type GetLocationsResponse struct {
	Count    int
	Next     string
	Previous string
	Results  []Location
}

func Map(config *Config, client *pokeapiclient.Client) (CallbackResponse, error) {
	url := "https://pokeapi.co/api/v2/location/"
	if config.NEXT_URL != nil {
		fmt.Println("Next URL is not nil")
		url = *config.NEXT_URL
	}
	response, err := client.HttpClient.Get(url)

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
	config.PREV_URL = &url
	client.Cache.Add(url, responseBytes)
	if marshalingError != nil {
		log.Fatalf("Failed to unmarshal response: %s\n", marshalingError)
	}
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return MapCommandResponse{locations.Results}, nil
}

func Mapb(config *Config, client *pokeapiclient.Client) ([]Location, error) {

	if config.PREV_URL == nil || *config.PREV_URL == "" {
		fmt.Println("There are no previous pages")
		return nil, errors.New("There are no previous pages")
	}

	url := *config.PREV_URL
	cachedBytes, exists := client.Cache.Get(url)

	if !exists {
		fmt.Println("No cached data!!!!")

		response, err := client.HttpClient.Get(url)
		if err != nil {
			return nil, errors.New("there was an issue with the API request")
		}
		if response == nil {
			return nil, errors.New("There was an issue with the API response")
		}
		body, _ := io.ReadAll(response.Body)
		if response.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", response.StatusCode, body)
		}
		responseBytes := []byte(body)
		locations, _ := Unmarshall(responseBytes)
		config.NEXT_URL = &locations.Next
		config.PREV_URL = &locations.Previous

		return locations.Results, nil
	} else {
		fmt.Println("Cache Hit")

		Unmarshall(cachedBytes)
	}

	return nil, nil
}

func Unmarshall(val []byte) (GetLocationsResponse, error) {
	var locations GetLocationsResponse
	marshalingError := json.Unmarshal(val, &locations)
	if marshalingError != nil {
		log.Fatalf("Failed to unmarshal response: %s\n", marshalingError)
	}
	return locations, nil
}
