package tools

import (
	"context"
	"fmt"
	"market_data_mcp_server/pkg/domain"

	"github.com/mark3labs/mcp-go/mcp"
)

type UserContextService interface {
	GetUserContext(userID string) (domain.UserContext, error)
	CreateUserContext(domain.UserContext) error
	UpdateUserContext(domain.UserContext) error
}

type GetUserContextRequest struct {
	UserID string `json:"user_id" jsonschema_description:"The id of the user to get the context for"`
}

type UserPortfolioHoldingSchema struct {
	AssetClass          string  `json:"asset_class" jsonschema_description:"Asset class of the holding. Valid values are: stock, etf, crypto, mutual_fund, bond, cash, real_estate, private_equity, commodities"`
	Symbol              string  `json:"symbol" jsonschema_description:"Symbol of the holding"`
	Name                string  `json:"name" jsonschema_description:"Name of the holding"`
	Quantity            float64 `json:"quantity" jsonschema_description:"Quantity of the holding(zero value means not known/given)"`
	PortfolioPercentage float64 `json:"portfolio_percentage" jsonschema_description:"Portfolio percentage of the holding(zero value means not known/given, 20 means 20%)"`
}

func (u UserPortfolioHoldingSchema) Validate() error {
	if u.AssetClass == "" {
		return fmt.Errorf("asset_class is required")
	}
	if u.Symbol == "" && u.Name == "" {
		return fmt.Errorf("you must define either symbol or name for all portfolio holdings")
	}

	if u.AssetClass != "stock" && u.AssetClass != "etf" && u.AssetClass != "crypto" && u.AssetClass != "mutual_fund" && u.AssetClass != "bond" && u.AssetClass != "cash" && u.AssetClass != "real_estate" && u.AssetClass != "private_equity" && u.AssetClass != "commodities" {
		return fmt.Errorf("asset_class valid values are: stock, etf, crypto, mutual_fund, bond, cash, real_estate, private_equity, commodities")
	}

	if u.Quantity < 0 {
		return fmt.Errorf("quantity must be a non-negative number")
	}

	if u.PortfolioPercentage < 0 || u.PortfolioPercentage > 100 {
		return fmt.Errorf("portfolio_percentage must be between 0 and 100")
	}

	return nil
}

type UserContextResponse struct {
	UserID        string                       `json:"user_id"`
	UserProfile   map[string]any               `json:"user_profile" jsonschema_description:"General information about the user"`
	UserPortfolio []UserPortfolioHoldingSchema `json:"user_portfolio"`
}

type GetUserContextTool struct {
	userContextService UserContextService
}

func NewGetUserContextTool(userContextService UserContextService) (*GetUserContextTool, error) {
	return &GetUserContextTool{
		userContextService: userContextService,
	}, nil
}

func (t *GetUserContextTool) HandleGetUserContext(ctx context.Context, req mcp.CallToolRequest, args GetUserContextRequest) (UserContextResponse, error) {
	if args.UserID == "" {
		return UserContextResponse{}, fmt.Errorf("user_id is required")
	}

	userContext, err := t.userContextService.GetUserContext(args.UserID)
	if err != nil {
		return UserContextResponse{}, err
	}

	portfolio := make([]UserPortfolioHoldingSchema, 0, len(userContext.UserPortfolio))
	for _, holding := range userContext.UserPortfolio {
		portfolio = append(portfolio, UserPortfolioHoldingSchema{
			AssetClass:          string(holding.AssetClass),
			Symbol:              holding.Symbol,
			Name:                holding.Name,
			Quantity:            holding.Quantity,
			PortfolioPercentage: holding.PortfolioPercentage,
		})
	}

	response := UserContextResponse{
		UserID:        userContext.UserID,
		UserProfile:   userContext.UserProfile,
		UserPortfolio: portfolio,
	}

	return response, nil
}

func (t *GetUserContextTool) GetTool() mcp.Tool {
	return mcp.NewTool("getUserContext",
		mcp.WithDescription("Get the user context including user profile and portfolio holdings"),
		mcp.WithInputSchema[GetUserContextRequest](),
		mcp.WithOutputSchema[UserContextResponse](),
	)
}

type UpdateUserContextRequest struct {
	UserID        string                       `json:"user_id" jsonschema_description:"The id of the user to update the context for"`
	UserProfile   map[string]any               `json:"user_profile" jsonschema_description:"General information about the user. Must provide the complete user profile as it will replace the existing one."`
	UserPortfolio []UserPortfolioHoldingSchema `json:"user_portfolio" jsonschema_description:"List of portfolio holdings. Must provide the complete portfolio as it will replace the existing one."`
}

type UpdateUserContextTool struct {
	userContextService UserContextService
}

func NewUpdateUserContextTool(userContextService UserContextService) (*UpdateUserContextTool, error) {
	return &UpdateUserContextTool{
		userContextService: userContextService,
	}, nil
}

func (t *UpdateUserContextTool) HandleUpdateUserContext(ctx context.Context, req mcp.CallToolRequest, args UpdateUserContextRequest) (UserContextResponse, error) {
	if args.UserID == "" {
		return UserContextResponse{}, fmt.Errorf("user_id is required")
	}

	// Validate portfolio holdings
	for _, holding := range args.UserPortfolio {
		if err := holding.Validate(); err != nil {
			return UserContextResponse{}, err
		}
	}

	// Convert request to domain.UserContext
	portfolioHoldings := make([]domain.UserPortfolioHolding, 0, len(args.UserPortfolio))
	for _, h := range args.UserPortfolio {
		portfolioHoldings = append(
			portfolioHoldings,
			domain.UserPortfolioHolding{
				AssetClass:          domain.AssetClass(h.AssetClass),
				Symbol:              h.Symbol,
				Name:                h.Name,
				Quantity:            h.Quantity,
				PortfolioPercentage: h.PortfolioPercentage,
			},
		)
	}

	userContext := domain.UserContext{
		UserID:        args.UserID,
		UserProfile:   args.UserProfile,
		UserPortfolio: portfolioHoldings,
	}

	err := t.userContextService.UpdateUserContext(userContext)
	if err != nil {
		return UserContextResponse{}, err
	}

	// Fetch the updated user context to return what's actually stored
	updatedUserContext, err := t.userContextService.GetUserContext(args.UserID)
	if err != nil {
		return UserContextResponse{}, err
	}

	// Return the updated user context
	portfolio := make([]UserPortfolioHoldingSchema, 0, len(updatedUserContext.UserPortfolio))
	for _, holding := range updatedUserContext.UserPortfolio {
		portfolio = append(portfolio, UserPortfolioHoldingSchema{
			AssetClass:          string(holding.AssetClass),
			Symbol:              holding.Symbol,
			Name:                holding.Name,
			Quantity:            holding.Quantity,
			PortfolioPercentage: holding.PortfolioPercentage,
		})
	}

	response := UserContextResponse{
		UserID:        updatedUserContext.UserID,
		UserProfile:   updatedUserContext.UserProfile,
		UserPortfolio: portfolio,
	}

	return response, nil
}

func (t *UpdateUserContextTool) GetTool() mcp.Tool {
	return mcp.NewTool("updateUserContext",
		mcp.WithDescription("Update the user context including user profile and portfolio holdings. Note: The provided context will completely replace the existing one, so the entire updated object must be provided."),
		mcp.WithInputSchema[UpdateUserContextRequest](),
		mcp.WithOutputSchema[UserContextResponse](),
	)
}
