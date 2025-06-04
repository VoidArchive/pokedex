package repl

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/voidarchive/pokedex/internal/pokeapi"
	"github.com/voidarchive/pokedex/internal/shared/constants"
)

const pokedexFilePath = "pokedex.json"

// SaveData encapsulates all data that needs to be persisted.
type SaveData struct {
	PokedexData map[string]pokeapi.UserPokemon `json:"pokedex"`
	PartyData   []pokeapi.UserPokemon          `json:"party"`
}

// savePokedex serializes the user's Pokedex and Party to a JSON file.
// This is typically called when the application exits.
func savePokedex(pokedex map[string]pokeapi.UserPokemon, party []pokeapi.UserPokemon) error {
	saveFile := SaveData{
		PokedexData: pokedex,
		PartyData:   party,
	}
	data, err := json.MarshalIndent(saveFile, "", "  ")
	if err != nil {
		return fmt.Errorf("%sfailed to marshal save data: %w%s", constants.ColorRed, err, constants.ColorReset)
	}

	err = os.WriteFile(pokedexFilePath, data, 0644) // 0644 provides read/write for owner, read for others.
	if err != nil {
		return fmt.Errorf("%sfailed to write save file: %w%s", constants.ColorRed, err, constants.ColorReset)
	}
	fmt.Printf("%sGame data saved.%s\n", constants.ColorGreen, constants.ColorReset)
	return nil
}

// loadPokedex deserializes the Pokedex and Party from a JSON file.
// If the file doesn't exist or is malformed, it returns an empty Pokedex and Party.
func loadPokedex() (map[string]pokeapi.UserPokemon, []pokeapi.UserPokemon, error) {
	pokedex := make(map[string]pokeapi.UserPokemon)
	var party []pokeapi.UserPokemon

	data, err := os.ReadFile(pokedexFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%sNo saved data found (%s). Starting fresh.%s\n", constants.ColorYellow, pokedexFilePath, constants.ColorReset)
			return pokedex, party, nil // Return empty structures, no error for non-existence
		}
		return nil, nil, fmt.Errorf("%sfailed to read save file: %w%s", constants.ColorRed, err, constants.ColorReset)
	}

	if len(data) == 0 {
		fmt.Printf("%sSave file (%s) is empty. Starting fresh.%s\n", constants.ColorYellow, pokedexFilePath, constants.ColorReset)
		return pokedex, party, nil // Return empty structures, no error for empty file
	}

	var saveData SaveData
	err = json.Unmarshal(data, &saveData)
	if err != nil {
		return nil, nil, fmt.Errorf("%sfailed to unmarshal save data: %w%s", constants.ColorRed, err, constants.ColorReset)
	}

	// Ensure Pokedex isn't nil if it was missing in JSON, though MarshalIndent should handle empty maps.
	if saveData.PokedexData == nil {
		saveData.PokedexData = make(map[string]pokeapi.UserPokemon)
	}
	// Party can be nil if empty, which is fine for an empty slice.

	fmt.Printf("%sGame data loaded from %s%s%s\n", constants.ColorGreen, constants.ColorYellow, pokedexFilePath, constants.ColorReset)
	return saveData.PokedexData, saveData.PartyData, nil
}
