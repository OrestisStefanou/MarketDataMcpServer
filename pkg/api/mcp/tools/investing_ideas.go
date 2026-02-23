package tools

import (
	"context"
	"fmt"
	"market_data_mcp_server/pkg/services"

	"github.com/mark3labs/mcp-go/mcp"
)

type GetInvestingIdeasRequest struct {
	// No input parameters needed
}

type InvestingIdea struct {
	IdeaID string `json:"idea_id"`
	Title  string `json:"title"`
}

type GetInvestingIdeasResponse struct {
	InvestingIdeas []InvestingIdea `json:"investing_ideas"`
}

type GetInvestingIdeaStocksRequest struct {
	IdeaID string `json:"idea_id" jsonschema_description:"Can be obtained from the getInvestingIdeas tool"`
}

type GetInvestingIdeaStocksResponse struct {
	Stocks []string `json:"stocks" jsonschema_description:"List of stocks(company names)"`
}

type GetInvestingIdeasTool struct {
	investingIdeasService services.InvestingIdeasService
}

func NewGetInvestingIdeasTool(investingIdeasService services.InvestingIdeasService) (*GetInvestingIdeasTool, error) {
	if investingIdeasService == nil {
		return nil, fmt.Errorf("investing ideas service is nil")
	}

	return &GetInvestingIdeasTool{
		investingIdeasService: investingIdeasService,
	}, nil
}

func (t *GetInvestingIdeasTool) HandleGetInvestingIdeas(ctx context.Context, req mcp.CallToolRequest, args GetInvestingIdeasRequest) (GetInvestingIdeasResponse, error) {
	investingIdeas, err := t.investingIdeasService.GetInvestingIdeas()
	if err != nil {
		return GetInvestingIdeasResponse{}, err
	}

	response := GetInvestingIdeasResponse{InvestingIdeas: make([]InvestingIdea, 0, len(investingIdeas))}

	for _, investingIdea := range investingIdeas {
		response.InvestingIdeas = append(response.InvestingIdeas, InvestingIdea{
			IdeaID: investingIdea.ID,
			Title:  investingIdea.Title,
		})
	}

	return response, nil
}

func (t *GetInvestingIdeasTool) GetTool() mcp.Tool {
	return mcp.NewTool("getInvestingIdeas",
		mcp.WithDescription("Get all investing ideas/themes (e.g. AI, Clean Energy, etc.)"),
		mcp.WithInputSchema[GetInvestingIdeasRequest](),
		mcp.WithOutputSchema[GetInvestingIdeasResponse](),
	)
}

type GetInvestingIdeaStocksTool struct {
	investingIdeasService services.InvestingIdeasService
}

func NewGetInvestingIdeaStocksTool(investingIdeasService services.InvestingIdeasService) (*GetInvestingIdeaStocksTool, error) {
	if investingIdeasService == nil {
		return nil, fmt.Errorf("investing ideas service is nil")
	}

	return &GetInvestingIdeaStocksTool{
		investingIdeasService: investingIdeasService,
	}, nil
}

func (t *GetInvestingIdeaStocksTool) HandleGetInvestingIdeaStocks(ctx context.Context, req mcp.CallToolRequest, args GetInvestingIdeaStocksRequest) (GetInvestingIdeaStocksResponse, error) {
	if args.IdeaID == "" {
		return GetInvestingIdeaStocksResponse{}, fmt.Errorf("idea_id is required")
	}

	investingIdeaStocks, err := t.investingIdeasService.GetInvestingIdeaStocks(args.IdeaID)
	if err != nil {
		return GetInvestingIdeaStocksResponse{}, err
	}

	response := GetInvestingIdeaStocksResponse{Stocks: investingIdeaStocks}

	return response, nil
}

func (t *GetInvestingIdeaStocksTool) GetTool() mcp.Tool {
	return mcp.NewTool("getInvestingIdeaStocks",
		mcp.WithDescription("Returns the stocks(company name) for the given investing idea/theme id"),
		mcp.WithInputSchema[GetInvestingIdeaStocksRequest](),
		mcp.WithOutputSchema[GetInvestingIdeaStocksResponse](),
	)
}
