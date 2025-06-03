package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		}, {
			input:    "  leading and trailing spaces  ",
			expected: []string{"leading", "and", "trailing", "spaces"},
		}, {
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("for input %q, expected length %d but got %d", c.input, len(c.expected), len(actual))
			continue
		}

		for i := range actual {
			word := actual[i]
			expected := c.expected[i]

			if word != expected {
				t.Errorf("for input %q, expected %q at index %d but got %q", c.input, c.expected[i], i, actual[i])
			}
		}

	}
}
