package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"net/http"
)

func scrapeEtfs() ([]domain.Etf, error) {
	url := "https://api.stockanalysis.com/api/screener/e/f?m=s&s=asc&c=s,n,assetClass,aum&i=etf"

	resp, err := http.Get(url)
	if err != nil {
		return []domain.Etf{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []domain.Etf{}, fmt.Errorf("Call to get etfs failed with status: %d", resp.StatusCode)
	}

	// Define an anonymous struct to match the JSON structure
	var apiResponse struct {
		Status int `json:"status"`
		Data   struct {
			Data []struct {
				S          string  `json:"s"`
				N          string  `json:"n"`
				AssetClass string  `json:"assetClass"`
				AUM        float64 `json:"aum"`
			} `json:"data"`
			ResultsCount int `json:"resultsCount"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return []domain.Etf{}, err
	}

	etfs := make([]domain.Etf, 0, len(apiResponse.Data.Data))
	for _, etfData := range apiResponse.Data.Data {
		etf := domain.Etf{
			Symbol:     etfData.S,
			Name:       etfData.N,
			AssetClass: etfData.AssetClass,
			Aum:        float32(etfData.AUM),
		}
		etfs = append(etfs, etf)
	}
	return etfs, nil
}
