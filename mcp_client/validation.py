import asyncio
from fastmcp import Client

# ANSI color codes
GREEN = "\033[92m"
RED = "\033[91m"
RESET = "\033[0m"

client = Client("http://localhost:8080/mcp")


def print_success(tool_name: str, params: dict):
    params_str = ", ".join(f"{k}={v!r}" for k, v in params.items()) if params else "no params"
    print(f"{GREEN}[PASS] {tool_name}({params_str}){RESET}")


def print_failure(tool_name: str, params: dict, error: str):
    params_str = ", ".join(f"{k}={v!r}" for k, v in params.items()) if params else "no params"
    print(f"{RED}[FAIL] {tool_name}({params_str}) => {error}{RESET}")


def validate_list_field(data: dict, field: str, item_validator=None) -> str | None:
    """Returns error string or None if valid."""
    if field not in data:
        return f"missing field '{field}'"
    if not isinstance(data[field], list):
        return f"field '{field}' is not a list"
    if item_validator and data[field]:
        err = item_validator(data[field][0])
        if err:
            return f"field '{field}[0]': {err}"
    return None


def validate_fields(data: dict, required_fields: list[str]) -> str | None:
    for field in required_fields:
        if field not in data:
            return f"missing field '{field}'"
    return None


async def call_tool(tool_name: str, params: dict, validator):
    await asyncio.sleep(1)
    try:
        result = await client.call_tool(name=tool_name, arguments=params)
        data = result.structured_content
        if data is None:
            print_failure(tool_name, params, "structured_content is None")
            return None
        err = validator(data)
        if err:
            print_failure(tool_name, params, err)
            return None
        print_success(tool_name, params)
        return data
    except Exception as e:
        print_failure(tool_name, params, str(e))
        return None


# --- Validators ---

def validate_stock_search(data):
    err = validate_list_field(data, "search_results")
    if err:
        return err
    if data["search_results"]:
        return validate_fields(data["search_results"][0], ["symbol", "company_name"])
    return None


def validate_etf_search(data):
    err = validate_list_field(data, "search_results")
    if err:
        return err
    if data["search_results"]:
        return validate_fields(data["search_results"][0], ["symbol", "etf_name", "asset_class"])
    return None


def validate_get_etf(data):
    return validate_fields(data, ["symbol", "description", "asset_class", "aum", "top_holdings"])


def validate_get_super_investors(data):
    err = validate_list_field(data, "super_investors")
    if err:
        return err
    if data["super_investors"]:
        return validate_fields(data["super_investors"][0], ["super_investor_name"])
    return None


def validate_get_super_investor_portfolio(data):
    err = validate_list_field(data, "holdings")
    if err:
        return err
    err = validate_list_field(data, "sector_analysis")
    if err:
        return err
    if data["holdings"]:
        return validate_fields(data["holdings"][0], ["stock", "portfolio_pct"])
    return None


def validate_get_market_news(data):
    # Response format: {'stock_symbol': [...news items...]}
    if not data:
        return "empty response"
    news_list = next(iter(data.values()))
    if not isinstance(news_list, list):
        return f"expected a list of news items, got {type(news_list).__name__}"
    if news_list:
        return validate_fields(news_list[0], ["url", "title", "source"])
    return None


def validate_get_sectors(data):
    err = validate_list_field(data, "sectors")
    if err:
        return err
    if data["sectors"]:
        return validate_fields(data["sectors"][0], ["name", "url_name"])
    return None


def validate_get_sector_stocks(data):
    err = validate_list_field(data, "sector_stocks")
    if err:
        return err
    if data["sector_stocks"]:
        return validate_fields(data["sector_stocks"][0], ["symbol", "company_name"])
    return None


def validate_get_stock_overview(data):
    return validate_fields(data, ["symbol", "stock_profile", "stock_financial_ratios"])


def validate_get_stock_financials(data):
    return validate_fields(data, ["symbol", "balance_sheets", "income_statements", "cash_flows"])


