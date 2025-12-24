package marketDataScraper

import (
	"encoding/json"
	"io"
	"market_data_mcp_server/pkg/domain"
	"net/http"
)

func scrapeSectors() ([]domain.Sector, error) {
	url := "https://stockanalysis.com/stocks/industry/sectors/__data.json"
	resp, err := http.Get(url)
	if err != nil {
		return []domain.Sector{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []domain.Sector{}, err
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return []domain.Sector{}, err
	}

	// Extract "nodes" from rawData
	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 3 {
		return []domain.Sector{}, err
	}

	// Access the second element in "nodes" which contains the data we are interested in
	nodeData, ok := nodes[2].(map[string]interface{})
	if !ok {
		return []domain.Sector{}, err
	}

	data, ok := nodeData["data"].([]interface{})
	if !ok {
		return []domain.Sector{}, err
	}

	dataMap, ok := data[0].(map[string]interface{})
	if !ok {
		return []domain.Sector{}, err
	}

	sectorsDataIndex, ok := dataMap["sectors"].(float64)
	if !ok {
		return []domain.Sector{}, err
	}

	sectorDataIndicesArray := data[int(sectorsDataIndex)].([]interface{})

	sectors := make([]domain.Sector, 0, len(sectorDataIndicesArray))
	for i := 0; i < len(sectorDataIndicesArray); i++ {
		sectorDataIndex := int(sectorDataIndicesArray[i].(float64))
		sectorData := data[sectorDataIndex].(map[string]interface{})
		sectorNameIndex := int(sectorData["sector_name"].(float64))
		sectorUrlNameIndex := int(sectorData["url"].(float64))
		numberOfStocksIndex := int(sectorData["stocks"].(float64))
		marketCapIndex := int(sectorData["marketCap"].(float64))
		dividendYieldIndex := int(sectorData["dividendYield"].(float64))
		peRatioIndex := int(sectorData["peRatio"].(float64))
		profitMarginIndex := int(sectorData["profitMargin"].(float64))
		oneYearChangeIndex := int(sectorData["ch1y"].(float64))

		sector := domain.Sector{
			Name:             data[sectorNameIndex].(string),
			UrlName:          data[sectorUrlNameIndex].(string),
			NumberOfStocks:   int(data[numberOfStocksIndex].(float64)),
			MarketCap:        float32(data[marketCapIndex].(float64)),
			DividendYieldPct: float32(data[dividendYieldIndex].(float64)),
			PeRatio:          float32(data[peRatioIndex].(float64)),
			ProfitMarginPct:  float32(data[profitMarginIndex].(float64)),
			OneYearChangePct: float32(data[oneYearChangeIndex].(float64)),
		}
		sectors = append(sectors, sector)
	}
	return sectors, nil
}
