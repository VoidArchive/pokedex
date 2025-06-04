package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

func FetchLocationArea(url string) (LocationAreaResponse, error) {
	var emptyResponse LocationAreaResponse

	if url == "" {
		return emptyResponse, fmt.Errorf("cannot fetch from empty URL")
	}
	res, err := http.Get(url)
	if err != nil {
		return emptyResponse, fmt.Errorf("failed to fetch location areas: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return emptyResponse, fmt.Errorf("failed to read response: %w", err)
	}
	var locationAreasRes LocationAreaResponse
	if err := json.Unmarshal(body, &locationAreasRes); err != nil {
		return emptyResponse, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return locationAreasRes, nil
}

func commandMap(cfg *Config) error {
	if cfg.NextLocationAreaURL == nil || *cfg.NextLocationAreaURL == "" {
		return fmt.Errorf("no next page of locations available")
	}

	locationData, err := FetchLocationArea(*cfg.NextLocationAreaURL)
	if err != nil {
		return err
	}
	for _, area := range locationData.Results {
		fmt.Println(area.Name)
	}
	cfg.NextLocationAreaURL = locationData.Next
	cfg.PrevLocationAreaURL = locationData.Previous

	return nil
}

func commandMapb(cfg *Config) error {
	if cfg.PrevLocationAreaURL == nil || *cfg.PrevLocationAreaURL == "" {
		return fmt.Errorf("you're on the first page of the locations")
	}
	locationData, err := FetchLocationArea(*cfg.PrevLocationAreaURL)
	if err != nil {
		return err
	}
	for _, area := range locationData.Results {
		fmt.Println(area.Name)
	}
	cfg.NextLocationAreaURL = locationData.Next
	cfg.PrevLocationAreaURL = locationData.Previous

	return nil
}