def validate_get_economic_indicator_time_series(data):
    err = validate_fields(data, ["indicator_name", "interval", "unit"])
    if err:
        return err
    return validate_list_field(data, "data")


def validate_get_commodity_time_series(data):
    err = validate_fields(data, ["commodity_name", "interval", "unit"])
    if err:
        return err
    return validate_list_field(data, "data")


def validate_search_cryptocurrencies(data):
    err = validate_list_field(data, "results")
    if err:
        return err
    if data["results"]:
        return validate_fields(data["results"][0], ["id", "name", "symbol"])
    return None


def validate_get_cryptocurrency_data_by_id(data):
    return validate_fields(data, ["id", "name", "symbol", "current_usd_price"])


def validate_get_cryptocurrency_news(data):
    err = validate_list_field(data, "news")
    if err:
        return err
    if data["news"]:
        return validate_fields(data["news"][0], ["url", "title"])
    return None


def validate_calculate_investment_future_value(data):
    return validate_fields(data, ["future_value"])


def validate_get_earnings_call_transcript(data):
    err = validate_fields(data, ["symbol", "year", "quarter"])
    if err:
        return err
    return validate_list_field(data, "earnings_call_transcripts")


def validate_get_insider_transactions(data):
    err = validate_fields(data, ["symbol", "year"])
    if err:
        return err
    return validate_list_field(data, "insider_transactions")


def validate_get_company_kpi_metrics(data):
    err = validate_fields(data, ["symbol"])
    if err:
        return err
    return validate_list_field(data, "kpi_metrics_categories")


def validate_get_investing_ideas(data):
    err = validate_list_field(data, "investing_ideas")
    if err:
        return err
    if data["investing_ideas"]:
        return validate_fields(data["investing_ideas"][0], ["idea_id", "title"])
    return None


def validate_get_investing_idea_stocks(data):
    return validate_list_field(data, "stocks")


def validate_get_currency_exchange_rate(data):
    return validate_fields(data, ["from_currency", "from_currency_name", "to_currency", "to_currency_name", "rate"])


STOCK_SYMBOLS = ["MSFT", "VRTX", "JPM", "BRK.B", "CAT", "TSLA", "LIN", "GOOGL", "WELL", "SHEL", "WMT", "NEE"]
ETF_SYMBOLS = ["VOO", "IEMG", "SLV", "EWJ"]


