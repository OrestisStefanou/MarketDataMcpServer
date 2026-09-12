import asyncio
from fastmcp import Client

# ANSI color codes
GREEN = "\033[92m"
RED = "\033[91m"
RESET = "\033[0m"

client = Client("http://localhost:8082/mcp")


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


async def call_tool_expect_failure(tool_name: str, params: dict):
    """For inputs the server must reject. A silent success here is the bug."""
    await asyncio.sleep(1)
    try:
        result = await client.call_tool(name=tool_name, arguments=params)
    except Exception:
        print_success(tool_name, params)
        return
    data = result.structured_content
    if data is None or result.is_error:
        print_success(tool_name, params)
        return
    print_failure(tool_name, params, f"expected the call to fail, got {data}")


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

def validate_limit(data: dict, field: str, limit: int | None) -> str | None:
    """A limit that returns limit+1 rows is an off-by-one no shape check would catch."""
    if limit is not None and len(data[field]) > limit:
        return f"returned {len(data[field])} results for a limit of {limit}"
    return None


def validate_stock_search(data, limit=None):
    err = validate_list_field(data, "search_results")
    if err:
        return err
    err = validate_limit(data, "search_results", limit)
    if err:
        return err
    if data["search_results"]:
        return validate_fields(data["search_results"][0], ["symbol", "company_name"])
    return None


def validate_etf_search(data, limit=None):
    err = validate_list_field(data, "search_results")
    if err:
        return err
    err = validate_limit(data, "search_results", limit)
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


def validate_get_sector_stocks(data, limit=None):
    err = validate_list_field(data, "sector_stocks")
    if err:
        return err
    err = validate_limit(data, "sector_stocks", limit)
    if err:
        return err
    if data["sector_stocks"]:
        return validate_fields(data["sector_stocks"][0], ["symbol", "company_name"])
    return None


def validate_get_stock_overview(data):
    return validate_fields(data, ["symbol", "stock_profile", "stock_financial_ratios"])


# An unmatched json tag in the Go structs unmarshals to zero rather than erroring, so a
# statement can look complete while its headline lines read 0. These are the lines no
# going concern reports as zero in a quarter it filed at all.
REQUIRED_NONZERO = {
    "balance_sheets": ["assets", "liabilities", "equity", "current_liabilities", "net_ppe", "retearn"],
    "income_statements": ["revenue", "gross_profit", "opinc", "netinc", "taxexp", "eps_dil"],
    "cash_flows": ["ncfo", "capex", "net_income_cf", "total_dep_amor_cf"],
}


def validate_get_stock_financials(data, limit=None):
    err = validate_fields(data, ["symbol", "balance_sheets", "income_statements", "cash_flows"])
    if err:
        return err

    for section, fields in REQUIRED_NONZERO.items():
        statements = data[section]
        if not statements:
            return f"field '{section}' is empty"
        if limit is not None and len(statements) > limit:
            return f"field '{section}' returned {len(statements)} statements for a limit of {limit}"
        # Across the returned quarters, not within one: a single quarter may legitimately
        # post a zero, but a field that is zero in every one of them is a broken mapping.
        for field in fields:
            if field not in statements[0]:
                return f"field '{section}[0]': missing field '{field}'"
            if all(not s.get(field) for s in statements):
                return f"field '{section}': '{field}' is zero in all {len(statements)} statements"
    return None


def validate_time_series_entries(data: dict, limit: int | None) -> str | None:
    """Shared checks for the FRED-backed time series tools."""
    entries = data["data"]
    if not entries:
        return "field 'data' is empty"
    if limit is not None and len(entries) > limit:
        return f"returned {len(entries)} entries for a limit of {limit}"

    err = validate_fields(entries[0], ["date", "value"])
    if err:
        return err

    # FRED publishes ascending dates but the tools present newest first. Getting this
    # backwards would silently serve the oldest observations.
    dates = [entry["date"] for entry in entries]
    if dates != sorted(dates, reverse=True):
        return f"entries are not newest first: {dates[:5]}"

    for entry in entries:
        try:
            float(entry["value"])
        except (TypeError, ValueError):
            return f"value {entry['value']!r} is not numeric"
    return None


def validate_get_economic_indicator_time_series(data, limit=None, value_range=None):
    err = validate_fields(data, ["indicator_name", "interval", "unit"])
    if err:
        return err
    err = validate_list_field(data, "data")
    if err:
        return err
    err = validate_time_series_entries(data, limit)
    if err:
        return err

    # Guards against an index-vs-percent regression: CPI is published as an index around
    # 330, but this tool has always returned an inflation rate.
    if value_range is not None:
        low, high = value_range
        newest = float(data["data"][0]["value"])
        if not low <= newest <= high:
            return f"newest value {newest} is outside the expected range {low} to {high}"
    return None


def validate_get_commodity_time_series(data, limit=None):
    err = validate_fields(data, ["commodity_name", "interval", "unit"])
    if err:
        return err
    err = validate_list_field(data, "data")
    if err:
        return err
    return validate_time_series_entries(data, limit)


