package pokeapi_test

import (
	"testing"

	"github.com/voidarchive/pokedex/internal/pokeapi"
)

func TestPokemonData_GetStat(t *testing.T) {
	pokemon := pokeapi.PokemonData{
		Name: "Pikachu",
		Stats: []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		}{
			{BaseStat: 35, Stat: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "hp"}},
			{BaseStat: 55, Stat: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "attack"}},
			{BaseStat: 40, Stat: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "defense"}},
		},
	}

	tests := []struct {
		name          string
		statName      string
		expectedValue int
		expectedFound bool
	}{
		{
			name:          "Get HP",
			statName:      "hp",
			expectedValue: 35,
			expectedFound: true,
		},
		{
			name:          "Get Attack",
			statName:      "attack",
			expectedValue: 55,
			expectedFound: true,
		},
		{
			name:          "Get Defense",
			statName:      "defense",
			expectedValue: 40,
			expectedFound: true,
		},
		{
			name:          "Get NonExistent Stat",
			statName:      "speed",
			expectedValue: 0,
			expectedFound: false,
		},
		{
			name:          "Get Stat with different case",
			statName:      "HP", // Test case insensitivity if desired (current GetStat is case-sensitive)
			expectedValue: 0,    // Assuming GetStat is case-sensitive as implemented
			expectedFound: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			value, found := pokemon.GetStat(tc.statName)
			if found != tc.expectedFound {
				t.Errorf("GetStat(%s) found = %v; want %v", tc.statName, found, tc.expectedFound)
			}
			if value != tc.expectedValue {
				t.Errorf("GetStat(%s) value = %d; want %d", tc.statName, value, tc.expectedValue)
			}
		})
	}
}
