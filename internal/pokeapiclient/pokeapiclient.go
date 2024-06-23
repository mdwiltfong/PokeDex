package pokeapiclient

import (
	"net/http"
	"time"

	"github.com/mdwiltfong/PokeDex/internal/pokecache"
)

type Client struct {
	Cache      *pokecache.Cache
	HttpClient http.Client
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		Cache:      pokecache.NewCache(cacheInterval),
		HttpClient: http.Client{},
	}
}
