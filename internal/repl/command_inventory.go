package repl

import (
	"fmt"

	"github.com/voidarchive/pokedex/internal/shared/constants"
)

// commandInventory displays the player's current inventory of items, including Pokeballs.
func commandInventory(cfg *Config, args ...string) error {
	if len(cfg.Inventory) == 0 {
		fmt.Printf("%sYour inventory is empty.%s\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}

	fmt.Printf("\n%sYour Inventory:%s\n", constants.ColorBrightCyan, constants.ColorReset)
	foundItems := false
	for itemNameKey, count := range cfg.Inventory {
		if count <= 0 { // Don't display items with zero count
			continue
		}
		foundItems = true
		// Check if it's a known Pokeball to use its color and proper name
		if ballType, ok := KnownPokeballs[itemNameKey]; ok {
			fmt.Printf("  %s- %s%s%s: %s%d%s\n", constants.ColorGray, ballType.Color, ballType.Name, constants.ColorReset, constants.ColorWhite, count, constants.ColorReset)
		} else {
			// For other items that might be added later
			fmt.Printf("  %s- %s%s%s: %s%d%s\n", constants.ColorGray, constants.ColorWhite, itemNameKey, constants.ColorReset, constants.ColorWhite, count, constants.ColorReset)
		}
	}

	if !foundItems {
		fmt.Printf("%sYour inventory has items, but all have a count of zero or less.%s\n", constants.ColorYellow, constants.ColorReset)
	}
	fmt.Println()
	return nil
}
