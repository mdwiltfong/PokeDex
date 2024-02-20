package main

import (
	"github.com/mdwiltfong/PokeDex/internal/pokeapiclient"
	"github.com/mdwiltfong/PokeDex/internal/utils"
)

func main() {
	cfg := &utils.Config{}
	client := pokeapiclient.NewClient(5000, 10000)
	utils.StartRepl(cfg, client)
}
