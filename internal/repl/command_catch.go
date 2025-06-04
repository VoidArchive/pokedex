package repl

import (
	"fmt"
	"math/rand"
)

func commandCatch(cfg *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a Pokemon name to catch")
	}
	pokemonName := args[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	// Check if already caught before making an API call
	if _, caught := cfg.Pokedex[pokemonName]; caught {
		fmt.Printf("%s has already been caught!\n", pokemonName)
		fmt.Println("You can inspect it using the 'inspect' command.") // Added hint for future inspect command
		return nil
	}

	pokemonData, err := cfg.PokeapiClient.FetchPokemon(pokemonName)
	if err != nil {
		// The FetchPokemon method already provides a good error message, including if not found.
		return err
	}

	// Catching logic:
	// Higher base experience makes it harder to catch.
	// We want the chance of success to decrease as base_experience increases.
	// Using a scale where a random number needs to be below a threshold that is
	// inversely related to base_experience.
	const maxRollValue = 400 // Defines the range of the random roll [0, maxRollValue-1]

	// Calculate a dynamic success threshold. Higher BE reduces this threshold.
	successThreshold := maxRollValue - pokemonData.BaseExperience

	// Set some reasonable floor and ceiling for the success threshold
	// to ensure a minimum chance to catch and a cap for very easy Pokemon.
	const minSuccessPoints = 40  // Corresponds to about 10% chance (40/400)
	const maxSuccessPoints = 380 // Corresponds to about 95% chance (380/400)

	if pokemonData.BaseExperience <= 0 { // Handle unlikely case of 0 or negative BE
		successThreshold = maxSuccessPoints // Make it very easy to catch
	} else {
		if successThreshold < minSuccessPoints {
			successThreshold = minSuccessPoints
		}
		if successThreshold > maxSuccessPoints {
			successThreshold = maxSuccessPoints
		}
	}

	roll := rand.Intn(maxRollValue)

	if roll < successThreshold {
		fmt.Printf("%s was caught!\n", pokemonData.Name)
		fmt.Println("You may now inspect it with the inspect command.")
		cfg.Pokedex[pokemonData.Name] = pokemonData
	} else {
		fmt.Printf("%s escaped!\n", pokemonData.Name)
	}

	return nil
}
