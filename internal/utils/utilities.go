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

func SanitizeInput(input string) []string {
	output := strings.TrimSpace(input)
	lowerCase := strings.ToLower(output)
	return strings.Split(lowerCase, " ")
}

type CallbackResponse interface {
	Response() interface{}
	Print()
}
type ExploreCommandResponse struct {
	Encounters []PokemonEncounter
}

func (h ExploreCommandResponse) Response() interface{} {
	return h.Encounters
}
func (h ExploreCommandResponse) Print() {
	for _, encounter := range h.Encounters {
		fmt.Println(encounter.Pokemon.Name)
	}
}

type HelpCommandResponse struct {
	CliCommandMapType
}

func (h HelpCommandResponse) Response() interface{} {
	return h.CliCommandMapType
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

func (h ExitCommandResponse) Response() interface{} {
	return h.Message
}

func (h ExitCommandResponse) Print() {
	fmt.Println("Okay! See you next time!")
}

type MapCommandResponse struct {
	Locations []Location
}

func (h MapCommandResponse) Response() interface{} {
	return h.Locations
}
func (h MapCommandResponse) Print() {
	for _, loc := range h.Locations {
		fmt.Println(loc.Name)
	}
}

type CallbackFunction func(*Config, *pokeapiclient.Client, string) (CallbackResponse, error)

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
			Callback:    HelpCommand,
		},
		"exit": {
			Name:        "exit",
			Description: "Exits the REPL",
			Callback:    ExitCommand,
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
		"explore": {
			Name:        "explore",
			Description: "Explore the possible pokemon encounters in an area",
			Callback:    Explore,
		},
	}

}
func HelpCommand(*Config, *pokeapiclient.Client, string) (CallbackResponse, error) {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for key, value := range CliCommandMap() {
		fmt.Println(key + ": " + value.Description)
	}
	fmt.Println("")
	return HelpCommandResponse{CliCommandMap()}, nil
}
func ExitCommand(*Config, *pokeapiclient.Client, string) (CallbackResponse, error) {
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
type PokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon"`
	VersionDetails []struct {
		EncounterDetails []struct {
			Chance          int   `json:"chance"`
			ConditionValues []any `json:"condition_values"`
			MaxLevel        int   `json:"max_level"`
			Method          struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"method"`
			MinLevel int `json:"min_level"`
		} `json:"encounter_details"`
		MaxChance int `json:"max_chance"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"version_details"`
}
type PokemonEncountersResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

func Map(config *Config, client *pokeapiclient.Client, commandInput string) (CallbackResponse, error) {
	url := "https://pokeapi.co/api/v2/location/"
	if config.NEXT_URL != nil {
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
	marshalingError := Unmarshall(responseBytes, &locations)

	config.NEXT_URL = &locations.Next
	config.PREV_URL = &url
	client.Cache.Add(url, responseBytes)
	if marshalingError != nil {
		log.Fatalf("Failed to unmarshal response: %s\n", marshalingError)
	}
	return MapCommandResponse{locations.Results}, nil
}

func Mapb(config *Config, client *pokeapiclient.Client, commandInput string) (CallbackResponse, error) {

	if config.PREV_URL == nil || *config.PREV_URL == "" {
		fmt.Println("There are no previous pages")
		return MapCommandResponse{}, errors.New("there are no previous pages")
	}

	url := *config.PREV_URL
	cachedBytes, exists := client.Cache.Get(url)
	var locations GetLocationsResponse
	if !exists {
		fmt.Println("No cached data!!!!")

		response, err := client.HttpClient.Get(url)
		if err != nil {
			return MapCommandResponse{}, errors.New("there was an issue with the API request")
		}
		if response == nil {
			return MapCommandResponse{}, errors.New("There was an issue with the API response")
		}
		body, _ := io.ReadAll(response.Body)
		if response.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", response.StatusCode, body)
		}
		responseBytes := []byte(body)
		error := Unmarshall(responseBytes, &locations)
		if error != nil {
			log.Fatalf("Failed to unmarshal response: %s\n", error)
		}
		config.NEXT_URL = &locations.Next
		config.PREV_URL = &locations.Previous

		return MapCommandResponse{locations.Results}, nil
	} else {
		fmt.Println("Cache Hit")

		error := Unmarshall(cachedBytes, &locations)
		if error != nil {
			log.Fatalf("Failed to unmarshal response: %s\n", error)
		}
		return MapCommandResponse{locations.Results}, nil
	}

}

func Explore(config *Config, client *pokeapiclient.Client, commandInput string) (CallbackResponse, error) {
	if commandInput == "" {
		errors.New("Please put in a location to explroe")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", commandInput)
	response, err := client.HttpClient.Get(url)
	if err != nil {
		errors.New(err.Error())
	}
	body, _ := io.ReadAll(response.Body)
	if response.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", response.StatusCode, body)
	}
	responseBytes := []byte(body)
	var encounter PokemonEncountersResponse
	unMarshallError := Unmarshall[PokemonEncountersResponse](responseBytes, &encounter)
	if unMarshallError != nil {
		log.Fatalf("Failed to unmarshal response: %s\n", unMarshallError)
	}
	return ExploreCommandResponse{encounter.PokemonEncounters}, nil
}

func Unmarshall[T GetLocationsResponse | PokemonEncountersResponse](val []byte, v *T) error {

	unmarshalError := json.Unmarshal(val, &v)
	if unmarshalError != nil {
		log.Fatalf("Failed to unmarshal response: %s\n", unmarshalError)
	}

	return nil
}
