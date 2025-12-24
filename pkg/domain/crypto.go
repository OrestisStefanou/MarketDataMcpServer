package domain

type Cryptocurrency struct {
	Id     string
	Name   string
	Symbol string
}

type CryptocurrencyData struct {
	Id                        string
	Name                      string
	Symbol                    string
	Description               string
	Whitepaper                string
	CurrentUsdPrice           float64
	MarketCapUsd              float64
	PriceChangePercentage24h  float64
	PriceChangePercentage7d   float64
	PriceChangePercentage14d  float64
	PriceChangePercentage30d  float64
	PriceChangePercentage60d  float64
	PriceChangePercentage200d float64
	PriceChangePercentage1y   float64
	TotalSupply               float64
	MaxSupply                 float64
}
