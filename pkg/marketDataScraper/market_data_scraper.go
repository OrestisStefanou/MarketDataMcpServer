package marketDataScraper

import (
	"fmt"
	"market_data_mcp_server/pkg/config"
	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/services"
	"time"
)

type MarketDataScraper struct {
}

// GetSectorStocks returns a list of stocks in a sector
// sector parameter should be the domain.Sector.UrlName value
func (mds MarketDataScraper) GetSectorStocks(sector string) ([]domain.SectorStock, error) {
	return scrapeSectorStocks(sector)
}

// GetSectors returns a list of sectors
func (mds MarketDataScraper) GetSectors() ([]domain.Sector, error) {
	return scrapeSectors()
}

// GetIndustryStocks returns a list of stocks in an industry
// industry parameter should be the domain.Industry.UrlName value
func (mds MarketDataScraper) GetIndustryStocks(industry string) ([]domain.IndustryStock, error) {
	return scrapeIndustryStocks(industry)
}

// GetIndustries returns a list of industries
func (mds MarketDataScraper) GetIndustries() ([]domain.Industry, error) {
	return scrapeIndustries()
}

// GetStockForecsat returns the forecast for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetStockForecast(symbol string) (domain.StockForecast, error) {
	return scrapeStockForecast(symbol)
}

// GetBalanceSheets returns a list of balance sheets for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetBalanceSheets(symbol string) ([]domain.BalanceSheet, error) {
	return scrapeBalanceSheets(symbol)
}

// GetIncomeStatements returns a list of income statements for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetIncomeStatements(symbol string) ([]domain.IncomeStatement, error) {
	return scrapeIncomeStatements(symbol)
}

// GetCashFlows returns a list of cash flows for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetCashFlows(symbol string) ([]domain.CashFlow, error) {
	return scrapeCashFlows(symbol)
}

// GetFinancialRatios returns a list of financial ratios for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetFinancialRatios(symbol string) ([]domain.FinancialRatios, error) {
	return scrapeFinancialRatios(symbol)
}

// GetEtfs returns a list of ETFs
func (mds MarketDataScraper) GetEtfs() ([]domain.Etf, error) {
	return scrapeEtfs()
}

// GetEtfOverview returns an overview of an ETF
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetEtfOverview(symbol string) (domain.EtfOverview, error) {
	return scrapeEtfOverview(symbol)
}

// GetStockProfile returns the profile of a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetStockProfile(symbol string) (domain.StockProfile, error) {
	return scrapeStockProfile(symbol)
}

// GetMarketNews returns the most recent news of the stock markets
func (mds MarketDataScraper) GetMarketNews() ([]domain.NewsArticle, error) {
	return scrapeMarketNews()
}

// GetStockNews returns the most recent news of the given stock symbol
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetStockNews(symbol string) ([]domain.NewsArticle, error) {
	return scrapeStockNews(symbol)
}

// GetTickers returns a list of Tickers(stock symbol and company name)
func (mds MarketDataScraper) GetTickers() ([]domain.Ticker, error) {
	return scrapeStockList()
}

// GetSuperInvestors returns a list of SuperInvestors (Name)
func (mds MarketDataScraper) GetSuperInvestors() ([]domain.SuperInvestor, error) {
	return scrapeSuperInvestors()
}

// GetSuperInvestorPortfolio returns the portfolio of the given super investor
func (mds MarketDataScraper) GetSuperInvestorPortfolio(superInvestorName string) (domain.SuperInvestorPortfolio, error) {
	return scrapeSuperInvestorPortfolio(superInvestorName)
}

func (mds MarketDataScraper) GetHistoricalPrices(ticker string, assetClass domain.AssetClass, period domain.Period) (domain.HistoricalPrices, error) {
	return scrapeHistoricalPrices(ticker, assetClass, period)
}

type MarketDataScraperWithCache struct {
	cache services.CacheService
	conf  config.Config
}

func NewMarketDataScraperWithCache(cache services.CacheService, conf config.Config) *MarketDataScraperWithCache {
	return &MarketDataScraperWithCache{cache: cache, conf: conf}
}

