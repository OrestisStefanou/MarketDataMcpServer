package tools

import (
	"context"
	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/services"

	"github.com/mark3labs/mcp-go/mcp"
)

type SearchStocksRequest struct {
	SearchString string `json:"search_string,omitempty" jsonschema_description:"Search string query"`
	Limit        int    `json:"limit,omitempty" jsonschema_description:"Maximum results" jsonschema:"minimum=1,default=100"`
}

type StockSearchResultSchema struct {
	Symbol      string `json:"symbol" jsonschema_description:"Stock symbol"`
	CompanyName string `json:"company_name" jsonschema_description:"Company name"`
}

type StockSearchResultsResponse struct {
	SearchResults []StockSearchResultSchema `json:"search_results" jsonschema_description:"Search results"`
}

type TickerService interface {
	GetTickers(filters services.TickerFilterOptions) ([]domain.Ticker, error)
}

type StockSearchTool struct {
	tickerService TickerService
}

func NewStockSearchTool(tickerService TickerService) (*StockSearchTool, error) {
	return &StockSearchTool{
		tickerService: tickerService,
	}, nil
}

func (h *StockSearchTool) HandleSearchStocks(ctx context.Context, req mcp.CallToolRequest, args SearchStocksRequest) (StockSearchResultsResponse, error) {
	if args.Limit == 0 {
		args.Limit = 100
	}
	tickerFilters := services.TickerFilterOptions{
		Limit:        args.Limit,
		SearchString: args.SearchString,
	}

	tickers, err := h.tickerService.GetTickers(tickerFilters)
	if err != nil {
		return StockSearchResultsResponse{}, err
	}

	response := StockSearchResultsResponse{
		SearchResults: make([]StockSearchResultSchema, 0, len(tickers)),
	}

	for _, t := range tickers {
		response.SearchResults = append(
			response.SearchResults,
			StockSearchResultSchema{
				Symbol:      t.Symbol,
				CompanyName: t.CompanyName,
			},
		)
	}

	return response, nil
}

func (h *StockSearchTool) GetTool() mcp.Tool {
	return mcp.NewTool("stockSearch",
		mcp.WithDescription("Search for a stock using the symbol or the company name"),
		mcp.WithInputSchema[SearchStocksRequest](),
		mcp.WithOutputSchema[StockSearchResultsResponse](),
	)
}
