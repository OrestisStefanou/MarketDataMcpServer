package main

import (
	"context"
	"log"
	alphavantage "market_data_mcp_server/pkg/alpha_vantage"
	"market_data_mcp_server/pkg/api/mcp/tools"
	coingecko "market_data_mcp_server/pkg/coin_gecko"
	"market_data_mcp_server/pkg/config"
	"market_data_mcp_server/pkg/marketDataScraper"
	"market_data_mcp_server/pkg/services"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	conf, _ := config.LoadConfig()

	// Initialize components
	logger := log.New(os.Stdout, "[MCP] ", log.LstdFlags)

	// Create middleware
	loggingMW := NewLoggingMiddleware(logger)

	mcpServer := server.NewMCPServer(
		"Market Data MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(false, true),
		server.WithPromptCapabilities(true),
		server.WithRecovery(),
		server.WithToolHandlerMiddleware(loggingMW.ToolMiddleware),
	)

	// Setup cache and data services
	cache, _ := services.NewBadgerCacheService()
	dataService := marketDataScraper.NewMarketDataScraperWithCache(cache, conf)

	alphaVantageClient, _ := alphavantage.NewAlphaVantageClientWithCache(conf.AlphaVantageApiKey, cache, conf.AlphaVantageCacheTtl)
	coinGeckoClient, _ := coingecko.NewCoinGeckoClientWithCache(conf.CoinGeckoApiKey, cache, conf.CoinGeckoCacheTtl)

	// Set up services
	tickerService, _ := services.NewTickerService(dataService)
	etfService, _ := services.NewEtfService(dataService)
	superInvestorService, _ := services.NewSuperInvestorService(dataService)
	cryptoService, _ := services.NewCryptoService(coinGeckoClient, alphaVantageClient)

	// Setup tools
	searchStocksTool, _ := tools.NewStockSearchTool(tickerService)
	searchEtfsTool, _ := tools.NewSearchEtfTool(etfService)
	getEtfTool, _ := tools.NewGetEtfTool(etfService)
	getSuperInvestorsTool, _ := tools.NewGetSuperInvestorsTool(superInvestorService)
	getSuperInvestorPortfolioTool, _ := tools.NewGetSuperInvestorPortfolioTool(superInvestorService)
	getMarketNewsTool, _ := tools.NewGetMarketNewsTool(dataService)
	getSectorsTool, _ := tools.NewGetSectorsTool(dataService)
	getSectorStocksTool, _ := tools.NewGetSectorStocksTool(dataService)
	getStockOverviewTool, _ := tools.NewGetStockOverviewTool(dataService)
	getStockFinancialsTool, _ := tools.NewGetStockFinancialsTool(dataService)
	getEconomicIndicatorTimeSeriesTool, _ := tools.NewGetEconomicIndicatorTimeSeriesTool(alphaVantageClient)
	getCommodityTimeSeriesTool, _ := tools.NewGetCommodityTimeSeriesTool(alphaVantageClient)
	searchCryptocurrenciesTool, _ := tools.NewSearchCryptocurrenciesTool(cryptoService)
	getCryptocurrencyDataByIdTool, _ := tools.NewGetCryptocurrencyDataByIdTool(cryptoService)
	getCryptocurrencyNewsTool, _ := tools.NewGetCryptocurrencyNewsTool(cryptoService)
	calculateInvestmentFutureValueTool, _ := tools.NewCalculateInvestmentFutureValueTool()
	getEarningsCallTranscriptTool, _ := tools.NewGetEarningsCallTranscriptTool(alphaVantageClient)
	getInsiderTransactionsTool, _ := tools.NewGetInsiderTransactionsTool(alphaVantageClient)

	// Add tools
	mcpServer.AddTool(
		searchStocksTool.GetTool(),
		mcp.NewStructuredToolHandler(searchStocksTool.HandleSearchStocks),
	)

	mcpServer.AddTool(
		searchEtfsTool.GetTool(),
		mcp.NewStructuredToolHandler(searchEtfsTool.HandleSearchEtfs),
	)

	mcpServer.AddTool(
		getEtfTool.GetTool(),
		mcp.NewStructuredToolHandler(getEtfTool.HandleGetEtf),
	)

	mcpServer.AddTool(
		getSuperInvestorsTool.GetTool(),
		mcp.NewStructuredToolHandler(getSuperInvestorsTool.HandleGetSuperInvestors),
	)

	mcpServer.AddTool(
		getSuperInvestorPortfolioTool.GetTool(),
		mcp.NewStructuredToolHandler(getSuperInvestorPortfolioTool.HandleGetSuperInvestorPortfolio),
	)

	mcpServer.AddTool(
		getMarketNewsTool.GetTool(),
		mcp.NewStructuredToolHandler(getMarketNewsTool.HandleGetNews),
	)

	mcpServer.AddTool(
		getSectorsTool.GetTool(),
		mcp.NewStructuredToolHandler(getSectorsTool.HandleGetSectors),
	)

	mcpServer.AddTool(
		getSectorStocksTool.GetTool(),
		mcp.NewStructuredToolHandler(getSectorStocksTool.HandleGetSectorStocks),
	)

	mcpServer.AddTool(
		getStockOverviewTool.GetTool(),
		mcp.NewStructuredToolHandler(getStockOverviewTool.HandleGetStockOverview),
	)

	mcpServer.AddTool(
		getStockFinancialsTool.GetTool(),
		mcp.NewStructuredToolHandler(getStockFinancialsTool.HandleGetStockFinancials),
	)

	mcpServer.AddTool(
		getEconomicIndicatorTimeSeriesTool.GetTool(),
		mcp.NewStructuredToolHandler(getEconomicIndicatorTimeSeriesTool.HandleGetEconomicIndicatorTimeSeries),
	)

	mcpServer.AddTool(
		getCommodityTimeSeriesTool.GetTool(),
		mcp.NewStructuredToolHandler(getCommodityTimeSeriesTool.HandleGetCommodityTimeSeries),
	)

	mcpServer.AddTool(
		searchCryptocurrenciesTool.GetTool(),
		mcp.NewStructuredToolHandler(searchCryptocurrenciesTool.HandleSearchCryptocurrencies),
	)

	mcpServer.AddTool(
		getCryptocurrencyDataByIdTool.GetTool(),
		mcp.NewStructuredToolHandler(getCryptocurrencyDataByIdTool.HandleGetCryptocurrencyDataById),
	)

	mcpServer.AddTool(
		getCryptocurrencyNewsTool.GetTool(),
		mcp.NewStructuredToolHandler(getCryptocurrencyNewsTool.HandleGetCryptocurrencyNews),
	)

	mcpServer.AddTool(
		calculateInvestmentFutureValueTool.GetTool(),
		mcp.NewStructuredToolHandler(calculateInvestmentFutureValueTool.HandleCalculateInvestmentFutureValue),
	)

	mcpServer.AddTool(
		getEarningsCallTranscriptTool.GetTool(),
		mcp.NewStructuredToolHandler(getEarningsCallTranscriptTool.HandleGetEarningsCallTranscript),
	)

	mcpServer.AddTool(
		getInsiderTransactionsTool.GetTool(),
		mcp.NewStructuredToolHandler(getInsiderTransactionsTool.HandleGetInsiderTransactions),
	)

	// Start the server
	startWithGracefulShutdown(mcpServer)
}

func startWithGracefulShutdown(mcpServer *server.MCPServer) {
	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var httpServer *server.StreamableHTTPServer

	// Start server in goroutine
	go func() {
		httpServer = server.NewStreamableHTTPServer(mcpServer)
		if err := httpServer.Start(":8080"); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
