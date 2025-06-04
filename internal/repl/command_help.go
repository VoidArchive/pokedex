package repl

import (
	"fmt"

	"github.com/voidarchive/pokedex/internal/shared/constants"
)

func commandHelp(cfg *Config, args ...string) error {
	fmt.Println()
	fmt.Printf("%sWelcome to the Pokedex!%s\n", constants.ColorBrightCyan, constants.ColorReset)
	fmt.Printf("%sUsage:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("  %s%s%s: %s%s%s\n", constants.ColorGreen, cmd.Name, constants.ColorReset, constants.ColorGray, cmd.Description, constants.ColorReset)
	}
	fmt.Println()
	return nil
}
