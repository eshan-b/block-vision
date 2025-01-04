// search.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

// Coin represents basic cryptocurrency data from the search API
type Coin struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	MarketCapRank int    `json:"market_cap_rank"`
	Thumb         string `json:"thumb"`
}

// APIResponse represents the structure of the search API response
type APIResponse struct {
	Coins []Coin `json:"coins"`
}

// CoinDetails contains detailed information about a specific cryptocurrency
type CoinDetails struct {
	Name                         string     `json:"name"`
	Symbol                       string     `json:"symbol"`
	MarketData                   MarketData `json:"market_data"`
	SentimentVotesUpPercentage   float64    `json:"sentiment_votes_up_percentage"`
	SentimentVotesDownPercentage float64    `json:"sentiment_votes_down_percentage"`
	Links                        Links      `json:"links"`
}

// MarketData contains price and other market-related details
type MarketData struct {
	CurrentPrice map[string]float64 `json:"current_price"`
}

// Links contains relevant external links for a cryptocurrency
type Links struct {
	Whitepaper string `json:"whitepaper"`
}

// Fetch coins from search endpoint based on query
func fetchCoins(query string) ([]list.Item, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/search?query=%s", query)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	items := make([]list.Item, len(apiResp.Coins))
	for i, coin := range apiResp.Coins {
		items[i] = item{
			title: fmt.Sprintf("%s (%s)", coin.Name, strings.ToUpper(coin.Symbol)),
			desc:  coin.ID,
			data:  coin,
		}
	}
	return items, nil
}

// Fetch detailed information for a specific coin
func fetchCoinDetails(id string) (CoinDetails, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		return CoinDetails{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return CoinDetails{}, err
	}

	var details CoinDetails
	if err := json.Unmarshal(body, &details); err != nil {
		return CoinDetails{}, err
	}
	return details, nil
}

// Item represents a list item in the TUI
type item struct {
	title string
	desc  string
	data  Coin
}

// Title returns the title of the item
func (i item) Title() string { return i.title }

// Description returns the description of the item
func (i item) Description() string { return i.desc }

// FilterValue is used for filtering the item in the list
func (i item) FilterValue() string { return i.title }
