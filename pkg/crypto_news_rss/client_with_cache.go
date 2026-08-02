package cryptonewsrss

import (
	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/services"
	"time"
)

// latestNewsCacheKey holds the merged, unfiltered feed. Caching before the symbol filter
// means asking about ten coins costs one pair of HTTP requests rather than ten.
const latestNewsCacheKey = "crypto_news_latest"

type CryptoNewsRssClientWithCache struct {
	cache           services.CacheService
	cacheTtlSeconds int
}

func NewCryptoNewsRssClientWithCache(cache services.CacheService, cacheTtlSeconds int) (*CryptoNewsRssClientWithCache, error) {
	return &CryptoNewsRssClientWithCache{cache: cache, cacheTtlSeconds: cacheTtlSeconds}, nil
}

func (c *CryptoNewsRssClientWithCache) GetCryptocurrencyNews(symbol string) (domain.CryptocurrencyNews, error) {
	articles, err := c.getLatestNews()
	if err != nil {
		return domain.CryptocurrencyNews{}, err
	}

	return selectForSymbol(symbol, articles), nil
}

func (c *CryptoNewsRssClientWithCache) getLatestNews() ([]domain.NewsArticle, error) {
	// Check if the data is in the cache
	var articles []domain.NewsArticle

	err := c.cache.Get(latestNewsCacheKey, &articles)
	if err == nil {
		return articles, nil
	}

	// If not in cache, get from the feeds
	cryptoNewsRssClient := CryptoNewsRssClient{}
	articles, err = cryptoNewsRssClient.GetLatestNews()
	if err != nil {
		return nil, err
	}

	// Set in cache
	err = c.cache.Set(latestNewsCacheKey, articles, time.Duration(c.cacheTtlSeconds)*time.Second)
	if err != nil {
		return nil, err
	}

	return articles, nil
}
