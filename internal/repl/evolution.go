package repl

import (
	"fmt"
	"time"

	"github.com/voidarchive/pokedex/internal/pokeapi"
	"github.com/voidarchive/pokedex/internal/shared/constants"
)

// TimeEvolutionThresholdNanos defines how long a Pokemon must be caught before being eligible for time-based evolution.
// For testing: 1 minute (60 * 1_000_000_000 nanos). For production, this would be much longer (e.g., 24 hours).
const TimeEvolutionThresholdNanos int64 = 1 * 60 * 1_000_000_000 // 1 minute

// CheckAndHandleEvolution attempts to evolve a Pokemon if it has met the criteria after leveling up or over time.
// It modifies cfg.Pokedex and cfg.Party directly if an evolution occurs.
// Returns true if an evolution happened, false otherwise, and an error if something went wrong during the process.
func CheckAndHandleEvolution(cfg *Config, pokemonName string) (bool, error) {
	userPokemon, exists := cfg.Pokedex[pokemonName]
	if !exists {
		return false, fmt.Errorf("%scannot check evolution for %s%s%s, not found in Pokedex%s", constants.ColorYellow, constants.ColorBrightRed, pokemonName, constants.ColorYellow, constants.ColorReset)
	}

	fmt.Printf("%sChecking if %s%s%s (Lvl %s%d%s) can evolve...%s\n", constants.ColorCyan, constants.ColorYellow, userPokemon.Name, constants.ColorCyan, constants.ColorYellow, userPokemon.Level, constants.ColorCyan, constants.ColorReset)
	species, err := cfg.PokeapiClient.FetchPokemonSpecies(userPokemon.Name)
	if err != nil {
		return false, fmt.Errorf("%scould not fetch species data for %s%s%s to check level-up evolution: %w%s", constants.ColorRed, constants.ColorBrightYellow, userPokemon.Name, constants.ColorRed, err, constants.ColorReset)
	}
	if species.EvolutionChain.URL == "" {
		fmt.Printf("  %s%s%s does not appear to have an evolution chain for level-up check.%s\n", constants.ColorGray, constants.ColorYellow, userPokemon.Name, constants.ColorReset)
	} else {
		evolutionChain, err := cfg.PokeapiClient.FetchEvolutionChain(species.EvolutionChain.URL)
		if err != nil {
			return false, fmt.Errorf("%scould not fetch evolution chain for %s%s%s: %w%s", constants.ColorRed, constants.ColorBrightYellow, userPokemon.Name, constants.ColorRed, err, constants.ColorReset)
		}

		evolvedToSpeciesName, evolutionByLevelTriggered := findPossibleLevelUpEvolution(userPokemon, evolutionChain.Chain)
		if evolutionByLevelTriggered {
			fmt.Printf("%sWhat? %s%s%s is evolving by %slevel-up%s!%s\n", constants.ColorBrightPurple, constants.ColorYellow, userPokemon.Name, constants.ColorBrightPurple, constants.ColorGreen, constants.ColorBrightPurple, constants.ColorReset)
			return performEvolution(cfg, userPokemon, evolvedToSpeciesName, "level-up")
		}
	}

	if userPokemon.CaughtTimestamp == 0 {
		fmt.Printf("  %sSkipping time-based evolution check for %s%s%s (no catch timestamp).%s\n", constants.ColorGray, constants.ColorYellow, userPokemon.Name, constants.ColorGray, constants.ColorReset)
		return false, nil
	}

	elapsedTime := time.Now().UnixNano() - userPokemon.CaughtTimestamp
	if elapsedTime > TimeEvolutionThresholdNanos {
		fmt.Printf("  %s%s%s has been with you for a while (%.2f minutes). Checking for %stime-based%s evolution...%s\n", constants.ColorCyan, constants.ColorYellow, userPokemon.Name, float64(elapsedTime)/float64(1*60*1_000_000_000), constants.ColorGreen, constants.ColorCyan, constants.ColorReset)

		if species.Name == "" || species.EvolutionChain.URL == "" {
			species, err = cfg.PokeapiClient.FetchPokemonSpecies(userPokemon.Name)
			if err != nil {
				return false, fmt.Errorf("%scould not fetch species data for %s%s%s for time-based evolution: %w%s", constants.ColorRed, constants.ColorBrightYellow, userPokemon.Name, constants.ColorRed, err, constants.ColorReset)
			}
		}
		if species.EvolutionChain.URL == "" {
			fmt.Printf("  %s%s%s does not appear to have an evolution chain for time-based check either.%s\n", constants.ColorGray, constants.ColorYellow, userPokemon.Name, constants.ColorReset)
			return false, nil
		}
		evolutionChain, err := cfg.PokeapiClient.FetchEvolutionChain(species.EvolutionChain.URL)
		if err != nil {
			return false, fmt.Errorf("%scould not fetch evolution chain for %s%s%s (time-based): %w%s", constants.ColorRed, constants.ColorBrightYellow, userPokemon.Name, constants.ColorRed, err, constants.ColorReset)
		}

		evolvedToSpeciesName, evolutionByTimeTriggered := findTimeBasedEvolutionCandidate(userPokemon, evolutionChain.Chain)
		if evolutionByTimeTriggered {
			fmt.Printf("%sWhat? %s%s%s is evolving by %stime%s!%s\n", constants.ColorBrightPurple, constants.ColorYellow, userPokemon.Name, constants.ColorBrightPurple, constants.ColorGreen, constants.ColorBrightPurple, constants.ColorReset)
			return performEvolution(cfg, userPokemon, evolvedToSpeciesName, "time")
		}
	}

	fmt.Printf("  %s%s%s is not ready to evolve yet or has no further evolutions meeting criteria.%s\n", constants.ColorGray, constants.ColorYellow, userPokemon.Name, constants.ColorReset)
	return false, nil
}

