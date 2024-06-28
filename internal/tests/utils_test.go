package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/mdwiltfong/PokeDex/internal/pokeapiclient"
	"github.com/mdwiltfong/PokeDex/internal/pokecache"
	"github.com/mdwiltfong/PokeDex/internal/utils"
)

func TestSanitizeInput(t *testing.T) {

	input := "  COMMAND INPUT   "
	output := utils.SanitizeInput(input)
	if output[0] != "command" && output[1] != "input" {
		t.Fatalf(`Command was: %s but expected Command \n Input was: %s but expected input`, output[0], output[1])

	}
}

func TestCliCommandMap(t *testing.T) {
	expectedCommands := []string{"help", "exit", "map", "mapb", "explore"}
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
	utils.Map(configInput, &clientInput, "")
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
	output1, _ := utils.Map(configInput, &clientInput, "")
	output2, _ := utils.Mapb(configInput, &clientInput, "")

	if isEqual(output1, output2) == false {
		t.Fatalf(`The two responses are not equal`)
	}
	cacheLength := clientInput.Cache.Length()
	if cacheLength > 1 {
		t.Fatalf(`The cache length is %v when it should be 1`, cacheLength)
	}
}

func TestCache(t *testing.T) {
	cache := pokecache.NewCache(5000)
	testInput := []byte{71, 111}
	cache.Add("Test", testInput)
	if cache.Length() != 1 {
		t.Fatalf(`Cache length is not one, but was %v instead`, cache.Length())
	}
	testOutput, _ := cache.Get("Test")
	if testOutput == nil {
		t.Fatalf("Cache is storing nil values")
	}
}
func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := pokecache.NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestExplore(t *testing.T) {
	configInput := &utils.Config{}
	clientInput := pokeapiclient.NewClient(50000, 10000)
	output, _ := utils.Explore(configInput, &clientInput, "canalave-city-area")
	if output.Response() == nil {
		t.Fatalf(`Explore returned nil response`)
	}
}

func TestExploreError404(t *testing.T) {
	configInput := &utils.Config{}
	clientInput := pokeapiclient.NewClient(50000, 10000)
	output, err := utils.Explore(configInput, &clientInput, "LOL")
	if output.Response() == nil {
		t.Fatalf(`Explore returned nil response`)
	}
	if err == nil {
		t.Fatalf("Error object should be nil but was: %s", err.Error())
	}
}

func TestExploreErrorNoInput(t *testing.T) {
	configInput := &utils.Config{}
	clientInput := pokeapiclient.NewClient(50000, 10000)
	_, err := utils.Explore(configInput, &clientInput, "")

	if err.Error() != "Please put in a location to explroe" {
		t.Fatalf("Error object should be nil but was: %s", err.Error())
	}
}
func TestExploreCache(t *testing.T) {
	configInput := &utils.Config{}
	clientInput := pokeapiclient.NewClient(50000, 10000)
	_, err := utils.Explore(configInput, &clientInput, "canalave-city-area")
	if err != nil {
		t.Fatalf("Error object should be nil but was: %s", err.Error())
	}
	_, err = utils.Explore(configInput, &clientInput, "canalave-city-area")
	if err != nil {
		t.Fatalf("Error object should be nil but was: %s", err.Error())
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

func isEqual(output1 utils.CallbackResponse, output2 utils.CallbackResponse) bool {
	locations1 := output1.Response().([]utils.Location)
	locations2 := output2.Response().([]utils.Location)
	for i, _ := range locations1 {
		if locations1[i].Name != locations2[i].Name {
			return false
		}
	}
	return true
}
