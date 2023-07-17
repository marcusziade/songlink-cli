package main

import (
	"fmt"
	"time"

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

	stopLoading := showLoadingIndicator()

	err = GetLinks(searchURL)
	if err != nil {
		return fmt.Errorf("error getting links: %w", err)
	}

	stopLoading <- true

	return nil
}

func showLoadingIndicator() chan bool {
	stopLoading := make(chan bool)
	go func() {
		chars := []string{"-", "\\", "|", "/"}
		i := 0
		for {
			select {
			case <-stopLoading:
				fmt.Print("\r")
				return
			default:
				fmt.Printf("\rLoading %s", chars[i])
				i = (i + 1) % len(chars)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	return stopLoading
}
