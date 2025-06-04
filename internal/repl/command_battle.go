package repl

import (
	"fmt"

	"github.com/voidarchive/pokedex/internal/battle"
	"github.com/voidarchive/pokedex/internal/shared/constants"
)

// commandBattle handles the 'battle' command from the REPL.
// It initiates a battle between a user's caught Pokemon and a specified opponent Pokemon.
func commandBattle(cfg *Config, args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("%susage: battle <your_pokemon_name> <opponent_pokemon_name>%s", constants.ColorYellow, constants.ColorReset)
	}
	playerPokemonName := args[0]
	opponentPokemonName := args[1]

	playerPokemon, caught := cfg.Pokedex[playerPokemonName]
	if !caught {
		return fmt.Errorf("%syou have not caught '%s%s%s' to battle with%s", constants.ColorYellow, constants.ColorBrightRed, playerPokemonName, constants.ColorYellow, constants.ColorReset)
	}

	fmt.Printf("%sFetching opponent %s%s%s for battle...%s\n", constants.ColorCyan, constants.ColorYellow, opponentPokemonName, constants.ColorCyan, constants.ColorReset)
	opponentPokemonData, err := cfg.PokeapiClient.FetchPokemon(opponentPokemonName)
	if err != nil {
		// Error from FetchPokemon is already descriptive, but we can color the wrapper message
		return fmt.Errorf("%scould not fetch opponent Pokemon '%s%s%s': %w%s", constants.ColorRed, constants.ColorYellow, opponentPokemonName, constants.ColorRed, err, constants.ColorReset)
	}

	// SimulateBattle now returns xpGained
	xpGained := battle.SimulateBattle(cfg.Randomizer, playerPokemon, opponentPokemonData)

	if xpGained > 0 {
		// Get the pokemon from Pokedex to update (it's a struct, so we operate on a copy then reassign)
		updatedPlayerPokemon := cfg.Pokedex[playerPokemonName] // Get a fresh copy
		leveledUp := updatedPlayerPokemon.AddXP(xpGained)      // AddXP modifies updatedPlayerPokemon directly
		cfg.Pokedex[playerPokemonName] = updatedPlayerPokemon  // Re-assign the modified Pokemon to the Pokedex

		// Update in party if present
		for i, p := range cfg.Party {
			if p.Name == playerPokemonName {
				cfg.Party[i] = updatedPlayerPokemon // Update the Pokemon in the party slot
				break
			}
		}

		if leveledUp {
			// The AddXP method already prints level up messages.
			// We could add more post-level up logic here if needed.
			fmt.Printf("%s%s%s's stats may have changed due to leveling up!%s\n", constants.ColorGreen, constants.ColorYellow, playerPokemonName, constants.ColorReset)

			// Check for evolution after leveling up
			evolved, err := CheckAndHandleEvolution(cfg, updatedPlayerPokemon.Name) // Pass the name of the (potentially) updated Pokemon
			if err != nil {
				// CheckAndHandleEvolution and performEvolution already color their errors, this is a fallback/wrapper
				fmt.Printf("%sError during evolution check for %s%s%s: %v%s\n", constants.ColorRed, constants.ColorYellow, updatedPlayerPokemon.Name, constants.ColorRed, err, constants.ColorReset)
			}
			if evolved {
				// If evolution occurred, updatedPlayerPokemon is now stale. cfg.Pokedex and cfg.Party have the new Pokemon.
				// The evolution messages are handled by CheckAndHandleEvolution.
				// We might need to refetch the evolved Pokemon if we need its new name here, but for now, not necessary.
				fmt.Printf("%s--- %s%s%s has evolved! ---%s\n", constants.ColorBrightPurple, constants.ColorYellow, playerPokemonName, constants.ColorBrightPurple, constants.ColorReset)
			}
		}
	}

	return nil
}
