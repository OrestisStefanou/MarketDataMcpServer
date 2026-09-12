package marketDataScraper

import (
	"math"
	"reflect"
	"testing"
)

// These tests hit stockanalysis.com. They guard against the upstream payload
// renaming or dropping fields, which unmarshals into silent zero values rather
// than an error. Run with -short to skip.

func requireNetwork(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping network test")
	}
}

// populatedFields reports which struct fields are non-zero in at least one record.
func populatedFields(records interface{}) (populated int, total int, allZero []string) {
	rv := reflect.ValueOf(records)
	if rv.Len() == 0 {
		return 0, 0, nil
	}
	rt := rv.Index(0).Type()
	seen := make([]bool, rt.NumField())
	for i := 0; i < rv.Len(); i++ {
		for f := 0; f < rt.NumField(); f++ {
			if !rv.Index(i).Field(f).IsZero() {
				seen[f] = true
			}
		}
	}
	for f := 0; f < rt.NumField(); f++ {
		if seen[f] {
			populated++
		} else {
			allZero = append(allZero, rt.Field(f).Name)
		}
	}
	return populated, rt.NumField(), allZero
}

// A generic large-cap should populate most non-sector-specific fields. A large
// drop here means upstream renamed keys and the json tags need updating.
func TestStatementsArePopulated(t *testing.T) {
	requireNetwork(t)

	cases := []struct {
		name string
		min  int
		fn   func() (interface{}, error)
	}{
		{"income", 30, func() (interface{}, error) { return scrapeIncomeStatements("aapl") }},
		{"balance", 35, func() (interface{}, error) { return scrapeBalanceSheets("aapl") }},
		{"cashflow", 30, func() (interface{}, error) { return scrapeCashFlows("aapl") }},
		{"ratios", 35, func() (interface{}, error) { return scrapeFinancialRatios("aapl") }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			recs, err := c.fn()
			if err != nil {
				t.Fatalf("scrape failed: %v", err)
			}
			got, total, allZero := populatedFields(recs)
			if total == 0 {
				t.Fatal("no records returned")
			}
			if got < c.min {
				t.Errorf("only %d/%d fields populated (want >= %d); all-zero: %v", got, total, c.min, allZero)
			}
		})
	}
}

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) <= 0.001*math.Max(1, math.Abs(b))
}

// Cross-field identities catch fields that are populated but mapped to the
// wrong upstream key.
func TestAccountingIdentities(t *testing.T) {
	requireNetwork(t)

	for _, sym := range []string{"aapl", "ftnt", "xom"} {
		t.Run(sym, func(t *testing.T) {
			bs, err := scrapeBalanceSheets(sym)
			if err != nil {
				t.Fatal(err)
			}
			is, err := scrapeIncomeStatements(sym)
			if err != nil {
				t.Fatal(err)
			}
			cf, err := scrapeCashFlows(sym)
			if err != nil {
				t.Fatal(err)
			}
			b, i, c := bs[0], is[0], cf[0]

			checks := []struct {
				name string
				got  float64
				want float64
			}{
				{"assets == liabilities + equity", b.Assets, b.Liabilitiesequity},
				{"liabilities + equity == assets", b.Liabilities + b.Equity, b.Assets},
				{"gross profit == revenue - cost of revenue", i.Gp, i.Revenue - i.Cor},
				{"operating income == gross profit - opex", i.Opinc, i.Gp - i.Opex},
				{"gross margin == gross profit / revenue", i.GrossMargin, i.Gp / i.Revenue},
				{"profit margin == to-company net income / revenue", i.ProfitMargin, i.NetincCompany / i.Revenue},
				// Upstream reports netinc and netinccmn both already net of minority
				// interest, and carries the minority interest line separately as a
				// disclosure rather than as a reconciling item. Verified on xom, whose
				// netinc and netinccmn are both 14,525m against a minority interest of
				// -356m. Subtracting it again would double-count.
				{"net income == to-company net income", i.Netinc, i.NetincCompany},
				{"free cash flow == operating cash flow + capex", c.Fcf, c.Ncfo + c.Capex},
			}
			for _, ch := range checks {
				if !approxEqual(ch.got, ch.want) {
					t.Errorf("%s: got %.4f, want %.4f", ch.name, ch.got, ch.want)
				}
			}

			// Upstream mirrors one of the two income statement net income lines
			// here, and which one varies by ticker.
			if !approxEqual(c.NetIncomeCF, i.Netinc) && !approxEqual(c.NetIncomeCF, i.NetincCompany) {
				t.Errorf("cash flow net income %.0f matches neither netinc %.0f nor netinc_company %.0f",
					c.NetIncomeCF, i.Netinc, i.NetincCompany)
			}
		})
	}
}

// Margins and yields are served as fractions (0.49 = 49%), which the schema
// descriptions promise. A vendor switch to percentages would silently inflate
// every margin by 100x.
func TestMarginsAreFractions(t *testing.T) {
	requireNetwork(t)

	is, err := scrapeIncomeStatements("aapl")
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range []struct {
		name string
		val  float64
	}{
		{"gross margin", is[0].GrossMargin},
		{"operating margin", is[0].OperatingMargin},
		{"profit margin", is[0].ProfitMargin},
		{"effective tax rate", is[0].Taxrate},
	} {
		if m.val <= 0 || m.val >= 1.5 {
			t.Errorf("%s = %v, expected a fraction", m.name, m.val)
		}
	}
}

// An absent financialData index must be an error, not a panic.
func TestMissingFinancialDataIsAnError(t *testing.T) {
	requireNetwork(t)

	// /financials/ serves the "overview" statement, which carries no financialData.
	_, err := scrapeFinancialStatementData("https://stockanalysis.com/stocks/ftnt/financials/__data.json?p=quarterly")
	if err == nil {
		t.Fatal("expected an error for a payload with no financialData")
	}
}
