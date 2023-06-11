package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SearchResult struct {
	Data []struct {
		AuthorName string `json:"author_name"`
		Title      string `json:"title"`
		Content    string `json:"content"`
	} `json:"data"`
}

func SearchCode(query, apiKey string) (*SearchResult, error) {
	baseURL := "https://api.grepper.com/v1/answers/search"
	params := url.Values{
		"query": {query},
	}
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(apiKey+":")))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var result SearchResult
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
}

func main() {
	query := "javascript loop array backwards"
	apiKey := "<INSERT API KEY>"
	results, err := SearchCode(query, apiKey)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Results:")
	for _, result := range results.Data {
		fmt.Printf("By %s\n", result.AuthorName)
		fmt.Println(result.Title)
		fmt.Println("--------------------------------")
		fmt.Println(result.Content)
		fmt.Println()
	}
}
