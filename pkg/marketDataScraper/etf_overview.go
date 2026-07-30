package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"strings"
)

func scrapeEtfOverview(symbol string) (domain.EtfOverview, error) {
	url := fmt.Sprintf("https://stockanalysis.com/etf/%s/__data.json", strings.ToLower(symbol))

	body, err := fetchData(url)
	if err != nil {
		return domain.EtfOverview{}, err
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return domain.EtfOverview{}, err
	}

	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 3 {
		return domain.EtfOverview{}, fmt.Errorf("unexpected structure in 'nodes'")
	}

	// The third node holds the page data, flattened into a single array where
	// every value is an index pointing at another entry of the same array
	nodeData, ok := nodes[2].(map[string]interface{})
	if !ok {
		return domain.EtfOverview{}, fmt.Errorf("unexpected structure in 'nodes[2]'")
	}

	data, ok := nodeData["data"].([]interface{})
	if !ok || len(data) == 0 {
		return domain.EtfOverview{}, fmt.Errorf("unexpected structure in 'data'")
	}

	overviewData, ok := data[0].(map[string]interface{})
	if !ok {
		return domain.EtfOverview{}, fmt.Errorf("unexpected structure in 'data[0]'")
	}

	performance := getMap(data, overviewData["performance"])
	holdingsTable := getMap(data, overviewData["holdingsTable"])
	holdings := getSlice(data, holdingsTable["holdings"])

	etfOverview := domain.EtfOverview{
		Symbol:           symbol,
		Description:      getString(data, overviewData["description"]),
		Aum:              getString(data, overviewData["aum"]),
		Nav:              getString(data, overviewData["nav"]),
		ExpenseRatio:     getString(data, overviewData["expenseRatio"]),
		PeRatio:          getString(data, overviewData["peRatio"]),
		Dps:              getString(data, overviewData["dps"]),
		DividendYield:    getString(data, overviewData["dividendYield"]),
		PayoutRatio:      getString(data, overviewData["payoutRatio"]),
		OneMonthReturn:   getFloat(data, performance["tr1m"]),
		OneYearReturn:    getFloat(data, performance["tr1y"]),
		YearToDateReturn: getFloat(data, performance["trYTD"]),
		FiveYearReturn:   getFloat(data, performance["cagr5y"]),
		TenYearReturn:    getFloat(data, performance["cagr10y"]),
		InceptionReturn:  getFloat(data, performance["cagrMAX"]),
		Beta:             getString(data, overviewData["beta"]),
		NumberOfHoldings: int32(getInt(data, overviewData["holdings"])),
		Website:          getString(data, overviewData["etf_website"]),
		TopHoldings:      make([]domain.EtfHolding, 0, len(holdings)),
	}

	for _, holding := range holdings {
		holdingData := getMap(data, holding)
		if len(holdingData) == 0 {
			continue
		}
		etfOverview.TopHoldings = append(etfOverview.TopHoldings, domain.EtfHolding{
			// Holding symbols come back link-prefixed, e.g. "$LLY"
			Symbol: strings.TrimPrefix(getString(data, holdingData["s"]), "$"),
			Name:   getString(data, holdingData["n"]),
			Weight: getString(data, holdingData["as"]),
		})
	}

	for _, row := range getSlice(data, overviewData["infoTable"]) {
		cells := getSlice(data, row)
		if len(cells) < 2 {
			continue
		}
		switch getString(data, cells[0]) {
		case "Asset Class":
			etfOverview.AssetClass = getString(data, cells[1])
		case "Category":
			etfOverview.Category = getString(data, cells[1])
		}
	}

	return etfOverview, nil
}

func getFloat(data []interface{}, index interface{}) float64 {
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
	return f
}

func getMap(data []interface{}, index interface{}) map[string]interface{} {
	idx, ok := index.(float64)
	if !ok {
		return nil
	}
	i := int(idx)
	if i < 0 || i >= len(data) {
		return nil
	}
	m, ok := data[i].(map[string]interface{})
	if !ok {
		return nil
	}
	return m
}

func getSlice(data []interface{}, index interface{}) []interface{} {
	idx, ok := index.(float64)
	if !ok {
		return nil
	}
	i := int(idx)
	if i < 0 || i >= len(data) {
		return nil
	}
	s, ok := data[i].([]interface{})
	if !ok {
		return nil
	}
	return s
}
