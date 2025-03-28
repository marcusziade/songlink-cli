package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/atotto/clipboard"
)

var (
	xFlag = flag.Bool("x", false, "Return the song.link URL without surrounding <>")
	dFlag = flag.Bool("d", false, "Return the song.link URL surrounded by <> and the Spotify URL")
	sFlag = flag.Bool("s", false, "Return only the Spotify URL")
)

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Execute     func(args []string) error
}

// Commands available in the application
var commands = []Command{
	{
		Name:        "search",
		Description: "Search for a song or album and get its links",
		Execute:     executeSearch,
	},
	{
		Name:        "config",
		Description: "Configure Apple Music API credentials",
		Execute:     executeConfig,
	},
}

func main() {
	// Define base flags
	flag.Parse()

	// Check if a subcommand is provided
	args := flag.Args()
	if len(args) > 0 {
		subcommand := args[0]
		
		// Find and execute the appropriate command
		for _, cmd := range commands {
			if cmd.Name == subcommand {
				err := cmd.Execute(args[1:])
				if err != nil {
					fmt.Println("An error occurred:", err)
					os.Exit(1)
				}
				return
			}
		}
		
		// If we get here, the subcommand wasn't recognized
		fmt.Printf("Unknown command: %s\n\n", subcommand)
		printUsage()
		os.Exit(1)
	}

	// No subcommand provided, run the default behavior
	err := runDefault()
	if err != nil {
		fmt.Println("An error occurred:", err)
		os.Exit(1)
	}
}

// executeSearch handles the search subcommand
func executeSearch(args []string) error {
	// Define search flags
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	typeFlag := searchCmd.String("type", "song", "Type of search: song, album, or both (default: song)")
	
	// Parse search flags
	if err := searchCmd.Parse(args); err != nil {
		return err
	}
	
	// Get search query
	searchArgs := searchCmd.Args()
	if len(searchArgs) == 0 {
		return fmt.Errorf("search query required")
	}
	
	query := searchArgs[0]
	
	// Determine search type
	var searchType SearchType
	switch *typeFlag {
	case "song":
		searchType = Song
	case "album":
		searchType = Album
	default:
		// Use Both to search for songs and albums
		searchType = Both
	}
	
	// Handle search
	return HandleSearch(query, searchType)
}

// executeConfig handles the config subcommand
func executeConfig(args []string) error {
	fmt.Println("Configuring Apple Music API credentials...")
	return RunOnboarding()
}

// runDefault runs the default behavior (process URL from clipboard)
func runDefault() error {
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

// printUsage prints usage information
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  songlink-cli [flags]                 Process URL from clipboard")
	fmt.Println("  songlink-cli search [flags] <query>  Search for a song or album")
	fmt.Println("  songlink-cli config                  Configure Apple Music API credentials")
	fmt.Println("\nFlags:")
	fmt.Println("  -x  Return the song.link URL without surrounding <>")
	fmt.Println("  -d  Return the song.link URL surrounded by <> and the Spotify URL")
	fmt.Println("  -s  Return only the Spotify URL")
	fmt.Println("\nSearch Flags:")
	fmt.Println("  -type=<type>  Type of search: song, album, or both (default: song)")
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
