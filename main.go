package main

import (
	"time"

	"github.com/voidarchive/pokedex/internal/pokeapi"
	"github.com/voidarchive/pokedex/internal/pokecache"
	"github.com/voidarchive/pokedex/internal/repl"
)

func main() {
	cacheInterval := 5 * time.Minute
	cache := pokecache.NewCache(cacheInterval)

	pokeAPIClient := pokeapi.NewClient(cache)
	repl.StartRepl(&pokeAPIClient, cache)
}
