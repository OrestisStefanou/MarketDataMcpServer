package coingecko

import (
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/services"
	"time"
)

type CoinGeckoClientWithCache struct {
	apiKey          string
	cache           services.CacheService
	cacheTtlSeconds int
}

func NewCoinGeckoClientWithCache(apiKey string, cache services.CacheService, cacheTtlSeconds int) (*CoinGeckoClientWithCache, error) {
	return &CoinGeckoClientWithCache{apiKey: apiKey, cache: cache, cacheTtlSeconds: cacheTtlSeconds}, nil
}

func (c *CoinGeckoClientWithCache) GetCryptocurrenciesList() ([]domain.Cryptocurrency, error) {
	// Check if the data is in the cache
	var cryptocurrenciesList []domain.Cryptocurrency

	key := "cryptocurrencies_list"
	err := c.cache.Get(key, &cryptocurrenciesList)
	if err == nil {
		return cryptocurrenciesList, nil
	}

	// If not in cache, get from API
	coinGeckoClient := CoinGeckoClient{apiKey: c.apiKey}
	cryptocurrenciesList, err = coinGeckoClient.GetCryptocurrenciesList()
	if err != nil {
		return nil, err
	}

	// Set in cache
	err = c.cache.Set(key, cryptocurrenciesList, time.Duration(c.cacheTtlSeconds)*time.Second)
	if err != nil {
		return nil, err
	}

	return cryptocurrenciesList, nil
}

func (c *CoinGeckoClientWithCache) GetCryptocurrencyDataById(id string) (domain.CryptocurrencyData, error) {
	// Check if the data is in the cache
	var cryptocurrencyData domain.CryptocurrencyData

	key := fmt.Sprintf("cryptocurrency_data_%s", id)
	err := c.cache.Get(key, &cryptocurrencyData)
	if err == nil {
		return cryptocurrencyData, nil
	}

	// If not in cache, get from API
	coinGeckoClient := CoinGeckoClient{apiKey: c.apiKey}
	cryptocurrencyData, err = coinGeckoClient.GetCryptocurrencyDataById(id)
	if err != nil {
		return domain.CryptocurrencyData{}, err
	}

	// Set in cache
	err = c.cache.Set(key, cryptocurrencyData, time.Duration(c.cacheTtlSeconds)*time.Second)
	if err != nil {
		return domain.CryptocurrencyData{}, err
	}

	return cryptocurrencyData, nil
}
