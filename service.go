package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type LinksResponse struct {
	PageUrl string `json:"pageUrl"`
}

func main() {
	LinksRequest("https://youtu.be/NObKVa0y9uo")
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
		fmt.Println(nonLocalURL)
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
