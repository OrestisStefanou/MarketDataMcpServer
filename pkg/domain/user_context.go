package domain

type AssetClass string

const (
	Stock  AssetClass = "stock"
	ETF    AssetClass = "etf"
	Crypto AssetClass = "crypto"
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
	CreatedAt     string // ISO 8601 format
	UpdatedAt     string // ISO 8601 format
}
