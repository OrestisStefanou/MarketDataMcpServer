package tools

import (
	"context"
	"fmt"
	"market_data_mcp_server/pkg/domain"

	"github.com/mark3labs/mcp-go/mcp"
)

type CryptoDataService interface {
	SearchCryptocurrencies(query string) ([]domain.Cryptocurrency, error)
	GetCryptocurrencyDataById(id string) (domain.CryptocurrencyData, error)
	GetCryptocurrencyNews(symbol string) ([]domain.NewsArticle, error)
}

type SearchCryptocurrenciesRequest struct {
	SearchQuery string `json:"search_query" jsonschema_description:"Search query(search by name, symbol)"`
	Limit       int    `json:"limit" jsonschema_description:"Number of results to return (default is 50)"`
}

type SearchCryptocurrencyResultSchema struct {
	Id     string `json:"id" jsonschema_description:"ID of the cryptocurrency"`
	Name   string `json:"name" jsonschema_description:"Name of the cryptocurrency"`
	Symbol string `json:"symbol" jsonschema_description:"Symbol of the cryptocurrency"`
}

type SearchCryptocurrenciesResponse struct {
	Results []SearchCryptocurrencyResultSchema `json:"results" jsonschema_description:"Search results"`
}

type SearchCryptocurrenciesTool struct {
	cryptoDataService CryptoDataService
}

func NewSearchCryptocurrenciesTool(cryptoDataService CryptoDataService) (*SearchCryptocurrenciesTool, error) {
	return &SearchCryptocurrenciesTool{
		cryptoDataService: cryptoDataService,
	}, nil
}

func (t *SearchCryptocurrenciesTool) HandleSearchCryptocurrencies(ctx context.Context, req mcp.CallToolRequest, args SearchCryptocurrenciesRequest) (SearchCryptocurrenciesResponse, error) {
	var err error

	cryptocurrencies, err := t.cryptoDataService.SearchCryptocurrencies(args.SearchQuery)
	if err != nil {
		return SearchCryptocurrenciesResponse{}, err
	}

	response := SearchCryptocurrenciesResponse{
		Results: make([]SearchCryptocurrencyResultSchema, 0, len(cryptocurrencies)),
	}

	limit := 50
	if args.Limit > 0 {
		limit = args.Limit
	}

	for i, cryptocurrency := range cryptocurrencies {
		if i >= limit {
			break
		}

		response.Results = append(response.Results, SearchCryptocurrencyResultSchema{
			Id:     cryptocurrency.Id,
			Name:   cryptocurrency.Name,
			Symbol: cryptocurrency.Symbol,
		})
	}

	return response, nil
}

func (t *SearchCryptocurrenciesTool) GetTool() mcp.Tool {
	return mcp.NewTool("searchCryptocurrencies",
		mcp.WithDescription("Search for cryptocurrencies by name or symbol."),
		mcp.WithInputSchema[SearchCryptocurrenciesRequest](),
		mcp.WithOutputSchema[SearchCryptocurrenciesResponse](),
	)
}

type GetCryptocurrencyDataByIdRequest struct {
	Id string `json:"id" jsonschema_description:"ID of the cryptocurrency"`
}

type GetCryptocurrencyDataByIdResponse struct {
	Id                        string  `json:"id" jsonschema_description:"ID of the cryptocurrency"`
	Name                      string  `json:"name" jsonschema_description:"Name of the cryptocurrency"`
	Symbol                    string  `json:"symbol" jsonschema_description:"Symbol of the cryptocurrency"`
	Description               string  `json:"description" jsonschema_description:"Description of the cryptocurrency"`
	Whitepaper                string  `json:"whitepaper" jsonschema_description:"Whitepaper of the cryptocurrency"`
	CurrentUsdPrice           float64 `json:"current_usd_price" jsonschema_description:"Current price of the cryptocurrency in USD"`
	MarketCapUsd              float64 `json:"market_cap_usd" jsonschema_description:"Market cap of the cryptocurrency in USD"`
	PriceChangePercentage24h  string  `json:"price_change_percentage_24h" jsonschema_description:"Price change percentage of the cryptocurrency in the last 24 hours"`
	PriceChangePercentage7d   string  `json:"price_change_percentage_7d" jsonschema_description:"Price change percentage of the cryptocurrency in the last 7 days"`
	PriceChangePercentage14d  string  `json:"price_change_percentage_14d" jsonschema_description:"Price change percentage of the cryptocurrency in the last 14 days"`
	PriceChangePercentage30d  string  `json:"price_change_percentage_30d" jsonschema_description:"Price change percentage of the cryptocurrency in the last 30 days"`
	PriceChangePercentage60d  string  `json:"price_change_percentage_60d" jsonschema_description:"Price change percentage of the cryptocurrency in the last 60 days"`
	PriceChangePercentage200d string  `json:"price_change_percentage_200d" jsonschema_description:"Price change percentage of the cryptocurrency in the last 200 days"`
	PriceChangePercentage1y   string  `json:"price_change_percentage_1y" jsonschema_description:"Price change percentage of the cryptocurrency in the last year"`
	TotalSupply               float64 `json:"total_supply" jsonschema_description:"Total supply of the cryptocurrency"`
	MaxSupply                 float64 `json:"max_supply" jsonschema_description:"Max supply of the cryptocurrency(zero means infinite)"`
}

