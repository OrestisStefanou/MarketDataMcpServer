package domain

import "time"

type Period string

const (
	Period1D Period = "1d"
	Period5D Period = "5d"
	Period1M Period = "1m"
	Period6M Period = "6m"
	Period1Y Period = "1y"
	Period5Y Period = "5y"
)

type Price struct {
	Date       time.Time
	ClosePrice float64
}

type HistoricalPrices struct {
	Period           Period
	Prices           []Price
	PercentageChange float64
}
