package domain

type InsiderTransaction struct {
	TransactionDate       string
	Ticker                string
	Executive             string
	ExecutiveTitle        string
	SecurityType          string
	AcquisitionOrDisposal string
	Shares                float64
	SharePrice            float64
}
