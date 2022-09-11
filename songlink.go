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

// The LinksResponse model
type LinksResponse struct {
	PageUrl string `json:"pageUrl"`
}

// Entrypoint of the app.
// Asks the user to paste and confirm a music service URL
// formats the input and passes it to the `LinksRequest` method
func main() {
	clipboard, err := clipboard.ReadAll()
	if err != nil {
		return
	}
	getLinks(clipboard)
}

// Takes a music service URL as input.
// checks if the response is succesful, decodes the json,
// copies the generated song.link URL to the clipboard and prints it to interface
func getLinks(searchURL string) {
	linksRes := LinksResponse{}

	response, err := http.Get(buildURL(searchURL))
	if err != nil {
		log.Fatal(err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(response.Body)
		err := decoder.Decode(&linksRes)
		if err != nil {
			log.Fatal("Error decoding response")
			return
		}
		nonLocalURL := strings.ReplaceAll(linksRes.PageUrl, "/fi", "")
		clipboard.WriteAll(nonLocalURL)
		fmt.Print("\nSuccess ✅\n", nonLocalURL, "\nSong.link URL copied to the clipboard\n\n")
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
