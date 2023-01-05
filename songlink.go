package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/atotto/clipboard"
)

type SonglinkResponse struct {
	PageURL         string          `json:"pageUrl"`
	LinksByPlatform LinksByPlatform `json:"linksByPlatform"`
}

type LinksByPlatform struct {
	Spotify PlatformMusic `json:"spotify"`
}

type PlatformMusic struct {
	URL string `json:"url"`
}

// Entrypoint of the app.
// Asks the user to paste and confirm a music service URL
// formats the input and passes it to the `LinksRequest` method
func main() {
	clipboard, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return
	}
	getLinks(clipboard)
}

// Used to fetch a song.link and Spotify URL for a music link in the clipboard.
// Copies two things to the clipboard:
// 1. A <> wrapped song.link url (This disables embedding in Discord)
// 2. A Spotify URL if the song is available on Spotify. This URL will enable the Spotify embed player.
func getLinks(searchURL string) {
	platform := PlatformMusic{
		URL: "",
	}
	links := LinksByPlatform{
		Spotify: platform,
	}
	linksResponse := SonglinkResponse{
		PageURL:         "",
		LinksByPlatform: links,
	}

	response, err := http.Get(buildURL(searchURL))
	if err != nil {
		log.Fatal(err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(response.Body)
		err := decoder.Decode(&linksResponse)
		if err != nil {
			log.Fatal("Error decoding response")
			return
		}
		
		nonLocalURL := strings.ReplaceAll(linksResponse.PageURL, "/fi", "")
		spotifyURL := linksResponse.LinksByPlatform.Spotify.URL
		outputString := fmt.Sprintf("<%s>\n\n%s", nonLocalURL, spotifyURL)

		clipboard.WriteAll(outputString)

		fmt.Print(
			"\nSuccess ✅\n",
			outputString,
			"\nCopied to the clipboard\n\n")
	} else {
		fmt.Println("\n❌", response.Status, "Check the search URL and retry.")
	}
}

// Takes in a music service URL and builds the song.link API query
func buildURL(searchURL string) string {
	url := url.URL{
		Scheme: "https",
		Host:   "api.song.link",
		Path:   "/v1-alpha.1/links",
	}
	values := url.Query()
	values.Add("url", searchURL)
	url.RawQuery = values.Encode()
	return url.String()
}
