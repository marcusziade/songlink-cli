package main

import (
   "context"
   "flag"
   "fmt"
   "os"
   "strings"
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
   {
       Name:        "download",
       Description: "Search for a song or album and download it as mp3 or mp4",
       Execute:     executeDownload,
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
   outFlag := searchCmd.String("out", "downloads", "Output directory for downloaded files")
   debugFlag := searchCmd.Bool("debug", false, "Enable debug logging during download")
	
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
   return HandleSearch(query, searchType, *outFlag, *debugFlag)
}

// executeConfig handles the config subcommand
func executeConfig(args []string) error {
	fmt.Println("Configuring Apple Music API credentials...")
   return RunOnboarding()
}

// executeDownload handles the download subcommand
func executeDownload(args []string) error {
   // Define download flags
   downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
   typeFlag := downloadCmd.String("type", "song", "Type of search: song, album, or both (default: song)")
   formatFlag := downloadCmd.String("format", "mp3", "Download format: mp3 or mp4 (default: mp3)")
   outFlag := downloadCmd.String("out", "downloads", "Output directory for downloaded files")
   debugFlag := downloadCmd.Bool("debug", false, "Enable debug logging (show yt-dlp/ffmpeg output)")

   // Parse flags
   if err := downloadCmd.Parse(args); err != nil {
       return err
   }

   // Get search query
   queryArgs := downloadCmd.Args()
   if len(queryArgs) == 0 {
       return fmt.Errorf("download query required")
   }
   query := strings.Join(queryArgs, " ")

   // Determine search type
   var searchType SearchType
   switch *typeFlag {
   case "song":
       searchType = Song
   case "album":
       searchType = Album
   default:
       searchType = Song
   }

   // Load config
   config, err := LoadConfig()
   if err != nil {
       return fmt.Errorf("error loading config: %w", err)
   }
   if !config.ConfigExists {
       fmt.Println("Apple Music API credentials not found. Let's set them up.")
       if err := RunOnboarding(); err != nil {
           return fmt.Errorf("error during onboarding: %w", err)
       }
       config, err = LoadConfig()
       if err != nil {
           return fmt.Errorf("error loading config after onboarding: %w", err)
       }
   }

   // Create music searcher
   searcher, err := NewMusicSearcher(config)
   if err != nil {
       return fmt.Errorf("error creating music searcher: %w", err)
   }

   // Search for music
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   results, err := searcher.Search(ctx, query, searchType)
   if err != nil {
       return fmt.Errorf("error searching: %w", err)
   }

   // Display results and select
   selected, err := DisplaySearchResults(results)
   if err != nil {
       return fmt.Errorf("error selecting result: %w", err)
   }
   fmt.Printf("\nSelected: %s - %s\n", selected.Name, selected.ArtistName)

   // Download track via YouTube
   fmt.Print("Downloading... ")
   path, err := DownloadTrack(selected.Name, selected.ArtistName, selected.ArtworkURL, *formatFlag, *outFlag, *debugFlag)
   if err != nil {
       return fmt.Errorf("download error: %w", err)
   }
   fmt.Printf("Done. Saved to %s\n", path)
   return nil
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
