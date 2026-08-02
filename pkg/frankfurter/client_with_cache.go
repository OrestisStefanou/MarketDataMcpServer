package frankfurter

import (
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/services"
	"time"
)

type FrankfurterClientWithCache struct {
	cache           services.CacheService
	cacheTtlSeconds int
}

func NewFrankfurterClientWithCache(cache services.CacheService, cacheTtlSeconds int) (*FrankfurterClientWithCache, error) {
	return &FrankfurterClientWithCache{cache: cache, cacheTtlSeconds: cacheTtlSeconds}, nil
}

func (c *FrankfurterClientWithCache) GetCurrencyExchangeRate(fromCurrency domain.Currency, toCurrency domain.Currency) (domain.CurrencyExchangeRate, error) {
	// Check if the data is in the cache
	var exchangeRate domain.CurrencyExchangeRate

	key := fmt.Sprintf("currency_exchange_rate_%s_%s", fromCurrency, toCurrency)
	err := c.cache.Get(key, &exchangeRate)
	if err == nil {
		return exchangeRate, nil
	}

	// If not in cache, get from API
	frankfurterClient := FrankfurterClient{}
	exchangeRate, err = frankfurterClient.GetCurrencyExchangeRate(fromCurrency, toCurrency)
	if err != nil {
		return domain.CurrencyExchangeRate{}, err
	}

	// Set in cache
	err = c.cache.Set(key, exchangeRate, time.Duration(c.cacheTtlSeconds)*time.Second)
	if err != nil {
		return domain.CurrencyExchangeRate{}, err
	}

	return exchangeRate, nil
}
