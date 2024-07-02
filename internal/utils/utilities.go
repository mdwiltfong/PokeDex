package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/mdwiltfong/PokeDex/internal/types"
)

func SanitizeInput(input string) []string {
	output := strings.TrimSpace(input)
	lowerCase := strings.ToLower(output)
	return strings.Split(lowerCase, " ")
}

func CliCommandMap() types.CliCommandMapType {

	return map[string]types.CliCommand{
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
		"catch": {
			Name:        "catch",
			Description: "Catch a pokemon",
			Callback:    Catch,
		},
	}

}

func HelpCommand(config *types.Config, dependency types.Dependency, commandInput string) (types.CallbackResponse, error) {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for key, value := range CliCommandMap() {
		fmt.Println(key + ": " + value.Description)
	}
	fmt.Println("")
	return types.HelpCommandResponse{CliCommandMapType: CliCommandMap()}, nil
}

func ExitCommand(config *types.Config, dep types.Dependency, commandInput string) (types.CallbackResponse, error) {
	return types.ExitCommandResponse{Message: "Okay! See you next time!"}, nil
}

func Map(config *types.Config, dep types.Dependency, commandInput string) (types.CallbackResponse, error) {
	url := "https://pokeapi.co/api/v2/location/"
	if config.NEXT_URL != nil {
		url = *config.NEXT_URL
	}
	response, err := config.Client.HttpClient.Get(url)

	if err != nil {
		errors.New("There was an issue with the API request")
	}
	body, _ := io.ReadAll(response.Body)
	if response.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", response.StatusCode, body)
	}
	responseBytes := []byte(body)
	var locations types.GetLocationsResponse
	marshalingError := Unmarshall(responseBytes, &locations)

	config.NEXT_URL = &locations.Next
	config.PREV_URL = &url
	config.Client.Cache.Add(url, responseBytes)
	if marshalingError != nil {
		log.Fatalf("Failed to unmarshal response: %s\n", marshalingError)
	}
	return types.MapCommandResponse{Locations: locations.Results}, nil
}

func Mapb(config *types.Config, dependency types.Dependency, commandInput string) (types.CallbackResponse, error) {

	if config.PREV_URL == nil || *config.PREV_URL == "" {
		fmt.Println("There are no previous pages")
		return types.MapCommandResponse{}, errors.New("there are no previous pages")
	}

	url := *config.PREV_URL
	var locations types.GetLocationsResponse
	cachedBytes, exists := config.Client.Cache.Get(url)
	if !exists {
		fmt.Println("No cached data!!!!")

		response, err := config.Client.HttpClient.Get(url)
		if err != nil {
			return types.MapCommandResponse{}, errors.New("there was an issue with the API request")
		}
		if response == nil {
			return types.MapCommandResponse{}, errors.New("There was an issue with the API response")
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

		return types.MapCommandResponse{Locations: locations.Results}, nil
	} else {
		fmt.Println("Cache Hit")

		marshalingError := Unmarshall(cachedBytes, &locations)
		if marshalingError != nil {
			log.Fatalf("Failed to unmarshal response: %s\n", marshalingError)
		}
		return types.MapCommandResponse{Locations: locations.Results}, nil
	}

}

func Explore(config *types.Config, dependency types.Dependency, commandInput string) (types.CallbackResponse, error) {
	if commandInput == "" {
		return types.ExploreCommandResponse{}, errors.New("Please put in a location to explore")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", commandInput)
	cachedBytes, exists := config.Client.Cache.Get(url)
	if !exists {
		response, err := config.Client.HttpClient.Get(url)
		if err != nil {
			return types.ExploreCommandResponse{}, errors.New("There was an issue retrieving the data:" + err.Error())
		}
		body, _ := io.ReadAll(response.Body)
		if response.StatusCode > 299 {
			if response.StatusCode == 404 {
				return types.ExploreCommandResponse{}, errors.New("Area was not found")
			}
			return types.ExploreCommandResponse{}, errors.New("There was an issue retrieving the data")
		}
		responseBytes := []byte(body)

		var encounter types.PokemonEncountersResponse
		unMarshallError := Unmarshall[types.PokemonEncountersResponse](responseBytes, &encounter)
		config.Client.Cache.Add(url, responseBytes)
		if unMarshallError != nil {
			log.Fatalf("Failed to unmarshal response: %s\n", unMarshallError)
			return types.ExploreCommandResponse{}, errors.New("There was an issue unmarshalling the data" + unMarshallError.Error())
		}
		return types.ExploreCommandResponse{Encounters: encounter.PokemonEncounters}, nil
	} else {
		fmt.Println("Cache hit")
		var encounter types.PokemonEncountersResponse
		unMarshallError := Unmarshall[types.PokemonEncountersResponse](cachedBytes, &encounter)
		if unMarshallError != nil {
			log.Fatalf("Failed to unmarshal response: %s\n", unMarshallError)
			return types.ExploreCommandResponse{}, errors.New("There was an issue unmarshalling the data" + unMarshallError.Error())
		}
		return types.ExploreCommandResponse{Encounters: encounter.PokemonEncounters}, nil
	}

}

func Catch(config *types.Config, dependency types.Dependency, commandInput string) (types.CallbackResponse, error) {
	if commandInput == "" {
		return types.ExploreCommandResponse{}, errors.New("Please enter a pokemon you'd like to catch")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", commandInput)
	cachedBytes, exists := config.Client.Cache.Get(url)
	var pokemonInformation types.PokemonInformation
	if !exists {
		response, err := config.Client.HttpClient.Get(url)
		if err != nil {
			return types.ExploreCommandResponse{}, errors.New("There was an issue retrieving the data:" + err.Error())
		}
		body, _ := io.ReadAll(response.Body)
		if response.StatusCode > 299 {
			if response.StatusCode == 404 {
				return types.ExploreCommandResponse{}, errors.New("Pokemon was not found")
			}
			return types.ExploreCommandResponse{}, errors.New("There was an issue retrieving the data")
		}
		responseBytes := []byte(body)

		unMarshallError := Unmarshall[types.PokemonInformation](responseBytes, &pokemonInformation)
		if unMarshallError != nil {
			log.Fatalf("Failed to unmarshal response: %s\n", unMarshallError)
			return types.ExploreCommandResponse{}, errors.New("There was an issue unmarshalling the data" + unMarshallError.Error())
		}
		config.Client.Cache.Add(url, responseBytes)

	} else {
		fmt.Println("Cache hit")
		unMarshallError := Unmarshall[types.PokemonInformation](cachedBytes, &pokemonInformation)
		if unMarshallError != nil {
			log.Fatalf("Failed to unmarshal response: %s\n", unMarshallError)
			return types.ExploreCommandResponse{}, errors.New("There was an issue unmarshalling the data" + unMarshallError.Error())
		}

	}
	randNum := dependency.RandInt(pokemonInformation.BaseExperience)
	chance := float64(randNum) / float64(pokemonInformation.BaseExperience)
	fmt.Printf("Chance of catching %s: %f\n", pokemonInformation.Name, chance)
	if chance > 0.5 {
		pokemonInformation.Caught = true
		config.Pokedex.AddPokemon(pokemonInformation)
	} else {
		pokemonInformation.Caught = false
	}

	return types.PokemonInformationResponse{Information: pokemonInformation}, nil
}

func Unmarshall[T types.GetLocationsResponse | types.PokemonEncountersResponse | types.PokemonInformation](val []byte, v *T) error {

	unmarshalError := json.Unmarshal(val, &v)
	if unmarshalError != nil {
		log.Fatalf("Failed to unmarshal response: %s\n", unmarshalError)
	}

	return nil
}
