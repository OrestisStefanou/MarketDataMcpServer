package polymarket

import (
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/services"
	"strings"
	"time"
)

type PolymarketClientWithCache struct {
	cache           services.CacheService
	cacheTtlSeconds int
}

func NewPolymarketClientWithCache(cache services.CacheService, cacheTtlSeconds int) (*PolymarketClientWithCache, error) {
	return &PolymarketClientWithCache{cache: cache, cacheTtlSeconds: cacheTtlSeconds}, nil
}

func (c *PolymarketClientWithCache) SearchEvents(query string, limit int) ([]domain.PredictionMarketEvent, error) {
	var events []domain.PredictionMarketEvent

	key := fmt.Sprintf("polymarket_events_%s_%d", strings.ToLower(strings.TrimSpace(query)), limit)
	if err := c.cache.Get(key, &events); err == nil {
		return events, nil
	}

	client := PolymarketClient{}
	events, err := client.SearchEvents(query, limit)
	if err != nil {
		return nil, err
	}

	if err := c.cache.Set(key, events, time.Duration(c.cacheTtlSeconds)*time.Second); err != nil {
		return nil, err
	}

	return events, nil
}