// GetSectorStocks returns a list of stocks in a sector
// sector parameter should be the domain.Sector.UrlName value
func (mds MarketDataScraperWithCache) GetSectorStocks(sector string) ([]domain.SectorStock, error) {
	// Check if the data is in the cache
	var sectorStocks []domain.SectorStock

	key := fmt.Sprintf("sector_stocks_%s", sector)
	err := mds.cache.Get(key, &sectorStocks)
	if err == nil {
		return sectorStocks, nil
	}

	sectorStocks, err = scrapeSectorStocks(sector)
	if err != nil {
		return nil, err
	}

	mds.cache.Set(key, sectorStocks, time.Duration(mds.conf.CacheTtl)*time.Second)
	return sectorStocks, nil
}

// GetSectors returns a list of sectors
func (mds MarketDataScraperWithCache) GetSectors() ([]domain.Sector, error) {
	// Check if the data is in the cache
	var sectors []domain.Sector

	key := "sectors"
	err := mds.cache.Get(key, &sectors)
	if err == nil {
		return sectors, nil
	}

	sectors, err = scrapeSectors()
	if err != nil {
		return nil, err
	}

	mds.cache.Set(key, sectors, time.Duration(mds.conf.CacheTtl)*time.Second)
	return sectors, nil
}

// GetIndustryStocks returns a list of stocks in an industry
// industry parameter should be the domain.Industry.UrlName value
func (mds MarketDataScraperWithCache) GetIndustryStocks(industry string) ([]domain.IndustryStock, error) {
	// Check if the data is in the cache
	var industryStocks []domain.IndustryStock

	key := fmt.Sprintf("industry_stocks_%s", industry)
	err := mds.cache.Get(key, &industryStocks)
	if err == nil {
		return industryStocks, nil
	}

	industryStocks, err = scrapeIndustryStocks(industry)
	if err != nil {
		return nil, err
	}

	mds.cache.Set(key, industryStocks, time.Duration(mds.conf.CacheTtl)*time.Second)
	return industryStocks, nil
}

// GetIndustries returns a list of industries
func (mds MarketDataScraperWithCache) GetIndustries() ([]domain.Industry, error) {
	// Check if the data is in the cache
	var industries []domain.Industry

	key := "industries"
	err := mds.cache.Get(key, &industries)
	if err == nil {
		return industries, nil
	}

	industries, err = scrapeIndustries()
	if err != nil {
		return nil, err
	}

	mds.cache.Set(key, industries, time.Duration(mds.conf.CacheTtl)*time.Second)
	return industries, nil
}

// GetStockForecsat returns the forecast for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraperWithCache) GetStockForecast(symbol string) (domain.StockForecast, error) {
	// Check if the data is in the cache
	var stockForecast domain.StockForecast

	key := fmt.Sprintf("stock_forecast_%s", symbol)
	err := mds.cache.Get(key, &stockForecast)
	if err == nil {
		return stockForecast, nil
	}

	stockForecast, err = scrapeStockForecast(symbol)
	if err != nil {
		return domain.StockForecast{}, err
	}

	mds.cache.Set(key, stockForecast, time.Duration(mds.conf.CacheTtl)*time.Second)
	return stockForecast, nil
}

// GetBalanceSheets returns a list of balance sheets for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraperWithCache) GetBalanceSheets(symbol string) ([]domain.BalanceSheet, error) {
	// Check if the data is in the cache
	var balanceSheets []domain.BalanceSheet

	key := fmt.Sprintf("balance_sheets_%s", symbol)
	err := mds.cache.Get(key, &balanceSheets)
	if err == nil {
		return balanceSheets, nil
	}

	balanceSheets, err = scrapeBalanceSheets(symbol)
	if err != nil {
		return nil, err
	}

	mds.cache.Set(key, balanceSheets, time.Duration(mds.conf.CacheTtl)*time.Second)
	return balanceSheets, nil
}

