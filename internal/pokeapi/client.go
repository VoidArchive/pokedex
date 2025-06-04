package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

// For detailed location area exploration
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

// PokemonData holds detailed information about a Pokemon
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
