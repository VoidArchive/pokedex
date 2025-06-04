package repl

import (
	"fmt"
	"strings"
	"time"

	"github.com/voidarchive/pokedex/internal/pokeapi"
	"github.com/voidarchive/pokedex/internal/shared/constants"
)

// ansiReset is already defined in command_catch.go, using ColorReset from colors.go is better for consistency.
// We'll assume ColorReset is available globally in the package or use repl.ColorReset if not.

func commandCatch(cfg *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("%susage: catch <pokemon_name> [pokeball_type (e.g., pokeball, greatball)]%s", constants.ColorYellow, constants.ColorReset)
	}
	pokemonName := args[0]

	ballToUseKey := DefaultBall // from types.go
	if len(args) > 1 {
		ballToUseKey = strings.ToLower(args[1])
	}

	chosenBall, ballExists := KnownPokeballs[ballToUseKey]
	if !ballExists {
		return fmt.Errorf("%sunknown pokeball type: %s%s%s. Available: pokeball, greatball, ultraball%s", constants.ColorYellow, constants.ColorBrightRed, ballToUseKey, constants.ColorYellow, constants.ColorReset)
	}

	// Check inventory
	if count, hasBall := cfg.Inventory[ballToUseKey]; !hasBall || count <= 0 {
		return fmt.Errorf("%syou don't have any %s%s%s left%s", constants.ColorYellow, chosenBall.Color, chosenBall.Name, constants.ColorYellow, constants.ColorReset)
	}

	// Decrement ball count (attempt is made)
	cfg.Inventory[ballToUseKey]--

	fmt.Printf("Throwing a %s%s%s at %s%s%s...\n", chosenBall.Color, chosenBall.Name, constants.ColorReset, constants.ColorYellow, pokemonName, constants.ColorReset)

	// Check if already caught before making an API call
	if _, caught := cfg.Pokedex[pokemonName]; caught {
		fmt.Printf("%s%s%s has already been caught!%s\n", constants.ColorYellow, constants.ColorBrightYellow, pokemonName, constants.ColorReset)
		fmt.Printf("%sYou can inspect it using the 'inspect' command.%s\n", constants.ColorGray, constants.ColorReset)
		// Note: The ball was used, so inventory remains decremented.
		return nil
	}

	pokemonData, err := cfg.PokeapiClient.FetchPokemon(pokemonName)
	if err != nil {
		return err // API client errors are already formatted or will be by the top-level handler
	}

	const maxRollValue = 400
	successThreshold := maxRollValue - pokemonData.BaseExperience

	// Apply Pokeball modifier
	successThreshold = int(float64(successThreshold) * chosenBall.CatchRateMod)

	const minSuccessPoints = 20  // Adjusted slightly if needed, or keep as is
	const maxSuccessPoints = 390 // Adjusted slightly, ensuring it's less than maxRollValue

	if pokemonData.BaseExperience <= 0 {
		successThreshold = maxSuccessPoints
	} else {
		if successThreshold < minSuccessPoints {
			successThreshold = minSuccessPoints
		}
		if successThreshold > maxSuccessPoints {
			successThreshold = maxSuccessPoints
		}
	}

	roll := cfg.Randomizer.Intn(maxRollValue)

	if roll < successThreshold {
		fmt.Printf("%s%s%s%s was caught!%s\n", chosenBall.Color, constants.ColorBrightGreen, pokemonData.Name, chosenBall.Color, constants.ColorReset)
		fmt.Printf("%sYou may now inspect it with the inspect command.%s\n", constants.ColorGray, constants.ColorReset)

		newUserPokemon := pokeapi.UserPokemon{
			PokemonData:     pokemonData,
			Level:           1,
			CurrentXP:       0,
			CaughtTimestamp: time.Now().UnixNano(), // Set caught timestamp
		}
		newUserPokemon.XPToNextLevel = newUserPokemon.CalculateNewXPToNextLevel() // Calculate based on its level

		cfg.Pokedex[newUserPokemon.Name] = newUserPokemon

		if len(cfg.Party) < MaxPartySize {
			cfg.Party = append(cfg.Party, newUserPokemon)
			fmt.Printf("%s%s%s has been added to your party!%s\n", chosenBall.Color, newUserPokemon.Name, constants.ColorReset, constants.ColorReset)
		} else {
			fmt.Printf("%s%s%s has been sent to your Pokedex storage as your party is full.%s\n", chosenBall.Color, newUserPokemon.Name, constants.ColorReset, constants.ColorReset)
		}

	} else {
		fmt.Printf("%sOh no! %s%s%s escaped!%s\n", constants.ColorRed, chosenBall.Color, pokemonData.Name, constants.ColorRed, constants.ColorReset)
	}

	return nil
}
