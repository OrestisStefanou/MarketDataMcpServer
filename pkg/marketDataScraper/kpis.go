package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"io"
	"market_data_mcp_server/pkg/domain"
	"net/http"
	"reflect"
	"strings"
)

// Root structure of the JSON
type root struct {
	Type  string `json:"type"`
	Nodes []node `json:"nodes"`
}

// Node structure
type node struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Output structures for clean JSON
type metricValue struct {
	Year  string      `json:"year"`
	Value interface{} `json:"value"`
}

type metric struct {
	ID       string        `json:"id"`
	Title    string        `json:"title"`
	Indented bool          `json:"indented,omitempty"`
	Format   string        `json:"format,omitempty"`
	Values   []metricValue `json:"values"`
}

type category struct {
	Name    string   `json:"name"`
	Metrics []metric `json:"metrics"`
}

type output struct {
	Title       string                 `json:"title"`
	Symbol      string                 `json:"symbol"`
	FiscalYears []string               `json:"fiscalYears"`
	Categories  []category             `json:"categories"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// resolveValue recursively resolves indexed values in the data structure
// with protection against infinite recursion
func resolveValue(val interface{}, dataArray []interface{}) interface{} {
	return resolveValueWithDepth(val, dataArray, 0, make(map[int]bool))
}

func resolveValueWithDepth(val interface{}, dataArray []interface{}, depth int, visited map[int]bool) interface{} {
	// Prevent infinite recursion - max depth of 100
	if depth > 100 {
		return val
	}

	switch v := val.(type) {
	case float64:
		// JSON numbers are float64, check if it's an integer index
		idx := int(v)
		if float64(idx) == v && idx > 0 && idx < len(dataArray) {
			// Check if we've already visited this index (circular reference)
			if visited[idx] {
				return val // Return the index itself to break the cycle
			}

			// Mark as visited
			newVisited := make(map[int]bool)
			for k, v := range visited {
				newVisited[k] = v
			}
			newVisited[idx] = true

			return resolveValueWithDepth(dataArray[idx], dataArray, depth+1, newVisited)
		}
		return val
	case map[string]interface{}:
		resolved := make(map[string]interface{})
		for key, value := range v {
			resolved[key] = resolveValueWithDepth(value, dataArray, depth+1, visited)
		}
		return resolved
	case []interface{}:
		resolved := make([]interface{}, len(v))
		for i, item := range v {
			resolved[i] = resolveValueWithDepth(item, dataArray, depth+1, visited)
		}
		return resolved
	default:
		return val
	}
}

func fetchData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return data, nil
}

func processData(data []byte) (*output, error) {
	var root root
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Navigate to the financial data node (third node)
	if len(root.Nodes) < 3 {
		return nil, fmt.Errorf("unexpected JSON structure: not enough nodes")
	}

	// Extract stock info from second node
	symbol := "UNKNOWN"
	companyName := "Unknown Company"

	if len(root.Nodes) >= 2 {
		var secondNodeData []interface{}
		if err := json.Unmarshal(root.Nodes[1].Data, &secondNodeData); err == nil && len(secondNodeData) > 0 {
			if stockData, ok := secondNodeData[0].(map[string]interface{}); ok {
				if infoIdx, ok := stockData["info"].(float64); ok {
					infoIndex := int(infoIdx)
					if infoIndex < len(secondNodeData) {
						if infoData, ok := secondNodeData[infoIndex].(map[string]interface{}); ok {
							// Extract symbol
							if symVal, ok := infoData["symbol"]; ok {
								if symIdx, isFloat := symVal.(float64); isFloat {
									idx := int(symIdx)
									if idx < len(secondNodeData) {
										symbol = fmt.Sprintf("%v", secondNodeData[idx])
									}
								} else {
									symbol = fmt.Sprintf("%v", symVal)
								}
							}

							// Extract name
							if nameVal, ok := infoData["name"]; ok {
								if nameIdx, isFloat := nameVal.(float64); isFloat {
									idx := int(nameIdx)
									if idx < len(secondNodeData) {
										companyName = fmt.Sprintf("%v", secondNodeData[idx])
									}
								} else {
									companyName = fmt.Sprintf("%v", nameVal)
								}
							}
						}
					}
				}
			}
		}
	}

	// Extract the data array from the third node
	var dataArray []interface{}
	if err := json.Unmarshal(root.Nodes[2].Data, &dataArray); err != nil {
		return nil, fmt.Errorf("failed to parse data array: %v", err)
	}

	if len(dataArray) == 0 {
		return nil, fmt.Errorf("data array is empty")
	}

	// Resolve the entire structure
	resolved := resolveValue(dataArray[0], dataArray).(map[string]interface{})

	// Extract title from resolved data
	title := fmt.Sprintf("%s Business Metrics & Revenue Breakdown", companyName)
	if titleVal, ok := resolved["title"]; ok {
		if titleStr, ok := titleVal.(string); ok {
			title = titleStr
		}
	}

	// Get financialData
	finDataRaw, ok := resolved["financialData"]
	if !ok {
		return nil, fmt.Errorf("financialData not found")
	}
	finData := finDataRaw.(map[string]interface{})

	// Extract key components
	fiscalYears := toStringSlice(finData["fiscalYear"])
	categoryNames := toStringSlice(finData["categoryNames"])
	metricOrder := finData["metricOrderByCategory"].(map[string]interface{})
	categories := finData["categories"].(map[string]interface{})
	mapList := toMapSlice(resolved["map"])

	// Get metadata
	detailsMeta := make(map[string]interface{})
	if details, ok := resolved["details"].(map[string]interface{}); ok {
		if source, ok := details["source"]; ok {
			detailsMeta["source"] = source
		}
		if lastDate, ok := details["lastTrailingDate"]; ok {
			detailsMeta["lastUpdated"] = lastDate
		}
		if fiscalYear, ok := details["fiscalYear"]; ok {
			detailsMeta["fiscalYearPeriod"] = fiscalYear
		}
		if fiscalYearShort, ok := details["fiscalYearShort"]; ok {
			detailsMeta["fiscalYearShort"] = fiscalYearShort
		}
	}

	// Build output structure
	output := &output{
		Title:       title,
		Symbol:      symbol,
		FiscalYears: fiscalYears,
		Categories:  make([]category, 0),
		Metadata:    detailsMeta,
	}

	// Create metric definitions lookup
	metricDefs := make(map[string]map[string]interface{})
	for _, entry := range mapList {
		if id, ok := entry["id"].(string); ok {
			metricDefs[id] = entry
		}
	}

	// Process each category
	for _, categoryName := range categoryNames {
		category := category{
			Name:    categoryName,
			Metrics: make([]metric, 0),
		}

		// Get metrics for this category
		categoryMetrics := toStringSlice(metricOrder[categoryName])
		categoryData := categories[categoryName].(map[string]interface{})

		// Process each metric
		for _, metricID := range categoryMetrics {
			metricDef, exists := metricDefs[metricID]
			if !exists {
				continue
			}

			metric := metric{
				ID:     metricID,
				Title:  metricDef["title"].(string),
				Values: make([]metricValue, 0),
			}

			// Check if indented
			if ind, ok := metricDef["indented"]; ok {
				if indStr, isString := ind.(string); isString && indStr != "" {
					metric.Indented = true
				} else if indBool, isBool := ind.(bool); isBool {
					metric.Indented = indBool
				}
			}

			// Add format if present
			if format, ok := metricDef["format"]; ok {
				if formatStr, ok := format.(string); ok {
					metric.Format = formatStr
				}
			}

			// Add values
			if values, ok := categoryData[metricID]; ok {
				if valSlice, ok := values.([]interface{}); ok {
					for i, val := range valSlice {
						if i < len(fiscalYears) {
							// Skip nil and [PRO] locked values
							if val == nil {
								continue
							}
							if strVal, ok := val.(string); ok && strVal == "[PRO]" {
								continue
							}
							metric.Values = append(metric.Values, metricValue{
								Year:  fiscalYears[i],
								Value: val,
							})
						}
					}
				}
			}

			category.Metrics = append(category.Metrics, metric)
		}

		output.Categories = append(output.Categories, category)
	}

	return output, nil
}

// toStringSlice converts an interface{} to []string
func toStringSlice(data interface{}) []string {
	if arr, ok := data.([]interface{}); ok {
		result := make([]string, len(arr))
		for i, v := range arr {
			result[i] = fmt.Sprintf("%v", v)
		}
		return result
	}
	return nil
}

// toMapSlice converts an interface{} to []map[string]interface{}
func toMapSlice(data interface{}) []map[string]interface{} {
	if arr, ok := data.([]interface{}); ok {
		result := make([]map[string]interface{}, 0, len(arr))
		for _, v := range arr {
			if m, ok := v.(map[string]interface{}); ok {
				result = append(result, m)
			}
		}
		return result
	}
	return nil
}

func scrapeCompanyKpiMetrics(stockSymbol string) (domain.CompanyKpiMetrics, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/financials/metrics/__data.json", strings.ToLower(stockSymbol))

	var data []byte
	var err error

	data, err = fetchData(url)
	if err != nil {
		return domain.CompanyKpiMetrics{}, err
	}

	// Process the data
	output, err := processData(data)
	if err != nil {
		return domain.CompanyKpiMetrics{}, err
	}

	// Convert to domain.CompanyKpiMetrics
	kpiMetrics := domain.CompanyKpiMetrics{
		StockSymbol: output.Symbol,
	}

	kpiCategories := make([]domain.KpiCategory, 0)

	for _, category := range output.Categories {
		kpiCategory := domain.KpiCategory{
			Name: category.Name,
		}

		kpiMetrics := make([]domain.KpiMetric, 0)

		for _, metric := range category.Metrics {
			metricValues := make([]domain.KpiMetricValue, 0)

			for _, value := range metric.Values {
				if value.Value == nil {
					continue
				}

				rv := reflect.ValueOf(value.Value)
				if (rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface || rv.Kind() == reflect.Slice || rv.Kind() == reflect.Map || rv.Kind() == reflect.Chan || rv.Kind() == reflect.Func) && rv.IsNil() {
					continue
				}
				metricValues = append(metricValues, domain.KpiMetricValue{
					Year:  value.Year,
					Value: value.Value,
				})
			}

			kpiMetric := domain.KpiMetric{
				Title:  metric.Title,
				Values: metricValues,
			}
			kpiMetrics = append(kpiMetrics, kpiMetric)
		}

		kpiCategory.Metrics = kpiMetrics
		kpiCategories = append(kpiCategories, kpiCategory)
	}

	kpiMetrics.KpiCategories = kpiCategories

	return kpiMetrics, nil
}
