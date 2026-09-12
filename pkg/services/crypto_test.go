package services

import (
	"errors"
	"testing"

	"market_data_mcp_server/pkg/domain"
)

// listStub serves a /coins/list ordering in which squatters precede the real asset,
// which is how CoinGecko actually returns it: the list carries no market cap rank.
type listStub struct {
	searchErr bool
}

func (s *listStub) GetCryptocurrenciesList() ([]domain.Cryptocurrency, error) {
	return []domain.Cryptocurrency{
		{Id: "aurasaylorfortniteyoda67inu", Name: "AuraSaylorFortniteYoda67Inu", Symbol: "bitcoin"},
		{Id: "batcat", Name: "batcat", Symbol: "btc"},
		{Id: "bitcoin", Name: "Bitcoin", Symbol: "btc"},
	}, nil
}

func (s *listStub) GetCryptocurrencyDataById(id string) (domain.CryptocurrencyData, error) {
	return domain.CryptocurrencyData{}, nil
}

func (s *listStub) SearchCryptocurrencies(query string) ([]domain.Cryptocurrency, error) {
	if s.searchErr {
		return nil, errors.New("upstream search unavailable")
	}
	return []domain.Cryptocurrency{{Id: "bitcoin", Name: "Bitcoin", Symbol: "btc"}}, nil
}

// The ranked upstream search is what separates the real asset from a symbol squatter.
func TestSearchCryptocurrenciesPrefersRankedUpstream(t *testing.T) {
	service, err := NewCryptoService(&listStub{}, nil)
	if err != nil {
		t.Fatalf("NewCryptoService: %v", err)
	}

	for _, query := range []string{"BTC", "Bitcoin"} {
		results, err := service.SearchCryptocurrencies(query)
		if err != nil {
			t.Fatalf("%s: %v", query, err)
		}
		if len(results) == 0 {
			t.Fatalf("%s: no results", query)
		}
		if results[0].Id != "bitcoin" {
			t.Errorf("%s: got %q, want bitcoin", query, results[0].Id)
		}
	}
}

// With the upstream search down, an exact id or name match must still outrank an exact
// symbol match that merely sits earlier in the list.
func TestSearchCryptocurrenciesFallbackRanksRatherThanReturningFirstExactHit(t *testing.T) {
	service, err := NewCryptoService(&listStub{searchErr: true}, nil)
	if err != nil {
		t.Fatalf("NewCryptoService: %v", err)
	}

	results, err := service.SearchCryptocurrencies("bitcoin")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(results) < 2 {
		t.Fatalf("got %d results, want the full ranked set", len(results))
	}
	if results[0].Id != "bitcoin" {
		t.Errorf("got %q first, want bitcoin (exact id beats an exact symbol squatter)", results[0].Id)
	}
}
