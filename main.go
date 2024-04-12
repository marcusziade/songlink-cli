package main

import (
    "flag"
    "fmt"
    "sync"
    "time"

    "github.com/atotto/clipboard"
)

var (
    xFlag = flag.Bool("x", false, "Return the song.link URL without surrounding <>")
    dFlag = flag.Bool("d", false, "Return the song.link URL surrounded by <> and the Spotify URL")
    sFlag = flag.Bool("s", false, "Return only the Spotify URL")
)

func main() {
    flag.Parse()

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

	var wg sync.WaitGroup
	wg.Add(1)
	stopLoading := make(chan bool)
	go func() {
		defer wg.Done()
		loadingIndicator(stopLoading)
	}()

	err = GetLinks(searchURL)
	if err != nil {
		return fmt.Errorf("error getting links: %w", err)
	}

	stopLoading <- true
	wg.Wait()

	return nil
}

func loadingIndicator(stop chan bool) {
	chars := []string{"-", "\\", "|", "/"}
	i := 0
	for {
		select {
		case <-stop:
			fmt.Print("\r")
			return
		default:
			fmt.Printf("\rLoading %s", chars[i])
			i = (i + 1) % len(chars)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
