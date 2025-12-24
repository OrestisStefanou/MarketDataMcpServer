package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"io"
	"market_data_mcp_server/pkg/domain"
	"net/http"
)

func scrapeIndustries() ([]domain.Industry, error) {
	url := "https://stockanalysis.com/stocks/industry/all/__data.json"
	resp, err := http.Get(url)
	if err != nil {
		return []domain.Industry{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []domain.Industry{}, err
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return []domain.Industry{}, err
	}

	// Extract "nodes" from rawData
	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 3 {
		return []domain.Industry{}, fmt.Errorf("unexpected structure in 'nodes'")
	}

	// Access the second element in "nodes" which contains the data we are interested in
	nodeData, ok := nodes[2].(map[string]interface{})
	if !ok {
		return []domain.Industry{}, fmt.Errorf("unexpected structure in 'nodes[2]'")
	}

	data, ok := nodeData["data"].([]interface{})
	if !ok {
		return []domain.Industry{}, fmt.Errorf("unexpected structure in 'data'")
	}

	dataMap, ok := data[0].(map[string]interface{})
	if !ok {
		return []domain.Industry{}, fmt.Errorf("unexpected structure in 'data[0]'")
	}

	industriesDataIndex, ok := dataMap["industries"].(float64)
	if !ok {
		return []domain.Industry{}, fmt.Errorf("unexpected structure for 'industries'")
	}

	industryDataIndicesArray := data[int(industriesDataIndex)].([]interface{})

	industries := make([]domain.Industry, 0, len(industryDataIndicesArray))
	for i := 0; i < len(industryDataIndicesArray); i++ {
		industryDataIndex := int(industryDataIndicesArray[i].(float64))
		industryData := data[industryDataIndex].(map[string]interface{})
		industryNameIndex := int(industryData["industry_name"].(float64))
		industryUrlNameIndex := int(industryData["url"].(float64))
		numberOfStocksIndex := int(industryData["stocks"].(float64))
		marketCapIndex := int(industryData["marketCap"].(float64))
		profitMarginIndex := int(industryData["profitMargin"].(float64))
		oneYearChangeIndex := int(industryData["ch1y"].(float64))

		// peRatio and dividendYield are handled differently because they could be missing from the industryData map
		var peRatio float32
		peRatioIndex, ok := industryData["peRatio"]
		if !ok {
			peRatio = 0
		} else {
			peRatioIndexInt := int(peRatioIndex.(float64))
			peRatio = float32(data[peRatioIndexInt].(float64))
		}

		var dividendYield float32
		dividendYieldIndex, ok := industryData["dividendYield"]
		if !ok {
			dividendYield = 0
		} else {
			dividendYieldIndexInt := int(dividendYieldIndex.(float64))
			dividendYield = float32(data[dividendYieldIndexInt].(float64))
		}

		industry := domain.Industry{
			Name:             data[industryNameIndex].(string),
			UrlName:          data[industryUrlNameIndex].(string),
			NumberOfStocks:   int(data[numberOfStocksIndex].(float64)),
			MarketCap:        float32(data[marketCapIndex].(float64)),
			DividendYieldPct: dividendYield,
			PeRatio:          peRatio,
			ProfitMarginPct:  float32(data[profitMarginIndex].(float64)),
			OneYearChangePct: float32(data[oneYearChangeIndex].(float64)),
		}
		industries = append(industries, industry)
	}

	return industries, nil
}
