package services

import (
	"market_data_mcp_server/pkg/domain"
	"strings"
)

type TickerDataService interface {
	GetTickers() ([]domain.Ticker, error)
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

func (f TickerFilterOptions) IsEmpty() bool {
	return f.Limit == 0 && f.Page == 0 && f.SearchString == ""
}

func (f TickerFilterOptions) HasSearchString() bool {
	return f.SearchString != ""
}

func (s TickerService) GetTickers(filters TickerFilterOptions) ([]domain.Ticker, error) {
	// TODO: Implement the page and limit filtering
	tickers, err := s.dataService.GetTickers()
	if err != nil {
		return nil, err
	}

	if filters.IsEmpty() {
		return tickers, nil
	}

	if filters.HasSearchString() {
		filteredTickers := make([]domain.Ticker, 0)
		for _, t := range tickers {
			search := strings.ToLower(filters.SearchString)
			symbol := strings.ToLower(t.Symbol)
			companyName := strings.ToLower(t.CompanyName)
			if strings.Contains(symbol, search) || strings.Contains(companyName, search) {
				filteredTickers = append(filteredTickers, t)
			}

		}
		return filteredTickers, nil
	}

	return tickers, nil
}
