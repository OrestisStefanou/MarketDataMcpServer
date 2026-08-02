package fred

import (
	"strings"
	"testing"
)

func TestParseSeriesCsvSkipsMissingObservations(t *testing.T) {
	// FRED writes a missing observation as an empty field, and as "." in some older series.
	body := strings.Join([]string{
		"observation_date,DGS10",
		"2026-07-01,4.21",
		"2026-07-02,",
		"2026-07-03,.",
		"2026-07-06,4.25",
	}, "\n")

	observations, err := parseSeriesCsv(strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(observations) != 2 {
		t.Fatalf("expected 2 observations, got %d: %+v", len(observations), observations)
	}
	if observations[0].Date != "2026-07-01" || observations[0].Value != "4.21" {
		t.Errorf("unexpected first observation: %+v", observations[0])
	}
	if observations[1].Date != "2026-07-06" || observations[1].Value != "4.25" {
		t.Errorf("unexpected second observation: %+v", observations[1])
	}
}

func TestParseSeriesCsvHeaderOnly(t *testing.T) {
	observations, err := parseSeriesCsv(strings.NewReader("observation_date,UNRATE\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(observations) != 0 {
		t.Fatalf("expected no observations, got %d", len(observations))
	}
}

// FRED returns ascending dates but the tools present and truncate newest first, so getting
// this backwards would silently serve the oldest observations.
func TestReverseProducesNewestFirst(t *testing.T) {
	reversed := reverse([]observation{
		{Date: "2024-01-01", Value: "1"},
		{Date: "2025-01-01", Value: "2"},
		{Date: "2026-01-01", Value: "3"},
	})

	if len(reversed) != 3 {
		t.Fatalf("expected 3 observations, got %d", len(reversed))
	}
	if reversed[0].Date != "2026-01-01" {
		t.Errorf("expected newest first, got %s", reversed[0].Date)
	}
	if reversed[2].Date != "2024-01-01" {
		t.Errorf("expected oldest last, got %s", reversed[2].Date)
	}
}

func TestToYearOverYearPercent(t *testing.T) {
	observations := []observation{
		{Date: "2025-01-01", Value: "100"},
		{Date: "2025-02-01", Value: "200"},
		{Date: "2026-01-01", Value: "103"},
		{Date: "2026-02-01", Value: "190"},
	}

	converted, err := toYearOverYearPercent(observations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The first year of an index series has no prior-year comparison and is dropped.
	if len(converted) != 2 {
		t.Fatalf("expected 2 converted observations, got %d: %+v", len(converted), converted)
	}
	if converted[0].Date != "2026-01-01" || converted[0].Value != "3.00" {
		t.Errorf("expected 2026-01-01 at 3.00, got %+v", converted[0])
	}
	if converted[1].Date != "2026-02-01" || converted[1].Value != "-5.00" {
		t.Errorf("expected 2026-02-01 at -5.00, got %+v", converted[1])
	}
}

// The prior-year value is looked up by date rather than by position, so a gap in the series
// must not shift the comparison window onto the wrong month.
func TestToYearOverYearPercentIsGapSafe(t *testing.T) {
	observations := []observation{
		{Date: "2025-01-01", Value: "100"},
		// 2025-02-01 is missing.
		{Date: "2025-03-01", Value: "300"},
		{Date: "2026-01-01", Value: "110"},
		{Date: "2026-03-01", Value: "330"},
	}

	converted, err := toYearOverYearPercent(observations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(converted) != 2 {
		t.Fatalf("expected 2 converted observations, got %d: %+v", len(converted), converted)
	}
	for _, obs := range converted {
		if obs.Value != "10.00" {
			t.Errorf("expected every comparison to be 10.00, got %+v", obs)
		}
	}
}

func TestToYearOverYearPercentWithoutPriorYear(t *testing.T) {
	_, err := toYearOverYearPercent([]observation{{Date: "2026-01-01", Value: "100"}})
	if err == nil {
		t.Fatal("expected an error when no observation has a prior-year comparison")
	}
}
