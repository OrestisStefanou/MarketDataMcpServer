package domain

type IndustryStock struct {
	Symbol      string
	CompanyName string
	MarketCap   float32
}

type Industry struct {
	Name             string
	UrlName          string
	NumberOfStocks   int
	MarketCap        float32
	DividendYieldPct float32
	PeRatio          float32
	ProfitMarginPct  float32
	OneYearChangePct float32
}
