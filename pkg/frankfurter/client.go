package frankfurter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/errors"
)

var (
	// frankfurterBaseUrl serves the European Central Bank's daily reference rates and needs
	// no API key. Note the .dev host: the older .app host permanently redirects here. It is
	// a var rather than a const only so that tests can point it at a stub server.
	frankfurterBaseUrl = "https://api.frankfurter.dev/v1"

	httpClient = &http.Client{Timeout: 10 * time.Second}
)

type FrankfurterClient struct{}

func NewFrankfurterClient() (*FrankfurterClient, error) {
	return &FrankfurterClient{}, nil
}

func (c *FrankfurterClient) GetCurrencyExchangeRate(fromCurrency domain.Currency, toCurrency domain.Currency) (domain.CurrencyExchangeRate, error) {
	fromCurrencyName, ok := domain.CurrencyCodeToNameMap[fromCurrency]
	if !ok {
		return domain.CurrencyExchangeRate{}, fmt.Errorf("unsupported currency: %s", fromCurrency)
	}

	toCurrencyName, ok := domain.CurrencyCodeToNameMap[toCurrency]
	if !ok {
		return domain.CurrencyExchangeRate{}, fmt.Errorf("unsupported currency: %s", toCurrency)
	}

	exchangeRate := domain.CurrencyExchangeRate{
		FromCurrency:     fromCurrency,
		FromCurrencyName: fromCurrencyName,
		ToCurrency:       toCurrency,
		ToCurrencyName:   toCurrencyName,
	}

	// Frankfurter has no rate for a currency against itself.
	if fromCurrency == toCurrency {
		exchangeRate.Rate = 1
		return exchangeRate, nil
	}

	url := fmt.Sprintf("%s/latest?base=%s&symbols=%s", frankfurterBaseUrl, fromCurrency, toCurrency)
	resp, err := httpClient.Get(url)
	if err != nil {
		return domain.CurrencyExchangeRate{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.CurrencyExchangeRate{}, errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("Call to get the %s to %s exchange rate failed", fromCurrency, toCurrency),
		}
	}

	var response latestResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return domain.CurrencyExchangeRate{}, errors.StreamError{
			Message: "failed to decode the Frankfurter exchange rate response",
			Err:     err,
		}
	}

	// Frankfurter drops symbols it does not publish instead of erroring, so an unchecked
	// map read here would return a rate of 0 that looks like a real answer.
	rate, ok := response.Rates[string(toCurrency)]
	if !ok {
		return domain.CurrencyExchangeRate{}, fmt.Errorf(
			"no published exchange rate from %s to %s", fromCurrency, toCurrency,
		)
	}

	exchangeRate.Rate = rate
	return exchangeRate, nil
}
