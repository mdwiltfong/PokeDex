package utils

import (
	"testing"

	"github.com/mdwiltfong/PokeDex/internal/pokeapiclient"
	"github.com/mdwiltfong/PokeDex/internal/utils"
)

func TestSanitizeInput(t *testing.T) {
	expected := "blah"
	input := "  BLAH   "
	output := utils.SanitizeInput(input)
	if output == "" || output != expected {
		t.Fatalf(`SanitizeInput(%v)=%v, expected %v`, input, output, expected)
	}
}

func TestCliCommandMap(t *testing.T) {
	expectedCommands := []string{"help", "exit", "map", "mapb"}
	outputCommands := utils.CliCommandMap()
	for key, _ := range outputCommands {
		output := contains(expectedCommands, key)
		if output == false {
			t.Fatalf(`The key %v is not expected`, key)
		}
	}
}

func TestMap(t *testing.T) {
	configInput := &utils.Config{}
	clientInput := pokeapiclient.NewClient(50000, 10000)
	utils.Map(configInput, &clientInput)
	_, exists := clientInput.Cache.Get("https://pokeapi.co/api/v2/location/")
	if exists == false {
		t.Fatalf(`Map did not store the url:%v`, "https://pokeapi.co/api/v2/location/")
	}
	cacheLength := clientInput.Cache.Length()
	if cacheLength > 1 {
		t.Fatalf(`Cache should be 1 but was %v instead`, cacheLength)
	}
	if configInput.NEXT_URL == nil {
		t.Fatalf(`The NEXT_URL should be set, but it was nill`)
	}
}

func TestMapb(t *testing.T) {
	configInput := &utils.Config{}
	clientInput := pokeapiclient.NewClient(50000, 10000)
	response1, _ := utils.Map(configInput, &clientInput)
	response2, _ := utils.Mapb(configInput, &clientInput)
	response1.Response()
	isEqual(&response1, &response2)
	cacheLength := clientInput.Cache.Length()
	if cacheLength > 1 {
		t.Fatalf(`The cache length is %v when it should be 1`, cacheLength)
	}

}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func isEqual(response1 *utils.CallbackResponse, response2 *utils.CallbackResponse) {

}
