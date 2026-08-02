package fred

import (
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/services"
	"time"
)

type FredClientWithCache struct {
	cache           services.CacheService
	cacheTtlSeconds int
}

func NewFredClientWithCache(cache services.CacheService, cacheTtlSeconds int) (*FredClientWithCache, error) {
	return &FredClientWithCache{cache: cache, cacheTtlSeconds: cacheTtlSeconds}, nil
}

func (c *FredClientWithCache) GetRealGdpTimeSeries(interval domain.EconomicIndicatorInterval) (domain.EconomicIndicatorTimeSeries, error) {
	key := fmt.Sprintf("real_gdp_%s", interval)
	return c.cachedEconomicIndicator(key, func(client FredClient) (domain.EconomicIndicatorTimeSeries, error) {
		return client.GetRealGdpTimeSeries(interval)
	})
}

func (c *FredClientWithCache) GetTreasuryYieldTimeSeries(maturity domain.TreasuryYieldMaturity) (domain.EconomicIndicatorTimeSeries, error) {
	key := fmt.Sprintf("treasury_yield_%s", maturity)
	return c.cachedEconomicIndicator(key, func(client FredClient) (domain.EconomicIndicatorTimeSeries, error) {
		return client.GetTreasuryYieldTimeSeries(maturity)
	})
}

func (c *FredClientWithCache) GetInterestRatesTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	return c.cachedEconomicIndicator("interest_rate", func(client FredClient) (domain.EconomicIndicatorTimeSeries, error) {
		return client.GetInterestRatesTimeSeries()
	})
}

func (c *FredClientWithCache) GetInflationTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	return c.cachedEconomicIndicator("inflation", func(client FredClient) (domain.EconomicIndicatorTimeSeries, error) {
		return client.GetInflationTimeSeries()
	})
}

func (c *FredClientWithCache) GetUnemploymentRateTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	return c.cachedEconomicIndicator("unemployment_rate", func(client FredClient) (domain.EconomicIndicatorTimeSeries, error) {
		return client.GetUnemploymentRateTimeSeries()
	})
}

func (c *FredClientWithCache) GetCommodityTimeSeries(commodity domain.Commodity) (domain.CommodityTimeSeries, error) {
	// Check if the data is in the cache
	var timeSeries domain.CommodityTimeSeries

	key := fmt.Sprintf("commodity_%s", commodity)
	err := c.cache.Get(key, &timeSeries)
	if err == nil {
		return timeSeries, nil
	}

	// If not in cache, get from API
	fredClient := FredClient{}
	timeSeries, err = fredClient.GetCommodityTimeSeries(commodity)
	if err != nil {
		return domain.CommodityTimeSeries{}, err
	}

	// Set in cache
	err = c.cache.Set(key, timeSeries, time.Duration(c.cacheTtlSeconds)*time.Second)
	if err != nil {
		return domain.CommodityTimeSeries{}, err
	}

	return timeSeries, nil
}

// cachedEconomicIndicator wraps the five economic indicator calls, which differ only in
// their cache key and which client method they delegate to.
func (c *FredClientWithCache) cachedEconomicIndicator(
	key string,
	fetch func(client FredClient) (domain.EconomicIndicatorTimeSeries, error),
) (domain.EconomicIndicatorTimeSeries, error) {
	// Check if the data is in the cache
	var timeSeries domain.EconomicIndicatorTimeSeries

	err := c.cache.Get(key, &timeSeries)
	if err == nil {
		return timeSeries, nil
	}

	// If not in cache, get from API
	timeSeries, err = fetch(FredClient{})
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, err
	}

	// Set in cache
	err = c.cache.Set(key, timeSeries, time.Duration(c.cacheTtlSeconds)*time.Second)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, err
	}

	return timeSeries, nil
}
