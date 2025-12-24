package marketDataScraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"market_data_mcp_server/pkg/domain"
	"net/http"
)

func scrapeFinancialStatementData(url string) ([]map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []map[string]interface{}{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return []map[string]interface{}{}, err
	}

	// Extract "nodes" from rawData
	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 3 {
		return []map[string]interface{}{}, errors.New("unexpected structure in 'nodes'")
	}

	// Access the second element in "nodes" which contains the data we are interested in
	nodeData, ok := nodes[2].(map[string]interface{})
	if !ok {
		return []map[string]interface{}{}, errors.New("unexpected structure in 'nodes[2]'")
	}

	data, ok := nodeData["data"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, errors.New("unexpected structure in 'data'")
	}

	dataMap, ok := data[0].(map[string]interface{})
	if !ok {
		return []map[string]interface{}{}, errors.New("unexpected structure in 'data[0]'")
	}

	financialDataIndex, ok := dataMap["financialData"].(float64)
	if !ok {
		return []map[string]interface{}{}, errors.New("unexpected structure for 'financialData'")
	}

	// Retrieve data map at financial data index
	statementDataMap, ok := data[int(financialDataIndex)].(map[string]interface{})
	if !ok {
		return []map[string]interface{}{}, errors.New("unexpected structure in balance sheet data map")
	}

	statement_data := make(map[string][]interface{})
	for field, fieldIndex := range statementDataMap {
		fieldIndexFloat, ok := fieldIndex.(float64)
		if !ok {
			return []map[string]interface{}{}, errors.New("unexpected index type in fieldIndex")
		}
		if fieldIndexFloat < 0 {
			continue
		}
		fieldValues := []interface{}{}
		for _, index := range data[int(fieldIndexFloat)].([]interface{}) {
			indexFloat, ok := index.(float64)
			if !ok {
				return []map[string]interface{}{}, errors.New("unexpected type in field values index")
			}
			fieldValues = append(fieldValues, data[int(indexFloat)])
		}
		statement_data[field] = fieldValues
	}

	// Converting the map of slices into a slice of maps to resemble final structure
	statement_data_slice := make([]map[string]interface{}, 0, len(statement_data["datekey"]))
	for i := 0; i < len(statement_data["datekey"]); i++ {
		record := make(map[string]interface{})
		for key, values := range statement_data {
			// There are some cases where the data is missing for a particular key
			if i >= len(values) {
				continue
			}
			record[key] = values[i]
		}
		statement_data_slice = append(statement_data_slice, record)
	}

	return statement_data_slice, nil
}

func scrapeBalanceSheets(symbol string) ([]domain.BalanceSheet, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/financials/balance-sheet/__data.json?p=quarterly", symbol)
	balanceSheetData, err := scrapeFinancialStatementData(url)
	if err != nil {
		return []domain.BalanceSheet{}, err
	}

	balanceSheets := make([]domain.BalanceSheet, 0, len(balanceSheetData))
	for i := 0; i < len(balanceSheetData); i++ {
		record := balanceSheetData[i]
		// Marshal the map to JSON
		jsonData, err := json.Marshal(record)
		if err != nil {
			return []domain.BalanceSheet{}, err
		}
		// Unmarshal the JSON data into an instance of balanceSheet
		var balanceSheetRecord domain.BalanceSheet
		err = json.Unmarshal(jsonData, &balanceSheetRecord)
		if err != nil {
			return []domain.BalanceSheet{}, err
		}
		balanceSheets = append(balanceSheets, balanceSheetRecord)
	}
	return balanceSheets, nil
}

func scrapeCashFlows(symbol string) ([]domain.CashFlow, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/financials/cash-flow-statement/__data.json?p=quarterly", symbol)
	cashFlowData, err := scrapeFinancialStatementData(url)
	if err != nil {
		return []domain.CashFlow{}, err
	}

	cashFlows := make([]domain.CashFlow, 0, len(cashFlowData))
	for i := 0; i < len(cashFlowData); i++ {
		record := cashFlowData[i]
		// Marshal the map to JSON
		jsonData, err := json.Marshal(record)
		if err != nil {
			return []domain.CashFlow{}, err
		}
		// Unmarshal the JSON data into an instance of balanceSheet
		var cashFlowRecord domain.CashFlow
		err = json.Unmarshal(jsonData, &cashFlowRecord)
		if err != nil {
			return []domain.CashFlow{}, err
		}
		cashFlows = append(cashFlows, cashFlowRecord)
	}
	return cashFlows, nil
}

func scrapeIncomeStatements(symbol string) ([]domain.IncomeStatement, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/financials/__data.json?p=quarterly", symbol)
	incomeStatementData, err := scrapeFinancialStatementData(url)
	if err != nil {
		return []domain.IncomeStatement{}, err
	}

	incomeStatements := make([]domain.IncomeStatement, 0, len(incomeStatementData))
	for i := 0; i < len(incomeStatementData); i++ {
		record := incomeStatementData[i]
		// Marshal the map to JSON
		jsonData, err := json.Marshal(record)
		if err != nil {
			return []domain.IncomeStatement{}, err
		}
		// Unmarshal the JSON data into an instance of balanceSheet
		var incomeStatementRecord domain.IncomeStatement
		err = json.Unmarshal(jsonData, &incomeStatementRecord)
		if err != nil {
			return []domain.IncomeStatement{}, err
		}
		incomeStatements = append(incomeStatements, incomeStatementRecord)
	}
	return incomeStatements, nil
}

func scrapeFinancialRatios(symbol string) ([]domain.FinancialRatios, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/financials/ratios/__data.json?p=quarterly", symbol)
	financialRatiosData, err := scrapeFinancialStatementData(url)
	if err != nil {
		return []domain.FinancialRatios{}, err
	}

	financialRatios := make([]domain.FinancialRatios, 0, len(financialRatiosData))
	for i := 0; i < len(financialRatiosData); i++ {
		record := financialRatiosData[i]
		// Marshal the map to JSON
		jsonData, err := json.Marshal(record)
		if err != nil {
			return []domain.FinancialRatios{}, err
		}
		// Unmarshal the JSON data into an instance of balanceSheet
		var financialRatiosRecord domain.FinancialRatios
		err = json.Unmarshal(jsonData, &financialRatiosRecord)
		if err != nil {
			return []domain.FinancialRatios{}, err
		}
		financialRatios = append(financialRatios, financialRatiosRecord)
	}

	return financialRatios, nil
}