type GetCryptocurrencyDataByIdTool struct {
	cryptoDataService CryptoDataService
}

func NewGetCryptocurrencyDataByIdTool(cryptoDataService CryptoDataService) (*GetCryptocurrencyDataByIdTool, error) {
	return &GetCryptocurrencyDataByIdTool{
		cryptoDataService: cryptoDataService,
	}, nil
}

func (t *GetCryptocurrencyDataByIdTool) HandleGetCryptocurrencyDataById(ctx context.Context, req mcp.CallToolRequest, args GetCryptocurrencyDataByIdRequest) (GetCryptocurrencyDataByIdResponse, error) {
	var err error

	cryptocurrencyData, err := t.cryptoDataService.GetCryptocurrencyDataById(args.Id)
	if err != nil {
		return GetCryptocurrencyDataByIdResponse{}, err
	}

	return GetCryptocurrencyDataByIdResponse{
		Id:                        cryptocurrencyData.Id,
		Name:                      cryptocurrencyData.Name,
		Symbol:                    cryptocurrencyData.Symbol,
		Description:               cryptocurrencyData.Description,
		Whitepaper:                cryptocurrencyData.Whitepaper,
		CurrentUsdPrice:           cryptocurrencyData.CurrentUsdPrice,
		MarketCapUsd:              cryptocurrencyData.MarketCapUsd,
		PriceChangePercentage24h:  fmt.Sprintf("%f %%", cryptocurrencyData.PriceChangePercentage24h),
		PriceChangePercentage7d:   fmt.Sprintf("%f %%", cryptocurrencyData.PriceChangePercentage7d),
		PriceChangePercentage14d:  fmt.Sprintf("%f %%", cryptocurrencyData.PriceChangePercentage14d),
		PriceChangePercentage30d:  fmt.Sprintf("%f %%", cryptocurrencyData.PriceChangePercentage30d),
		PriceChangePercentage60d:  fmt.Sprintf("%f %%", cryptocurrencyData.PriceChangePercentage60d),
		PriceChangePercentage200d: fmt.Sprintf("%f %%", cryptocurrencyData.PriceChangePercentage200d),
		PriceChangePercentage1y:   fmt.Sprintf("%f %%", cryptocurrencyData.PriceChangePercentage1y),
		TotalSupply:               cryptocurrencyData.TotalSupply,
		MaxSupply:                 cryptocurrencyData.MaxSupply,
	}, nil
}

func (t *GetCryptocurrencyDataByIdTool) GetTool() mcp.Tool {
	return mcp.NewTool("getCryptocurrencyDataById",
		mcp.WithDescription("Get the data of the given cryptocurrency by ID."),
		mcp.WithInputSchema[GetCryptocurrencyDataByIdRequest](),
		mcp.WithOutputSchema[GetCryptocurrencyDataByIdResponse](),
	)
}

type GetCryptocurrencyNewsRequest struct {
	Symbol string `json:"symbol" jsonschema_description:"Symbol of the cryptocurrency"`
}

type GetCryptocurrencyNewsResponse struct {
	News []NewsArticleSchema `json:"news" jsonschema_description:"A list of crypto news articles"`
}

type GetCryptocurrencyNewsTool struct {
	cryptoDataService CryptoDataService
}

func NewGetCryptocurrencyNewsTool(cryptoDataService CryptoDataService) (*GetCryptocurrencyNewsTool, error) {
	return &GetCryptocurrencyNewsTool{
		cryptoDataService: cryptoDataService,
	}, nil
}

func (t *GetCryptocurrencyNewsTool) HandleGetCryptocurrencyNews(ctx context.Context, req mcp.CallToolRequest, args GetCryptocurrencyNewsRequest) (GetCryptocurrencyNewsResponse, error) {
	news, err := t.cryptoDataService.GetCryptocurrencyNews(args.Symbol)
	if err != nil {
		return GetCryptocurrencyNewsResponse{}, err
	}

	response := GetCryptocurrencyNewsResponse{
		News: make([]NewsArticleSchema, 0, len(news)),
	}

	for _, article := range news {
		response.News = append(response.News, NewsArticleSchema{
			Url:    article.Url,
			Image:  article.Image,
			Title:  article.Title,
			Text:   article.Text,
			Source: article.Source,
			Time:   article.Time,
		})
	}

	return response, nil
}

func (t *GetCryptocurrencyNewsTool) GetTool() mcp.Tool {
	return mcp.NewTool("getCryptocurrencyNews",
		mcp.WithDescription("Get the news of the given cryptocurrency by symbol."),
		mcp.WithInputSchema[GetCryptocurrencyNewsRequest](),
		mcp.WithOutputSchema[GetCryptocurrencyNewsResponse](),
	)
}
