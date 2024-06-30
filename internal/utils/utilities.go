package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"

	"github.com/mdwiltfong/PokeDex/internal/pokeapiclient"
)

type Config struct {
	PREV_URL *string
	NEXT_URL *string
	Client   *pokeapiclient.Client
	Pokedex  Pokedex
}
type Pokedex map[string]PokemonInformation

func (p Pokedex) AddPokemon(pokemon PokemonInformation) {
	_, exists := p[pokemon.Name]
	if !exists {
		p[pokemon.Name] = pokemon
	}
}

func (p Pokedex) GetPokemon(name string) (PokemonInformation, error) {
	pokemon, exists := p[name]
	if !exists {
		return PokemonInformation{}, errors.New("Pokemon not found")
	}
	return pokemon, nil
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

type PokemonInformationResponse struct {
	Information PokemonInformation
}

func (h PokemonInformationResponse) Response() interface{} {
	return h.Information
}
func (h PokemonInformationResponse) Print() {
	fmt.Printf("Throwing a Pokeball at %s\n", h.Information.Name)
	if h.Information.Caught {
		fmt.Printf("You caught %s!\n", h.Information.Name)
	} else {
		fmt.Printf("Oh no! %s got away!\n", h.Information.Name)
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

type CallbackFunction func(*Config, string) (CallbackResponse, error)

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
		"catch": {
			Name:        "catch",
			Description: "Catch a pokemon",
			Callback:    Catch,
		},
	}

}

func HelpCommand(*Config, string) (CallbackResponse, error) {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for key, value := range CliCommandMap() {
		fmt.Println(key + ": " + value.Description)
	}
	fmt.Println("")
	return HelpCommandResponse{CliCommandMap()}, nil
}

func ExitCommand(*Config, string) (CallbackResponse, error) {
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
type PokemonInformation struct {
	Caught    bool `json:"caught"`
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []any  `json:"past_abilities"`
	PastTypes     []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
	} `json:"past_types"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       any    `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  any    `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      any    `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale any    `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func Map(config *Config, commandInput string) (CallbackResponse, error) {
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
	var locations GetLocationsResponse
	marshalingError := Unmarshall(responseBytes, &locations)

	config.NEXT_URL = &locations.Next
	config.PREV_URL = &url
	config.Client.Cache.Add(url, responseBytes)
	if marshalingError != nil {
		log.Fatalf("Failed to unmarshal response: %s\n", marshalingError)
	}
	return MapCommandResponse{locations.Results}, nil
}

func Mapb(config *Config, commandInput string) (CallbackResponse, error) {

	if config.PREV_URL == nil || *config.PREV_URL == "" {
		fmt.Println("There are no previous pages")
		return MapCommandResponse{}, errors.New("there are no previous pages")
	}

	url := *config.PREV_URL
	var locations GetLocationsResponse
	cachedBytes, exists := config.Client.Cache.Get(url)
	if !exists {
		fmt.Println("No cached data!!!!")

		response, err := config.Client.HttpClient.Get(url)
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

		marshalingError := Unmarshall(cachedBytes, &locations)
		if marshalingError != nil {
			log.Fatalf("Failed to unmarshal response: %s\n", marshalingError)
		}
		return MapCommandResponse{locations.Results}, nil
	}

}

func Explore(config *Config, commandInput string) (CallbackResponse, error) {
	if commandInput == "" {
		return ExploreCommandResponse{}, errors.New("Please put in a location to explore")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", commandInput)
	cachedBytes, exists := config.Client.Cache.Get(url)
	if !exists {
		response, err := config.Client.HttpClient.Get(url)
		if err != nil {
			return ExploreCommandResponse{}, errors.New("There was an issue retrieving the data:" + err.Error())
		}
		body, _ := io.ReadAll(response.Body)
		if response.StatusCode > 299 {
			if response.StatusCode == 404 {
				return ExploreCommandResponse{}, errors.New("Area was not found")
			}
			return ExploreCommandResponse{}, errors.New("There was an issue retrieving the data")
		}
		responseBytes := []byte(body)

		var encounter PokemonEncountersResponse
		unMarshallError := Unmarshall[PokemonEncountersResponse](responseBytes, &encounter)
		config.Client.Cache.Add(url, responseBytes)
		if unMarshallError != nil {
			log.Fatalf("Failed to unmarshal response: %s\n", unMarshallError)
			return ExploreCommandResponse{}, errors.New("There was an issue unmarshalling the data" + unMarshallError.Error())
		}
		return ExploreCommandResponse{encounter.PokemonEncounters}, nil
	} else {
		fmt.Println("Cache hit")
		var encounter PokemonEncountersResponse
		unMarshallError := Unmarshall[PokemonEncountersResponse](cachedBytes, &encounter)
		if unMarshallError != nil {
			log.Fatalf("Failed to unmarshal response: %s\n", unMarshallError)
			return ExploreCommandResponse{}, errors.New("There was an issue unmarshalling the data" + unMarshallError.Error())
		}
		return ExploreCommandResponse{encounter.PokemonEncounters}, nil
	}

}

func Catch(config *Config, commandInput string) (CallbackResponse, error) {
	if commandInput == "" {
		return ExploreCommandResponse{}, errors.New("Please enter a pokemon you'd like to catch")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", commandInput)
	cachedBytes, exists := config.Client.Cache.Get(url)
	var pokemonInformation PokemonInformation
	if !exists {
		response, err := config.Client.HttpClient.Get(url)
		if err != nil {
			return ExploreCommandResponse{}, errors.New("There was an issue retrieving the data:" + err.Error())
		}
		body, _ := io.ReadAll(response.Body)
		if response.StatusCode > 299 {
			if response.StatusCode == 404 {
				return ExploreCommandResponse{}, errors.New("Pokemon was not found")
			}
			return ExploreCommandResponse{}, errors.New("There was an issue retrieving the data")
		}
		responseBytes := []byte(body)

		unMarshallError := Unmarshall(responseBytes, &pokemonInformation)
		if unMarshallError != nil {
			log.Fatalf("Failed to unmarshal response: %s\n", unMarshallError)
			return ExploreCommandResponse{}, errors.New("There was an issue unmarshalling the data" + unMarshallError.Error())
		}
		config.Client.Cache.Add(url, responseBytes)

	} else {
		fmt.Println("Cache hit")
		unMarshallError := Unmarshall(cachedBytes, &pokemonInformation)
		if unMarshallError != nil {
			log.Fatalf("Failed to unmarshal response: %s\n", unMarshallError)
			return ExploreCommandResponse{}, errors.New("There was an issue unmarshalling the data" + unMarshallError.Error())
		}

	}
	randNum := rand.Intn(pokemonInformation.BaseExperience)
	chance := float64(randNum) / float64(pokemonInformation.BaseExperience)
	fmt.Printf("Chance of catching %s: %f\n", pokemonInformation.Name, chance)
	if chance > 0.5 {
		pokemonInformation.Caught = true
		config.Pokedex.AddPokemon(pokemonInformation)
	} else {
		pokemonInformation.Caught = false
	}

	return PokemonInformationResponse{pokemonInformation}, nil
}

func Unmarshall[T GetLocationsResponse | PokemonEncountersResponse | PokemonInformation](val []byte, v *T) error {

	unmarshalError := json.Unmarshal(val, &v)
	if unmarshalError != nil {
		log.Fatalf("Failed to unmarshal response: %s\n", unmarshalError)
	}

	return nil
}
