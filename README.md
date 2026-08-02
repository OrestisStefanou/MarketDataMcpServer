# MarketDataMcpServer

MarketDataMcpServer is a Model Context Protocol (MCP) server that provides comprehensive market data for stocks, ETFs, cryptocurrencies, and economic indicators. It runs without any API keys, drawing on public sources such as FRED, the European Central Bank, SEC EDGAR, CoinGecko, and Polymarket.

## Features

- **Stock Market Data**: Search for stocks, get detailed company overviews, and access financial statements.
- **ETF Analysis**: Search for ETFs and view their holdings and detailed information.
- **Crypto Tracking**: Search for cryptocurrencies, get real-time data, and stay updated with crypto news.
- **Super Investor Insights**: Track institutional "Super Investors" and their portfolios.
- **Economic Indicators**: Access time series data for key economic indicators (GDP, Inflation, etc.) and commodities (Oil, Gas, etc.), sourced from FRED.
- **Insider Activity**: Read Form 4 insider transactions straight from SEC EDGAR.
- **Market Intelligence**: Get the latest market news and sector performances.
- **Prediction Markets**: Search Polymarket for events by free-text query and get current implied odds for their markets.

## Prerequisites

- [Go](https://go.dev/) 1.21 or higher.

**No API keys are required.** The server runs with zero configuration. A CoinGecko key is optional and only raises rate limits.

## Data sources

| Source | Used for | Key required |
| --- | --- | --- |
| [FRED](https://fred.stlouisfed.org/) | Economic indicators, commodities | No |
| [Frankfurter](https://frankfurter.dev/) (ECB reference rates) | Currency exchange rates | No |
| [SEC EDGAR](https://www.sec.gov/edgar) | Insider transactions (Form 4) | No, but a contact address is required in the User-Agent |
| stockanalysis.com | Stocks, ETFs, sectors, financials, news | No |
| dataroma.com | Super investor portfolios | No |
| [CoinGecko](https://www.coingecko.com/en/api) | Cryptocurrency search and data | Optional |
| CoinDesk, CoinTelegraph RSS | Cryptocurrency news | No |
| [Polymarket Gamma](https://docs.polymarket.com/) | Prediction market odds | No |

## Configuration

Every setting has a working default. To override any of them, create a `.env` file in the root directory:

```env
PORT=8080

# Optional. Raises CoinGecko's rate limits; the server works without it.
COIN_GECKO_API_KEY=

# The SEC asks automated clients to identify themselves with a contact address.
# Replace the default with your own email. Note that EDGAR rejects a User-Agent
# containing a URL.
SEC_EDGAR_USER_AGENT=MarketDataMcpServer/1.0 (you@example.com)

# Cache TTL (in seconds)
CACHE_TTL=3600
COIN_GECKO_CACHE_TTL=3600
POLYMARKET_CACHE_TTL=300
FRED_CACHE_TTL=86400
FRANKFURTER_CACHE_TTL=3600
SEC_EDGAR_CACHE_TTL=86400
CRYPTO_NEWS_CACHE_TTL=900
```

## Getting Started

### Installation

1. Clone the repository.
2. Install dependencies:
   ```bash
   make install
   ```

### Running the Server

Start the MCP server on port `8080`:
```bash
make run_mcp_server
```

### Building

Build the binary:
```bash
make build_mcp_server
```

### Testing

Unit tests cover the response parsers:
```bash
go test ./...
```

`mcp_client/validation.py` is an end-to-end check that calls every registered tool against a
running server. It expects the server on port `8082`:
```bash
PORT=8082 make run_mcp_server
uv run --with fastmcp python mcp_client/validation.py
```

## Available Tools

| Tool | Description |
| --- | --- |
| `stockSearch` | Search for stock tickers based on keywords. |
| `etfSearch` | Search for ETFs based on keywords. |
| `getETF` | Get detailed information and holdings for a specific ETF. |
| `getSuperInvestors` | List tracked institutional super investors. |
| `getSuperInvestorPortfolio` | Get the portfolio holdings of a specific super investor. |
| `getMarketNews` | Get the latest global market news. |
| `getSectors` | Get a list of market sectors and their performance. |
| `getSectorStocks` | Get the top stocks for a specific sector. |
| `getStockOverview` | Get a comprehensive overview of a company (valuation, growth, etc.). |
| `getStockFinancials` | Get financial statements (Income Statement, Balance Sheet, Cash Flow). |
| `getEconomicIndicatorTimeSeries` | Get historical data for economic indicators (GDP, Inflation, Interest Rate, Unemployment, Treasury Yield) from FRED. Treasury yields are daily; the rest are monthly or quarterly. Inflation is reported as a year-over-year rate, not a CPI index level. |
| `getCommodityTimeSeries` | Get historical data for commodities (Crude Oil, Natural Gas, Copper, Aluminum, Wheat, Corn, Sugar, Coffee) from FRED. Oil and gas are daily; the rest are monthly IMF global prices published with a one to two month lag. |
| `searchCryptocurrencies` | Search for cryptocurrencies on CoinGecko. |
| `getCryptocurrencyDataById` | Get detailed real-time data for a specific cryptocurrency. |
| `getCryptocurrencyNews` | Get the latest news related to a cryptocurrency. Returns `matched_symbol`, which is false when no recent article mentioned the symbol and the response is general crypto news instead. |
| `calculateInvestmentFutureValue` | Calculate the future value of an investment based on initial amount, annual return, and years. |
| `getInsiderTransactions` | Get the insider transactions of the stock with the given symbol and year, parsed from SEC EDGAR Form 4 filings. Limited to the years EDGAR keeps in its recent submissions index. |
| `getCompanyKpiMetrics` | Get the KPI metrics(revenue breakdown, revenue by geography etc) of the stock with the given symbol. |
| `getCurrencyExchangeRate` | Get the exchange rate between two currencies (USD, EUR, GBP, JPY, CHF, CAD, AUD). |
| `getPolymarketEventOdds` | Search Polymarket for prediction-market events matching a free-text query and return matching events with their markets and implied outcome probabilities. |
