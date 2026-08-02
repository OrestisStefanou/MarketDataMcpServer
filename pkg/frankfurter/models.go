package frankfurter

// latestResponse is the shape of GET /v1/latest.
//
// Rates only contains the symbols Frankfurter actually supports. An unsupported symbol is
// omitted silently rather than reported as an error, so every read of this map must check
// for presence.
type latestResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}
