package alphavantage

import (
	"encoding/json"
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

type AlphaVantageClient struct {
	apiKey string
}

const alphaVantageBaseURL = "https://www.alphavantage.co/query"

func NewAlphaVantageClient(apiKey string) (*AlphaVantageClient, error) {
	return &AlphaVantageClient{apiKey: apiKey}, nil
}

func (c *AlphaVantageClient) GetRealGdpTimeSeries(interval domain.EconomicIndicatorInterval) (domain.EconomicIndicatorTimeSeries, error) {
	// Map domain interval to API interval
	apiInterval := "annual"
	if interval == domain.QuarterlyEconomicIndicatorInterval {
		apiInterval = "quarterly"
	}

	// Build URL with query parameters
	requestUrl, err := url.Parse(alphaVantageBaseURL)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to parse base URL: %v", err),
		}
	}

	q := requestUrl.Query()
	q.Set("function", string(RealGDP))
	q.Set("interval", apiInterval)
	q.Set("apikey", c.apiKey)
	requestUrl.RawQuery = q.Encode()

	// Create HTTP request
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to send HTTP request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Parse JSON response
	var apiResponse EconomicIndicatorTimeSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.JSONMarshalError{
			Message: "failed to decode JSON response",
			Err:     err,
		}
	}

	// Map API response to domain model
	domainInterval := domain.AnnualEconomicIndicatorInterval
	if apiResponse.Interval == "quarterly" {
		domainInterval = domain.QuarterlyEconomicIndicatorInterval
	}

	// Map unit - API returns "billions of dollars" for Real GDP
	domainUnit := domain.BillionsOfDollarsEconomicIndicatorUnit

	// Map data entries
	data := make([]domain.EconomicIndicatorTimeSeriesEntry, len(apiResponse.Data))
	for i, entry := range apiResponse.Data {
		data[i] = domain.EconomicIndicatorTimeSeriesEntry{
			Date:  entry.Date,
			Value: entry.Value,
		}
	}

	return domain.EconomicIndicatorTimeSeries{
		Name:     domain.RealGDP,
		Interval: domainInterval,
		Unit:     domainUnit,
		Data:     data,
	}, nil
}

// GetTreasuryYieldTimeSeries returns the monthly treasury yield of the given maturity
func (c *AlphaVantageClient) GetTreasuryYieldTimeSeries(
	maturity domain.TreasuryYieldMaturity,
) (domain.EconomicIndicatorTimeSeries, error) {
	apiInterval := "monthly"

	// Build URL with query parameters
	requestUrl, err := url.Parse(alphaVantageBaseURL)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to parse base URL: %v", err),
		}
	}

	var maturityParam string
	switch maturity {
	case domain.ThreeMonthTreasuryYieldMaturity:
		maturityParam = "3month"
	case domain.TwoYearTreasuryYieldMaturity:
		maturityParam = "2year"
	case domain.FiveYearTreasuryYieldMaturity:
		maturityParam = "5year"
	case domain.TenYearTreasuryYieldMaturity:
		maturityParam = "10year"
	case domain.ThirtyYearTreasuryYieldMaturity:
		maturityParam = "30year"
	default:
		maturityParam = "10year"
	}

	q := requestUrl.Query()
	q.Set("function", string(TreasuryYield))
	q.Set("interval", apiInterval)
	q.Set("apikey", c.apiKey)
	q.Set("maturity", maturityParam)
	requestUrl.RawQuery = q.Encode()

	// Create HTTP request
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to send HTTP request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Parse JSON response
	var apiResponse EconomicIndicatorTimeSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.JSONMarshalError{
			Message: "failed to decode JSON response",
			Err:     err,
		}
	}

	domainInterval := domain.MonthlyEconomicIndicatorInterval
	domainUnit := domain.PercentEconomicIndicatorUnit

	// Map data entries
	data := make([]domain.EconomicIndicatorTimeSeriesEntry, len(apiResponse.Data))
	for i, entry := range apiResponse.Data {
		data[i] = domain.EconomicIndicatorTimeSeriesEntry{
			Date:  entry.Date,
			Value: entry.Value,
		}
	}

	return domain.EconomicIndicatorTimeSeries{
		Name:     domain.TreasuryYield,
		Interval: domainInterval,
		Unit:     domainUnit,
		Data:     data,
	}, nil
}

