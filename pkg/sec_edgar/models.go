package secedgar

import "encoding/json"

// companyTickersResponse is the shape of https://www.sec.gov/files/company_tickers.json,
// a map keyed by an arbitrary row index.
type companyTickersResponse map[string]struct {
	Cik    json.Number `json:"cik_str"`
	Ticker string      `json:"ticker"`
	Title  string      `json:"title"`
}

// Form4Filing identifies one Form 4 in a company's filing history.
type Form4Filing struct {
	AccessionNumber string
	FilingDate      string
	PrimaryDocument string
}

// Submissions is the subset of a company's EDGAR submissions index that we need.
//
// EDGAR stores the recent filings as parallel arrays rather than a list of objects, and
// caps them at roughly the 1,000 most recent. Older filings live in separate archive files
// that we deliberately do not page through.
type Submissions struct {
	Filings struct {
		Recent struct {
			Form            []string `json:"form"`
			FilingDate      []string `json:"filingDate"`
			AccessionNumber []string `json:"accessionNumber"`
			PrimaryDocument []string `json:"primaryDocument"`
		} `json:"recent"`
	} `json:"filings"`
}

// Form4Filings returns every Form 4 in the submissions index, newest first.
func (s Submissions) Form4Filings() []Form4Filing {
	recent := s.Filings.Recent

	filings := make([]Form4Filing, 0, len(recent.Form))
	for i, form := range recent.Form {
		if form != "4" {
			continue
		}
		if i >= len(recent.FilingDate) || i >= len(recent.AccessionNumber) || i >= len(recent.PrimaryDocument) {
			break
		}
		filings = append(filings, Form4Filing{
			AccessionNumber: recent.AccessionNumber[i],
			FilingDate:      recent.FilingDate[i],
			PrimaryDocument: recent.PrimaryDocument[i],
		})
	}

	return filings
}

// EarliestFilingDate reports the oldest filing date held in the submissions index, which is
// how far back a query can reach before it needs the archive files.
func (s Submissions) EarliestFilingDate() string {
	earliest := ""
	for _, date := range s.Filings.Recent.FilingDate {
		if earliest == "" || date < earliest {
			earliest = date
		}
	}
	return earliest
}
