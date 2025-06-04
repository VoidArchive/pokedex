package repl

import (
	"fmt"
	"os"

	"github.com/voidarchive/pokedex/internal/shared/constants"
)

func commandExit(cfg *Config, args ...string) error {
	fmt.Printf("%sClosing the Pokedex... Goodbye!%s\n", constants.ColorGreen, constants.ColorReset)
	if err := savePokedex(cfg.Pokedex, cfg.Party); err != nil {
		// savePokedex already returns a colored error string, but we might want to ensure the whole message is structured.
		fmt.Fprintf(os.Stderr, "%sError saving game data on exit: %s%s\n", constants.ColorRed, err.Error(), constants.ColorReset)
	}
	os.Exit(0)
	return nil // This line is never reached due to os.Exit(0)
}
