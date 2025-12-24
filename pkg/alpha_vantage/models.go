package alphavantage

type EconomicIndicator string

const (
	RealGDP          EconomicIndicator = "REAL_GDP"
	TreasuryYield    EconomicIndicator = "TREASURY_YIELD"
	FederalFundsRate EconomicIndicator = "FEDERAL_FUNDS_RATE"
	Inflation        EconomicIndicator = "INFLATION"
	UnemploymentRate EconomicIndicator = "UNEMPLOYMENT"
)

type TimeSeriesEntry struct {
	Date  string `json:"date"`
	Value string `json:"value"`
}

type EconomicIndicatorTimeSeriesResponse struct {
	Name     string            `json:"name"`
	Interval string            `json:"interval"`
	Unit     string            `json:"unit"`
	Data     []TimeSeriesEntry `json:"data"`
}

type Commodity string

const (
	WTI         Commodity = "WTI"
	NATURAL_GAS Commodity = "NATURAL_GAS"
	COPPER      Commodity = "COPPER"
	ALUMINIUM   Commodity = "ALUMINUM"
	WHEAT       Commodity = "WHEAT"
	CORN        Commodity = "CORN"
	SUGAR       Commodity = "SUGAR"
	COFFEE      Commodity = "COFFEE"
)

type CommodityTimeSeriesResponse struct {
	Name     string            `json:"name"`
	Interval string            `json:"interval"`
	Unit     string            `json:"unit"`
	Data     []TimeSeriesEntry `json:"data"`
}

type GetNewsResponse struct {
	Items                    string         `json:"items"`
	SentimentScoreDefinition string         `json:"sentiment_score_definition"`
	RelevanceScoreDefinition string         `json:"relevance_score_definition"`
	Feed                     []NewsFeedItem `json:"feed"`
}

type NewsFeedItem struct {
	Title                 string            `json:"title"`
	URL                   string            `json:"url"`
	TimePublished         string            `json:"time_published"`
	Authors               []string          `json:"authors"`
	Summary               string            `json:"summary"`
	BannerImage           string            `json:"banner_image"`
	Source                string            `json:"source"`
	CategoryWithinSource  string            `json:"category_within_source"`
	SourceDomain          string            `json:"source_domain"`
	Topics                []Topic           `json:"topics"`
	OverallSentimentScore float64           `json:"overall_sentiment_score"`
	OverallSentimentLabel string            `json:"overall_sentiment_label"`
	TickerSentiment       []TickerSentiment `json:"ticker_sentiment"`
}

type Topic struct {
	Topic          string `json:"topic"`
	RelevanceScore string `json:"relevance_score"`
}

type TickerSentiment struct {
	Ticker               string `json:"ticker"`
	RelevanceScore       string `json:"relevance_score"`
	TickerSentimentScore string `json:"ticker_sentiment_score"`
	TickerSentimentLabel string `json:"ticker_sentiment_label"`
}
