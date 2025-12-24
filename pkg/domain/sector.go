package domain

type Sector struct {
	Name             string
	UrlName          string
	NumberOfStocks   int
	MarketCap        float32
	DividendYieldPct float32
	PeRatio          float32
	ProfitMarginPct  float32
	OneYearChangePct float32
}

type SectorStock struct {
	Symbol      string
	CompanyName string
	MarketCap   float32
}
