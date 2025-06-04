package repl

import (
	"fmt"

	"github.com/voidarchive/pokedex/internal/shared/constants"
)

func commandMapf(cfg *Config, args ...string) error {
	if cfg.NextLocationAreaURL == nil || *cfg.NextLocationAreaURL == "" {
		cfg.CurrentAreaChoices = nil
		return fmt.Errorf("%sno next page of locations available%s", constants.ColorYellow, constants.ColorReset)
	}

	locationData, err := cfg.PokeapiClient.FetchLocationAreas(*cfg.NextLocationAreaURL)
	if err != nil {
		cfg.CurrentAreaChoices = nil
		return err // Error from FetchLocationAreas should be colored by the client or caller
	}

	fmt.Printf("\n%sLocation Areas:%s\n", constants.ColorCyan, constants.ColorReset)
	cfg.CurrentAreaChoices = locationData.Results // Store for explore command
	for i, area := range cfg.CurrentAreaChoices {
		fmt.Printf("  %s%d%s: %s%s%s\n", constants.ColorYellow, i+1, constants.ColorReset, constants.ColorWhite, area.Name, constants.ColorReset)
	}
	fmt.Printf("\n%sType 'explore <number>' or 'explore <full_area_name>' to see Pokemon in an area.%s\n", constants.ColorGray, constants.ColorReset)

	cfg.NextLocationAreaURL = locationData.Next
	cfg.PrevLocationAreaURL = locationData.Previous

	return nil
}

func commandMapb(cfg *Config, args ...string) error {
	if cfg.PrevLocationAreaURL == nil || *cfg.PrevLocationAreaURL == "" {
		cfg.CurrentAreaChoices = nil
		return fmt.Errorf("%syou're on the first page of the locations%s", constants.ColorYellow, constants.ColorReset)
	}
	locationData, err := cfg.PokeapiClient.FetchLocationAreas(*cfg.PrevLocationAreaURL)
	if err != nil {
		cfg.CurrentAreaChoices = nil
		return err
	}

	fmt.Printf("\n%sLocation Areas:%s\n", constants.ColorCyan, constants.ColorReset)
	cfg.CurrentAreaChoices = locationData.Results // Store for explore command
	for i, area := range cfg.CurrentAreaChoices {
		fmt.Printf("  %s%d%s: %s%s%s\n", constants.ColorYellow, i+1, constants.ColorReset, constants.ColorWhite, area.Name, constants.ColorReset)
	}
	fmt.Printf("\n%sType 'explore <number>' or 'explore <full_area_name>' to see Pokemon in an area.%s\n", constants.ColorGray, constants.ColorReset)

	cfg.NextLocationAreaURL = locationData.Next
	cfg.PrevLocationAreaURL = locationData.Previous

	return nil
}
