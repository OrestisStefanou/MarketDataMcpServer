package alphavantage

import (
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/services"
	"time"
)

type AlphaVantageClientWithCache struct {
	apiKey          string
	cache           services.CacheService
	cacheTtlSeconds int
}

func NewAlphaVantageClientWithCache(apiKey string, cache services.CacheService, cacheTtlSeconds int) (*AlphaVantageClientWithCache, error) {
	return &AlphaVantageClientWithCache{apiKey: apiKey, cache: cache, cacheTtlSeconds: cacheTtlSeconds}, nil
}

func (c *AlphaVantageClientWithCache) GetRealGdpTimeSeries(interval domain.EconomicIndicatorInterval) (domain.EconomicIndicatorTimeSeries, error) {
	// Check if the data is in the cache
	var economicIndicatorTimeSeries domain.EconomicIndicatorTimeSeries

	key := fmt.Sprintf("real_gdp_%s", interval)
	err := c.cache.Get(key, &economicIndicatorTimeSeries)
	if err == nil {
		return economicIndicatorTimeSeries, nil
	}

	// If not in cache, get from API
	alphaVantageClient := AlphaVantageClient{apiKey: c.apiKey}
	economicIndicatorTimeSeries, err = alphaVantageClient.GetRealGdpTimeSeries(interval)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, err
	}

	// Set in cache
	c.cache.Set(key, economicIndicatorTimeSeries, time.Duration(c.cacheTtlSeconds)*time.Second)

	return economicIndicatorTimeSeries, nil
}

func (c *AlphaVantageClientWithCache) GetTreasuryYieldTimeSeries(maturity domain.TreasuryYieldMaturity) (domain.EconomicIndicatorTimeSeries, error) {
	// Check if the data is in the cache
	var economicIndicatorTimeSeries domain.EconomicIndicatorTimeSeries

	key := fmt.Sprintf("treasury_yield_%s", maturity)
	err := c.cache.Get(key, &economicIndicatorTimeSeries)
	if err == nil {
		return economicIndicatorTimeSeries, nil
	}

	// If not in cache, get from API
	alphaVantageClient := AlphaVantageClient{apiKey: c.apiKey}
	economicIndicatorTimeSeries, err = alphaVantageClient.GetTreasuryYieldTimeSeries(maturity)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, err
	}

	// Set in cache
	c.cache.Set(key, economicIndicatorTimeSeries, time.Duration(c.cacheTtlSeconds)*time.Second)

	return economicIndicatorTimeSeries, nil
}

func (c *AlphaVantageClientWithCache) GetInterestRatesTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	// Check if the data is in the cache
	var economicIndicatorTimeSeries domain.EconomicIndicatorTimeSeries

	key := "interest_rate"
	err := c.cache.Get(key, &economicIndicatorTimeSeries)
	if err == nil {
		return economicIndicatorTimeSeries, nil
	}

	// If not in cache, get from API
	alphaVantageClient := AlphaVantageClient{apiKey: c.apiKey}
	economicIndicatorTimeSeries, err = alphaVantageClient.GetInterestRatesTimeSeries()
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, err
	}

	// Set in cache
	c.cache.Set(key, economicIndicatorTimeSeries, time.Duration(c.cacheTtlSeconds)*time.Second)

	return economicIndicatorTimeSeries, nil
}

func (c *AlphaVantageClientWithCache) GetInflationTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	// Check if the data is in the cache
	var economicIndicatorTimeSeries domain.EconomicIndicatorTimeSeries

	key := "inflation"
	err := c.cache.Get(key, &economicIndicatorTimeSeries)
	if err == nil {
		return economicIndicatorTimeSeries, nil
	}

	// If not in cache, get from API
	alphaVantageClient := AlphaVantageClient{apiKey: c.apiKey}
	economicIndicatorTimeSeries, err = alphaVantageClient.GetInflationTimeSeries()
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, err
	}

	// Set in cache
	c.cache.Set(key, economicIndicatorTimeSeries, time.Duration(c.cacheTtlSeconds)*time.Second)

	return economicIndicatorTimeSeries, nil
}

func (c *AlphaVantageClientWithCache) GetUnemploymentRateTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	// Check if the data is in the cache
	var economicIndicatorTimeSeries domain.EconomicIndicatorTimeSeries

	key := "unemployment_rate"
	err := c.cache.Get(key, &economicIndicatorTimeSeries)
	if err == nil {
		return economicIndicatorTimeSeries, nil
	}

	// If not in cache, get from API
	alphaVantageClient := AlphaVantageClient{apiKey: c.apiKey}
	economicIndicatorTimeSeries, err = alphaVantageClient.GetUnemploymentRateTimeSeries()
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, err
	}

	// Set in cache
	c.cache.Set(key, economicIndicatorTimeSeries, time.Duration(c.cacheTtlSeconds)*time.Second)

	return economicIndicatorTimeSeries, nil
}

func (c *AlphaVantageClientWithCache) GetCommodityTimeSeries(commodity domain.Commodity) (domain.CommodityTimeSeries, error) {
	// Check if the data is in the cache
	var commodityTimeSeries domain.CommodityTimeSeries

	key := fmt.Sprintf("commodity_%s", commodity)
	err := c.cache.Get(key, &commodityTimeSeries)
	if err == nil {
		return commodityTimeSeries, nil
	}

	// If not in cache, get from API
	alphaVantageClient := AlphaVantageClient{apiKey: c.apiKey}
	commodityTimeSeries, err = alphaVantageClient.GetCommodityTimeSeries(commodity)
	if err != nil {
		return domain.CommodityTimeSeries{}, err
	}

	// Set in cache
	c.cache.Set(key, commodityTimeSeries, time.Duration(c.cacheTtlSeconds)*time.Second)

	return commodityTimeSeries, nil
}

func (c *AlphaVantageClientWithCache) GetCryptocurrencyNews(symbol string) ([]domain.NewsArticle, error) {
	// Check if the data is in the cache
	var newsArticles []domain.NewsArticle

	key := fmt.Sprintf("cryptocurrency_news_%s", symbol)
	err := c.cache.Get(key, &newsArticles)
	if err == nil {
		return newsArticles, nil
	}

	// If not in cache, get from API
	alphaVantageClient := AlphaVantageClient{apiKey: c.apiKey}
	newsArticles, err = alphaVantageClient.GetCryptocurrencyNews(symbol)
	if err != nil {
		return nil, err
	}

	// Set in cache
	c.cache.Set(key, newsArticles, time.Duration(c.cacheTtlSeconds)*time.Second)

	return newsArticles, nil
}
