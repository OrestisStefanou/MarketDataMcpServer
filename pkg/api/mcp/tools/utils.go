package tools

import (
	"context"
	"math"

	"github.com/mark3labs/mcp-go/mcp"
)

type CalculateInvestmentFutureValueRequest struct {
	InitialInvestment float64 `json:"initial_investment" jsonschema_description:"Initial investment amount"`
	AnnualReturn      float64 `json:"annual_return" jsonschema_description:"Annual return percentage(10 means 10%)"`
	Years             int     `json:"years" jsonschema_description:"Number of years"`
}

type CalculateInvestmentFutureValueResponse struct {
	FutureValue float64 `json:"future_value" jsonschema_description:"Future value of the investment"`
}

type CalculateInvestmentFutureValueTool struct{}

func NewCalculateInvestmentFutureValueTool() (*CalculateInvestmentFutureValueTool, error) {
	return &CalculateInvestmentFutureValueTool{}, nil
}

func (t *CalculateInvestmentFutureValueTool) HandleCalculateInvestmentFutureValue(ctx context.Context, req mcp.CallToolRequest, args CalculateInvestmentFutureValueRequest) (CalculateInvestmentFutureValueResponse, error) {
	futureValue := args.InitialInvestment * math.Pow(1+args.AnnualReturn/100, float64(args.Years))
	return CalculateInvestmentFutureValueResponse{
		FutureValue: futureValue,
	}, nil
}

func (t *CalculateInvestmentFutureValueTool) GetTool() mcp.Tool {
	return mcp.NewTool("calculateInvestmentFutureValue",
		mcp.WithDescription("Calculate the future value of an investment"),
		mcp.WithInputSchema[CalculateInvestmentFutureValueRequest](),
		mcp.WithOutputSchema[CalculateInvestmentFutureValueResponse](),
	)
}
