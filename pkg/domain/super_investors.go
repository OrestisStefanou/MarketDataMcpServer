package domain

type SuperInvestorPortfolioHolding struct {
	Stock          string
	PortfolioPct   string
	RecentActivity string
	Shares         string
	Value          string
}

type SuperInvestorPortfolioSectorAnalysis struct {
	Sector       string
	PortfolioPct string
}

type SuperInvestor struct {
	Name string
}

type SuperInvestorPortfolio struct {
	Holdings       []SuperInvestorPortfolioHolding
	SectorAnalysis []SuperInvestorPortfolioSectorAnalysis
}
