package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"net/http"
	"time"
)

func scrapeHistoricalPrices(ticker string, assetClass domain.AssetClass, period domain.Period) (domain.HistoricalPrices, error) {
	var assetClassPrefix string
	var periodPrefix string

	switch assetClass {
	case domain.Stock:
		assetClassPrefix = "s"
	case domain.ETF:
		assetClassPrefix = "e"
	}

	switch period {
	case domain.Period1D:
		periodPrefix = "1D"
	case domain.Period5D:
		periodPrefix = "5D"
	case domain.Period1M:
		periodPrefix = "1M"
	case domain.Period6M:
		periodPrefix = "6M"
	case domain.Period1Y:
		periodPrefix = "1Y"
	case domain.Period5Y:
		periodPrefix = "5Y"
	}

	url := fmt.Sprintf("https://stockanalysis.com/api/charts/%s/%s/%s/l", assetClassPrefix, ticker, periodPrefix)
	resp, err := http.Get(url)
	if err != nil {
		return domain.HistoricalPrices{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.HistoricalPrices{}, fmt.Errorf("Call to get historical prices failed with status: %d", resp.StatusCode)
	}

	// Define an anonymous struct to match the JSON structure
	var apiResponse struct {
		Status int `json:"status"`
		Data   []struct {
			T int64   `json:"t"`
			C float64 `json:"c"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return domain.HistoricalPrices{}, err
	}

	prices := make([]domain.Price, 0, len(apiResponse.Data))
	for _, price := range apiResponse.Data {
		price := domain.Price{
			Date:       time.Unix(price.T, 0),
			ClosePrice: price.C,
		}
		prices = append(prices, price)
	}

	firstPrice := prices[0].ClosePrice
	lastPrice := prices[len(prices)-1].ClosePrice
	percentChange := ((lastPrice - firstPrice) / firstPrice) * 100

	return domain.HistoricalPrices{
		Period:           period,
		Prices:           prices,
		PercentageChange: percentChange,
	}, nil
}
