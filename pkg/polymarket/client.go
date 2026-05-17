package polymarket

import (
	"encoding/json"
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"net/http"
	"net/url"
	"strconv"
)

const (
	gammaBaseURL = "https://gamma-api.polymarket.com"
	sourceName   = "polymarket"
)

type PolymarketClient struct{}

func NewPolymarketClient() (*PolymarketClient, error) {
	return &PolymarketClient{}, nil
}

func (c *PolymarketClient) SearchEvents(query string, limit int) ([]domain.PredictionMarketEvent, error) {
	if limit <= 0 {
		limit = 5
	}

	requestUrl, err := url.Parse(gammaBaseURL + "/public-search")
	if err != nil {
		return nil, fmt.Errorf("failed to parse polymarket base url: %w", err)
	}

	q := requestUrl.Query()
	q.Set("q", query)
	q.Set("limit_per_type", strconv.Itoa(limit))
	q.Set("events_status", "active")
	requestUrl.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create polymarket request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call polymarket: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("polymarket returned non-200 status: %s", resp.Status)
	}

	var searchResp polymarketSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode polymarket response: %w", err)
	}

	events := make([]domain.PredictionMarketEvent, 0, len(searchResp.Events))
	for i, ev := range searchResp.Events {
		if i >= limit {
			break
		}
		events = append(events, mapEvent(ev))
	}

	return events, nil
}

func mapEvent(ev polymarketEvent) domain.PredictionMarketEvent {
	markets := make([]domain.PredictionMarket, 0, len(ev.Markets))
	for _, m := range ev.Markets {
		markets = append(markets, mapMarket(m))
	}

	return domain.PredictionMarketEvent{
		Id:          ev.Id,
		Slug:        ev.Slug,
		Title:       ev.Title,
		Description: ev.Description,
		EndDate:     ev.EndDate,
		Volume:      ev.Volume,
		Active:      ev.Active,
		Closed:      ev.Closed,
		Source:      sourceName,
		Markets:     markets,
	}
}

func mapMarket(m polymarketMarket) domain.PredictionMarket {
	names := parseStringArray(m.Outcomes)
	prices := parseStringArray(m.OutcomePrices)

	outcomes := make([]domain.PredictionMarketOutcome, 0, len(names))
	for i, name := range names {
		probability := 0.0
		if i < len(prices) {
			if p, err := strconv.ParseFloat(prices[i], 64); err == nil {
				probability = p
			}
		}
		outcomes = append(outcomes, domain.PredictionMarketOutcome{
			Name:        name,
			Probability: probability,
		})
	}

	volume := 0.0
	if m.Volume != "" {
		if v, err := strconv.ParseFloat(m.Volume, 64); err == nil {
			volume = v
		}
	}

	return domain.PredictionMarket{
		Question: m.Question,
		Slug:     m.Slug,
		Volume:   volume,
		EndDate:  m.EndDate,
		Active:   m.Active,
		Closed:   m.Closed,
		Outcomes: outcomes,
	}
}

// Polymarket's Gamma API returns outcomes / outcomePrices as JSON-encoded strings
// (e.g. "[\"Yes\",\"No\"]"). Decode them into Go slices.
func parseStringArray(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return nil
	}
	return out
}