// GetIncomeStatements returns a list of income statements for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraperWithCache) GetIncomeStatements(symbol string) ([]domain.IncomeStatement, error) {
	// Check if the data is in the cache
	var incomeStatements []domain.IncomeStatement

	key := fmt.Sprintf("income_statements_%s", symbol)
	err := mds.cache.Get(key, &incomeStatements)
	if err == nil {
		return incomeStatements, nil
	}

	incomeStatements, err = scrapeIncomeStatements(symbol)
	if err != nil {
		return nil, err
	}

	mds.cache.Set(key, incomeStatements, time.Duration(mds.conf.CacheTtl)*time.Second)
	return incomeStatements, nil
}

// GetCashFlows returns a list of cash flows for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraperWithCache) GetCashFlows(symbol string) ([]domain.CashFlow, error) {
	// Check if the data is in the cache
	var cashFlows []domain.CashFlow

	key := fmt.Sprintf("cash_flows_%s", symbol)
	err := mds.cache.Get(key, &cashFlows)
	if err == nil {
		return cashFlows, nil
	}

	cashFlows, err = scrapeCashFlows(symbol)
	if err != nil {
		return nil, err
	}

	mds.cache.Set(key, cashFlows, time.Duration(mds.conf.CacheTtl)*time.Second)
	return cashFlows, nil
}

// GetFinancialRatios returns a list of financial ratios for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraperWithCache) GetFinancialRatios(symbol string) ([]domain.FinancialRatios, error) {
	// Check if the data is in the cache
	var financialRatios []domain.FinancialRatios

	key := fmt.Sprintf("financial_ratios_%s", symbol)
	err := mds.cache.Get(key, &financialRatios)
	if err == nil {
		return financialRatios, nil
	}

	financialRatios, err = scrapeFinancialRatios(symbol)
	if err != nil {
		return nil, err
	}

	mds.cache.Set(key, financialRatios, time.Duration(mds.conf.CacheTtl)*time.Second)
	return financialRatios, nil
}

// GetEtfs returns a list of ETFs
func (mds MarketDataScraperWithCache) GetEtfs() ([]domain.Etf, error) {
	// Check if the data is in the cache
	var etfs []domain.Etf

	key := "etfs"
	err := mds.cache.Get(key, &etfs)
	if err == nil {
		return etfs, nil
	}

	etfs, err = scrapeEtfs()
	if err != nil {
		return nil, err
	}

	mds.cache.Set(key, etfs, time.Duration(mds.conf.CacheTtl)*time.Second)
	return etfs, nil
}

// GetEtfOverview returns an overview of an ETF
// symbol parameter should be in lowercase
func (mds MarketDataScraperWithCache) GetEtfOverview(symbol string) (domain.EtfOverview, error) {
	// Check if the data is in the cache
	var etfOverview domain.EtfOverview

	key := fmt.Sprintf("etf_overview_%s", symbol)
	err := mds.cache.Get(key, &etfOverview)
	if err == nil {
		return etfOverview, nil
	}

	etfOverview, err = scrapeEtfOverview(symbol)
	if err != nil {
		return domain.EtfOverview{}, err
	}

	mds.cache.Set(key, etfOverview, time.Duration(mds.conf.CacheTtl)*time.Second)
	return etfOverview, nil
}

// GetStockProfile returns the profile of a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraperWithCache) GetStockProfile(symbol string) (domain.StockProfile, error) {
	// Check if the data is in the cache
	var stockProfile domain.StockProfile

	key := fmt.Sprintf("stock_profile_%s", symbol)
	err := mds.cache.Get(key, &stockProfile)
	if err == nil {
		return stockProfile, nil
	}

	stockProfile, err = scrapeStockProfile(symbol)
	if err != nil {
		return domain.StockProfile{}, err
	}

	mds.cache.Set(key, stockProfile, time.Duration(mds.conf.CacheTtl)*time.Second)
	return stockProfile, nil
}

