package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"net/http"
)

func scrapeEtfs() ([]domain.Etf, error) {
	url := "https://stockanalysis.com/_api/endpoints/screener/table?type=e&m=s&s=asc&c=s,n,assetClass,aum,expenseRatio&cn=10000&i=etf"

	resp, err := http.Get(url)
	if err != nil {
		return []domain.Etf{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []domain.Etf{}, fmt.Errorf("Call to get etfs failed with status: %d", resp.StatusCode)
	}

	var apiResponse struct {
		Status int `json:"status"`
		Data   struct {
			Data []struct {
				S            string  `json:"s"`
				N            string  `json:"n"`
				AssetClass   string  `json:"assetClass"`
				AUM          float64 `json:"aum"`
				ExpenseRatio float64 `json:"expenseRatio"`
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
			Symbol:       etfData.S,
			Name:         etfData.N,
			AssetClass:   etfData.AssetClass,
			Aum:          float32(etfData.AUM),
			ExpenseRatio: float32(etfData.ExpenseRatio),
		}
		etfs = append(etfs, etf)
	}
	return etfs, nil
}
