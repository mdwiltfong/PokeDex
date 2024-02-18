package main

import (
	"github.com/mdwiltfong/PokeDex/internal/utils"
)

func main() {
	cfg := &utils.Config{}
	utils.StartRepl(cfg)
}
