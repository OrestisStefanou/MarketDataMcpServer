package domain

type PredictionMarketOutcome struct {
	Name        string  `json:"name"`
	Probability float64 `json:"probability"`
}

type PredictionMarket struct {
	Question string                    `json:"question"`
	Slug     string                    `json:"slug"`
	Volume   float64                   `json:"volume"`
	EndDate  string                    `json:"end_date"`
	Active   bool                      `json:"active"`
	Closed   bool                      `json:"closed"`
	Outcomes []PredictionMarketOutcome `json:"outcomes"`
}

type PredictionMarketEvent struct {
	Id          string             `json:"id"`
	Slug        string             `json:"slug"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	EndDate     string             `json:"end_date"`
	Volume      float64            `json:"volume"`
	Active      bool               `json:"active"`
	Closed      bool               `json:"closed"`
	Source      string             `json:"source"`
	Markets     []PredictionMarket `json:"markets"`
}