// GetInterestRatesTimeSeries returns the monthly interest rate time series
func (c *AlphaVantageClient) GetInterestRatesTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	apiInterval := "monthly"

	// Build URL with query parameters
	requestUrl, err := url.Parse(alphaVantageBaseURL)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to parse base URL: %v", err),
		}
	}

	q := requestUrl.Query()
	q.Set("function", string(FederalFundsRate))
	q.Set("interval", apiInterval)
	q.Set("apikey", c.apiKey)
	requestUrl.RawQuery = q.Encode()

	// Create HTTP request
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to send HTTP request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Parse JSON response
	var apiResponse EconomicIndicatorTimeSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.JSONMarshalError{
			Message: "failed to decode JSON response",
			Err:     err,
		}
	}

	domainInterval := domain.MonthlyEconomicIndicatorInterval
	domainUnit := domain.PercentEconomicIndicatorUnit

	// Map data entries
	data := make([]domain.EconomicIndicatorTimeSeriesEntry, len(apiResponse.Data))
	for i, entry := range apiResponse.Data {
		data[i] = domain.EconomicIndicatorTimeSeriesEntry{
			Date:  entry.Date,
			Value: entry.Value,
		}
	}

	return domain.EconomicIndicatorTimeSeries{
		Name:     domain.InterestRate,
		Interval: domainInterval,
		Unit:     domainUnit,
		Data:     data,
	}, nil
}

// GetInflationTimeSeries returns the annual inflation time series
func (c *AlphaVantageClient) GetInflationTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	// Build URL with query parameters
	requestUrl, err := url.Parse(alphaVantageBaseURL)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to parse base URL: %v", err),
		}
	}

	q := requestUrl.Query()
	q.Set("function", string(Inflation))
	q.Set("apikey", c.apiKey)
	requestUrl.RawQuery = q.Encode()

	// Create HTTP request
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to send HTTP request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Parse JSON response
	var apiResponse EconomicIndicatorTimeSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.JSONMarshalError{
			Message: "failed to decode JSON response",
			Err:     err,
		}
	}

	domainInterval := domain.AnnualEconomicIndicatorInterval
	domainUnit := domain.PercentEconomicIndicatorUnit

	// Map data entries
	data := make([]domain.EconomicIndicatorTimeSeriesEntry, len(apiResponse.Data))
	for i, entry := range apiResponse.Data {
		data[i] = domain.EconomicIndicatorTimeSeriesEntry{
			Date:  entry.Date,
			Value: entry.Value,
		}
	}

	return domain.EconomicIndicatorTimeSeries{
		Name:     domain.Inflation,
		Interval: domainInterval,
		Unit:     domainUnit,
		Data:     data,
	}, nil
}

// GetUnemploymentRateTimeSeries returns the monthly unemployment rate time series
func (c *AlphaVantageClient) GetUnemploymentRateTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	// Build URL with query parameters
	requestUrl, err := url.Parse(alphaVantageBaseURL)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to parse base URL: %v", err),
		}
	}

	q := requestUrl.Query()
	q.Set("function", string(UnemploymentRate))
	q.Set("apikey", c.apiKey)
	requestUrl.RawQuery = q.Encode()

	// Create HTTP request
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to send HTTP request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Parse JSON response
	var apiResponse EconomicIndicatorTimeSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.JSONMarshalError{
			Message: "failed to decode JSON response",
			Err:     err,
		}
	}

	domainInterval := domain.MonthlyEconomicIndicatorInterval
	domainUnit := domain.PercentEconomicIndicatorUnit

	// Map data entries
	data := make([]domain.EconomicIndicatorTimeSeriesEntry, len(apiResponse.Data))
	for i, entry := range apiResponse.Data {
		data[i] = domain.EconomicIndicatorTimeSeriesEntry{
			Date:  entry.Date,
			Value: entry.Value,
		}
	}

	return domain.EconomicIndicatorTimeSeries{
		Name:     domain.UnemploymentRate,
		Interval: domainInterval,
		Unit:     domainUnit,
		Data:     data,
	}, nil
}

