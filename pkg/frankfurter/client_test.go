package frankfurter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"market_data_mcp_server/pkg/domain"
)

// Frankfurter omits a symbol it does not publish rather than reporting an error. Reading the
// map without checking for presence would return a rate of 0 that looks like a real answer,
// so this is the most important behaviour in the package.
func TestGetCurrencyExchangeRateRejectsAMissingRate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(latestResponse{
			Amount: 1,
			Base:   "USD",
			Date:   "2026-08-01",
			Rates:  map[string]float64{},
		})
	}))
	defer server.Close()

	withBaseUrl(t, server.URL)

	if _, err := (&FrankfurterClient{}).GetCurrencyExchangeRate(domain.USD, domain.GBP); err == nil {
		t.Fatal("expected an error when the response omits the requested rate")
	}
}

func TestGetCurrencyExchangeRate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(latestResponse{
			Amount: 1,
			Base:   "EUR",
			Date:   "2026-08-01",
			Rates:  map[string]float64{"USD": 1.1485},
		})
	}))
	defer server.Close()

	withBaseUrl(t, server.URL)

	rate, err := (&FrankfurterClient{}).GetCurrencyExchangeRate(domain.EUR, domain.USD)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rate.Rate != 1.1485 {
		t.Errorf("expected 1.1485, got %v", rate.Rate)
	}
	// The API returns no currency names, so they come from the domain map.
	if rate.FromCurrencyName != domain.Euro || rate.ToCurrencyName != domain.UnitedStatesDollar {
		t.Errorf("unexpected currency names: %+v", rate)
	}
}

// Frankfurter has no rate for a currency against itself.
func TestGetCurrencyExchangeRateForTheSameCurrency(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("the same currency should not require an HTTP call")
	}))
	defer server.Close()

	withBaseUrl(t, server.URL)

	rate, err := (&FrankfurterClient{}).GetCurrencyExchangeRate(domain.USD, domain.USD)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rate.Rate != 1 {
		t.Errorf("expected a rate of 1, got %v", rate.Rate)
	}
}

// AED was dropped when the FX source moved to the ECB, which does not publish it.
func TestGetCurrencyExchangeRateRejectsUnsupportedCurrencies(t *testing.T) {
	if _, err := (&FrankfurterClient{}).GetCurrencyExchangeRate("AED", domain.USD); err == nil {
		t.Fatal("expected an error for a currency outside the supported set")
	}
	if _, err := (&FrankfurterClient{}).GetCurrencyExchangeRate(domain.USD, "AED"); err == nil {
		t.Fatal("expected an error for a currency outside the supported set")
	}
}

func withBaseUrl(t *testing.T, url string) {
	t.Helper()

	original := frankfurterBaseUrl
	frankfurterBaseUrl = url
	t.Cleanup(func() { frankfurterBaseUrl = original })
}
