package repl

import (
	"fmt"
)

func commandPokedex(cfg *Config, args ...string) error {
	if len(cfg.Pokedex) == 0 {
		fmt.Println("Your Pokedex is empty. Go catch some Pokemon!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range cfg.Pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}