async def main():
    async with client:
        await client.ping()

        # 1. stockSearch
        await call_tool("stockSearch", {"search_string": "Microsoft", "limit": 5}, validate_stock_search)

        # 2. etfSearch
        await call_tool("etfSearch", {"search_string": "Vanguard", "limit": 5}, validate_etf_search)

        # 3. getETF - all ETF symbols
        for symbol in ETF_SYMBOLS:
            await call_tool("getETF", {"etf_symbol": symbol}, validate_get_etf)

        # 4. getSuperInvestors
        super_investors_data = await call_tool("getSuperInvestors", {}, validate_get_super_investors)

        # 5. getSuperInvestorPortfolio - use first super investor from previous result
        super_investor_name = "Warren Buffett"
        if super_investors_data and super_investors_data.get("super_investors"):
            super_investor_name = super_investors_data["super_investors"][0]["super_investor_name"]
        await call_tool(
            "getSuperInvestorPortfolio",
            {"super_investor_name": super_investor_name},
            validate_get_super_investor_portfolio,
        )

        # 6. getMarketNews (general)
        await call_tool("getMarketNews", {}, validate_get_market_news)

        # 7. getMarketNews - all stock symbols
        for symbol in STOCK_SYMBOLS:
            await call_tool("getMarketNews", {"stock_symbol": symbol}, validate_get_market_news)

        # 8. getSectors
        sectors_data = await call_tool("getSectors", {}, validate_get_sectors)

        # 9. getSectorStocks - use url_name from getSectors result
        sector_url_name = "technology"
        if sectors_data and sectors_data.get("sectors"):
            sector_url_name = sectors_data["sectors"][0]["url_name"]
        await call_tool(
            "getSectorStocks",
            {"url_name": sector_url_name, "limit": 10},
            validate_get_sector_stocks,
        )

        # 10. getStockOverview - all stock symbols
        for symbol in STOCK_SYMBOLS:
            await call_tool("getStockOverview", {"stock_symbol": symbol}, validate_get_stock_overview)

        # 11. getStockFinancials - all stock symbols
        for symbol in STOCK_SYMBOLS:
            await call_tool(
                "getStockFinancials",
                {
                    "stock_symbol": symbol,
                    "include_balance_sheets": True,
                    "include_income_statements": True,
                    "include_cash_flows": True,
                    "limit": 4,
                },
                validate_get_stock_financials,
            )

        # 12. getEconomicIndicatorTimeSeries - Inflation
        await call_tool(
            "getEconomicIndicatorTimeSeries",
            {"indicator_name": "Inflation", "limit": 5},
            validate_get_economic_indicator_time_series,
        )

        # 13. getEconomicIndicatorTimeSeries - TreasuryYield
        await call_tool(
            "getEconomicIndicatorTimeSeries",
            {"indicator_name": "TreasuryYield", "treasury_yield_maturity": "10Y", "limit": 5},
            validate_get_economic_indicator_time_series,
        )

        # 14. getCommodityTimeSeries
        await call_tool(
            "getCommodityTimeSeries",
            {"commodity_name": "CrudeOil", "limit": 5},
            validate_get_commodity_time_series,
        )

        # 15. searchCryptocurrencies
        crypto_data = await call_tool(
            "searchCryptocurrencies",
            {"search_query": "Bitcoin", "limit": 5},
            validate_search_cryptocurrencies,
        )

        # 16. getCryptocurrencyDataById - use id from searchCryptocurrencies result
        crypto_id = "bitcoin"
        if crypto_data and crypto_data.get("results"):
            crypto_id = crypto_data["results"][0]["id"]
        await call_tool(
            "getCryptocurrencyDataById",
            {"id": crypto_id},
            validate_get_cryptocurrency_data_by_id,
        )

        # 17. getCryptocurrencyNews - BTC and ETH
        for symbol in ["BTC", "ETH"]:
            await call_tool("getCryptocurrencyNews", {"symbol": symbol}, validate_get_cryptocurrency_news)

        # 18. calculateInvestmentFutureValue
        await call_tool(
            "calculateInvestmentFutureValue",
            {"initial_investment": 10000.0, "annual_return": 8.0, "years": 10},
            validate_calculate_investment_future_value,
        )

        # 19. getEarningsCallTranscript
        await call_tool(
            "getEarningsCallTranscript",
            {"stock_symbol": "MSFT", "year": 2024, "quarter": "Q1"},
            validate_get_earnings_call_transcript,
        )

        # 20. getInsiderTransactions
        await call_tool(
            "getInsiderTransactions",
            {"stock_symbol": "TSLA", "year": 2024},
            validate_get_insider_transactions,
        )

        # 21. getCompanyKpiMetrics - all stock symbols
        for symbol in STOCK_SYMBOLS:
            await call_tool("getCompanyKpiMetrics", {"stock_symbol": symbol}, validate_get_company_kpi_metrics)

        # 22. getInvestingIdeas
        investing_ideas_data = await call_tool("getInvestingIdeas", {}, validate_get_investing_ideas)

        # 23. getInvestingIdeaStocks - use idea_id from getInvestingIdeas result
        if investing_ideas_data and investing_ideas_data.get("investing_ideas"):
            idea_id = investing_ideas_data["investing_ideas"][0]["idea_id"]
            await call_tool(
                "getInvestingIdeaStocks",
                {"idea_id": idea_id},
                validate_get_investing_idea_stocks,
            )
        else:
            # Skip with a note if no ideas available
            print(f"{RED}[SKIP] getInvestingIdeaStocks - no investing ideas available{RESET}")

        # 24. getCurrencyExchangeRate
        await call_tool(
            "getCurrencyExchangeRate",
            {"from_currency": "EUR", "to_currency": "USD"},
            validate_get_currency_exchange_rate,
        )


asyncio.run(main())
