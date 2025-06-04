package repl

import (
	"fmt"
	"strconv"

	"github.com/voidarchive/pokedex/internal/shared/constants"
)

func commandExplore(cfg *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("%syou must provide a location area name or number (if choices are available from 'map')%s", constants.ColorYellow, constants.ColorReset)
	}

	areaIdentifier := args[0]
	areaNameToExplore := ""

	// Try to parse as a number first
	if areaNum, err := strconv.Atoi(areaIdentifier); err == nil {
		// Successfully parsed as a number, check if it's a valid choice
		if cfg.CurrentAreaChoices != nil && areaNum > 0 && areaNum <= len(cfg.CurrentAreaChoices) {
			areaNameToExplore = cfg.CurrentAreaChoices[areaNum-1].Name
			fmt.Printf("%sChoosing area #%d: %s%s%s\n", constants.ColorGray, areaNum, constants.ColorYellow, areaNameToExplore, constants.ColorReset)
		} else if cfg.CurrentAreaChoices != nil {
			// It was a number, but not a valid choice from the current list
			return fmt.Errorf("%sinvalid area number: %s%d%s. Please choose from %s1%s to %s%d%s, or provide a full area name. Or, run 'map' again for updated choices%s", constants.ColorYellow, constants.ColorBrightRed, areaNum, constants.ColorYellow, constants.ColorGreen, constants.ColorYellow, constants.ColorGreen, len(cfg.CurrentAreaChoices), constants.ColorYellow, constants.ColorReset)
		} else {
			// It was a number, but no choices are currently loaded (e.g., 'map' hasn't been run successfully)
			// In this case, we assume the number might be part of an area name (e.g. "area-123")
			areaNameToExplore = areaIdentifier
		}
	} else {
		// Not a number, treat as a full area name
		areaNameToExplore = areaIdentifier
	}

	if areaNameToExplore == "" {
		return fmt.Errorf("%sno valid area specified for exploration. Run 'map' to see available areas%s", constants.ColorYellow, constants.ColorReset)
	}

	fmt.Printf("\n%sExploring %s%s%s...%s\n", constants.ColorCyan, constants.ColorYellow, areaNameToExplore, constants.ColorCyan, constants.ColorReset)

	locationDetail, err := cfg.PokeapiClient.FetchLocationAreaDetail(areaNameToExplore)
	if err != nil {
		return err // Error from FetchLocationAreaDetail should be colored by the client or caller
	}

	fmt.Printf("%sFound Pokemon:%s\n", constants.ColorGreen, constants.ColorReset)
	if len(locationDetail.PokemonEncounters) == 0 {
		fmt.Printf("  %sNo Pokemon found in this area.%s\n", constants.ColorGray, constants.ColorReset)
	} else {
		for _, encounter := range locationDetail.PokemonEncounters {
			fmt.Printf("  %s- %s%s%s\n", constants.ColorGray, constants.ColorWhite, encounter.Pokemon.Name, constants.ColorReset)
		}
	}

	// Random Encounter Logic
	if len(locationDetail.PokemonEncounters) > 0 {
		// Select one Pokemon randomly from the list of possible encounters in this area
		randomIndex := cfg.Randomizer.Intn(len(locationDetail.PokemonEncounters))
		encounteredPokemon := locationDetail.PokemonEncounters[randomIndex].Pokemon

		fmt.Printf("\n%sA wild %s%s%s has appeared!%s\n", constants.ColorBrightYellow, constants.ColorYellow, encounteredPokemon.Name, constants.ColorBrightYellow, constants.ColorReset)
		// Future: Add options like battle, catch, run here.
		// For now, just the appearance is the "encounter".
	}

	// Consider clearing CurrentAreaChoices if we want to force 'map' before every 'explore <number>'
	// cfg.CurrentAreaChoices = nil

	return nil
}
