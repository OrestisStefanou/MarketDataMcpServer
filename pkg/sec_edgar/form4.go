package secedgar

import (
	"encoding/xml"
	"strconv"
	"strings"

	"market_data_mcp_server/pkg/domain"
)

// valueField wraps the <value> element Form 4 uses for almost every leaf, since each one
// can also carry a footnote reference alongside the value.
type valueField struct {
	Value string `xml:"value"`
}

type reportingOwner struct {
	Id struct {
		Name string `xml:"rptOwnerName"`
	} `xml:"reportingOwnerId"`
	Relationship struct {
		IsDirector        string `xml:"isDirector"`
		IsOfficer         string `xml:"isOfficer"`
		IsTenPercentOwner string `xml:"isTenPercentOwner"`
		IsOther           string `xml:"isOther"`
		OfficerTitle      string `xml:"officerTitle"`
	} `xml:"reportingOwnerRelationship"`
}

// form4Transaction covers both the non-derivative and derivative tables, which share the
// fields we care about.
type form4Transaction struct {
	SecurityTitle   valueField `xml:"securityTitle"`
	TransactionDate valueField `xml:"transactionDate"`
	Amounts         struct {
		Shares               valueField `xml:"transactionShares"`
		PricePerShare        valueField `xml:"transactionPricePerShare"`
		AcquiredDisposedCode valueField `xml:"transactionAcquiredDisposedCode"`
	} `xml:"transactionAmounts"`
}

type ownershipDocument struct {
	XMLName xml.Name `xml:"ownershipDocument"`
	Issuer  struct {
		Name          string `xml:"issuerName"`
		TradingSymbol string `xml:"issuerTradingSymbol"`
	} `xml:"issuer"`
	ReportingOwners []reportingOwner `xml:"reportingOwner"`
	// Holdings (nonDerivativeHolding, derivativeHolding) are deliberately excluded: they
	// report a standing position rather than a transaction.
	NonDerivativeTransactions []form4Transaction `xml:"nonDerivativeTable>nonDerivativeTransaction"`
	DerivativeTransactions    []form4Transaction `xml:"derivativeTable>derivativeTransaction"`
}

// parseForm4 turns the raw XML of a single Form 4 filing into insider transactions.
func parseForm4(raw []byte) ([]domain.InsiderTransaction, error) {
	var document ownershipDocument
	if err := xml.Unmarshal(raw, &document); err != nil {
		return nil, err
	}

	executive, executiveTitle := ownerDetails(document.ReportingOwners)
	ticker := strings.TrimSpace(document.Issuer.TradingSymbol)

	all := append(document.NonDerivativeTransactions, document.DerivativeTransactions...)
	transactions := make([]domain.InsiderTransaction, 0, len(all))
	for _, transaction := range all {
		date := strings.TrimSpace(transaction.TransactionDate.Value)
		if date == "" {
			continue
		}

		transactions = append(transactions, domain.InsiderTransaction{
			TransactionDate:       date,
			Ticker:                ticker,
			Executive:             executive,
			ExecutiveTitle:        executiveTitle,
			SecurityType:          strings.TrimSpace(transaction.SecurityTitle.Value),
			AcquisitionOrDisposal: strings.TrimSpace(transaction.Amounts.AcquiredDisposedCode.Value),
			Shares:                parseFloat(transaction.Amounts.Shares.Value),
			SharePrice:            parseFloat(transaction.Amounts.PricePerShare.Value),
		})
	}

	return transactions, nil
}

// ownerDetails collapses the reporting owners into a single name and title. A Form 4 is
// usually filed by one person, but joint filings do happen.
func ownerDetails(owners []reportingOwner) (string, string) {
	names := make([]string, 0, len(owners))
	titles := make([]string, 0, len(owners))

	for _, owner := range owners {
		if name := strings.TrimSpace(owner.Id.Name); name != "" {
			names = append(names, name)
		}
		if title := ownerTitle(owner); title != "" {
			titles = append(titles, title)
		}
	}

	return strings.Join(names, ", "), strings.Join(titles, ", ")
}

// ownerTitle prefers the free-text officer title, falling back to the relationship flags so
// that directors and ten percent owners are still described.
func ownerTitle(owner reportingOwner) string {
	relationship := owner.Relationship

	if title := strings.TrimSpace(relationship.OfficerTitle); title != "" {
		return title
	}

	roles := make([]string, 0, 4)
	if isFlagSet(relationship.IsOfficer) {
		roles = append(roles, "Officer")
	}
	if isFlagSet(relationship.IsDirector) {
		roles = append(roles, "Director")
	}
	if isFlagSet(relationship.IsTenPercentOwner) {
		roles = append(roles, "10% Owner")
	}
	if isFlagSet(relationship.IsOther) {
		roles = append(roles, "Other")
	}

	return strings.Join(roles, ", ")
}

// isFlagSet handles the two spellings EDGAR uses for booleans.
func isFlagSet(value string) bool {
	value = strings.TrimSpace(value)
	return value == "1" || strings.EqualFold(value, "true")
}

// parseFloat returns 0 for the values Form 4 legitimately leaves blank, such as the price
// of a gifted or inherited share.
func parseFloat(value string) float64 {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return 0
	}
	return parsed
}