def validate_search_cryptocurrencies(data, limit=None, expected_top_id=None):
    err = validate_list_field(data, "results")
    if err:
        return err
    err = validate_limit(data, "results", limit)
    if err:
        return err
    if not data["results"]:
        return "field 'results' is empty"
    err = validate_fields(data["results"][0], ["id", "name", "symbol"])
    if err:
        return err

    # /coins/list carries no market cap rank, so a symbol squatter can outrank the real
    # asset: a search for "BTC" once returned the memecoin "batcat". Shape checks alone
    # pass that, and the id then feeds getCryptocurrencyDataById as if it were correct.
    if expected_top_id is not None and data["results"][0]["id"] != expected_top_id:
        return f"expected {expected_top_id!r} ranked first, got {data['results'][0]['id']!r}"
    return None


def validate_get_cryptocurrency_data_by_id(data):
    return validate_fields(data, ["id", "name", "symbol", "current_usd_price"])


def validate_get_cryptocurrency_news(data):
    err = validate_list_field(data, "news")
    if err:
        return err
    # matched_symbol distinguishes real coverage from the general-news fallback. Without it
    # the two are indistinguishable to the caller.
    if "matched_symbol" not in data:
        return "missing field 'matched_symbol'"
    if data["news"]:
        return validate_fields(data["news"][0], ["url", "title"])
    return None


def validate_calculate_investment_future_value(data):
    return validate_fields(data, ["future_value"])


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


def validate_get_currency_exchange_rate(data, expected_rate=None):
    err = validate_fields(data, ["from_currency", "from_currency_name", "to_currency", "to_currency_name", "rate"])
    if err:
        return err

    # Frankfurter omits an unsupported symbol rather than erroring, so a rate of 0 is the
    # shape a missed presence check takes.
    rate = float(data["rate"])
    if rate <= 0:
        return f"rate {rate} is not a usable exchange rate"
    if expected_rate is not None and abs(rate - expected_rate) > 1e-9:
        return f"expected a rate of {expected_rate}, got {rate}"
    return None


def validate_get_polymarket_event_odds(data):
    err = validate_list_field(data, "events")
    if err:
        return err
    if data["events"]:
        event = data["events"][0]
        err = validate_fields(event, ["id", "slug", "title", "source", "markets"])
        if err:
            return err
        if event["markets"]:
            market = event["markets"][0]
            err = validate_fields(market, ["question", "outcomes"])
            if err:
                return err
            if market["outcomes"]:
                return validate_fields(market["outcomes"][0], ["name", "probability"])
    return None


STOCK_SYMBOLS = ["MSFT", "VRTX", "JPM", "BRK.B", "CAT", "TSLA", "LIN", "GOOGL", "WELL", "SHEL", "WMT", "NEE"]
ETF_SYMBOLS = ["VOO", "IEMG", "SLV", "EWJ"]
TREASURY_MATURITIES = ["3m", "2Y", "5Y", "10Y", "30Y"]
COMMODITIES = ["CrudeOil", "NaturalGas", "Copper", "Aluminum", "Wheat", "Corn", "Sugar", "Coffee"]

# Every series id is looked up in a map, so a typo in any single one is an otherwise silent
# 404. The same goes for a tool that fails to register.
EXPECTED_TOOLS = {
    "calculateInvestmentFutureValue",
    "etfSearch",
    "getCommodityTimeSeries",
    "getCompanyKpiMetrics",
    "getCryptocurrencyDataById",
    "getCryptocurrencyNews",
    "getCurrencyExchangeRate",
    "getETF",
    "getEconomicIndicatorTimeSeries",
    "getInsiderTransactions",
    "getMarketNews",
    "getPolymarketEventOdds",
    "getSectorStocks",
    "getSectors",
    "getStockFinancials",
    "getStockOverview",
    "getSuperInvestorPortfolio",
    "getSuperInvestors",
    "searchCryptocurrencies",
    "stockSearch",
}


async def validate_tool_set():
    tools = {tool.name for tool in await client.list_tools()}
    missing = EXPECTED_TOOLS - tools
    unexpected = tools - EXPECTED_TOOLS
    if missing or unexpected:
        print_failure("list_tools", {}, f"missing={sorted(missing)} unexpected={sorted(unexpected)}")
    else:
        print_success("list_tools", {"count": len(tools)})