// performEvolution centralizes the logic to execute an evolution once a candidate is found.
func performEvolution(cfg *Config, originalPokemon pokeapi.UserPokemon, evolvedSpeciesName string, method string) (bool, error) {
	evolvedPokemonData, err := cfg.PokeapiClient.FetchPokemon(evolvedSpeciesName)
	if err != nil {
		return false, fmt.Errorf("%scould not fetch data for evolved form %s%s%s: %w%s", constants.ColorRed, constants.ColorBrightYellow, evolvedSpeciesName, constants.ColorRed, err, constants.ColorReset)
	}

	newEvolvedUserPokemon := pokeapi.UserPokemon{
		PokemonData:     evolvedPokemonData,
		Level:           originalPokemon.Level,
		CurrentXP:       0,
		CaughtTimestamp: time.Now().UnixNano(),
	}
	newEvolvedUserPokemon.XPToNextLevel = newEvolvedUserPokemon.CalculateNewXPToNextLevel()

	fmt.Printf("%sCongratulations! Your %s%s%s evolved into %s%s%s by %s%s%s!%s\n",
		constants.ColorBrightGreen, constants.ColorYellow, originalPokemon.Name, constants.ColorBrightGreen,
		constants.ColorBrightYellow, newEvolvedUserPokemon.Name, constants.ColorBrightGreen,
		constants.ColorGreen, method, constants.ColorBrightGreen, constants.ColorReset)

	delete(cfg.Pokedex, originalPokemon.Name)
	cfg.Pokedex[newEvolvedUserPokemon.Name] = newEvolvedUserPokemon

	for i, p := range cfg.Party {
		if p.Name == originalPokemon.Name {
			cfg.Party[i] = newEvolvedUserPokemon
			fmt.Printf("  %s%s%s was updated in your party.%s\n", constants.ColorGray, constants.ColorBrightYellow, newEvolvedUserPokemon.Name, constants.ColorReset)
			break
		}
	}
	return true, nil
}

// findPossibleLevelUpEvolution (renamed from findPossibleEvolution)
func findPossibleLevelUpEvolution(currentPokemon pokeapi.UserPokemon, chainLink pokeapi.CorrectedChainLink) (evolvesToName string, triggered bool) {
	if chainLink.Species.Name == currentPokemon.Name {
		for _, evolution := range chainLink.EvolvesTo {
			for _, detail := range evolution.EvolutionDetails {
				if detail.Trigger.Name == "level-up" && detail.MinLevel != nil && currentPokemon.Level >= *detail.MinLevel {
					return evolution.Species.Name, true
				}
			}
		}
		return "", false
	}
	for _, nextLink := range chainLink.EvolvesTo {
		evolvesToName, triggered = findPossibleLevelUpEvolution(currentPokemon, nextLink)
		if triggered {
			return evolvesToName, true
		}
	}
	return "", false
}

// findTimeBasedEvolutionCandidate checks if the current Pokemon has any evolution it can go to,
// simplifying the time-based rule to "any next stage is eligible after enough time".
func findTimeBasedEvolutionCandidate(currentPokemon pokeapi.UserPokemon, chainLink pokeapi.CorrectedChainLink) (evolvesToName string, canEvolveByTime bool) {
	// First, find the current Pokemon in the chain.
	if chainLink.Species.Name == currentPokemon.Name {
		// If it has any evolutions listed, pick the first one as a candidate for time-based evolution.
		// This is a major simplification and doesn't check *how* it evolves, only *that* it can.
		if len(chainLink.EvolvesTo) > 0 {
			// We should ensure this evolution isn't the same species (e.g. some special cases or data errors)
			if chainLink.EvolvesTo[0].Species.Name != currentPokemon.Name {
				return chainLink.EvolvesTo[0].Species.Name, true
			}
		}
		return "", false // No further evolutions from this specific species
	}

	// If current link is not our Pokemon, recursively check deeper in the chain.
	for _, nextLink := range chainLink.EvolvesTo {
		evolvesToName, canEvolveByTime = findTimeBasedEvolutionCandidate(currentPokemon, nextLink)
		if canEvolveByTime {
			return evolvesToName, true // Evolution found in a deeper branch
		}
	}
	return "", false
}
