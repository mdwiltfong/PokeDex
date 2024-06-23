package utils

import (
	"testing"

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

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