async def main():
    async with client:
        await client.ping()

        # 0. the registered tool set
        await validate_tool_set()

        # 1. stockSearch
        await call_tool(
            "stockSearch",
            {"search_string": "Microsoft", "limit": 5},
            lambda data: validate_stock_search(data, limit=5),
        )

        # 2. etfSearch
        await call_tool(
            "etfSearch",
            {"search_string": "Vanguard", "limit": 5},
            lambda data: validate_etf_search(data, limit=5),
        )

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
            lambda data: validate_get_sector_stocks(data, limit=10),
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
                lambda data: validate_get_stock_financials(data, limit=4),
            )

        # 12. getEconomicIndicatorTimeSeries - Inflation must be a rate, not a CPI index
        await call_tool(
            "getEconomicIndicatorTimeSeries",
            {"indicator_name": "Inflation", "limit": 5},
            lambda data: validate_get_economic_indicator_time_series(data, limit=5, value_range=(-20, 20)),
        )

        # 13. getEconomicIndicatorTimeSeries - the remaining indicators
        for indicator in ["InterestRate", "UnemploymentRate", "RealGDP"]:
            value_range = (0, 100) if indicator != "RealGDP" else None
            await call_tool(
                "getEconomicIndicatorTimeSeries",
                {"indicator_name": indicator, "limit": 5},
                lambda data, value_range=value_range: validate_get_economic_indicator_time_series(
                    data, limit=5, value_range=value_range
                ),
            )

        # 14. getEconomicIndicatorTimeSeries - every treasury maturity
        for maturity in TREASURY_MATURITIES:
            await call_tool(
                "getEconomicIndicatorTimeSeries",
                {"indicator_name": "TreasuryYield", "treasury_yield_maturity": maturity, "limit": 5},
                lambda data: validate_get_economic_indicator_time_series(data, limit=5, value_range=(0, 25)),
            )

        # 15. getCommodityTimeSeries - every commodity
        for commodity in COMMODITIES:
            await call_tool(
                "getCommodityTimeSeries",
                {"commodity_name": commodity, "limit": 5},
                lambda data: validate_get_commodity_time_series(data, limit=5),
            )

        # A limit larger than the series must clamp rather than panic.
        await call_tool(
            "getCommodityTimeSeries",
            {"commodity_name": "Copper", "limit": 100000},
            validate_get_commodity_time_series,
        )

        # 16. searchCryptocurrencies
        crypto_data = await call_tool(
            "searchCryptocurrencies",
            {"search_query": "Bitcoin", "limit": 5},
            lambda data: validate_search_cryptocurrencies(data, limit=5, expected_top_id="bitcoin"),
        )

        # The symbol is the query a squatter is most likely to win.
        await call_tool(
            "searchCryptocurrencies",
            {"search_query": "BTC", "limit": 5},
            lambda data: validate_search_cryptocurrencies(data, limit=5, expected_top_id="bitcoin"),
        )

        # 17. getCryptocurrencyDataById - use id from searchCryptocurrencies result
        crypto_id = "bitcoin"
        if crypto_data and crypto_data.get("results"):
            crypto_id = crypto_data["results"][0]["id"]
        await call_tool(
            "getCryptocurrencyDataById",
            {"id": crypto_id},
            validate_get_cryptocurrency_data_by_id,
        )

        # 18. getCryptocurrencyNews - ZZZQQQ exercises the general-news fallback
        for symbol in ["BTC", "ETH", "SOL", "ADA", "ZZZQQQ"]:
            await call_tool("getCryptocurrencyNews", {"symbol": symbol}, validate_get_cryptocurrency_news)

        # 19. calculateInvestmentFutureValue
        await call_tool(
            "calculateInvestmentFutureValue",
            {"initial_investment": 10000.0, "annual_return": 8.0, "years": 10},
            validate_calculate_investment_future_value,
        )

        # 20. getInsiderTransactions - including a share class written with a dot
        for symbol in ["TSLA", "BRK.B"]:
            await call_tool(
                "getInsiderTransactions",
                {"stock_symbol": symbol, "year": 2024},
                validate_get_insider_transactions,
            )

        # A year older than the EDGAR submissions index must say so, not return an empty list.
        await call_tool_expect_failure("getInsiderTransactions", {"stock_symbol": "TSLA", "year": 1999})

        # 21. getCompanyKpiMetrics - all stock symbols
        for symbol in STOCK_SYMBOLS:
            await call_tool("getCompanyKpiMetrics", {"stock_symbol": symbol}, validate_get_company_kpi_metrics)

        # 22. getCurrencyExchangeRate
        await call_tool(
            "getCurrencyExchangeRate",
            {"from_currency": "EUR", "to_currency": "USD"},
            validate_get_currency_exchange_rate,
        )

        # A currency against itself never reaches the API.
        await call_tool(
            "getCurrencyExchangeRate",
            {"from_currency": "USD", "to_currency": "USD"},
            lambda data: validate_get_currency_exchange_rate(data, expected_rate=1.0),
        )

        # AED was dropped along with Alpha Vantage: the ECB does not publish it.
        await call_tool_expect_failure(
            "getCurrencyExchangeRate", {"from_currency": "AED", "to_currency": "USD"}
        )

        # 23. getPolymarketEventOdds
        await call_tool(
            "getPolymarketEventOdds",
            {"event_query": "election", "limit": 3},
            validate_get_polymarket_event_odds,
        )


asyncio.run(main())
