// trending.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type TrendingResponse struct {
	Coins []struct {
		Item struct {
			ID            string  `json:"id"`
			Name          string  `json:"name"`
			Symbol        string  `json:"symbol"`
			MarketCapRank int     `json:"market_cap_rank"`
			Price         float64 `json:"price_btc"`
			MarketCap     string  `json:"market_cap"`
			TotalVolume   string  `json:"total_volume"`
			Data          struct {
				PriceChange24h map[string]float64 `json:"price_change_percentage_24h"`
			} `json:"data"`
		} `json:"item"`
	} `json:"coins"`
}

// InitializeTrendingTable creates and returns a configured table model
func InitializeTrendingTable() table.Model {
	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "Coin", Width: 20},
		{Title: "Symbol", Width: 10},
		{Title: "Price (BTC)", Width: 15},
		{Title: "24h Change (USD)", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	// Set table styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func fetchTrendingCryptos() ([]table.Row, error) {
	resp, err := http.Get("https://api.coingecko.com/api/v3/search/trending")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trending data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var trendingData TrendingResponse
	if err := json.Unmarshal(body, &trendingData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	var rows []table.Row
	for _, coin := range trendingData.Coins {
		priceChange := coin.Item.Data.PriceChange24h["usd"]
		var priceChangeColor lipgloss.Color
		if priceChange > 0 {
			priceChangeColor = lipgloss.Color("2") // Green for positive change
		} else {
			priceChangeColor = lipgloss.Color("1") // Red for negative change
		}

		priceChangeStyled := lipgloss.NewStyle().Foreground(priceChangeColor).Render(fmt.Sprintf("%.2f%%", priceChange))

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", coin.Item.MarketCapRank),
			coin.Item.Name,
			coin.Item.Symbol,
			fmt.Sprintf("%.8f", coin.Item.Price),
			priceChangeStyled,
		})
	}

	return rows, nil
}
