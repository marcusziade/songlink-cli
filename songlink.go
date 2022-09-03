package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
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
	fmt.Print("Enter search URL...\n")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return
	}
	input = strings.TrimSuffix(input, "\n")
	LinksRequest(input)
}

// Takes a music service URL as input.
// checks if the response is succesful, decodes the json,
// copies the generated song.link URL to the clipboard and prints it to interface
func LinksRequest(searchURL string) {
	linksRes := LinksResponse{}

	response, err := http.Get(buildURL(searchURL))
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(response.Body)
		err := decoder.Decode(&linksRes)
		if err != nil {
			panic(err)
		}
		nonLocalURL := strings.ReplaceAll(linksRes.PageUrl, "/fi", "")
		clipboard.WriteAll(nonLocalURL)
		fmt.Print("\nSuccess ✅\n", nonLocalURL, "\nSong.link URL copied to the clipboard")
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
