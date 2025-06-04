package repl

import (
	"fmt"

	"github.com/voidarchive/pokedex/internal/shared/constants"
)

func commandPokedex(cfg *Config, args ...string) error {
	if len(cfg.Pokedex) == 0 {
		fmt.Printf("%sYour Pokedex is empty. Go catch some Pokemon!%s\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}

	fmt.Printf("\n%sYour Pokedex:%s\n", constants.ColorBrightCyan, constants.ColorReset)
	for name := range cfg.Pokedex {
		// Could add more info here like Level if desired, similar to party view
		fmt.Printf("  %s- %s%s%s\n", constants.ColorGray, constants.ColorWhite, name, constants.ColorReset)
	}
	fmt.Println()
	return nil
}
