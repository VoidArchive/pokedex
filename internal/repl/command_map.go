package repl

import (
	"fmt"
)

func commandMapf(cfg *Config, args ...string) error {
	if cfg.NextLocationAreaURL == nil || *cfg.NextLocationAreaURL == "" {
		return fmt.Errorf("no next page of locations available")
	}

	locationData, err := cfg.PokeapiClient.FetchLocationAreas(*cfg.NextLocationAreaURL)
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

func commandMapb(cfg *Config, args ...string) error {
	if cfg.PrevLocationAreaURL == nil || *cfg.PrevLocationAreaURL == "" {
		return fmt.Errorf("you're on the first page of the locations")
	}
	locationData, err := cfg.PokeapiClient.FetchLocationAreas(*cfg.PrevLocationAreaURL)
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
