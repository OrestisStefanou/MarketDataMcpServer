package domain

// InsiderTransactions carries one company-year of insider dealing.
//
// Truncated records whether the year held more filings than a single request will fetch.
// It matters because an under-reported year reads exactly like a quiet one: a partial view
// of 2021 Tesla filings would understate insider selling rather than look obviously wrong.
type InsiderTransactions struct {
	Symbol         string
	Year           int
	Transactions   []InsiderTransaction
	FilingsFound   int
	FilingsFetched int
	Truncated      bool
}

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
