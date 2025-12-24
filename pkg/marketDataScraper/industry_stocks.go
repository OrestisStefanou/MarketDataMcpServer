package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"io"
	"market_data_mcp_server/pkg/domain"
	"net/http"
)

func scrapeIndustryStocks(industry string) ([]domain.IndustryStock, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/industry/%s/__data.json", industry)
	resp, err := http.Get(url)
	if err != nil {
		return []domain.IndustryStock{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []domain.IndustryStock{}, err
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return []domain.IndustryStock{}, err
	}

	// Extract "nodes" from rawData
	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 3 {
		return []domain.IndustryStock{}, fmt.Errorf("unexpected structure in 'nodes'")
	}

	// Access the second element in "nodes" which contains the data we are interested in
	nodeData, ok := nodes[2].(map[string]interface{})
	if !ok {
		return []domain.IndustryStock{}, fmt.Errorf("unexpected structure in 'nodes[2]'")
	}

	data, ok := nodeData["data"].([]interface{})
	if !ok {
		return []domain.IndustryStock{}, fmt.Errorf("unexpected structure in 'data'")
	}

	dataMap, ok := data[0].(map[string]interface{})
	if !ok {
		return []domain.IndustryStock{}, fmt.Errorf("unexpected structure in 'data[0]'")
	}

	stocksArrayDataIndex, ok := dataMap["data"].(float64)
	if !ok {
		return []domain.IndustryStock{}, fmt.Errorf("unexpected structure for 'data'")
	}

	stocksDataIndicesArray := data[int(stocksArrayDataIndex)].([]interface{})
	if !ok {
		return []domain.IndustryStock{}, fmt.Errorf("unexpected structure for 'data'")
	}

	stocks := make([]domain.IndustryStock, 0, len(stocksDataIndicesArray))
	for i := 0; i < len(stocksDataIndicesArray); i++ {
		stockDataIndex := int(stocksDataIndicesArray[i].(float64))
		stockData := data[stockDataIndex].(map[string]interface{})
		stockSymbolIndex := int(stockData["s"].(float64))
		stockCompanyNameIndex := int(stockData["n"].(float64))
		stockMarketCapIndex := int(stockData["marketCap"].(float64))
		stock := domain.IndustryStock{
			Symbol:      data[stockSymbolIndex].(string),
			CompanyName: data[stockCompanyNameIndex].(string),
			MarketCap:   float32(data[stockMarketCapIndex].(float64)),
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}
