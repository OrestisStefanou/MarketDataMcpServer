package services

import (
	"market_data_mcp_server/pkg/domain"
	"strings"
)

type EtfDataService interface {
	GetEtfs() ([]domain.Etf, error)
	GetEtfOverview(symbol string) (domain.EtfOverview, error)
}

type EtfService struct {
	dataService EtfDataService
}

func NewEtfService(dataService EtfDataService) (*EtfService, error) {
	return &EtfService{dataService: dataService}, nil
}

type EtfFilterOptions struct {
	SearchString string
}

func (f EtfFilterOptions) IsEmpty() bool {
	return f.SearchString == ""
}

func (f EtfFilterOptions) HasSearchString() bool {
	return f.SearchString != ""
}

func (s EtfService) GetEtfs(filters EtfFilterOptions) ([]domain.Etf, error) {
	etfs, err := s.dataService.GetEtfs()
	if err != nil {
		return nil, err
	}

	if filters.IsEmpty() {
		return etfs, nil
	}

	if filters.HasSearchString() {
		filteredEtfs := make([]domain.Etf, 0)
		for _, e := range etfs {
			search := strings.ToLower(filters.SearchString)
			symbol := strings.ToLower(e.Symbol)
			etfName := strings.ToLower(e.Name)
			if strings.Contains(symbol, search) || strings.Contains(etfName, search) {
				filteredEtfs = append(filteredEtfs, e)
			}

		}
		return filteredEtfs, nil
	}

	return etfs, nil
}

func (s EtfService) GetEtf(etfSymbol string) (domain.EtfOverview, error) {
	return s.dataService.GetEtfOverview(etfSymbol)
}
