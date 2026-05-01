package tools

import (
	"context"
	"fmt"
	"market_data_mcp_server/pkg/domain"

	"github.com/mark3labs/mcp-go/mcp"
)

type CurrencyExchangeService interface {
	GetCurrencyExchangeRate(fromCurrency domain.Currency, toCurrency domain.Currency) (domain.CurrencyExchangeRate, error)
}

type GetCurrencyExchangeRateRequest struct {
	FromCurrency string `json:"from_currency" jsonschema_description:"Currency code of the currency to convert from" jsonschema:"enum=AED,enum=USD,enum=EUR,enum=GBP,enum=JPY,enum=CHF,enum=CAD,enum=AUD" jsonschema_required:"true"`
	ToCurrency   string `json:"to_currency" jsonschema_description:"Currency code of the currency to convert to" jsonschema:"enum=AED,enum=USD,enum=EUR,enum=GBP,enum=JPY,enum=CHF,enum=CAD,enum=AUD" jsonschema_required:"true"`
}

type GetCurrencyExchangeRateResponse struct {
	FromCurrency     string `json:"from_currency" jsonschema_description:"Currency code of the currency to convert from"`
	FromCurrencyName string `json:"from_currency_name" jsonschema_description:"Name of the currency to convert from"`
	ToCurrency       string `json:"to_currency" jsonschema_description:"Currency code of the currency to convert to"`
	ToCurrencyName   string `json:"to_currency_name" jsonschema_description:"Name of the currency to convert to"`
	Rate             string `json:"rate" jsonschema_description:"Exchange rate"`
}

type GetCurrencyExchangeRateTool struct {
	currencyExchangeService CurrencyExchangeService
}

func NewGetCurrencyExchangeRateTool(currencyExchangeService CurrencyExchangeService) (*GetCurrencyExchangeRateTool, error) {
	return &GetCurrencyExchangeRateTool{
		currencyExchangeService: currencyExchangeService,
	}, nil
}

func (t *GetCurrencyExchangeRateTool) HandleGetCurrencyExchangeRate(ctx context.Context, req mcp.CallToolRequest, args GetCurrencyExchangeRateRequest) (GetCurrencyExchangeRateResponse, error) {
	if _, ok := domain.CurrencyCodeToNameMap[domain.Currency(args.FromCurrency)]; !ok {
		return GetCurrencyExchangeRateResponse{}, fmt.Errorf("invalid from_currency: %s", args.FromCurrency)
	}
	if _, ok := domain.CurrencyCodeToNameMap[domain.Currency(args.ToCurrency)]; !ok {
		return GetCurrencyExchangeRateResponse{}, fmt.Errorf("invalid to_currency: %s", args.ToCurrency)
	}

	rate, err := t.currencyExchangeService.GetCurrencyExchangeRate(domain.Currency(args.FromCurrency), domain.Currency(args.ToCurrency))
	if err != nil {
		return GetCurrencyExchangeRateResponse{}, err
	}

	return GetCurrencyExchangeRateResponse{
		FromCurrency:     string(rate.FromCurrency),
		FromCurrencyName: string(rate.FromCurrencyName),
		ToCurrency:       string(rate.ToCurrency),
		ToCurrencyName:   string(rate.ToCurrencyName),
		Rate:             fmt.Sprintf("%v", rate.Rate),
	}, nil
}

func (t *GetCurrencyExchangeRateTool) GetTool() mcp.Tool {
	return mcp.NewTool("getCurrencyExchangeRate",
		mcp.WithDescription("Get the exchange rate between two currencies."),
		mcp.WithInputSchema[GetCurrencyExchangeRateRequest](),
		mcp.WithOutputSchema[GetCurrencyExchangeRateResponse](),
	)
}
