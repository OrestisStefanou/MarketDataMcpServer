package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"net/http"
	"net/url"
)

func scrapeStockSearch(query string) ([]domain.Ticker, error) {
	endpoint := fmt.Sprintf("https://stockanalysis.com/api/search?q=%s", url.QueryEscape(query))

	resp, err := http.Get(endpoint)
	if err != nil {
		return []domain.Ticker{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []domain.Ticker{}, fmt.Errorf("Call to stock search failed with status: %d", resp.StatusCode)
	}

	var apiResponse struct {
		Status int `json:"status"`
		Data   []struct {
			S string `json:"s"`
			T string `json:"t"`
			N string `json:"n"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return []domain.Ticker{}, err
	}

	tickers := make([]domain.Ticker, 0, len(apiResponse.Data))
	for _, row := range apiResponse.Data {
		if row.T != "s" {
			continue
		}
		tickers = append(tickers, domain.Ticker{
			Symbol:      row.S,
			CompanyName: row.N,
		})
	}
	return tickers, nil
}