// GetCommodityTimeSeries returns the monthly time series for the given commodity
func (c *AlphaVantageClient) GetCommodityTimeSeries(commodity domain.Commodity) (domain.CommodityTimeSeries, error) {
	var function string
	var unit domain.CommodityUnit

	switch commodity {
	case domain.CrudeOil:
		function = string(WTI)
		unit = domain.DollarsPerBarrelCommodityUnit
	case domain.NaturalGas:
		function = string(NATURAL_GAS)
		unit = domain.DollarsPerMillionBTUCommodityUnit
	case domain.Copper:
		function = string(COPPER)
		unit = domain.DollarsPerMetricTonCommodityUnit
	case domain.Aluminum:
		function = string(ALUMINIUM)
		unit = domain.DollarsPerMetricTonCommodityUnit
	case domain.Wheat:
		function = string(WHEAT)
		unit = domain.DollarsPerMetricTonCommodityUnit
	case domain.Corn:
		function = string(CORN)
		unit = domain.DollarsPerMetricTonCommodityUnit
	case domain.Sugar:
		function = string(SUGAR)
		unit = domain.CentsPerPoundCommodityUnit
	case domain.Coffee:
		function = string(COFFEE)
		unit = domain.CentsPerPoundCommodityUnit
	default:
		return domain.CommodityTimeSeries{}, fmt.Errorf("unsupported commodity: %v", commodity)
	}

	apiResponse, err := c.fetchCommodityData(function, "monthly")
	if err != nil {
		return domain.CommodityTimeSeries{}, err
	}

	data := make([]domain.CommodityTimeSeriesEntry, len(apiResponse.Data))
	for i, entry := range apiResponse.Data {
		data[i] = domain.CommodityTimeSeriesEntry{
			Date:  entry.Date,
			Value: entry.Value,
		}
	}

	return domain.CommodityTimeSeries{
		Name:     commodity,
		Interval: domain.MonthlyCommodityInterval,
		Unit:     unit,
		Data:     data,
	}, nil
}

func (c *AlphaVantageClient) fetchCommodityData(function string, interval string) (CommodityTimeSeriesResponse, error) {
	requestUrl, err := url.Parse(alphaVantageBaseURL)
	if err != nil {
		return CommodityTimeSeriesResponse{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to parse base URL: %v", err),
		}
	}

	q := requestUrl.Query()
	q.Set("function", function)
	q.Set("interval", interval)
	q.Set("apikey", c.apiKey)
	requestUrl.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return CommodityTimeSeriesResponse{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return CommodityTimeSeriesResponse{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to send HTTP request: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CommodityTimeSeriesResponse{}, &errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	var apiResponse CommodityTimeSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return CommodityTimeSeriesResponse{}, &errors.JSONMarshalError{
			Message: "failed to decode JSON response",
			Err:     err,
		}
	}

	return apiResponse, nil
}

func (c *AlphaVantageClient) GetCryptocurrencyNews(symbol string) ([]domain.NewsArticle, error) {
	// Build URL with query parameters
	requestUrl, err := url.Parse(alphaVantageBaseURL)
	if err != nil {
		return nil, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to parse base URL: %v", err),
		}
	}

	q := requestUrl.Query()
	q.Set("function", "NEWS_SENTIMENT")
	q.Set("tickers", fmt.Sprintf("CRYPTO:%s", strings.ToUpper(symbol)))
	q.Set("apikey", c.apiKey)
	requestUrl.RawQuery = q.Encode()

	// Create HTTP request
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return nil, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to send HTTP request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, &errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Parse JSON response
	var apiResponse GetNewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, &errors.JSONMarshalError{
			Message: "failed to decode JSON response",
			Err:     err,
		}
	}

	// Map API response to domain model
	var newsArticles []domain.NewsArticle
	for _, article := range apiResponse.Feed {
		newsArticles = append(newsArticles, domain.NewsArticle{
			Title:  article.Title,
			Url:    article.URL,
			Time:   article.TimePublished,
			Image:  article.BannerImage,
			Source: article.Source,
			Text:   article.Summary,
		})
	}

	return newsArticles, nil
}

func (c *AlphaVantageClient) GetEarningsCallTranscript(symbol string, year int, quarter domain.Quarter) ([]domain.EarningsCallTranscript, error) {
	requestUrl, err := url.Parse(alphaVantageBaseURL)
	if err != nil {
		return nil, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to parse base URL: %v", err),
		}
	}

	q := requestUrl.Query()
	q.Set("function", "EARNINGS_CALL_TRANSCRIPT")
	q.Set("symbol", symbol)
	q.Set("quarter", fmt.Sprintf("%d%s", year, quarter))
	q.Set("apikey", c.apiKey)
	requestUrl.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return nil, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to send HTTP request: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	var apiResponse GetEarningsCallTranscriptResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, &errors.JSONMarshalError{
			Message: "failed to decode JSON response",
			Err:     err,
		}
	}

	var earningsCallTranscript []domain.EarningsCallTranscript
	for _, transcript := range apiResponse.Transcript {
		earningsCallTranscript = append(earningsCallTranscript, domain.EarningsCallTranscript{
			Speaker:   transcript.Speaker,
			Title:     transcript.Title,
			Content:   transcript.Content,
			Sentiment: transcript.Sentiment,
		})
	}

	return earningsCallTranscript, nil
}
