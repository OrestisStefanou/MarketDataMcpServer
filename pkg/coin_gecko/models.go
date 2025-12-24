package coingecko

type CoinGeckoCoin struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type CoinGeckoCoinData struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`

	Description struct {
		En string `json:"en"`
	} `json:"description"`

	Links struct {
		Whitepaper string `json:"whitepaper"`
	} `json:"links"`

	MarketData struct {
		CurrentPrice struct {
			USD float64 `json:"usd"`
		} `json:"current_price"`

		MarketCap struct {
			USD float64 `json:"usd"`
		} `json:"market_cap"`

		PriceChangePercentage24h  float64 `json:"price_change_percentage_24h"`
		PriceChangePercentage7d   float64 `json:"price_change_percentage_7d"`
		PriceChangePercentage14d  float64 `json:"price_change_percentage_14d"`
		PriceChangePercentage30d  float64 `json:"price_change_percentage_30d"`
		PriceChangePercentage60d  float64 `json:"price_change_percentage_60d"`
		PriceChangePercentage200d float64 `json:"price_change_percentage_200d"`
		PriceChangePercentage1y   float64 `json:"price_change_percentage_1y"`

		TotalSupply float64  `json:"total_supply"`
		MaxSupply   *float64 `json:"max_supply"`
	} `json:"market_data"`
}
