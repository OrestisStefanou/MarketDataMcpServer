package secedgar

import (
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/services"
	"time"
)

const (
	// The ticker to CIK mapping is a large file that changes only when companies list or
	// delist, and the submissions index only when a company files something new.
	tickerMapCacheTtl   = 7 * 24 * time.Hour
	submissionsCacheTtl = 6 * time.Hour
)

type SecEdgarClientWithCache struct {
	userAgent       string
	cache           services.CacheService
	cacheTtlSeconds int
}

func NewSecEdgarClientWithCache(userAgent string, cache services.CacheService, cacheTtlSeconds int) (*SecEdgarClientWithCache, error) {
	return &SecEdgarClientWithCache{userAgent: userAgent, cache: cache, cacheTtlSeconds: cacheTtlSeconds}, nil
}

func (c *SecEdgarClientWithCache) GetInsiderTransactions(symbol string, year int) (domain.InsiderTransactions, error) {
	// Check if the data is in the cache
	var insiderTransactions domain.InsiderTransactions

	key := fmt.Sprintf("insider_transactions_%s_%d", symbol, year)
	err := c.cache.Get(key, &insiderTransactions)
	if err == nil {
		return insiderTransactions, nil
	}

	// If not in cache, get from API. The ticker map and the submissions index are cached
	// separately, so a second year for the same company only pays for the filings.
	secEdgarClient, err := NewSecEdgarClient(c.userAgent)
	if err != nil {
		return domain.InsiderTransactions{}, err
	}

	cik, err := c.getCikBySymbol(secEdgarClient, symbol)
	if err != nil {
		return domain.InsiderTransactions{}, err
	}

	submissions, err := c.getSubmissions(secEdgarClient, cik)
	if err != nil {
		return domain.InsiderTransactions{}, err
	}

	filings, total, err := SelectForm4Filings(submissions, year)
	if err != nil {
		return domain.InsiderTransactions{}, err
	}

	insiderTransactions = domain.InsiderTransactions{
		Symbol:         symbol,
		Year:           year,
		Transactions:   FilterByYear(secEdgarClient.FetchForm4Transactions(cik, filings), year),
		FilingsFound:   total,
		FilingsFetched: len(filings),
		Truncated:      len(filings) < total,
	}

	// Set in cache
	err = c.cache.Set(key, insiderTransactions, time.Duration(c.cacheTtlSeconds)*time.Second)
	if err != nil {
		return domain.InsiderTransactions{}, err
	}

	return insiderTransactions, nil
}

func (c *SecEdgarClientWithCache) getCikBySymbol(client *SecEdgarClient, symbol string) (int, error) {
	var tickerToCik map[string]int

	key := "sec_edgar_ticker_cik_map"
	err := c.cache.Get(key, &tickerToCik)
	if err != nil {
		tickerToCik, err = client.GetTickerToCikMap()
		if err != nil {
			return 0, err
		}
		if err := c.cache.Set(key, tickerToCik, tickerMapCacheTtl); err != nil {
			return 0, err
		}
	}

	return lookupCik(tickerToCik, symbol)
}

func (c *SecEdgarClientWithCache) getSubmissions(client *SecEdgarClient, cik int) (Submissions, error) {
	var submissions Submissions

	key := fmt.Sprintf("sec_edgar_submissions_%d", cik)
	err := c.cache.Get(key, &submissions)
	if err == nil {
		return submissions, nil
	}

	submissions, err = client.GetSubmissions(cik)
	if err != nil {
		return Submissions{}, err
	}

	if err := c.cache.Set(key, submissions, submissionsCacheTtl); err != nil {
		return Submissions{}, err
	}

	return submissions, nil
}
