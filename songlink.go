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

type LinksResponse struct {
	PageUrl string `json:"pageUrl"`
}

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
