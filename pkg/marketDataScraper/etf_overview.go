package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"net/http"
)

func scrapeEtfOverview(symbol string) (domain.EtfOverview, error) {
	url := fmt.Sprintf("https://api.stockanalysis.com/api/symbol/e/%s/overview", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return domain.EtfOverview{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.EtfOverview{}, err
	}

	// Define an anonymous struct to match the JSON structure
	var apiResponse struct {
		Status int `json:"status"`
		Data   struct {
			Aum           string     `json:"aum"`
			Nav           string     `json:"nav"`
			ExpenseRatio  string     `json:"expenseRatio"`
			Description   string     `json:"description"`
			PeRatio       string     `json:"peRatio"`
			Dps           string     `json:"dps"`
			DividendYield string     `json:"dividendYield"`
			PayoutRatio   string     `json:"payoutRatio"`
			Ch1y          string     `json:"ch1y"`
			Beta          string     `json:"beta"`
			Holdings      int32      `json:"holdings"`
			EtfWebsite    string     `json:"etf_website"`
			InfoTable     [][]string `json:"infoTable"`
			HoldingsTable struct {
				Count    int `json:"count"`
				Holdings []struct {
					S  string `json:"s"`
					N  string `json:"n"`
					As string `json:"as"`
				} `json:"holdings"`
			} `json:"holdingsTable"`
			Performance struct {
				Tr1m    float64 `json:"tr1m"`
				TrYTD   float64 `json:"trYTD"`
				Tr1y    float64 `json:"tr1y"`
				Cagr5y  float64 `json:"cagr5y"`
				Cagr10y float64 `json:"cagr10y"`
				CagrMAX float64 `json:"cagrMAX"`
			} `json:"performance"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return domain.EtfOverview{}, err
	}

	etfOverview := domain.EtfOverview{
		Symbol:           symbol,
		Description:      apiResponse.Data.Description,
		Aum:              apiResponse.Data.Aum,
		Nav:              apiResponse.Data.Nav,
		ExpenseRatio:     apiResponse.Data.ExpenseRatio,
		PeRatio:          apiResponse.Data.PeRatio,
		Dps:              apiResponse.Data.Dps,
		DividendYield:    apiResponse.Data.DividendYield,
		PayoutRatio:      apiResponse.Data.PayoutRatio,
		OneMonthReturn:   apiResponse.Data.Performance.Tr1m,
		OneYearReturn:    apiResponse.Data.Performance.Tr1y,
		YearToDateReturn: apiResponse.Data.Performance.TrYTD,
		FiveYearReturn:   apiResponse.Data.Performance.Cagr5y,
		TenYearReturn:    apiResponse.Data.Performance.Cagr10y,
		InceptionReturn:  apiResponse.Data.Performance.CagrMAX,
		Beta:             apiResponse.Data.Beta,
		NumberOfHoldings: apiResponse.Data.Holdings,
		Website:          apiResponse.Data.EtfWebsite,
		TopHoldings:      make([]domain.EtfHolding, 0, len(apiResponse.Data.HoldingsTable.Holdings)),
	}

	for _, holding := range apiResponse.Data.HoldingsTable.Holdings {
		etfHolding := domain.EtfHolding{
			Symbol: holding.S,
			Name:   holding.N,
			Weight: holding.As,
		}
		etfOverview.TopHoldings = append(etfOverview.TopHoldings, etfHolding)
	}

	var assetClass string
	for _, info := range apiResponse.Data.InfoTable {
		if info[0] == "Asset Class" {
			assetClass = info[1]
			break
		}
	}
	etfOverview.AssetClass = assetClass

	var category string
	for _, info := range apiResponse.Data.InfoTable {
		if info[0] == "Category" {
			category = info[1]
			break
		}
	}
	etfOverview.Category = category

	return etfOverview, nil
}
