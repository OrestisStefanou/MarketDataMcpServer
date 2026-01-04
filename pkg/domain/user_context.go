package domain

type AssetClass string

const (
	Stock         AssetClass = "stock"
	ETF           AssetClass = "etf"
	Crypto        AssetClass = "crypto"
	MutualFund    AssetClass = "mutual_fund"
	Bond          AssetClass = "bond"
	Cash          AssetClass = "cash"
	RealEstate    AssetClass = "real_estate"
	PrivateEquity AssetClass = "private_equity"
	Commodities   AssetClass = "commodities"
)

type UserPortfolioHolding struct {
	AssetClass          AssetClass
	Symbol              string
	Name                string
	Quantity            float64
	PortfolioPercentage float64
}

type UserContext struct {
	UserID        string
	UserProfile   map[string]any
	UserPortfolio []UserPortfolioHolding
}
