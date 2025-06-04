package battle_test

import (
	"math/rand"
	"testing"

	"github.com/voidarchive/pokedex/internal/battle"
)

func TestCalculateDamage(t *testing.T) {
	// Use a fixed seed for the test's random number generator to ensure reproducible test outcomes.
	// This local RNG is then passed to battle.CalculateDamage.
	rngSource := rand.NewSource(1) // Fixed seed for determinism
	r := rand.New(rngSource)

	tests := []struct {
		name        string
		attack      int
		defense     int
		expectedMin int
		expectedMax int
	}{
		{
			name:        "Normal case",
			attack:      50,
			defense:     30,
			expectedMin: 20,
			expectedMax: 29,
		},
		{
			name:        "High attack, low defense",
			attack:      100,
			defense:     10,
			expectedMin: 80,
			expectedMax: 119,
		},
		{
			name:        "Low attack, high defense",
			attack:      10,
			defense:     100,
			expectedMin: 1,
			expectedMax: 1,
		},
		{
			name:        "Zero defense",
			attack:      50,
			defense:     0,
			expectedMin: 72,
			expectedMax: 107,
		},
		{
			name:        "Zero attack",
			attack:      0,
			defense:     50,
			expectedMin: 1,
			expectedMax: 1,
		},
		{
			name:        "High defense makes damage 1",
			attack:      20,
			defense:     500,
			expectedMin: 1,
			expectedMax: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// The 100 iterations are to test the *output range* of CalculateDamage
			// given fixed inputs and its internal randomness, using the deterministic RNG.
			for i := range 100 {
				damage := battle.CalculateDamage(r, tc.attack, tc.defense)
				if damage < tc.expectedMin || damage > tc.expectedMax {
					t.Errorf("battle.CalculateDamage(r, %d, %d) = %d; want between %d and %d (iteration %d)", tc.attack, tc.defense, damage, tc.expectedMin, tc.expectedMax, i)
					break
				}
			}
		})
	}
}
