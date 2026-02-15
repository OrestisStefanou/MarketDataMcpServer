# MarketDataMcpServer

MarketDataMcpServer is a Model Context Protocol (MCP) server that provides comprehensive market data for stocks, ETFs, cryptocurrencies, and economic indicators. It leverages multiple data providers like Alpha Vantage and CoinGecko to offer high-quality financial information to AI models.

## Features

- **Stock Market Data**: Search for stocks, get detailed company overviews, and access financial statements.
- **ETF Analysis**: Search for ETFs and view their holdings and detailed information.
- **Crypto Tracking**: Search for cryptocurrencies, get real-time data, and stay updated with crypto news.
- **Super Investor Insights**: Track institutional "Super Investors" and their portfolios.
- **Economic Indicators**: Access time series data for key economic indicators (GDP, Inflation, etc.) and commodities (Oil, Gas, etc.).
- **Market Intelligence**: Get the latest market news and sector performances.

## Prerequisites

- [Go](https://go.dev/) 1.21 or higher.
- API keys for:
  - [Alpha Vantage](https://www.alphavantage.co/support/#api-key)
  - [CoinGecko](https://www.coingecko.com/en/api)

## Configuration

The server is configured via environment variables. Create a `.env` file in the root directory:

```env
# API Keys
ALPHA_VANTAGE_API_KEY=your_alpha_vantage_key
COIN_GECKO_API_KEY=your_coin_gecko_key

# Cache TTL (in seconds)
CACHE_TTL=3600
ALPHA_VANTAGE_CACHE_TTL=3600
COIN_GECKO_CACHE_TTL=3600
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

## Available Tools

| Tool | Description |
| --- | --- |
| `search_stocks` | Search for stock tickers based on keywords. |
| `search_etfs` | Search for ETFs based on keywords. |
| `get_etf` | Get detailed information and holdings for a specific ETF. |
| `get_super_investors` | List tracked institutional super investors. |
| `get_super_investor_portfolio` | Get the portfolio holdings of a specific super investor. |
| `get_market_news` | Get the latest global market news. |
| `get_sectors` | Get a list of market sectors and their performance. |
| `get_sector_stocks` | Get the top stocks for a specific sector. |
| `get_stock_overview` | Get a comprehensive overview of a company (valuation, growth, etc.). |
| `get_stock_financials` | Get financial statements (Income Statement, Balance Sheet, Cash Flow). |
| `get_economic_indicator_time_series` | Get historical data for economic indicators (e.g., GDP, Inflation). |
| `get_commodity_time_series` | Get historical data for commodities (e.g., Crude Oil, Natural Gas). |
| `search_cryptocurrencies` | Search for cryptocurrencies on CoinGecko. |
| `get_cryptocurrency_data_by_id` | Get detailed real-time data for a specific cryptocurrency. |
| `get_cryptocurrency_news` | Get the latest news related to cryptocurrencies. |
| `calculate_investment_future_value` | Calculate the future value of an investment based on initial amount, annual return, and years. |
| `getEarningsCallTranscript` | Get the earnings call transcript of the stock with the given symbol, year, and quarter. |
| `getInsiderTransactions` | Get the insider transactions of the stock with the given symbol and year. |


## Available Prompts

- `investment_advisor`: Sets up the AI as a professional financial advisor, providing it with the necessary context and constraints to provide safe and accurate analysis.
