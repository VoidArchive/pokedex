package repl

import (
	"fmt"
	"strings"

	"github.com/voidarchive/pokedex/internal/shared/constants"
)

func commandParty(cfg *Config, args ...string) error {
	if len(cfg.Party) == 0 {
		fmt.Printf("%sYour party is empty.%s\n", constants.ColorYellow, constants.ColorReset)
		fmt.Printf("%sCaught Pokemon will be added to your party if there is space.%s\n", constants.ColorGray, constants.ColorReset)
		return nil
	}

	fmt.Printf("\n%sYour Party:%s\n", constants.ColorBrightCyan, constants.ColorReset)
	for i, p := range cfg.Party {
		fmt.Printf("  %sSlot %d:%s %s%s%s (Lvl %s%d%s) - XP: %s%d%s/%s%d%s\n",
			constants.ColorYellow, i+1, constants.ColorReset,
			constants.ColorWhite, p.Name, constants.ColorReset,
			constants.ColorBrightCyan, p.Level, constants.ColorReset,
			constants.ColorBrightCyan, p.CurrentXP, constants.ColorReset,
			constants.ColorCyan, p.XPToNextLevel, constants.ColorReset)

		var typeStrings []string
		for _, t := range p.PokemonData.Types {
			// Could add specific colors per type here in the future
			typeStrings = append(typeStrings, fmt.Sprintf("%s%s%s", constants.ColorPurple, t.Type.Name, constants.ColorReset))
		}
		fmt.Printf("      %sTypes:%s %s\n", constants.ColorGreen, constants.ColorReset, strings.Join(typeStrings, ", "))
	}
	fmt.Println()
	return nil
}
