package repl

import (
	"math/rand"

	"github.com/voidarchive/pokedex/internal/pokeapi"
)

const MaxPartySize = 6

// UserPokemon has been moved to internal/pokeapi/client.go

// PokeballType defines the properties of a Pokeball.
type PokeballType struct {
	Name         string
	CatchRateMod float64 // e.g., 1.0 for Poke Ball, 1.5 for Great Ball, 2.0 for Ultra Ball
	Color        string  // ANSI color code for display
}

// KnownPokeballs maps ball names to their properties.
// This can be expanded with more ball types.
var KnownPokeballs = map[string]PokeballType{
	"pokeball":  {Name: "Poke Ball", CatchRateMod: 1.0, Color: "\033[0;37m"},  // White
	"greatball": {Name: "Great Ball", CatchRateMod: 1.5, Color: "\033[0;34m"}, // Blue
	"ultraball": {Name: "Ultra Ball", CatchRateMod: 2.0, Color: "\033[0;33m"}, // Yellow/Gold
	// Add Master Ball later if desired: {Name: "Master Ball", CatchRateMod: 255.0, Color: "\033[0;35m"}, // Purple
}

const DefaultBall = "pokeball"

type Config struct {
	NextLocationAreaURL *string
	PrevLocationAreaURL *string
	PokeapiClient       PokeapiClient
	Cache               Pokecache
	Pokedex             map[string]pokeapi.UserPokemon
	Party               []pokeapi.UserPokemon
	Inventory           map[string]int // Item name -> count (e.g., "pokeball" -> 10)
	Randomizer          *rand.Rand
	CurrentAreaChoices  []pokeapi.LocationArea // For 'map' command to store choices for 'explore'
}

type PokeapiClient interface {
	FetchLocationAreas(url string) (pokeapi.LocationAreasResponse, error)
	FetchLocationAreaDetail(areaName string) (pokeapi.LocationAreaDetail, error)
	FetchPokemon(pokemonName string) (pokeapi.PokemonData, error)
	FetchPokemonSpecies(pokemonNameOrID string) (pokeapi.PokemonSpecies, error)
	FetchEvolutionChain(url string) (pokeapi.EvolutionChainResponse, error)
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
