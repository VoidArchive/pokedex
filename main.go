package main

import "strings"

func main() {
	println("Hello, World!")
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
