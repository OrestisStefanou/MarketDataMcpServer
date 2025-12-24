package services

import "market_data_mcp_server/pkg/domain"

type SuperInvestorDataService interface {
	GetSuperInvestors() ([]domain.SuperInvestor, error)
	GetSuperInvestorPortfolio(superInvestorName string) (domain.SuperInvestorPortfolio, error)
}

type SuperInvestorService struct {
	dataService SuperInvestorDataService
}

func NewSuperInvestorService(dataService SuperInvestorDataService) (*SuperInvestorService, error) {
	return &SuperInvestorService{dataService: dataService}, nil
}

func (s SuperInvestorService) GetSuperInvestors() ([]domain.SuperInvestor, error) {
	return s.dataService.GetSuperInvestors()
}

func (s SuperInvestorService) GetSuperInvestorPortfolio(superInvestorName string) (domain.SuperInvestorPortfolio, error) {
	return s.dataService.GetSuperInvestorPortfolio(superInvestorName)
}
