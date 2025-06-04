package repl

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/voidarchive/pokedex/internal/pokeapi"
	"github.com/voidarchive/pokedex/internal/shared/constants"
)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"map": {
			Name:        "map",
			Description: "Displays the next 20 Pokemon location areas",
			Callback:    commandMapf,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays the previous 20 Pokemon location areas",
			Callback:    commandMapb,
		},
		"explore": {
			Name:        "explore <location_area_name>",
			Description: "Explore a location area for Pokemon",
			Callback:    commandExplore,
		},
		"catch": {
			Name:        "catch <pokemon_name>",
			Description: "Attempt to catch a Pokemon",
			Callback:    commandCatch,
		},
		"inspect": {
			Name:        "inspect <pokemon_name>",
			Description: "View details of a caught Pokemon",
			Callback:    commandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "View all Pokemon in your Pokedex",
			Callback:    commandPokedex,
		},
		"party": {
			Name:        "party",
			Description: "View the Pokemon in your active party",
			Callback:    commandParty,
		},
		"inventory": {
			Name:        "inventory",
			Description: "View your items, including Pokeballs",
			Callback:    commandInventory,
		},
		"battle": {
			Name:        "battle <your_pokemon> <opponent_pokemon>",
			Description: "Simulate a battle between one of your Pokemon and an opponent",
			Callback:    commandBattle,
		},
	}
}

func stringToPtr(s string) *string {
	return &s
}

func StartRepl(pokeapiClient PokeapiClient, cache Pokecache) {
	loadedPokedex, loadedParty, err := loadPokedex()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError loading saved data: %v. Starting fresh.%s\n", constants.ColorRed, err, constants.ColorReset)
		loadedPokedex = make(map[string]pokeapi.UserPokemon)
		loadedParty = []pokeapi.UserPokemon{}
	} else {
		if loadedParty == nil {
			loadedParty = []pokeapi.UserPokemon{}
		}
	}

	initialInventory := map[string]int{
		"pokeball":  10,
		"greatball": 5,
	}

	cfg := &Config{
		NextLocationAreaURL: stringToPtr(pokeapi.BaseURL + "/location-area"),
		PrevLocationAreaURL: nil,
		PokeapiClient:       pokeapiClient,
		Cache:               cache,
		Pokedex:             loadedPokedex,
		Party:               loadedParty,
		Inventory:           initialInventory,
		Randomizer:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          fmt.Sprintf("%sPokedex > %s", constants.ColorCyan, constants.ColorReset),
		HistoryFile:     "/tmp/pokedex_history.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Printf("%sError creating readline instance: %v%s\n", constants.ColorRed, err, constants.ColorReset)
		return
	}
	defer rl.Close()

	commands := getCommands()
	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt { // Ctrl+C
			fmt.Printf("\n%sInterrupt received, saving game data and exiting...%s\n", constants.ColorYellow, constants.ColorReset)
			if saveErr := savePokedex(cfg.Pokedex, cfg.Party); saveErr != nil {
				fmt.Fprintf(os.Stderr, "%sError saving game data on interrupt: %v%s\n", constants.ColorRed, saveErr, constants.ColorReset)
			}
			return
		} else if err == io.EOF { // Ctrl+D
			fmt.Printf("\n%sEOF received, saving game data and exiting...%s\n", constants.ColorYellow, constants.ColorReset)
			if saveErr := savePokedex(cfg.Pokedex, cfg.Party); saveErr != nil {
				fmt.Fprintf(os.Stderr, "%sError saving game data on EOF: %v%s\n", constants.ColorRed, saveErr, constants.ColorReset)
			}
			return
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "%sError reading input: %v%s\n", constants.ColorRed, err, constants.ColorReset)
			continue
		}

		if line == "" {
			continue
		}

		userInput := CleanInput(line)
		if len(userInput) == 0 {
			continue
		}

		commandName := userInput[0]
		args := []string{}
		if len(userInput) > 1 {
			args = userInput[1:]
		}

		command, exists := commands[commandName]
		if exists {
			err := command.Callback(cfg, args...)
			if err != nil {
				fmt.Printf("%s%v%s\n", constants.ColorRed, err, constants.ColorReset)
			}
		} else {
			fmt.Printf("%sUnknown command: %s%s%s\n", constants.ColorRed, commandName, constants.ColorRed, constants.ColorReset)
		}
	}
}

func CleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
