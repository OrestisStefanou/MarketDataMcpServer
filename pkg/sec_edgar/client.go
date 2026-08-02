package secedgar

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/errors"
)

const (
	secTickersUrl           = "https://www.sec.gov/files/company_tickers.json"
	secSubmissionsUrlFormat = "https://data.sec.gov/submissions/CIK%010d.json"
	secArchivesUrlFormat    = "https://www.sec.gov/Archives/edgar/data/%d/%s/%s"

	// maxForm4FilingsPerRequest bounds how many filings a single tool call will fetch, since
	// each filing is its own HTTP request. It is set above the busiest year of a heavily
	// traded issuer, and any year that still exceeds it is reported as truncated rather than
	// silently short.
	maxForm4FilingsPerRequest = 120
	form4FetchConcurrency     = 4
)

// xslPrefixPattern strips the rendering-stylesheet prefix EDGAR puts on primaryDocument
// (xslF345X05/, xslF345X03/ and friends) to get at the raw XML.
var xslPrefixPattern = regexp.MustCompile(`^xsl[^/]*/`)

// SecEdgarClient reads insider transactions from the SEC's EDGAR full-text archive.
// EDGAR needs no API key but does require a descriptive User-Agent on every request.
type SecEdgarClient struct {
	userAgent  string
	httpClient *http.Client
}

func NewSecEdgarClient(userAgent string) (*SecEdgarClient, error) {
	return &SecEdgarClient{
		userAgent:  userAgent,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (c *SecEdgarClient) GetInsiderTransactions(symbol string, year int) (domain.InsiderTransactions, error) {
	cik, err := c.GetCikBySymbol(symbol)
	if err != nil {
		return domain.InsiderTransactions{}, err
	}

	submissions, err := c.GetSubmissions(cik)
	if err != nil {
		return domain.InsiderTransactions{}, err
	}

	filings, total, err := SelectForm4Filings(submissions, year)
	if err != nil {
		return domain.InsiderTransactions{}, err
	}

	return domain.InsiderTransactions{
		Symbol:         symbol,
		Year:           year,
		Transactions:   FilterByYear(c.FetchForm4Transactions(cik, filings), year),
		FilingsFound:   total,
		FilingsFetched: len(filings),
		Truncated:      len(filings) < total,
	}, nil
}

// GetCikBySymbol resolves a ticker to its EDGAR central index key.
func (c *SecEdgarClient) GetCikBySymbol(symbol string) (int, error) {
	tickerToCik, err := c.GetTickerToCikMap()
	if err != nil {
		return 0, err
	}

	cik, err := lookupCik(tickerToCik, symbol)
	if err != nil {
		return 0, err
	}

	return cik, nil
}

// GetTickerToCikMap downloads EDGAR's full ticker to CIK mapping. It is a large file that
// changes rarely, so callers are expected to cache it.
func (c *SecEdgarClient) GetTickerToCikMap() (map[string]int, error) {
	body, err := c.get(secTickersUrl)
	if err != nil {
		return nil, err
	}

	var response companyTickersResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.StreamError{Message: "failed to decode the SEC ticker map", Err: err}
	}

	tickerToCik := make(map[string]int, len(response))
	for _, company := range response {
		cik, err := company.Cik.Int64()
		if err != nil {
			continue
		}
		tickerToCik[strings.ToUpper(company.Ticker)] = int(cik)
	}

	return tickerToCik, nil
}

func (c *SecEdgarClient) GetSubmissions(cik int) (Submissions, error) {
	body, err := c.get(fmt.Sprintf(secSubmissionsUrlFormat, cik))
	if err != nil {
		return Submissions{}, err
	}

	var submissions Submissions
	if err := json.Unmarshal(body, &submissions); err != nil {
		return Submissions{}, errors.StreamError{
			Message: fmt.Sprintf("failed to decode the SEC submissions index for CIK %d", cik),
			Err:     err,
		}
	}

	return submissions, nil
}

// FetchForm4Transactions downloads and parses the given filings concurrently. A filing that
// fails to download or parse is skipped rather than failing the whole call, since one
// malformed legacy document should not lose the rest of the year.
func (c *SecEdgarClient) FetchForm4Transactions(cik int, filings []Form4Filing) []domain.InsiderTransaction {
	results := make([][]domain.InsiderTransaction, len(filings))

	var waitGroup sync.WaitGroup
	semaphore := make(chan struct{}, form4FetchConcurrency)

	for i, filing := range filings {
		waitGroup.Add(1)
		go func(index int, filing Form4Filing) {
			defer waitGroup.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			transactions, err := c.fetchForm4(cik, filing)
			if err != nil {
				return
			}
			results[index] = transactions
		}(i, filing)
	}

	waitGroup.Wait()

	transactions := make([]domain.InsiderTransaction, 0, len(filings)*2)
	for _, result := range results {
		transactions = append(transactions, result...)
	}

	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].TransactionDate > transactions[j].TransactionDate
	})

	return transactions
}

