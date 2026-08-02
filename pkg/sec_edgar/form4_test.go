package secedgar

import (
	"os"
	"testing"
)

// The fixture is a real Form 4: Elon Musk's 2026-06-16 filing for Tesla, which contains both
// a derivative option exercise and a non-derivative tax withholding.
func TestParseForm4(t *testing.T) {
	raw, err := os.ReadFile("testdata/form4_tsla.xml")
	if err != nil {
		t.Fatalf("failed to read the fixture: %v", err)
	}

	transactions, err := parseForm4(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Two non-derivative rows (the option exercise and the tax withholding) plus the
	// matching derivative row for the option itself.
	if len(transactions) != 3 {
		t.Fatalf("expected 3 transactions, got %d: %+v", len(transactions), transactions)
	}

	for _, transaction := range transactions {
		if transaction.Ticker != "TSLA" {
			t.Errorf("expected ticker TSLA, got %q", transaction.Ticker)
		}
		if transaction.Executive != "Musk Elon" {
			t.Errorf("expected executive Musk Elon, got %q", transaction.Executive)
		}
		if transaction.ExecutiveTitle != "CEO" {
			t.Errorf("expected title CEO, got %q", transaction.ExecutiveTitle)
		}
		if transaction.TransactionDate != "2026-06-16" {
			t.Errorf("expected 2026-06-16, got %q", transaction.TransactionDate)
		}
	}

	acquired := transactions[0]
	if acquired.AcquisitionOrDisposal != "A" {
		t.Errorf("expected an acquisition, got %q", acquired.AcquisitionOrDisposal)
	}
	if acquired.Shares != 303960630 {
		t.Errorf("expected 303960630 shares, got %v", acquired.Shares)
	}
	if acquired.SharePrice != 23.34 {
		t.Errorf("expected a price of 23.34, got %v", acquired.SharePrice)
	}
	if acquired.SecurityType != "Common Stock" {
		t.Errorf("expected Common Stock, got %q", acquired.SecurityType)
	}

	disposed := transactions[1]
	if disposed.AcquisitionOrDisposal != "D" {
		t.Errorf("expected a disposal, got %q", disposed.AcquisitionOrDisposal)
	}
	if disposed.SharePrice != 404.66 {
		t.Errorf("expected a price of 404.66, got %v", disposed.SharePrice)
	}

	// The derivative row carries its own security title, which is what makes it
	// self-describing next to the common stock rows.
	derivative := transactions[2]
	if derivative.SecurityType != "Non-Qualified Stock Option (right to buy)" {
		t.Errorf("expected the option security title, got %q", derivative.SecurityType)
	}
	if derivative.SharePrice != 0 {
		t.Errorf("expected no price on the option exercise, got %v", derivative.SharePrice)
	}
}

func TestParseForm4RejectsGarbage(t *testing.T) {
	if _, err := parseForm4([]byte("not xml at all")); err == nil {
		t.Fatal("expected an error for malformed XML")
	}
}

func TestOwnerTitleFallsBackToRelationshipFlags(t *testing.T) {
	tests := []struct {
		name     string
		owner    reportingOwner
		expected string
	}{
		{"officer title wins", ownerWith("1", "0", "0", "Chief Financial Officer"), "Chief Financial Officer"},
		{"director flag", ownerWith("0", "1", "0", ""), "Director"},
		{"ten percent owner", ownerWith("0", "0", "1", ""), "10% Owner"},
		{"combined roles", ownerWith("1", "1", "0", ""), "Officer, Director"},
		{"no relationship at all", ownerWith("0", "0", "0", ""), ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ownerTitle(test.owner); got != test.expected {
				t.Errorf("expected %q, got %q", test.expected, got)
			}
		})
	}
}

// EDGAR writes booleans as either 1/0 or true/false.
func TestIsFlagSet(t *testing.T) {
	for _, value := range []string{"1", "true", "TRUE", " true "} {
		if !isFlagSet(value) {
			t.Errorf("expected %q to be set", value)
		}
	}
	for _, value := range []string{"0", "false", "", "  "} {
		if isFlagSet(value) {
			t.Errorf("expected %q to be unset", value)
		}
	}
}

// Form 4 legitimately omits a price for gifts and inherited shares.
func TestParseFloatDefaultsToZero(t *testing.T) {
	if got := parseFloat(""); got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
	if got := parseFloat("12.5"); got != 12.5 {
		t.Errorf("expected 12.5, got %v", got)
	}
}

func ownerWith(isOfficer, isDirector, isTenPercentOwner, officerTitle string) reportingOwner {
	var owner reportingOwner
	owner.Relationship.IsOfficer = isOfficer
	owner.Relationship.IsDirector = isDirector
	owner.Relationship.IsTenPercentOwner = isTenPercentOwner
	owner.Relationship.OfficerTitle = officerTitle
	return owner
}
