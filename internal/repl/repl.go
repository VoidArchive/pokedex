package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/voidarchive/pokedex/internal/pokeapi"
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
			Name:        "explore",
			Description: "Explore a location area for Pokemon",
			Callback:    commandExplore,
		},
		"catch": {
			Name:        "catch <pokemon_name>",
			Description: "Attempt to catch a Pokemon",
			Callback:    commandCatch,
		},
	}
}

func stringToPtr(s string) *string {
	return &s
}

func StartRepl(pokeapiClient PokeapiClient, cache Pokecache) {
	cfg := &Config{
		NextLocationAreaURL: stringToPtr(pokeapi.BaseURL + "/location-area"),
		PrevLocationAreaURL: nil,
		PokeapiClient:       pokeapiClient,
		Cache:               cache,
		Pokedex:             make(map[string]pokeapi.PokemonData),
	}

	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading input:", err)
			}
			break
		}

		userInput := CleanInput(scanner.Text())
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
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func CleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
