package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/voidarchive/pokedex/internal/shared/constants"
)

const BaseURL = "https://pokeapi.co/api/v2"

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreasResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaDetail struct {
	ID                int                `json:"id"`
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonData struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

type UserPokemon struct {
	PokemonData           // Embeds all fields from PokemonData (Name, Stats, Types, etc.)
	Level           int   `json:"level"`
	CurrentXP       int   `json:"current_xp"`
	XPToNextLevel   int   `json:"xp_to_next_level"`
	CaughtTimestamp int64 `json:"caught_timestamp"` // Unix nanoseconds when caught
	// TODO: Potentially add other fields later, like current HP, status conditions, moveset, etc.
}

// CalculateNewXPToNextLevel provides a basic formula for determining XP for the next level.
// This can be adjusted for different leveling curves.
// Exported for potential use elsewhere if needed, but primarily a helper for AddXP.
func (up *UserPokemon) CalculateNewXPToNextLevel() int {
	// Example formula: (currentLevel^2 * 20) + 100. Min level 1 for calc.
	level := up.Level
	if level < 1 {
		level = 1 // Should not happen if Pokemon start at level 1+
	}
	return (level * level * 20) + 100
}

// AddXP adds experience points to the Pokemon and handles leveling up.
// It returns true if the Pokemon leveled up, false otherwise.
// This method modifies the UserPokemon instance it's called on.
func (up *UserPokemon) AddXP(xpGained int) bool {
	if xpGained <= 0 {
		return false
	}

	// Use BrightGreen for the Pokemon's name and Yellow for the XP amount.
	fmt.Printf("%s%s%s gained %s%d%s XP!\n", constants.ColorBrightGreen, up.Name, constants.ColorReset, constants.ColorYellow, xpGained, constants.ColorReset)
	up.CurrentXP += xpGained
	var leveledUpThisCycle bool = false
	var anyLevelUp bool = false

	for up.CurrentXP >= up.XPToNextLevel {
		leveledUpThisCycle = true
		anyLevelUp = true
		up.Level++
		xpForThisLevel := up.XPToNextLevel
		up.CurrentXP -= xpForThisLevel
		up.XPToNextLevel = up.CalculateNewXPToNextLevel()
		// Use BrightPurple for the congratulations message, BrightGreen for name, BrightCyan for level and new XP status.
		fmt.Printf("%sCongratulations! %s%s%s grew to Level %s%d%s! (XP: %s%d%s/%s%d%s)\n",
			constants.ColorBrightPurple,
			constants.ColorBrightGreen, up.Name, constants.ColorReset,
			constants.ColorBrightCyan, up.Level, constants.ColorReset,
			constants.ColorBrightCyan, up.CurrentXP, constants.ColorReset,
			constants.ColorCyan, up.XPToNextLevel, constants.ColorReset)

		if up.CurrentXP < 0 {
			up.CurrentXP = 0
		}
	}

	if leveledUpThisCycle {
		// If a level up occurred, ensure XP doesn't exceed the new cap if it was a multi-level up ending exactly.
		// This case is subtle: if CurrentXP became 0 after the last level up, and XPToNextLevel is now N, it's correct.
		// The loop condition `up.CurrentXP >= up.XPToNextLevel` handles continuing level-ups.
	}
	return anyLevelUp
}

// GetStat retrieves a specific stat value by name for a Pokemon.
// Returns the stat value and true if found, otherwise 0 and false.
func (pd *PokemonData) GetStat(statName string) (int, bool) {
	for _, s := range pd.Stats {
		if s.Stat.Name == statName {
			return s.BaseStat, true
		}
	}
	return 0, false
}

type CacheInterface interface {
	Add(key string, val []byte)
	Get(key string) ([]byte, bool)
}

type Client struct {
	httpClient http.Client
	cache      CacheInterface
}

func NewClient(cache CacheInterface) Client {
	return Client{
		httpClient: http.Client{},
		cache:      cache,
	}
}

func (c *Client) FetchLocationAreas(url string) (LocationAreasResponse, error) {
	var emptyResponse LocationAreasResponse

	if url == "" {
		return emptyResponse, fmt.Errorf("cannot fetch from empty URL")
	}
	if cacheData, ok := c.cache.Get(url); ok {
		log.Printf("Cache hit for URL: %s\n", url)
		var locationAreasRes LocationAreasResponse
		if err := json.Unmarshal(cacheData, &locationAreasRes); err != nil {
			log.Printf("Error unmarshalling cached data for %s: %v.", url, err)
		} else {
			return locationAreasRes, nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return emptyResponse, err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return emptyResponse, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return emptyResponse, fmt.Errorf("failed to read response: %w", err)
	}

	// Cache the response
	c.cache.Add(url, body)

	var locationAreasRes LocationAreasResponse
	if err := json.Unmarshal(body, &locationAreasRes); err != nil {
		return emptyResponse, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return locationAreasRes, nil
}

func (c *Client) FetchLocationAreaDetail(areaName string) (LocationAreaDetail, error) {
	var emptyResponse LocationAreaDetail

	url := fmt.Sprintf("%s/location-area/%s", BaseURL, areaName)

	// Check cache first
	if cacheData, ok := c.cache.Get(url); ok {
		log.Printf("Cache hit for URL: %s\n", url)
		var locationAreaDetail LocationAreaDetail
		if err := json.Unmarshal(cacheData, &locationAreaDetail); err != nil {
			log.Printf("Error unmarshalling cached data for %s: %v.", url, err)
		} else {
			return locationAreaDetail, nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return emptyResponse, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return emptyResponse, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return emptyResponse, fmt.Errorf("location area '%s' not found", areaName)
	}

	if res.StatusCode > 299 { // Broader check for non-2xx status codes
		return emptyResponse, fmt.Errorf("API request failed with status: %d for URL %s", res.StatusCode, url)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return emptyResponse, fmt.Errorf("failed to read response: %w", err)
	}

	// Cache the response
	c.cache.Add(url, body)

	var locationAreaDetail LocationAreaDetail
	if err := json.Unmarshal(body, &locationAreaDetail); err != nil {
		return emptyResponse, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return locationAreaDetail, nil
}

// FetchPokemon retrieves detailed information about a specific Pokemon by its name.
func (c *Client) FetchPokemon(pokemonName string) (PokemonData, error) {
	var emptyResponse PokemonData
	url := fmt.Sprintf("%s/pokemon/%s", BaseURL, pokemonName)

	if data, ok := c.cache.Get(url); ok {
		log.Printf("Cache hit for Pokemon: %s\n", pokemonName)
		var pokemonData PokemonData
		if err := json.Unmarshal(data, &pokemonData); err != nil {
			return emptyResponse, fmt.Errorf("failed to unmarshal cached pokemon data for %s: %w", pokemonName, err)
		}
		return pokemonData, nil
	}
	log.Printf("Cache miss for Pokemon: %s. Fetching from API: %s\n", pokemonName, url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return emptyResponse, fmt.Errorf("could not create request for %s: %w", pokemonName, err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return emptyResponse, fmt.Errorf("request to fetch %s failed: %w", pokemonName, err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return emptyResponse, fmt.Errorf("pokemon '%s' not found", pokemonName)
	}
	if res.StatusCode > 299 {
		bodyBytes, _ := io.ReadAll(res.Body) // Try to read body for more error info
		return emptyResponse, fmt.Errorf("API request for %s failed with status %d: %s", pokemonName, res.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return emptyResponse, fmt.Errorf("failed to read response body for %s: %w", pokemonName, err)
	}

	var pokemonData PokemonData
	if err := json.Unmarshal(body, &pokemonData); err != nil {
		return emptyResponse, fmt.Errorf("failed to unmarshal pokemon data for %s: %w. Body: %s", pokemonName, err, string(body))
	}

	c.cache.Add(url, body)
	return pokemonData, nil
}

// PokemonSpecies represents data from the /pokemon-species/{id_or_name}/ endpoint.
// We primarily need it to get the URL for the evolution chain.
type PokemonSpecies struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	EvolutionChain struct {
		URL string `json:"url"`
	} `json:"evolution_chain"`
	// Other fields like flavor_text_entries, generation, etc., are available but not immediately needed.
}

// EvolutionDetail contains the specifics of how a Pokemon evolves, e.g., min_level.
// There can be multiple evolution details (e.g., level up + hold item), so it's an array.
type EvolutionDetail struct {
	MinLevel *int `json:"min_level"` // Pointer because it can be null
	Trigger  struct {
		Name string `json:"name"`
	} `json:"trigger"` // e.g., "level-up"
	// Other details: item, gender, held_item, known_move, known_move_type, location, min_affection,
	// min_beauty, min_happiness, needs_overworld_rain, party_species, party_type, relative_physical_stats,
	// time_of_day, trade_species, turn_upside_down. For now, we only care about min_level and level-up trigger.
}

// ChainLink represents one link in an evolution chain.
// It contains the species in this link and how it evolves to the next species in the chain.
type ChainLink struct {
	Species          Pokemon           `json:"species"`           // Re-using the Pokemon struct (Name, URL) from above
	EvolvesTo        []ChainLink       `json:"evolves_to"`        // Nested evolutions
	EvolutionDetails []EvolutionDetail `json:"evolution_details"` // Details for evolution TO this species from previous -- this is incorrect, should be for evolution TO NEXT species
}

// Corrected ChainLink - EvolutionDetails are for the evolution from THIS species to the ones in EvolvesTo
// However, the API structure is: EvolvesTo contains a list of ChainLinks, and EACH of those has its own evolution_details for *how it evolved from the parent*.
// So, the details are actually on the child link that describes its evolution from the parent.
// Let's adjust the ChainLink struct to reflect the API's structure more accurately for parsing logic.
// The API structure is: current_species -> evolves_to (list) -> [ {evolution_details_for_this_evolution, species_it_becomes, its_own_evolves_to_list} ]

// Revised ChainLink structure to better match PokeAPI: evolution_details are part of the transition TO the species in `EvolvesTo`.
// The API gives us: current_species_info, and a list of potential evolutions (EvolvesTo).
// Each item in `EvolvesTo` is itself a `ChainLink` which also contains the evolution_details that LED to it.

// Simpler approach: The `ChainLink` represents a specific species in the chain.
// `EvolvesTo` is a list of *further* `ChainLink`s that *this* species can evolve into.
// The `EvolutionDetails` on a `ChainLink` within an `EvolvesTo` list describe how the *parent* evolves into *that specific ChainLink species*.

type CorrectedChainLink struct { // Renaming to avoid confusion during thought process, will consolidate
	IsBaby           bool                 `json:"is_baby"`
	Species          Pokemon              `json:"species"`
	EvolutionDetails []EvolutionDetail    `json:"evolution_details"` // How this species evolved from its predecessor
	EvolvesTo        []CorrectedChainLink `json:"evolves_to"`        // Further evolutions from this species
}

// EvolutionChainResponse is the top-level structure for an evolution chain API response.
type EvolutionChainResponse struct {
	ID    int                `json:"id"`
	Chain CorrectedChainLink `json:"chain"` // The start of the evolution chain (the base form)
}

func (c *Client) FetchPokemonSpecies(pokemonNameOrID string) (PokemonSpecies, error) {
	var emptyResponse PokemonSpecies
	url := fmt.Sprintf("%s/pokemon-species/%s", BaseURL, pokemonNameOrID)

	if data, ok := c.cache.Get(url); ok {
		log.Printf("Cache hit for Pokemon Species: %s\\n", pokemonNameOrID)
		var speciesData PokemonSpecies
		if err := json.Unmarshal(data, &speciesData); err != nil {
			return emptyResponse, fmt.Errorf("failed to unmarshal cached species data for %s: %w", pokemonNameOrID, err)
		}
		return speciesData, nil
	}
	log.Printf("Cache miss for Pokemon Species: %s. Fetching from API: %s\\n", pokemonNameOrID, url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return emptyResponse, fmt.Errorf("could not create request for species %s: %w", pokemonNameOrID, err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return emptyResponse, fmt.Errorf("request to fetch species %s failed: %w", pokemonNameOrID, err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return emptyResponse, fmt.Errorf("pokemon species '%s' not found", pokemonNameOrID)
	}
	if res.StatusCode > 299 {
		bodyBytes, _ := io.ReadAll(res.Body)
		return emptyResponse, fmt.Errorf("API request for species %s failed with status %d: %s", pokemonNameOrID, res.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return emptyResponse, fmt.Errorf("failed to read response body for species %s: %w", pokemonNameOrID, err)
	}

	var speciesData PokemonSpecies
	if err := json.Unmarshal(body, &speciesData); err != nil {
		return emptyResponse, fmt.Errorf("failed to unmarshal species data for %s: %w. Body: %s", pokemonNameOrID, err, string(body))
	}

	c.cache.Add(url, body)
	return speciesData, nil
}

func (c *Client) FetchEvolutionChain(url string) (EvolutionChainResponse, error) {
	var emptyResponse EvolutionChainResponse

	if url == "" {
		return emptyResponse, fmt.Errorf("cannot fetch evolution chain from empty URL")
	}

	if data, ok := c.cache.Get(url); ok {
		log.Printf("Cache hit for Evolution Chain URL: %s\\n", url)
		var chainData EvolutionChainResponse
		if err := json.Unmarshal(data, &chainData); err != nil {
			return emptyResponse, fmt.Errorf("failed to unmarshal cached evolution chain data for %s: %w", url, err)
		}
		return chainData, nil
	}
	log.Printf("Cache miss for Evolution Chain URL: %s. Fetching from API.\\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return emptyResponse, fmt.Errorf("could not create request for evolution chain %s: %w", url, err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return emptyResponse, fmt.Errorf("request to fetch evolution chain %s failed: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return emptyResponse, fmt.Errorf("evolution chain at '%s' not found", url)
	}
	if res.StatusCode > 299 {
		bodyBytes, _ := io.ReadAll(res.Body)
		return emptyResponse, fmt.Errorf("API request for evolution chain %s failed with status %d: %s", url, res.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return emptyResponse, fmt.Errorf("failed to read response body for evolution chain %s: %w", url, err)
	}

	var chainData EvolutionChainResponse
	if err := json.Unmarshal(body, &chainData); err != nil {
		return emptyResponse, fmt.Errorf("failed to unmarshal evolution chain data from %s: %w. Body: %s", url, err, string(body))
	}

	c.cache.Add(url, body)
	return chainData, nil
}
