package marketDataScraper

import (
	"encoding/json"
	"market_data_mcp_server/pkg/domain"
	"net/http"
)

func scrapeStockList() ([]domain.Ticker, error) {
	url := "https://stockanalysis.com/api/screener/s/f?m=s&s=asc&c=s,n&i=stocks"

	resp, err := http.Get(url)
	if err != nil {
		return []domain.Ticker{}, err
	}
	defer resp.Body.Close()

	// Define an anonymous struct to match the JSON structure
	var apiResponse struct {
		Status int `json:"status"`
		Data   struct {
			Data []struct {
				S string `json:"s"`
				N string `json:"n"`
			} `json:"data"`
			ResultsCount int `json:"resultsCount"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return []domain.Ticker{}, err
	}

	tickers := make([]domain.Ticker, 0, len(apiResponse.Data.Data))
	for _, tickerData := range apiResponse.Data.Data {
		ticker := domain.Ticker{
			Symbol:      tickerData.S,
			CompanyName: tickerData.N,
		}
		tickers = append(tickers, ticker)
	}

	return tickers, nil
}
