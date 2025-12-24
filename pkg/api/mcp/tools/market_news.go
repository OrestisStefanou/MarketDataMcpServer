package tools

import (
	"context"
	"market_data_mcp_server/pkg/domain"

	"github.com/mark3labs/mcp-go/mcp"
)

type MarketNewsService interface {
	GetMarketNews() ([]domain.NewsArticle, error)
	GetStockNews(symbol string) ([]domain.NewsArticle, error)
}

type NewsArticleSchema struct {
	Url    string `json:"url" jsonschema_description:"News article URL"`
	Image  string `json:"image" jsonschema_description:"News article image URL"`
	Title  string `json:"title" jsonschema_description:"News article title"`
	Text   string `json:"text" jsonschema_description:"News article text content"`
	Source string `json:"source" jsonschema_description:"News article source"`
	Time   string `json:"time" jsonschema_description:"News article publication time"`
}

type GetMarketNewsRequest struct {
	StockSymbol string `json:"stock_symbol,omitempty" jsonschema_description:"The symbol of the stock to get the news for. Leave empty to get general market news"`
}

type GetMarketNewsResponse struct {
	News []NewsArticleSchema `json:"stock_symbol,omitempty" jsonschema_description:"A list of market news articles"`
}

type GetMarketNewsTool struct {
	newsService MarketNewsService
}

func NewGetMarketNewsTool(newsService MarketNewsService) (*GetMarketNewsTool, error) {
	return &GetMarketNewsTool{
		newsService: newsService,
	}, nil
}

func (t *GetMarketNewsTool) HandleGetNews(ctx context.Context, req mcp.CallToolRequest, args GetMarketNewsRequest) (GetMarketNewsResponse, error) {
	var news []domain.NewsArticle
	var err error

	if args.StockSymbol != "" {
		news, err = t.newsService.GetStockNews(args.StockSymbol)
	} else {
		news, err = t.newsService.GetMarketNews()
	}
	if err != nil {
		return GetMarketNewsResponse{}, err
	}

	response := GetMarketNewsResponse{
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

func (t *GetMarketNewsTool) GetTool() mcp.Tool {
	return mcp.NewTool("getMarketNews",
		mcp.WithDescription("Get the latest market news of a stock or of the market in general"),
		mcp.WithInputSchema[GetMarketNewsRequest](),
		mcp.WithOutputSchema[GetMarketNewsResponse](),
	)
}
