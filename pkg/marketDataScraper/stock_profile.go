package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"io"
	"market_data_mcp_server/pkg/domain"
	"net/http"
)

func scrapeStockProfile(symbol string) (domain.StockProfile, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/company/__data.json", symbol)

	resp, err := http.Get(url)
	if err != nil {
		return domain.StockProfile{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.StockProfile{}, err
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return domain.StockProfile{}, err
	}

	// Extract "nodes" from rawData
	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 3 {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure in 'nodes'")
	}

	// Access the second element in "nodes" which contains the data we are interested in
	nodeData, ok := nodes[2].(map[string]interface{})
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure in 'nodes[2]'")
	}

	data, ok := nodeData["data"].([]interface{})
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure in 'data'")
	}

	dataMap, ok := data[0].(map[string]interface{})
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure in 'data[0]'")
	}

	descriptionIndex, ok := dataMap["description"].(float64)
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure for 'description'")
	}

	profileIndex, ok := dataMap["profile"].(float64)
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure for 'profile'")
	}

	if int(profileIndex) < 0 || int(profileIndex) >= len(data) {
		return domain.StockProfile{}, fmt.Errorf("profile index out of bounds")
	}
	stockProfileData, ok := data[int(profileIndex)].(map[string]interface{})
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure for stockProfileData")
	}

	industryDataIndex, ok := stockProfileData["industry"].(float64)
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure for 'industry'")
	}

	if int(industryDataIndex) < 0 || int(industryDataIndex) >= len(data) {
		return domain.StockProfile{}, fmt.Errorf("industry index out of bounds")
	}
	industryData, ok := data[int(industryDataIndex)].(map[string]interface{})
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure for industryData")
	}
	industryNameIndex, _ := industryData["value"].(float64)

	sectorDataIndex, ok := stockProfileData["sector"].(float64)
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure for 'sector'")
	}
	if int(sectorDataIndex) < 0 || int(sectorDataIndex) >= len(data) {
		return domain.StockProfile{}, fmt.Errorf("sector index out of bounds")
	}
	sectorData, ok := data[int(sectorDataIndex)].(map[string]interface{})
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure for sectorData")
	}
	sectorNameIndex, _ := sectorData["value"].(float64)

	stockNameINdex, _ := stockProfileData["name"].(float64)
	stockCountryIndex, _ := stockProfileData["country"].(float64)
	stockFoundedIndex, _ := stockProfileData["founded"].(float64)
	stockIpoDateIndex, _ := stockProfileData["ipoDate"].(float64)
	stockCeoIndex, _ := stockProfileData["ceo"].(float64)

	return domain.StockProfile{
		Name:        getString(data, stockNameINdex),
		Description: getString(data, descriptionIndex),
		Country:     getString(data, stockCountryIndex),
		Founded:     getInt(data, stockFoundedIndex),
		IpoDate:     getString(data, stockIpoDateIndex),
		Industry:    getString(data, industryNameIndex),
		Sector:      getString(data, sectorNameIndex),
		Ceo:         getString(data, stockCeoIndex),
	}, nil
}

func getString(data []interface{}, index interface{}) string {
	idx, ok := index.(float64)
	if !ok {
		return ""
	}
	i := int(idx)
	if i < 0 || i >= len(data) {
		return ""
	}
	s, ok := data[i].(string)
	if !ok {
		return ""
	}
	return s
}

func getInt(data []interface{}, index interface{}) int {
	idx, ok := index.(float64)
	if !ok {
		return 0
	}
	i := int(idx)
	if i < 0 || i >= len(data) {
		return 0
	}
	f, ok := data[i].(float64)
	if !ok {
		return 0
	}
	return int(f)
}
