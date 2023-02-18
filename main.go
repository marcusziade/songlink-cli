package main

import (
	"fmt"

	"github.com/atotto/clipboard"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("An error occurred:", err)
	}
}

func run() error {
	searchURL, err := clipboard.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading clipboard: %w", err)
	}
	err = GetLinks(searchURL)
	if err != nil {
		return fmt.Errorf("error getting links: %w", err)
	}
	return nil
}
