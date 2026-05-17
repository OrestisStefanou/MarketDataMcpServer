package tools

import (
	"context"
	"market_data_mcp_server/pkg/domain"

	"github.com/mark3labs/mcp-go/mcp"
)

type PredictionMarketsService interface {
	SearchEvents(query string, limit int) ([]domain.PredictionMarketEvent, error)
}

type GetPolymarketEventOddsRequest struct {
	EventQuery string `json:"event_query" jsonschema_description:"Free-text event name to search for on Polymarket (e.g. 'US presidential election 2028', 'Bitcoin price end of year')" jsonschema_required:"true"`
	Limit      int    `json:"limit" jsonschema_description:"Maximum number of matching events to return (default 5, max 20)"`
}

type PredictionMarketOutcomeSchema struct {
	Name        string  `json:"name" jsonschema_description:"Outcome name (e.g. 'Yes', 'No', or a candidate / option name)"`
	Probability float64 `json:"probability" jsonschema_description:"Implied probability of the outcome, between 0 and 1"`
}

type PredictionMarketSchema struct {
	Question string                          `json:"question" jsonschema_description:"The market question"`
	Slug     string                          `json:"slug" jsonschema_description:"Polymarket slug of the market"`
	Volume   float64                         `json:"volume" jsonschema_description:"All-time market volume in USD"`
	EndDate  string                          `json:"end_date" jsonschema_description:"Resolution / end date of the market in ISO 8601 format"`
	Active   bool                            `json:"active" jsonschema_description:"Whether the market is currently active"`
	Closed   bool                            `json:"closed" jsonschema_description:"Whether the market is closed (no longer accepting trades)"`
	Outcomes []PredictionMarketOutcomeSchema `json:"outcomes" jsonschema_description:"Possible outcomes with their implied probabilities"`
}

type PredictionMarketEventSchema struct {
	Id          string                   `json:"id" jsonschema_description:"Event ID on the source platform"`
	Slug        string                   `json:"slug" jsonschema_description:"Event slug on the source platform"`
	Title       string                   `json:"title" jsonschema_description:"Event title"`
	Description string                   `json:"description" jsonschema_description:"Event description"`
	EndDate     string                   `json:"end_date" jsonschema_description:"Event end date in ISO 8601 format"`
	Volume      float64                  `json:"volume" jsonschema_description:"All-time event volume in USD"`
	Active      bool                     `json:"active" jsonschema_description:"Whether the event is currently active"`
	Closed      bool                     `json:"closed" jsonschema_description:"Whether the event is closed"`
	Source      string                   `json:"source" jsonschema_description:"Name of the prediction market source (e.g. 'polymarket')"`
	Markets     []PredictionMarketSchema `json:"markets" jsonschema_description:"Markets that belong to this event, each with its outcomes and odds"`
}

type GetPolymarketEventOddsResponse struct {
	Events []PredictionMarketEventSchema `json:"events" jsonschema_description:"Ranked list of matching events with their markets and odds"`
}

type GetPolymarketEventOddsTool struct {
	predictionMarketsService PredictionMarketsService
}

func NewGetPolymarketEventOddsTool(predictionMarketsService PredictionMarketsService) (*GetPolymarketEventOddsTool, error) {
	return &GetPolymarketEventOddsTool{
		predictionMarketsService: predictionMarketsService,
	}, nil
}

func (t *GetPolymarketEventOddsTool) HandleGetPolymarketEventOdds(ctx context.Context, req mcp.CallToolRequest, args GetPolymarketEventOddsRequest) (GetPolymarketEventOddsResponse, error) {
	limit := args.Limit
	if limit <= 0 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	events, err := t.predictionMarketsService.SearchEvents(args.EventQuery, limit)
	if err != nil {
		return GetPolymarketEventOddsResponse{}, err
	}

	response := GetPolymarketEventOddsResponse{
		Events: make([]PredictionMarketEventSchema, 0, len(events)),
	}

	for _, event := range events {
		markets := make([]PredictionMarketSchema, 0, len(event.Markets))
		for _, m := range event.Markets {
			// Skip markets without live odds: closed markets are already resolved
			// (probabilities are 0/1, not real odds), and inactive ones are
			// placeholder slots Polymarket creates ahead of time.
			if m.Closed || !m.Active {
				continue
			}
			outcomes := make([]PredictionMarketOutcomeSchema, 0, len(m.Outcomes))
			for _, o := range m.Outcomes {
				outcomes = append(outcomes, PredictionMarketOutcomeSchema{
					Name:        o.Name,
					Probability: o.Probability,
				})
			}
			markets = append(markets, PredictionMarketSchema{
				Question: m.Question,
				Slug:     m.Slug,
				Volume:   m.Volume,
				EndDate:  m.EndDate,
				Active:   m.Active,
				Closed:   m.Closed,
				Outcomes: outcomes,
			})
		}

		response.Events = append(response.Events, PredictionMarketEventSchema{
			Id:          event.Id,
			Slug:        event.Slug,
			Title:       event.Title,
			Description: event.Description,
			EndDate:     event.EndDate,
			Volume:      event.Volume,
			Active:      event.Active,
			Closed:      event.Closed,
			Source:      event.Source,
			Markets:     markets,
		})
	}

	return response, nil
}

func (t *GetPolymarketEventOddsTool) GetTool() mcp.Tool {
	return mcp.NewTool("getPolymarketEventOdds",
		mcp.WithDescription("Search Polymarket for prediction-market events matching a free-text query and return the matching events along with their markets and implied outcome probabilities (odds)."),
		mcp.WithInputSchema[GetPolymarketEventOddsRequest](),
		mcp.WithOutputSchema[GetPolymarketEventOddsResponse](),
	)
}