func (c *SecEdgarClient) fetchForm4(cik int, filing Form4Filing) ([]domain.InsiderTransaction, error) {
	accession := strings.ReplaceAll(filing.AccessionNumber, "-", "")
	document := xslPrefixPattern.ReplaceAllString(filing.PrimaryDocument, "")

	body, err := c.get(fmt.Sprintf(secArchivesUrlFormat, cik, accession, document))
	if err != nil {
		return nil, err
	}

	return parseForm4(body)
}

// get issues a rate-limited GET carrying the User-Agent the SEC requires. Without one that
// names a contact address EDGAR replies 403. Note that EDGAR also rejects a User-Agent
// containing a URL, so the contact must be an email address.
//
// Accept-Encoding is deliberately left unset so that net/http negotiates gzip and
// transparently decompresses the response.
func (c *SecEdgarClient) get(url string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", c.userAgent)

	limiter.wait()

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("Call to SEC EDGAR %s failed", url),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.StreamError{Message: "failed to read the SEC EDGAR response", Err: err}
	}

	return body, nil
}

// SelectForm4Filings narrows a submissions index to the Form 4s worth fetching for a year.
//
// The window runs to mid February of the following year because a transaction in late
// December is often filed in the new year. Filings are newest first, so the cap keeps the
// most recent.
// It also reports the total number of filings in the window, so the caller can tell whether
// the cap dropped any.
func SelectForm4Filings(submissions Submissions, year int) ([]Form4Filing, int, error) {
	windowStart := fmt.Sprintf("%d-01-01", year)
	windowEnd := fmt.Sprintf("%d-02-15", year+1)

	filings := make([]Form4Filing, 0, maxForm4FilingsPerRequest)
	total := 0
	for _, filing := range submissions.Form4Filings() {
		if filing.FilingDate < windowStart || filing.FilingDate > windowEnd {
			continue
		}
		total++
		if len(filings) < maxForm4FilingsPerRequest {
			filings = append(filings, filing)
		}
	}

	if total == 0 {
		// Distinguish "nothing was filed" from "the index does not reach back that far",
		// because EDGAR only keeps the most recent filings in this index.
		earliest := submissions.EarliestFilingDate()
		if earliest != "" && windowEnd < earliest {
			return nil, 0, fmt.Errorf(
				"SEC EDGAR only indexes filings back to %s for this company, so %d is out of range",
				earliest, year,
			)
		}
	}

	return filings, total, nil
}

// lookupCik resolves a ticker, accounting for EDGAR writing share classes with a hyphen
// (BRK-B) where quote feeds normally use a dot (BRK.B).
func lookupCik(tickerToCik map[string]int, symbol string) (int, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	if cik, ok := tickerToCik[symbol]; ok {
		return cik, nil
	}
	if cik, ok := tickerToCik[strings.ReplaceAll(symbol, ".", "-")]; ok {
		return cik, nil
	}

	return 0, fmt.Errorf("no SEC EDGAR filer found for symbol %s", symbol)
}

// FilterByYear keeps only the transactions that actually took place in the requested year.
// The filing window is deliberately wider than the year, so this trims the overhang.
func FilterByYear(transactions []domain.InsiderTransaction, year int) []domain.InsiderTransaction {
	prefix := strconv.Itoa(year)

	filtered := make([]domain.InsiderTransaction, 0, len(transactions))
	for _, transaction := range transactions {
		if strings.HasPrefix(transaction.TransactionDate, prefix) {
			filtered = append(filtered, transaction)
		}
	}

	return filtered
}
