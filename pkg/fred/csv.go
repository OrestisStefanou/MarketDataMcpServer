package fred

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"market_data_mcp_server/pkg/errors"
)

const fredCsvBaseUrl = "https://fred.stlouisfed.org/graph/fredgraph.csv"

var httpClient = &http.Client{Timeout: 15 * time.Second}

// observation is a single FRED data point, in the ascending date order FRED publishes.
type observation struct {
	Date  string
	Value string
}

// fetchSeries downloads a FRED series as CSV and returns its observations in ascending
// date order. FRED marks a missing observation with an empty field (and, in some older
// series, a "."), so both are skipped rather than emitted as a blank value.
func fetchSeries(seriesId string, lookbackYears int) ([]observation, error) {
	url := fmt.Sprintf("%s?id=%s", fredCsvBaseUrl, seriesId)
	if lookbackYears > 0 {
		start := time.Now().AddDate(-lookbackYears, 0, 0)
		url = fmt.Sprintf("%s&cosd=%s", url, start.Format("2006-01-02"))
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("Call to get FRED series %s failed", seriesId),
		}
	}

	observations, err := parseSeriesCsv(resp.Body)
	if err != nil {
		return nil, errors.StreamError{
			Message: fmt.Sprintf("failed to read FRED series %s", seriesId),
			Err:     err,
		}
	}

	if len(observations) == 0 {
		return nil, fmt.Errorf("FRED series %s returned no observations", seriesId)
	}

	return observations, nil
}

// parseSeriesCsv reads FRED's two-column CSV, dropping the header and any observation FRED
// publishes as missing.
func parseSeriesCsv(body io.Reader) ([]observation, error) {
	reader := csv.NewReader(body)

	// The header is "observation_date,<SERIES_ID>".
	if _, err := reader.Read(); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	observations := make([]observation, 0, 512)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) < 2 {
			continue
		}

		// A missing observation is an empty field, and in some older series a ".".
		value := strings.TrimSpace(record[1])
		if value == "" || value == "." {
			continue
		}

		observations = append(observations, observation{
			Date:  strings.TrimSpace(record[0]),
			Value: value,
		})
	}

	return observations, nil
}

// reverse flips observations into newest-first order, which is the order the MCP tools
// present and truncate with their limit argument.
func reverse(observations []observation) []observation {
	reversed := make([]observation, len(observations))
	for i, obs := range observations {
		reversed[len(observations)-1-i] = obs
	}
	return reversed
}

// toYearOverYearPercent converts an index-level series (such as CPI) into percent change
// against the same period a year earlier. The prior-year value is looked up by date rather
// than by position, so a gap in the series cannot silently shift the comparison window.
func toYearOverYearPercent(observations []observation) ([]observation, error) {
	valuesByDate := make(map[string]float64, len(observations))
	for _, obs := range observations {
		value, err := strconv.ParseFloat(obs.Value, 64)
		if err != nil {
			continue
		}
		valuesByDate[obs.Date] = value
	}

	converted := make([]observation, 0, len(observations))
	for _, obs := range observations {
		current, ok := valuesByDate[obs.Date]
		if !ok {
			continue
		}

		date, err := time.Parse("2006-01-02", obs.Date)
		if err != nil {
			continue
		}

		prior, ok := valuesByDate[date.AddDate(-1, 0, 0).Format("2006-01-02")]
		if !ok || prior == 0 {
			continue
		}

		converted = append(converted, observation{
			Date:  obs.Date,
			Value: strconv.FormatFloat((current/prior-1)*100, 'f', 2, 64),
		})
	}

	if len(converted) == 0 {
		return nil, fmt.Errorf("not enough observations to compute a year over year change")
	}

	return converted, nil
}
