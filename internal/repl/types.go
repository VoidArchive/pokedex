package repl

import "github.com/voidarchive/pokedex/internal/pokeapi"

type Config struct {
	NextLocationAreaURL *string
	PrevLocationAreaURL *string
	PokeapiClient       PokeapiClient
	Cache               Pokecache
	Pokedex             map[string]pokeapi.PokemonData
}

type PokeapiClient interface {
	FetchLocationAreas(url string) (pokeapi.LocationAreasResponse, error)
	FetchLocationAreaDetail(areaName string) (pokeapi.LocationAreaDetail, error)
	FetchPokemon(pokemonName string) (pokeapi.PokemonData, error)
}

type Pokecache interface {
	Add(key string, val []byte)
	Get(key string) ([]byte, bool)
}

type cliCommand struct {
	Name        string
	Description string
	Callback    func(*Config, ...string) error
}
