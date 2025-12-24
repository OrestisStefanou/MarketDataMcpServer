package tools

import (
	"context"
	"fmt"
	"market_data_mcp_server/pkg/domain"

	"github.com/mark3labs/mcp-go/mcp"
)

type SuperInvestorsService interface {
	GetSuperInvestors() ([]domain.SuperInvestor, error)
	GetSuperInvestorPortfolio(superInvestorName string) (domain.SuperInvestorPortfolio, error)
}

type SuperInvestorSchema struct {
	Name string `json:"super_investor_name" jsonschema_description:"The name of the super investor(Portfolio Manager - Firm)"`
}

type GetSuperInvestorsRequest struct {
	// No input parameters required
}

type GetSuperInvestorsResponse struct {
	SuperInvestors []SuperInvestorSchema `json:"super_investors" jsonschema_description:"A list with the names of the super investors(Portfolio Managers - Firms) names"`
}

type GetSuperInvestorsTool struct {
	superInvestorsService SuperInvestorsService
}

func NewGetSuperInvestorsTool(superInvestorsService SuperInvestorsService) (*GetSuperInvestorsTool, error) {
	return &GetSuperInvestorsTool{
		superInvestorsService: superInvestorsService,
	}, nil
}

func (t *GetSuperInvestorsTool) HandleGetSuperInvestors(ctx context.Context, req mcp.CallToolRequest, args GetSuperInvestorsRequest) (GetSuperInvestorsResponse, error) {
	superInvestors, err := t.superInvestorsService.GetSuperInvestors()
	if err != nil {
		return GetSuperInvestorsResponse{}, err
	}

	response := GetSuperInvestorsResponse{
		SuperInvestors: make([]SuperInvestorSchema, 0, len(superInvestors)),
	}

	for _, si := range superInvestors {
		response.SuperInvestors = append(
			response.SuperInvestors,
			SuperInvestorSchema{Name: si.Name},
		)
	}

	return response, nil
}

func (t *GetSuperInvestorsTool) GetTool() mcp.Tool {
	return mcp.NewTool("getSuperInvestors",
		mcp.WithDescription("Get a list of all super investors (Portfolio Managers - Firms)"),
		mcp.WithInputSchema[GetSuperInvestorsRequest](),
		mcp.WithOutputSchema[GetSuperInvestorsResponse](),
	)
}

type SuperInvestorPortfolioHoldingSchema struct {
	Stock          string `json:"stock" jsonschema_description:"Stock name"`
	PortfolioPct   string `json:"portfolio_pct" jsonschema_description:"Percentage of portfolio"`
	RecentActivity string `json:"recent_activity" jsonschema_description:"Recent activity"`
	Shares         string `json:"shares" jsonschema_description:"Number of shares held"`
	Value          string `json:"value" jsonschema_description:"Value of holding"`
}

type SuperInvestorPortfolioSectorAnalysisSchema struct {
	Sector       string `json:"sector" jsonschema_description:"Sector name"`
	PortfolioPct string `json:"portfolio_pct" jsonschema_description:"Percentage of portfolio"`
}

type GetSuperInvestorPortfolioRequest struct {
	SuperInvestorName string `json:"super_investor_name" jsonschema_description:"The name of the super investor (Portfolio Manager - Firm) to get the portfolio for"`
}

type GetSuperInvestorPortfolioResponse struct {
	Holdings       []SuperInvestorPortfolioHoldingSchema        `json:"holdings" jsonschema_description:"List of portfolio holdings"`
	SectorAnalysis []SuperInvestorPortfolioSectorAnalysisSchema `json:"sector_analysis" jsonschema_description:"Sector analysis of the portfolio"`
}

type GetSuperInvestorPortfolioTool struct {
	superInvestorsService SuperInvestorsService
}

func NewGetSuperInvestorPortfolioTool(superInvestorsService SuperInvestorsService) (*GetSuperInvestorPortfolioTool, error) {
	return &GetSuperInvestorPortfolioTool{
		superInvestorsService: superInvestorsService,
	}, nil
}

func (t *GetSuperInvestorPortfolioTool) HandleGetSuperInvestorPortfolio(ctx context.Context, req mcp.CallToolRequest, args GetSuperInvestorPortfolioRequest) (GetSuperInvestorPortfolioResponse, error) {
	if args.SuperInvestorName == "" {
		return GetSuperInvestorPortfolioResponse{}, fmt.Errorf("super_investor_name is required")
	}

	portfolio, err := t.superInvestorsService.GetSuperInvestorPortfolio(args.SuperInvestorName)
	if err != nil {
		return GetSuperInvestorPortfolioResponse{}, err
	}

	holdings := make([]SuperInvestorPortfolioHoldingSchema, 0, len(portfolio.Holdings))
	for _, holding := range portfolio.Holdings {
		holdings = append(holdings, SuperInvestorPortfolioHoldingSchema{
			Stock:          holding.Stock,
			PortfolioPct:   holding.PortfolioPct,
			RecentActivity: holding.RecentActivity,
			Shares:         holding.Shares,
			Value:          holding.Value,
		})
	}

	sectorAnalysis := make([]SuperInvestorPortfolioSectorAnalysisSchema, 0, len(portfolio.SectorAnalysis))
	for _, sector := range portfolio.SectorAnalysis {
		sectorAnalysis = append(sectorAnalysis, SuperInvestorPortfolioSectorAnalysisSchema{
			Sector:       sector.Sector,
			PortfolioPct: sector.PortfolioPct,
		})
	}

	response := GetSuperInvestorPortfolioResponse{
		Holdings:       holdings,
		SectorAnalysis: sectorAnalysis,
	}

	return response, nil
}

func (t *GetSuperInvestorPortfolioTool) GetTool() mcp.Tool {
	return mcp.NewTool("getSuperInvestorPortfolio",
		mcp.WithDescription("Get the portfolio of a super investor (Portfolio Manager - Firm) including holdings and sector analysis"),
		mcp.WithInputSchema[GetSuperInvestorPortfolioRequest](),
		mcp.WithOutputSchema[GetSuperInvestorPortfolioResponse](),
	)
}
