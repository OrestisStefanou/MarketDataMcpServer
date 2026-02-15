package tools

import (
	"context"
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

type GetInsiderTransactionsRequest struct {
	StockSymbol string `json:"stock_symbol" jsonschema_description:"Symbol of the stock to get data for"`
	Year        int    `json:"year" jsonschema_description:"Year of the insider transactions"`
}

type InsiderTransaction struct {
	TransactionDate       string  `json:"transaction_date"`
	Ticker                string  `json:"ticker"`
	Executive             string  `json:"executive"`
	ExecutiveTitle        string  `json:"executive_title"`
	SecurityType          string  `json:"security_type"`
	AcquisitionOrDisposal string  `json:"acquisition_or_disposal"`
	Shares                float64 `json:"shares"`
	SharePrice            float64 `json:"share_price"`
}

type GetInsiderTransactionsResponse struct {
	Symbol              string               `json:"symbol" jsonschema_description:"Symbol of the stock"`
	Year                int                  `json:"year" jsonschema_description:"Year of the insider transactions"`
	InsiderTransactions []InsiderTransaction `json:"insider_transactions" jsonschema_description:"Insider transactions of the stock"`
}

type InsiderTransactionsService interface {
	GetInsiderTransactions(symbol string) ([]domain.InsiderTransaction, error)
}

type GetInsiderTransactionsTool struct {
	insiderTransactionsService InsiderTransactionsService
}

func NewGetInsiderTransactionsTool(insiderTransactionsService InsiderTransactionsService) (*GetInsiderTransactionsTool, error) {
	if insiderTransactionsService == nil {
		return nil, fmt.Errorf("insiderTransactionsService is required")
	}
	return &GetInsiderTransactionsTool{insiderTransactionsService: insiderTransactionsService}, nil
}

func (t *GetInsiderTransactionsTool) HandleGetInsiderTransactions(ctx context.Context, req mcp.CallToolRequest, args GetInsiderTransactionsRequest) (GetInsiderTransactionsResponse, error) {
	if args.StockSymbol == "" {
		return GetInsiderTransactionsResponse{}, fmt.Errorf("stock_symbol is required")
	}

	if args.Year == 0 {
		return GetInsiderTransactionsResponse{}, fmt.Errorf("year is required")
	}

	insiderTransactions, err := t.insiderTransactionsService.GetInsiderTransactions(args.StockSymbol)
	if err != nil {
		return GetInsiderTransactionsResponse{}, err
	}

	// Filter by year
	var filteredTransactions []domain.InsiderTransaction
	yearStr := fmt.Sprintf("%d", args.Year)
	for _, transaction := range insiderTransactions {
		if strings.HasPrefix(transaction.TransactionDate, yearStr) {
			filteredTransactions = append(filteredTransactions, transaction)
		}
	}

	// Convert to InsiderTransaction
	insiderTransactionsResponse := make([]InsiderTransaction, 0)
	for _, transaction := range filteredTransactions {
		insiderTransactionsResponse = append(insiderTransactionsResponse, InsiderTransaction{
			TransactionDate:       transaction.TransactionDate,
			Ticker:                transaction.Ticker,
			Executive:             transaction.Executive,
			ExecutiveTitle:        transaction.ExecutiveTitle,
			SecurityType:          transaction.SecurityType,
			AcquisitionOrDisposal: transaction.AcquisitionOrDisposal,
			Shares:                transaction.Shares,
			SharePrice:            transaction.SharePrice,
		})
	}

	return GetInsiderTransactionsResponse{
		Symbol:              args.StockSymbol,
		Year:                args.Year,
		InsiderTransactions: insiderTransactionsResponse,
	}, nil
}

func (t *GetInsiderTransactionsTool) GetTool() mcp.Tool {
	return mcp.NewTool("getInsiderTransactions",
		mcp.WithDescription("Get the insider transactions of the stock with the given symbol."),
		mcp.WithInputSchema[GetInsiderTransactionsRequest](),
		mcp.WithOutputSchema[GetInsiderTransactionsResponse](),
	)
}
