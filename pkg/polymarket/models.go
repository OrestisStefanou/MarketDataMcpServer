package polymarket

type polymarketMarket struct {
	Id            string  `json:"id"`
	Question      string  `json:"question"`
	Slug          string  `json:"slug"`
	Volume        string  `json:"volume"`
	EndDate       string  `json:"endDate"`
	Active        bool    `json:"active"`
	Closed        bool    `json:"closed"`
	Outcomes      string  `json:"outcomes"`
	OutcomePrices string  `json:"outcomePrices"`
}

type polymarketEvent struct {
	Id          string             `json:"id"`
	Slug        string             `json:"slug"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	EndDate     string             `json:"endDate"`
	Volume      float64            `json:"volume"`
	Active      bool               `json:"active"`
	Closed      bool               `json:"closed"`
	Markets     []polymarketMarket `json:"markets"`
}

type polymarketSearchResponse struct {
	Events []polymarketEvent `json:"events"`
}
