package apiclient

import (
	"fmt"
	"net/http"
)

type HTTPMethod int

const (
	GET    HTTPMethod = iota
	PUT    HTTPMethod = iota
	POST   HTTPMethod = iota
	DELETE HTTPMethod = iota
)

func PokemonHttpRequest(endpoint string, method HTTPMethod) {
	baseUrl := "https://pokeapi.co/api/v2"
	switch method {
	case GET:
		http.Get(baseUrl + endpoint)
	default:
		fmt.Print("That HTTP method is not supported")
	}

}
