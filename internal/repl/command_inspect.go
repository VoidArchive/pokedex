package repl

import (
	"fmt"

	"github.com/voidarchive/pokedex/internal/shared/constants"
)

func commandInspect(cfg *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("%syou must provide a Pokemon name to inspect%s", constants.ColorYellow, constants.ColorReset)
	}
	pokemonName := args[0]

	pokemon, caught := cfg.Pokedex[pokemonName]
	if !caught {
		return fmt.Errorf("%syou have not caught %s%s%s%s", constants.ColorYellow, constants.ColorBrightRed, pokemonName, constants.ColorYellow, constants.ColorReset)
	}

	fmt.Printf("\n%s--- Inspecting: %s%s%s ---%s\n", constants.ColorCyan, constants.ColorBrightYellow, pokemon.Name, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %sName:%s %s%s%s\n", constants.ColorGreen, constants.ColorReset, constants.ColorWhite, pokemon.Name, constants.ColorReset)
	fmt.Printf("  %sHeight:%s %s%d%s\n", constants.ColorGreen, constants.ColorReset, constants.ColorWhite, pokemon.Height, constants.ColorReset)
	fmt.Printf("  %sWeight:%s %s%d%s\n", constants.ColorGreen, constants.ColorReset, constants.ColorWhite, pokemon.Weight, constants.ColorReset)

	fmt.Printf("  %sLevel:%s %s%d%s\n", constants.ColorGreen, constants.ColorReset, constants.ColorBrightCyan, pokemon.Level, constants.ColorReset)
	fmt.Printf("  %sXP:%s %s%d%s/%s%d%s\n", constants.ColorGreen, constants.ColorReset, constants.ColorBrightCyan, pokemon.CurrentXP, constants.ColorReset, constants.ColorCyan, pokemon.XPToNextLevel, constants.ColorReset)

	fmt.Printf("  %sStats:%s\n", constants.ColorGreen, constants.ColorReset)
	for _, stat := range pokemon.Stats {
		fmt.Printf("    %s-%s%s: %s%d%s\n", constants.ColorBlue, stat.Stat.Name, constants.ColorReset, constants.ColorWhite, stat.BaseStat, constants.ColorReset)
	}
	fmt.Printf("  %sTypes:%s\n", constants.ColorGreen, constants.ColorReset)
	for _, typeInfo := range pokemon.Types {
		// We can add specific colors per type later if desired
		fmt.Printf("    %s- %s%s%s\n", constants.ColorPurple, constants.ColorWhite, typeInfo.Type.Name, constants.ColorReset)
	}
	fmt.Println(constants.ColorCyan + "------------------------" + constants.ColorReset)

	return nil
}
