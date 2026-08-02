package secedgar

import (
	"testing"

	"market_data_mcp_server/pkg/domain"
)

// EDGAR writes share classes with a hyphen where quote feeds use a dot.
func TestLookupCikNormalisesShareClasses(t *testing.T) {
	tickerToCik := map[string]int{"TSLA": 1318605, "BRK-B": 1067983}

	for _, symbol := range []string{"TSLA", "tsla", " TSLA "} {
		if cik, err := lookupCik(tickerToCik, symbol); err != nil || cik != 1318605 {
			t.Errorf("expected 1318605 for %q, got %d (%v)", symbol, cik, err)
		}
	}

	for _, symbol := range []string{"BRK-B", "BRK.B", "brk.b"} {
		if cik, err := lookupCik(tickerToCik, symbol); err != nil || cik != 1067983 {
			t.Errorf("expected 1067983 for %q, got %d (%v)", symbol, cik, err)
		}
	}

	if _, err := lookupCik(tickerToCik, "NOTREAL"); err == nil {
		t.Error("expected an error for an unknown symbol")
	}
}

// The prefix varies across filings (xslF345X05, xslF345X03), so it cannot be stripped with
// a literal.
func TestXslPrefixStripping(t *testing.T) {
	tests := map[string]string{
		"xslF345X05/tm123_4seq1.xml": "tm123_4seq1.xml",
		"xslF345X03/edgardoc.xml":    "edgardoc.xml",
		"xslF345X06/doc4.xml":        "doc4.xml",
		"primary_doc.xml":            "primary_doc.xml",
	}

	for input, expected := range tests {
		if got := xslPrefixPattern.ReplaceAllString(input, ""); got != expected {
			t.Errorf("stripping %q: expected %q, got %q", input, expected, got)
		}
	}
}

func TestSelectForm4Filings(t *testing.T) {
	submissions := submissionsWith(
		[]string{"4", "4", "8-K", "4", "4"},
		[]string{"2026-01-20", "2025-06-15", "2025-03-01", "2025-01-10", "2019-05-02"},
	)

	// The window runs into mid February of the next year, because a December transaction is
	// often filed in January.
	filings, total, err := SelectForm4Filings(submissions, 2025)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 3 {
		t.Errorf("expected 3 filings in the window, got %d", total)
	}
	if len(filings) != 3 {
		t.Fatalf("expected 3 filings for 2025, got %d: %+v", len(filings), filings)
	}
	if filings[0].FilingDate != "2026-01-20" {
		t.Errorf("expected the newest filing first, got %s", filings[0].FilingDate)
	}

	// The 8-K must never be selected.
	for _, filing := range filings {
		if filing.AccessionNumber == "acc-2" {
			t.Error("selected a non Form 4 filing")
		}
	}
}

// A year the index cannot reach must be reported, not silently returned as "no filings".
func TestSelectForm4FilingsOutOfRange(t *testing.T) {
	submissions := submissionsWith([]string{"4", "4"}, []string{"2026-01-20", "2018-05-07"})

	if _, _, err := SelectForm4Filings(submissions, 1999); err == nil {
		t.Fatal("expected an error for a year older than the index")
	}

	// A year inside the index with no filings is simply empty, not an error.
	filings, _, err := SelectForm4Filings(submissions, 2021)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filings) != 0 {
		t.Fatalf("expected no filings, got %d", len(filings))
	}
}

// A busy year must report that it was cut short. Silently returning a partial year would
// read as a quiet year, understating insider activity.
func TestSelectForm4FilingsReportsTruncation(t *testing.T) {
	count := maxForm4FilingsPerRequest + 25
	forms := make([]string, count)
	dates := make([]string, count)
	for i := range forms {
		forms[i] = "4"
		dates[i] = "2025-06-15"
	}

	filings, total, err := SelectForm4Filings(submissionsWith(forms, dates), 2025)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filings) != maxForm4FilingsPerRequest {
		t.Fatalf("expected the request to be capped at %d, got %d", maxForm4FilingsPerRequest, len(filings))
	}
	// total counts the whole window, which is what makes the truncation visible.
	if total != count {
		t.Fatalf("expected the window total to be %d, got %d", count, total)
	}
	if len(filings) >= total {
		t.Fatal("expected the caller to be able to detect truncation")
	}
}

// The filing window is wider than the year, so the overhang has to be trimmed by
// transaction date.
func TestFilterByYear(t *testing.T) {
	transactions := []domain.InsiderTransaction{
		{TransactionDate: "2025-12-30"},
		{TransactionDate: "2026-01-05"},
		{TransactionDate: "2025-01-02"},
	}

	filtered := FilterByYear(transactions, 2025)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 transactions in 2025, got %d: %+v", len(filtered), filtered)
	}
}

func submissionsWith(forms []string, filingDates []string) Submissions {
	var submissions Submissions
	recent := &submissions.Filings.Recent

	recent.Form = forms
	recent.FilingDate = filingDates
	for i := range forms {
		recent.AccessionNumber = append(recent.AccessionNumber, "acc-"+string(rune('0'+i%10)))
		recent.PrimaryDocument = append(recent.PrimaryDocument, "xslF345X05/doc.xml")
	}

	return submissions
}
