package domain

type StockProfile struct {
	Name        string
	Description string
	Country     string
	Founded     int
	IpoDate     string
	Industry    string
	Sector      string
	Ceo         string
}

type StockEstimation struct {
	Date          string
	Eps           float64
	EpsGrowth     float64
	FiscalQuarter string
	FiscalYear    string
	Revenue       float64
	RevenueGrowth float64
}

type StockTargetPrc struct {
	Average float32
	High    float32
	Low     float32
	Median  float32
}

type StockForecast struct {
	Estimations []StockEstimation
	TargetPrice StockTargetPrc
}

type Ticker struct {
	Symbol      string
	CompanyName string
}

type EarningsCallTranscript struct {
	Speaker   string
	Title     string
	Content   string
	Sentiment string
}

type Quarter string

const (
	Q1Quarter Quarter = "Q1"
	Q2Quarter Quarter = "Q2"
	Q3Quarter Quarter = "Q3"
	Q4Quarter Quarter = "Q4"
)
