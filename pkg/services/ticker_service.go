package services

import (
	"fmt"
	"market_data_mcp_server/pkg/domain"
)

type TickerDataService interface {
	GetTickers(query string) ([]domain.Ticker, error)
}

type TickerService struct {
	dataService TickerDataService
}

func NewTickerService(dataService TickerDataService) (*TickerService, error) {
	return &TickerService{
		dataService: dataService,
	}, nil
}

type TickerFilterOptions struct {
	Limit        int
	Page         int
	SearchString string
}

func (s TickerService) GetTickers(filters TickerFilterOptions) ([]domain.Ticker, error) {
	if filters.SearchString == "" {
		return nil, fmt.Errorf("search_string is required")
	}

	tickers, err := s.dataService.GetTickers(filters.SearchString)
	if err != nil {
		return nil, err
	}

	if filters.Limit > 0 && len(tickers) > filters.Limit {
		tickers = tickers[:filters.Limit]
	}

	return tickers, nil
}
