package main

import (
	"github.com/mdwiltfong/PokeDex/internal/pokeapiclient"
	"github.com/mdwiltfong/PokeDex/internal/utils"
)

func main() {
	client := pokeapiclient.NewClient(50000, 100000)
	cfg := &utils.Config{
		PREV_URL: nil,
		NEXT_URL: nil,
		Client:   &client,
	}
	utils.StartRepl(cfg)
}
