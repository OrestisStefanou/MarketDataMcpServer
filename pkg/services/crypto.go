package services

import (
	"market_data_mcp_server/pkg/domain"
	"sort"
	"strings"
)

type ICryptoDataService interface {
	GetCryptocurrenciesList() ([]domain.Cryptocurrency, error)
	GetCryptocurrencyDataById(id string) (domain.CryptocurrencyData, error)
	SearchCryptocurrencies(query string) ([]domain.Cryptocurrency, error)
}

type CryptoNewsSource interface {
	GetCryptocurrencyNews(symbol string) (domain.CryptocurrencyNews, error)
}

type CryptoService struct {
	cryptoDataService ICryptoDataService
	cryptoNewsSource  CryptoNewsSource
}

func NewCryptoService(cryptoDataService ICryptoDataService, cryptoNewsSource CryptoNewsSource) (*CryptoService, error) {
	return &CryptoService{cryptoDataService: cryptoDataService, cryptoNewsSource: cryptoNewsSource}, nil
}

func (s *CryptoService) GetCryptocurrenciesList() ([]domain.Cryptocurrency, error) {
	return s.cryptoDataService.GetCryptocurrenciesList()
}

func (s *CryptoService) GetCryptocurrencyDataById(id string) (domain.CryptocurrencyData, error) {
	return s.cryptoDataService.GetCryptocurrencyDataById(id)
}

// SearchCryptocurrencies prefers the upstream ranked search, which orders by market
// cap. The local scan below is only a fallback: /coins/list carries no rank, so ties on
// an exact symbol match resolve by list position, and squatters win those ties.
func (s *CryptoService) SearchCryptocurrencies(query string) ([]domain.Cryptocurrency, error) {
	searchResults, err := s.cryptoDataService.SearchCryptocurrencies(query)
	if err == nil && len(searchResults) > 0 {
		return searchResults, nil
	}

	return s.searchCryptocurrenciesLocally(query)
}

// searchCryptocurrenciesLocally ranks every match rather than returning the first exact
// hit, so the ordering stays stable when the upstream search is unavailable.
func (s *CryptoService) searchCryptocurrenciesLocally(query string) ([]domain.Cryptocurrency, error) {
	cryptocurrenciesList, err := s.cryptoDataService.GetCryptocurrenciesList()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)

	type rankedResult struct {
		cryptocurrency domain.Cryptocurrency
		rank           int
	}

	var ranked []rankedResult
	for _, cryptocurrency := range cryptocurrenciesList {
		cryptoName := strings.ToLower(cryptocurrency.Name)
		cryptoSymbol := strings.ToLower(cryptocurrency.Symbol)
		cryptoId := strings.ToLower(cryptocurrency.Id)

		rank := -1
		switch {
		case cryptoId == query:
			rank = 0
		case cryptoName == query:
			rank = 1
		case cryptoSymbol == query:
			rank = 2
		case strings.HasPrefix(cryptoName, query) || strings.HasPrefix(cryptoId, query):
			rank = 3
		case strings.Contains(cryptoName, query) || strings.Contains(cryptoSymbol, query) || strings.Contains(cryptoId, query):
			rank = 4
		}

		if rank >= 0 {
			ranked = append(ranked, rankedResult{cryptocurrency: cryptocurrency, rank: rank})
		}
	}

	sort.SliceStable(ranked, func(i, j int) bool { return ranked[i].rank < ranked[j].rank })

	searchResults := make([]domain.Cryptocurrency, 0, len(ranked))
	for _, result := range ranked {
		searchResults = append(searchResults, result.cryptocurrency)
	}

	return searchResults, nil
}

func (s *CryptoService) GetCryptocurrencyNews(symbol string) (domain.CryptocurrencyNews, error) {
	return s.cryptoNewsSource.GetCryptocurrencyNews(symbol)
}
