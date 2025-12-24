package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"io"
	"market_data_mcp_server/pkg/domain"
	"net/http"
)

func scrapeStockForecast(symbol string) (domain.StockForecast, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/forecast/__data.json", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return domain.StockForecast{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.StockForecast{}, err
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return domain.StockForecast{}, err
	}

	// Accessing nodes data
	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 3 {
		return domain.StockForecast{}, fmt.Errorf("invalid response structure: nodes missing or too short")
	}

	nodeData, ok := nodes[2].(map[string]interface{})
	if !ok {
		return domain.StockForecast{}, fmt.Errorf("invalid response structure: node[2] is not a map")
	}

	data, ok := nodeData["data"].([]interface{})
	if !ok || len(data) == 0 {
		return domain.StockForecast{}, fmt.Errorf("invalid response structure: missing or empty data")
	}

	dataMap, ok := data[0].(map[string]interface{})
	if !ok {
		return domain.StockForecast{}, fmt.Errorf("invalid response structure: data[0] is not a map")
	}

	// Estimates Scraping
	quarterlyEstimatesData := make(map[string][]interface{})

	estimatesIdxVal, ok := dataMap["estimates"]
	if !ok {
		return domain.StockForecast{}, fmt.Errorf("missing estimates index")
	}
	estimatesDataIndex := int(estimatesIdxVal.(float64))

	estimatesMap, ok := data[estimatesDataIndex].(map[string]interface{})
	if !ok {
		return domain.StockForecast{}, fmt.Errorf("invalid estimates data")
	}

	tableIdxVal, ok := estimatesMap["table"]
	if !ok {
		return domain.StockForecast{}, fmt.Errorf("missing table index")
	}
	estimatesTableDataIndex := int(tableIdxVal.(float64))

	tableMap, ok := data[estimatesTableDataIndex].(map[string]interface{})
	if !ok {
		return domain.StockForecast{}, fmt.Errorf("invalid table data")
	}

	quarterlyIdxVal, ok := tableMap["quarterly"]
	if !ok {
		return domain.StockForecast{}, fmt.Errorf("missing quarterly index")
	}
	quarterlyEstimatesDataIndex := int(quarterlyIdxVal.(float64))

	quarterlyEstimatesDataMap, ok := data[quarterlyEstimatesDataIndex].(map[string]interface{})
	if !ok {
		return domain.StockForecast{}, fmt.Errorf("invalid quarterly estimates data")
	}

	for estimationField, estimationFieldIdx := range quarterlyEstimatesDataMap {
		if estimationField == "lastDate" {
			continue
		}
		idx, ok := estimationFieldIdx.(float64)
		if !ok {
			continue
		}
		fieldData, ok := data[int(idx)].([]interface{})
		if !ok {
			continue
		}

		var estimationFieldValues []interface{}
		for _, fieldValueIndex := range fieldData {
			fvIdx, ok := fieldValueIndex.(float64)
			if !ok {
				continue
			}
			if int(fvIdx) >= len(data) {
				continue
			}
			fieldValue := data[int(fvIdx)]
			if fieldValue == "[PRO]" {
				continue
			}
			estimationFieldValues = append(estimationFieldValues, fieldValue)
		}
		quarterlyEstimatesData[estimationField] = estimationFieldValues
	}

	// Prepare estimations_doc in the same format as in Python
	var estimationsDoc []map[string]interface{}
	epsData, ok := quarterlyEstimatesData["eps"]
	if ok {
		for i := 0; i < len(epsData); i++ {
			record := make(map[string]interface{})
			for key, values := range quarterlyEstimatesData {
				if i < len(values) {
					record[key] = values[i]
				}
			}
			estimationsDoc = append(estimationsDoc, record)
		}
	}

	// Create a slice of StockEstimation structs from estimationsDoc
	estimations := make([]domain.StockEstimation, 0, len(estimationsDoc))
	for _, record := range estimationsDoc {
		var date string
		if record["dates"] != nil {
			if d, ok := record["dates"].(string); ok {
				date = d
			}
		}

		var eps float64
		if record["eps"] != nil {
			if e, ok := record["eps"].(float64); ok {
				eps = e
			}
		}

		var epsGrowth float64
		if record["epsGrowth"] != nil {
			if eg, ok := record["epsGrowth"].(float64); ok {
				epsGrowth = eg
			}
		}

		var revenue float64
		if record["revenue"] != nil {
			if r, ok := record["revenue"].(float64); ok {
				revenue = r
			}
		}

		var revenueGrowth float64
		if record["revenueGrowth"] != nil {
			if rg, ok := record["revenueGrowth"].(float64); ok {
				revenueGrowth = rg
			}
		}

		var fiscalQuarter string
		if fq, ok := record["fiscalQuarter"].(string); ok {
			fiscalQuarter = fq
		}

		var fiscalYear string
		if fy, ok := record["fiscalYear"].(string); ok {
			fiscalYear = fy
		}

		estimation := domain.StockEstimation{
			Date:          date,
			Eps:           eps,
			EpsGrowth:     epsGrowth,
			FiscalQuarter: fiscalQuarter,
			FiscalYear:    fiscalYear,
			Revenue:       revenue,
			RevenueGrowth: revenueGrowth,
		}
		estimations = append(estimations, estimation)
	}

	// Target Price Scraping
	targetsIdxVal, ok := dataMap["targets"]
	if !ok {
		return domain.StockForecast{Estimations: estimations}, nil
	}
	targetDataIndex := int(targetsIdxVal.(float64))

	targetDataMap, ok := data[targetDataIndex].(map[string]interface{})
	if !ok {
		return domain.StockForecast{Estimations: estimations}, nil
	}

	targetPriceKeys := []string{"average", "high", "low", "median"}
	targetPriceDoc := make(map[string]interface{})
	for _, targetKey := range targetPriceKeys {
		if val, ok := targetDataMap[targetKey]; ok {
			targetValueIndex := int(val.(float64))
			if targetValueIndex < len(data) {
				targetPriceDoc[targetKey] = data[targetValueIndex]
			}
		}
	}

	// Create a StockTargetPrc struct from targetPriceDoc
	var avg, high, low, median float32
	if v, ok := targetPriceDoc["average"].(float64); ok {
		avg = float32(v)
	}
	if v, ok := targetPriceDoc["high"].(float64); ok {
		high = float32(v)
	}
	if v, ok := targetPriceDoc["low"].(float64); ok {
		low = float32(v)
	}
	if v, ok := targetPriceDoc["median"].(float64); ok {
		median = float32(v)
	}

	targetPrice := domain.StockTargetPrc{
		Average: avg,
		High:    high,
		Low:     low,
		Median:  median,
	}

	return domain.StockForecast{
		Estimations: estimations,
		TargetPrice: targetPrice,
	}, nil
}
