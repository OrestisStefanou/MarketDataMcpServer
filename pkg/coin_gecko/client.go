package coingecko

import (
	"encoding/json"
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"net/http"
	"net/url"
)

const coinGeckoBaseURL = "https://api.coingecko.com/api/v3"

type CoinGeckoClient struct {
	apiKey string
}

func NewCoinGeckoClient(apiKey string) (*CoinGeckoClient, error) {
	return &CoinGeckoClient{apiKey: apiKey}, nil
}

func (c *CoinGeckoClient) GetCryptocurrenciesList() ([]domain.Cryptocurrency, error) {
	requestUrl := fmt.Sprintf("%s/coins/list", coinGeckoBaseURL)

	// The api key is optional. CoinGecko serves these endpoints anonymously at a lower
	// rate limit, so the header is only set when a key is configured.
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}
	if c.apiKey != "" {
		req.Header.Set("x-cg-demo-api-key", c.apiKey)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get coins list: %s", resp.Status)
	}

	// Parse the response
	var coinsList []CoinGeckoCoin
	if err := json.NewDecoder(resp.Body).Decode(&coinsList); err != nil {
		return nil, err
	}

	var cryptocurrenciesList []domain.Cryptocurrency
	for _, coin := range coinsList {
		cryptocurrenciesList = append(cryptocurrenciesList, domain.Cryptocurrency{
			Id:     coin.Id,
			Name:   coin.Name,
			Symbol: coin.Symbol,
		})
	}

	return cryptocurrenciesList, nil
}

func (c *CoinGeckoClient) GetCryptocurrencyDataById(id string) (domain.CryptocurrencyData, error) {
	requestUrl := fmt.Sprintf("%s/coins/%s", coinGeckoBaseURL, id)

	// The api key is optional. CoinGecko serves these endpoints anonymously at a lower
	// rate limit, so the header is only set when a key is configured.
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return domain.CryptocurrencyData{}, err
	}
	if c.apiKey != "" {
		req.Header.Set("x-cg-demo-api-key", c.apiKey)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return domain.CryptocurrencyData{}, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return domain.CryptocurrencyData{}, fmt.Errorf("failed to get coin data: %s", resp.Status)
	}

	// Parse the response
	var coinData CoinGeckoCoinData
	if err := json.NewDecoder(resp.Body).Decode(&coinData); err != nil {
		return domain.CryptocurrencyData{}, err
	}

	var cryptocurrencyData domain.CryptocurrencyData
	cryptocurrencyData.Id = coinData.ID
	cryptocurrencyData.Name = coinData.Name
	cryptocurrencyData.Symbol = coinData.Symbol
	cryptocurrencyData.Description = coinData.Description.En
	cryptocurrencyData.Whitepaper = coinData.Links.Whitepaper
	cryptocurrencyData.CurrentUsdPrice = coinData.MarketData.CurrentPrice.USD
	cryptocurrencyData.MarketCapUsd = coinData.MarketData.MarketCap.USD
	cryptocurrencyData.PriceChangePercentage24h = coinData.MarketData.PriceChangePercentage24h
	cryptocurrencyData.PriceChangePercentage7d = coinData.MarketData.PriceChangePercentage7d
	cryptocurrencyData.PriceChangePercentage14d = coinData.MarketData.PriceChangePercentage14d
	cryptocurrencyData.PriceChangePercentage30d = coinData.MarketData.PriceChangePercentage30d
	cryptocurrencyData.PriceChangePercentage60d = coinData.MarketData.PriceChangePercentage60d
	cryptocurrencyData.PriceChangePercentage200d = coinData.MarketData.PriceChangePercentage200d
	cryptocurrencyData.PriceChangePercentage1y = coinData.MarketData.PriceChangePercentage1y
	cryptocurrencyData.TotalSupply = coinData.MarketData.TotalSupply

	if coinData.MarketData.MaxSupply != nil {
		cryptocurrencyData.MaxSupply = *coinData.MarketData.MaxSupply
	}

	return cryptocurrencyData, nil
}

// SearchCryptocurrencies delegates to CoinGecko's own search, which ranks by market
// cap. Scanning /coins/list instead matches unranked symbol squatters first: a query
// for "btc" hits the memecoin "batcat" long before bitcoin.
func (c *CoinGeckoClient) SearchCryptocurrencies(query string) ([]domain.Cryptocurrency, error) {
	requestUrl := fmt.Sprintf("%s/search?query=%s", coinGeckoBaseURL, url.QueryEscape(query))

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}
	if c.apiKey != "" {
		req.Header.Set("x-cg-demo-api-key", c.apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search coins: %s", resp.Status)
	}

	var searchResponse CoinGeckoSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	cryptocurrencies := make([]domain.Cryptocurrency, 0, len(searchResponse.Coins))
	for _, coin := range searchResponse.Coins {
		cryptocurrencies = append(cryptocurrencies, domain.Cryptocurrency{
			Id:     coin.Id,
			Name:   coin.Name,
			Symbol: coin.Symbol,
		})
	}

	return cryptocurrencies, nil
}
