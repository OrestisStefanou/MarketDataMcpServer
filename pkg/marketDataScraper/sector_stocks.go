package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"io"
	"market_data_mcp_server/pkg/domain"
	"net/http"
)

func scrapeSectorStocks(sector string) ([]domain.SectorStock, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/sector/%s/__data.json", sector)
	resp, err := http.Get(url)
	if err != nil {
		return []domain.SectorStock{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []domain.SectorStock{}, err
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return []domain.SectorStock{}, err
	}

	// Extract "nodes" from rawData
	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 2 {
		return []domain.SectorStock{}, fmt.Errorf("unexpected structure in 'nodes'")
	}

	// Access the second element in "nodes" which contains the data we are interested in
	nodeData, ok := nodes[2].(map[string]interface{})
	if !ok {
		return []domain.SectorStock{}, fmt.Errorf("unexpected structure in 'nodes[1]'")
	}

	data, ok := nodeData["data"].([]interface{})
	if !ok {
		return []domain.SectorStock{}, fmt.Errorf("unexpected structure in 'data'")
	}

	dataMap, ok := data[0].(map[string]interface{})
	if !ok {
		return []domain.SectorStock{}, fmt.Errorf("unexpected structure in 'data[0]'")
	}

	stocksArrayDataIndex, ok := dataMap["data"].(float64)
	if !ok {
		return []domain.SectorStock{}, fmt.Errorf("unexpected structure for 'data'")
	}

	stocksDataIndicesArray := data[int(stocksArrayDataIndex)].([]interface{})
	if !ok {
		return []domain.SectorStock{}, fmt.Errorf("unexpected structure for 'data[stocksArrayDataIndex]'")
	}

	stocks := make([]domain.SectorStock, 0, len(stocksDataIndicesArray))
	for i := 0; i < len(stocksDataIndicesArray); i++ {
		stockDataIndex := int(stocksDataIndicesArray[i].(float64))
		stockData := data[stockDataIndex].(map[string]interface{})
		stockSymbolIndex := int(stockData["s"].(float64))
		stockCompanyNameIndex := int(stockData["n"].(float64))

		var marketCap float32
		stockMarketCapIndex, ok := stockData["marketCap"].(float64)
		if !ok {
			// Do nothing in this case since there are times that this fied doesn't exist
		} else {
			marketCap = float32(data[int(stockMarketCapIndex)].(float64))
		}
		stock := domain.SectorStock{
			Symbol:      data[stockSymbolIndex].(string),
			CompanyName: data[stockCompanyNameIndex].(string),
			MarketCap:   marketCap,
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}