// GetMarketNews returns the most recent news of the stock markets
func (mds MarketDataScraperWithCache) GetMarketNews() ([]domain.NewsArticle, error) {
	// Check if the data is in the cache
	var marketNews []domain.NewsArticle

	key := "market_news"
	err := mds.cache.Get(key, &marketNews)
	if err == nil {
		return marketNews, nil
	}

	marketNews, err = scrapeMarketNews()
	if err != nil {
		return nil, err
	}

	mds.cache.Set(key, marketNews, time.Duration(mds.conf.CacheTtl)*time.Second)
	return marketNews, nil
}

// GetStockNews returns the most recent news of the given stock symbol
// symbol parameter should be in lowercase
func (mds MarketDataScraperWithCache) GetStockNews(symbol string) ([]domain.NewsArticle, error) {
	// Check if the data is in the cache
	var stockNews []domain.NewsArticle

	key := fmt.Sprintf("stock_news_%s", symbol)
	err := mds.cache.Get(key, &stockNews)
	if err == nil {
		return stockNews, nil
	}

	stockNews, err = scrapeStockNews(symbol)
	if err != nil {
		return nil, err
	}

	mds.cache.Set(key, stockNews, time.Duration(mds.conf.CacheTtl)*time.Second)
	return stockNews, nil
}

// GetTickers returns a list of Tickers(stock symbol and company name)
func (mds MarketDataScraperWithCache) GetTickers() ([]domain.Ticker, error) {
	// Check if the data is in the cache
	var tickers []domain.Ticker

	key := "tickers"
	err := mds.cache.Get(key, &tickers)
	if err == nil {
		return tickers, nil
	}

	tickers, err = scrapeStockList()
	if err != nil {
		return nil, err
	}

	mds.cache.Set(key, tickers, time.Duration(mds.conf.CacheTtl)*time.Second)
	return tickers, nil
}

// GetSuperInvestors returns a list of SuperInvestors (Name)
func (mds MarketDataScraperWithCache) GetSuperInvestors() ([]domain.SuperInvestor, error) {
	// Check if the data is in the cache
	var superInvestors []domain.SuperInvestor

	key := "super_investors"
	err := mds.cache.Get(key, &superInvestors)
	if err == nil {
		return superInvestors, nil
	}

	superInvestors, err = scrapeSuperInvestors()
	if err != nil {
		return nil, err
	}

	mds.cache.Set(key, superInvestors, time.Duration(mds.conf.CacheTtl)*time.Second)
	return superInvestors, nil
}

// GetSuperInvestorPortfolio returns the portfolio of the given super investor
func (mds MarketDataScraperWithCache) GetSuperInvestorPortfolio(superInvestorName string) (domain.SuperInvestorPortfolio, error) {
	// Check if the data is in the cache
	var superInvestorPortfolio domain.SuperInvestorPortfolio

	key := fmt.Sprintf("super_investor_portfolio_%s", superInvestorName)
	err := mds.cache.Get(key, &superInvestorPortfolio)
	if err == nil {
		return superInvestorPortfolio, nil
	}

	superInvestorPortfolio, err = scrapeSuperInvestorPortfolio(superInvestorName)
	if err != nil {
		return domain.SuperInvestorPortfolio{}, err
	}

	mds.cache.Set(key, superInvestorPortfolio, time.Duration(mds.conf.CacheTtl)*time.Second)
	return superInvestorPortfolio, nil
}

func (mds MarketDataScraperWithCache) GetHistoricalPrices(ticker string, assetClass domain.AssetClass, period domain.Period) (domain.HistoricalPrices, error) {
	// Check if the data is in the cache
	var historicalPrices domain.HistoricalPrices

	key := fmt.Sprintf("historical_prices_%s_%s_%s", ticker, assetClass, period)
	err := mds.cache.Get(key, &historicalPrices)
	if err == nil {
		return historicalPrices, nil
	}

	historicalPrices, err = scrapeHistoricalPrices(ticker, assetClass, period)
	if err != nil {
		return domain.HistoricalPrices{}, err
	}

	mds.cache.Set(key, historicalPrices, time.Duration(mds.conf.CacheTtl)*time.Second)
	return historicalPrices, nil
}
