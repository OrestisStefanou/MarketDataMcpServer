package services

import (
	"market_data_mcp_server/pkg/domain"
	"strings"
)

type ICryptoDataService interface {
	GetCryptocurrenciesList() ([]domain.Cryptocurrency, error)
	GetCryptocurrencyDataById(id string) (domain.CryptocurrencyData, error)
}

type CryptoNewsSource interface {
	GetCryptocurrencyNews(symbol string) ([]domain.NewsArticle, error)
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

func (s *CryptoService) SearchCryptocurrencies(query string) ([]domain.Cryptocurrency, error) {
	cryptocurrenciesList, err := s.cryptoDataService.GetCryptocurrenciesList()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	var searchResults []domain.Cryptocurrency
	for _, cryptocurrency := range cryptocurrenciesList {
		cryptoName := strings.ToLower(cryptocurrency.Name)
		cryptoSymbol := strings.ToLower(cryptocurrency.Symbol)
		cryptoId := strings.ToLower(cryptocurrency.Id)
		if strings.Contains(cryptoName, query) {
			searchResults = append(searchResults, cryptocurrency)
			if cryptoName == query {
				return []domain.Cryptocurrency{cryptocurrency}, nil
			}
			continue
		}
		if strings.Contains(cryptoSymbol, query) {
			searchResults = append(searchResults, cryptocurrency)
			if cryptoSymbol == query {
				return []domain.Cryptocurrency{cryptocurrency}, nil
			}
			continue
		}
		if strings.Contains(cryptoId, query) {
			searchResults = append(searchResults, cryptocurrency)
			if cryptoId == query {
				return []domain.Cryptocurrency{cryptocurrency}, nil
			}
			continue
		}
	}

	return searchResults, nil
}

func (s *CryptoService) GetCryptocurrencyNews(symbol string) ([]domain.NewsArticle, error) {
	return s.cryptoNewsSource.GetCryptocurrencyNews(symbol)
}
