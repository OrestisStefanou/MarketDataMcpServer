package tools

import (
	"context"
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/services"

	"github.com/mark3labs/mcp-go/mcp"
)

type SearchEtfRequest struct {
	SearchString string `json:"search_string,omitempty" jsonschema_description:"Search string query"`
	Limit        int    `json:"limit,omitempty" jsonschema_description:"Maximum results" jsonschema:"minimum=1,default=100"`
}

type EtfSearchResultSchema struct {
	Symbol     string  `json:"symbol" jsonschema_description:"ETF symbol"`
	EtfName    string  `json:"etf_name" jsonschema_description:"ETF name"`
	AssetClass string  `json:"asset_class" jsonschema_description:"ETF asset class"`
	Aum        float32 `json:"aum" jsonschema_description:"ETF assets under management"`
}

type EtfSearchResultsResponse struct {
	SearchResults []EtfSearchResultSchema `json:"search_results" jsonschema_description:"Search results"`
}

type EtfService interface {
	GetEtfs(filters services.EtfFilterOptions) ([]domain.Etf, error)
	GetEtf(etfSymbol string) (domain.EtfOverview, error)
}

type SearchEtfTool struct {
	etfService EtfService
}

func NewSearchEtfTool(etfService EtfService) (*SearchEtfTool, error) {
	return &SearchEtfTool{
		etfService: etfService,
	}, nil
}

func (t *SearchEtfTool) HandleSearchEtfs(ctx context.Context, req mcp.CallToolRequest, args SearchEtfRequest) (EtfSearchResultsResponse, error) {
	if args.Limit == 0 {
		args.Limit = 100
	}
	tickerFilters := services.EtfFilterOptions{
		SearchString: args.SearchString,
	}

	etfs, err := t.etfService.GetEtfs(tickerFilters)
	if err != nil {
		return EtfSearchResultsResponse{}, err
	}

	response := EtfSearchResultsResponse{
		SearchResults: make([]EtfSearchResultSchema, 0, len(etfs)),
	}

	for i, e := range etfs {
		if i > args.Limit {
			break
		}
		response.SearchResults = append(
			response.SearchResults,
			EtfSearchResultSchema{Symbol: e.Symbol, EtfName: e.Name, AssetClass: e.AssetClass, Aum: e.Aum},
		)
	}

	return response, nil
}

func (t *SearchEtfTool) GetTool() mcp.Tool {
	return mcp.NewTool("etfSearch",
		mcp.WithDescription("Search for an ETF using the symbol or the ETF name"),
		mcp.WithInputSchema[SearchEtfRequest](),
		mcp.WithOutputSchema[EtfSearchResultsResponse](),
	)
}

type EtfHoldingSchema struct {
	Symbol string `json:"symbol" jsonschema_description:"Holding symbol"`
	Name   string `json:"name" jsonschema_description:"Holding name"`
	Weight string `json:"weight" jsonschema_description:"Holding weight percentage in the ETF"`
}

type GetEtfResponse struct {
	Symbol           string             `json:"symbol" jsonschema_description:"ETF symbol"`
	Description      string             `json:"description" jsonschema_description:"ETF description"`
	AssetClass       string             `json:"asset_class" jsonschema_description:"ETF asset class"`
	Category         string             `json:"category" jsonschema_description:"ETF category"`
	Aum              string             `json:"aum" jsonschema_description:"Assets under management"`
	Nav              string             `json:"nav" jsonschema_description:"Net asset value"`
	ExpenseRatio     string             `json:"expense_ratio" jsonschema_description:"Annual expense ratio"`
	PeRatio          string             `json:"pe_ratio" jsonschema_description:"Price-to-earnings ratio"`
	Dps              string             `json:"dps" jsonschema_description:"Dividends per share"`
	DividendYield    string             `json:"dividend_yield" jsonschema_description:"Dividend yield percentage"`
	PayoutRatio      string             `json:"payout_ratio" jsonschema_description:"Dividend payout ratio"`
	OneMonthReturn   float64            `json:"one_month_return" jsonschema_description:"One month return percentage"`
	OneYearReturn    float64            `json:"one_year_return" jsonschema_description:"One year return percentage"`
	YearToDateReturn float64            `json:"year_to_date_return" jsonschema_description:"Year-to-date return percentage"`
	FiveYearReturn   float64            `json:"five_year_return" jsonschema_description:"Five year return percentage"`
	TenYearReturn    float64            `json:"ten_year_return" jsonschema_description:"Ten year return percentage"`
	InceptionReturn  float64            `json:"inception_return" jsonschema_description:"Return since inception percentage"`
	Beta             string             `json:"beta" jsonschema_description:"Beta coefficient measuring volatility"`
	NumberOfHoldings int32              `json:"number_of_holdings" jsonschema_description:"Total number of holdings in the ETF"`
	Website          string             `json:"website" jsonschema_description:"ETF website URL"`
	TopHoldings      []EtfHoldingSchema `json:"top_holdings" jsonschema_description:"List of top holdings in the ETF"`
}

type GetEtfRequest struct {
	EtfSymbol string `json:"etf_symbol" jsonschema_description:"Symbol of the ETF to return"`
}

type GetEtfTool struct {
	etfService EtfService
}

func NewGetEtfTool(etfService EtfService) (*GetEtfTool, error) {
	return &GetEtfTool{
		etfService: etfService,
	}, nil
}

func (t *GetEtfTool) HandleGetEtf(ctx context.Context, req mcp.CallToolRequest, args GetEtfRequest) (GetEtfResponse, error) {
	if args.EtfSymbol == "" {
		return GetEtfResponse{}, fmt.Errorf("etf_symbol is required")
	}

	etf, err := t.etfService.GetEtf(args.EtfSymbol)
	if err != nil {
		return GetEtfResponse{}, err
	}

	topHoldings := make([]EtfHoldingSchema, 0, len(etf.TopHoldings))
	for _, holding := range etf.TopHoldings {
		topHoldings = append(topHoldings, EtfHoldingSchema{
			Symbol: holding.Symbol,
			Name:   holding.Name,
			Weight: holding.Weight,
		})
	}

	response := GetEtfResponse{
		Symbol:           etf.Symbol,
		Description:      etf.Description,
		AssetClass:       etf.AssetClass,
		Category:         etf.Category,
		Aum:              etf.Aum,
		Nav:              etf.Nav,
		ExpenseRatio:     etf.ExpenseRatio,
		PeRatio:          etf.PeRatio,
		Dps:              etf.Dps,
		DividendYield:    etf.DividendYield,
		PayoutRatio:      etf.PayoutRatio,
		OneMonthReturn:   etf.OneMonthReturn,
		OneYearReturn:    etf.OneYearReturn,
		YearToDateReturn: etf.YearToDateReturn,
		FiveYearReturn:   etf.FiveYearReturn,
		TenYearReturn:    etf.TenYearReturn,
		InceptionReturn:  etf.InceptionReturn,
		Beta:             etf.Beta,
		NumberOfHoldings: etf.NumberOfHoldings,
		Website:          etf.Website,
		TopHoldings:      topHoldings,
	}

	return response, nil
}

func (t *GetEtfTool) GetTool() mcp.Tool {
	return mcp.NewTool("getETF",
		mcp.WithDescription("Get an ETF using it's symbol"),
		mcp.WithInputSchema[GetEtfRequest](),
		mcp.WithOutputSchema[GetEtfResponse](),
	)
}
