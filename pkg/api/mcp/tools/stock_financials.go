package tools

import (
	"context"
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

type StockFinancialsService interface {
	GetBalanceSheets(symbol string) ([]domain.BalanceSheet, error)
	GetIncomeStatements(symbol string) ([]domain.IncomeStatement, error)
	GetCashFlows(symbol string) ([]domain.CashFlow, error)
}

type GetStockFinancialsRequest struct {
	StockSymbol             string `json:"stock_symbol" jsonschema_description:"Symbol of the stock to get data for"`
	IncludeBalanceSheets    bool   `json:"include_balance_sheets" jsonschema_description:"If true, balance sheets of the stock will be included in the response."`
	IncludeIncomeStatements bool   `json:"include_income_statements" jsonschema_description:"If true, income statements of the stock will be included in the response."`
	IncludeCashFlows        bool   `json:"include_cash_flows" jsonschema_description:"If true, cash flows of the stock will be included in the response"`
	Limit                   int    `json:"limit" jsonschema_description:"How many financial statements to return starting from the most recent one. Will return all if not given or if it's zero"`
}

type BalanceSheetSchema struct {
	Datekey                        string  `json:"datekey" jsonschema_description:"Date key"`
	FiscalYear                     string  `json:"fiscal_year" jsonschema_description:"Fiscal year"`
	FiscalQuarter                  string  `json:"fiscal_quarter" jsonschema_description:"Fiscal quarter"`
	Cashneq                        float64 `json:"cashneq" jsonschema_description:"Cash and cash equivalents (zero value means not known)"`
	Investmentsc                   float64 `json:"investmentsc" jsonschema_description:"Short term investments (zero value means not known)"`
	Totalcash                      float64 `json:"totalcash" jsonschema_description:"Total cash (zero value means not known)"`
	CashGrowth                     float64 `json:"cash_growth" jsonschema_description:"Cash growth percentage (zero value means not known)"`
	AccountsReceivable             float64 `json:"accounts_receivable" jsonschema_description:"Accounts receivable (zero value means not known)"`
	OtherReceivables               float64 `json:"other_receivables" jsonschema_description:"Other receivables (zero value means not known)"`
	Receivables                    float64 `json:"receivables" jsonschema_description:"Total receivables (zero value means not known)"`
	Inventory                      float64 `json:"inventory" jsonschema_description:"Inventory (zero value means not known)"`
	RestrictedCash                 float64 `json:"restricted_cash" jsonschema_description:"Restricted cash (zero value means not known)"`
	Othercurrent                   float64 `json:"othercurrent" jsonschema_description:"Other current assets (zero value means not known)"`
	Assetsc                        float64 `json:"assetsc" jsonschema_description:"Total current assets (zero value means not known)"`
	NetPPE                         float64 `json:"net_ppe" jsonschema_description:"Net property, plant, and equipment (zero value means not known)"`
	Investmentsnc                  float64 `json:"investmentsnc" jsonschema_description:"Non-current investments (zero value means not known)"`
	Goodwill                       float64 `json:"goodwill" jsonschema_description:"Goodwill (zero value means not known)"`
	OtherIntangibles               float64 `json:"other_intangibles" jsonschema_description:"Other intangible assets (zero value means not known)"`
	Othernoncurrent                float64 `json:"othernoncurrent" jsonschema_description:"Other non-current assets (zero value means not known)"`
	Assets                         float64 `json:"assets" jsonschema_description:"Total assets (zero value means not known)"`
	AccountsPayable                float64 `json:"accounts_payable" jsonschema_description:"Accounts payable (zero value means not known)"`
	AccruedExpenses                float64 `json:"accrued_expenses" jsonschema_description:"Accrued expenses (zero value means not known)"`
	Debtc                          float64 `json:"debtc" jsonschema_description:"Current debt (zero value means not known)"`
	CurrentPortDebt                float64 `json:"current_port_debt" jsonschema_description:"Current portion of long-term debt (zero value means not known)"`
	CurrentCapLeases               float64 `json:"current_cap_leases" jsonschema_description:"Current capital leases (zero value means not known)"`
	CurrentIncomeTaxesPayable      float64 `json:"current_income_taxes_payable" jsonschema_description:"Current income taxes payable (zero value means not known)"`
	CurrentUnearnedRevenue         float64 `json:"current_unearned_revenue" jsonschema_description:"Current unearned revenue (zero value means not known)"`
	OtherCurrentLiabilities        float64 `json:"other_current_liabilities" jsonschema_description:"Other current liabilities (zero value means not known)"`
	CurrentLiabilities             float64 `json:"current_liabilities" jsonschema_description:"Total current liabilities (zero value means not known)"`
	Debtnc                         float64 `json:"debtnc" jsonschema_description:"Non-current debt (zero value means not known)"`
	CapitalLeases                  float64 `json:"capital_leases" jsonschema_description:"Long-term capital leases (zero value means not known)"`
	LongTermUnearnedRevenue        float64 `json:"long_term_unearned_revenue" jsonschema_description:"Long-term unearned revenue (zero value means not known)"`
	LongTermDeferredTaxLiabilities float64 `json:"long_term_deferred_tax_liabilities" jsonschema_description:"Long-term deferred tax liabilities (zero value means not known)"`
	Otherliabilitiesnoncurrent     float64 `json:"otherliabilitiesnoncurrent" jsonschema_description:"Other non-current liabilities (zero value means not known)"`
	Liabilities                    float64 `json:"liabilities" jsonschema_description:"Total liabilities (zero value means not known)"`
	CommonStock                    float64 `json:"common_stock" jsonschema_description:"Common stock (zero value means not known)"`
	Retearn                        float64 `json:"retearn" jsonschema_description:"Retained earnings (zero value means not known)"`
	OtherEquity                    float64 `json:"other_equity" jsonschema_description:"Other equity (zero value means not known)"`
	Equity                         float64 `json:"equity" jsonschema_description:"Total equity (zero value means not known)"`
	Liabilitiesequity              float64 `json:"liabilitiesequity" jsonschema_description:"Total liabilities and equity (zero value means not known)"`
	SharesOutFilingDate            float64 `json:"shares_out_filing_date" jsonschema_description:"Shares outstanding at filing date (zero value means not known)"`
	SharesOutTotalCommon           float64 `json:"shares_out_total_common" jsonschema_description:"Total common shares outstanding (zero value means not known)"`
	Bvps                           float64 `json:"bvps" jsonschema_description:"Book value per share (zero value means not known)"`
	TangibleBookValue              float64 `json:"tangible_book_value" jsonschema_description:"Tangible book value (zero value means not known)"`
	TangibleBookValuePerShare      float64 `json:"tangible_book_value_per_share" jsonschema_description:"Tangible book value per share (zero value means not known)"`
	Debt                           float64 `json:"debt" jsonschema_description:"Total debt (zero value means not known)"`
	Netcash                        float64 `json:"netcash" jsonschema_description:"Net cash (zero value means not known)"`
	NetCashGrowth                  float64 `json:"net_cash_growth" jsonschema_description:"Net cash growth percentage (zero value means not known)"`
	Netcashpershare                float64 `json:"netcashpershare" jsonschema_description:"Net cash per share (zero value means not known)"`
	Workingcapital                 float64 `json:"workingcapital" jsonschema_description:"Working capital (zero value means not known)"`
	Land                           float64 `json:"land" jsonschema_description:"Land (zero value means not known)"`
	Machinery                      float64 `json:"machinery" jsonschema_description:"Machinery (zero value means not known)"`
	LeaseholdImprovements          float64 `json:"leasehold_improvements" jsonschema_description:"Leasehold improvements (zero value means not known)"`
	TradingAssetSecurities         float64 `json:"trading_asset_securities" jsonschema_description:"Trading asset securities (zero value means not known)"`
}

type CashFlowSchema struct {
	Datekey                  string  `json:"datekey" jsonschema_description:"Date key"`
	FiscalYear               string  `json:"fiscal_year" jsonschema_description:"Fiscal year"`
	FiscalQuarter            string  `json:"fiscal_quarter" jsonschema_description:"Fiscal quarter"`
	NetIncomeCF              float64 `json:"net_income_cf" jsonschema_description:"Net income from cash flow (zero value means not known)"`
	TotalDepAmorCF           float64 `json:"total_dep_amor_cf" jsonschema_description:"Total depreciation and amortization (zero value means not known)"`
	Sbcomp                   float64 `json:"sbcomp" jsonschema_description:"Stock-based compensation (zero value means not known)"`
	ChangeAR                 float64 `json:"change_ar" jsonschema_description:"Change in accounts receivable (zero value means not known)"`
	ChangeInventory          float64 `json:"change_inventory" jsonschema_description:"Change in inventory (zero value means not known)"`
	ChangeAP                 float64 `json:"change_ap" jsonschema_description:"Change in accounts payable (zero value means not known)"`
	ChangeUnearnedRev        float64 `json:"change_unearned_rev" jsonschema_description:"Change in unearned revenue (zero value means not known)"`
	ChangeOtherNetOperAssets float64 `json:"change_other_net_oper_assets" jsonschema_description:"Change in other net operating assets (zero value means not known)"`
	OtherOperating           float64 `json:"other_operating" jsonschema_description:"Other operating activities (zero value means not known)"`
	Ncfo                     float64 `json:"ncfo" jsonschema_description:"Net cash from operating activities (zero value means not known)"`
	OcfGrowth                float64 `json:"ocf_growth" jsonschema_description:"Operating cash flow growth percentage (zero value means not known)"`
	Capex                    float64 `json:"capex" jsonschema_description:"Capital expenditures (zero value means not known)"`
	CashAcquisition          float64 `json:"cash_acquisition" jsonschema_description:"Cash used for acquisitions (zero value means not known)"`
	SalePurchaseIntangibles  float64 `json:"sale_purchase_intangibles" jsonschema_description:"Sale or purchase of intangible assets (zero value means not known)"`
	InvestInSecurities       float64 `json:"invest_in_securities" jsonschema_description:"Investment in securities (zero value means not known)"`
	OtherInvesting           float64 `json:"other_investing" jsonschema_description:"Other investing activities (zero value means not known)"`
	Ncfi                     float64 `json:"ncfi" jsonschema_description:"Net cash from investing activities (zero value means not known)"`
	DebtIssuedShortTerm      float64 `json:"debt_issued_short_term" jsonschema_description:"Short-term debt issued (zero value means not known)"`
	DebtIssuedLongTerm       float64 `json:"debt_issued_long_term" jsonschema_description:"Long-term debt issued (zero value means not known)"`
	DebtIssuedTotal          float64 `json:"debt_issued_total" jsonschema_description:"Total debt issued (zero value means not known)"`
	DebtRepaidShortTerm      float64 `json:"debt_repaid_short_term" jsonschema_description:"Short-term debt repaid (zero value means not known)"`
	DebtRepaidLongTerm       float64 `json:"debt_repaid_long_term" jsonschema_description:"Long-term debt repaid (zero value means not known)"`
	DebtRepaidTotal          float64 `json:"debt_repaid_total" jsonschema_description:"Total debt repaid (zero value means not known)"`
	NetDebtIssued            float64 `json:"net_debt_issued" jsonschema_description:"Net debt issued (zero value means not known)"`
	CommonIssued             float64 `json:"common_issued" jsonschema_description:"Common stock issued (zero value means not known)"`
	CommonRepurchased        float64 `json:"common_repurchased" jsonschema_description:"Common stock repurchased (zero value means not known)"`
	CommonDividendCF         float64 `json:"common_dividend_cf" jsonschema_description:"Common stock dividends paid (zero value means not known)"`
	OtherFinancing           float64 `json:"other_financing" jsonschema_description:"Other financing activities (zero value means not known)"`
	Ncff                     float64 `json:"ncff" jsonschema_description:"Net cash from financing activities (zero value means not known)"`
	Ncf                      float64 `json:"ncf" jsonschema_description:"Net change in cash (zero value means not known)"`
	Fcf                      float64 `json:"fcf" jsonschema_description:"Free cash flow (zero value means not known)"`
	FcfGrowth                float64 `json:"fcf_growth" jsonschema_description:"Free cash flow growth percentage (zero value means not known)"`
	FcfMargin                float64 `json:"fcf_margin" jsonschema_description:"Free cash flow margin percentage (zero value means not known)"`
	Fcfps                    float64 `json:"fcfps" jsonschema_description:"Free cash flow per share (zero value means not known)"`
	LeveredFCF               float64 `json:"levered_fcf" jsonschema_description:"Levered free cash flow (zero value means not known)"`
	UnleveredFCF             float64 `json:"unlevered_fcf" jsonschema_description:"Unlevered free cash flow (zero value means not known)"`
	CashInterestPaid         float64 `json:"cash_interest_paid" jsonschema_description:"Cash interest paid (zero value means not known)"`
	CashTaxesPaid            float64 `json:"cash_taxes_paid" jsonschema_description:"Cash taxes paid (zero value means not known)"`
	ChangeNetWorkingCapital  float64 `json:"change_net_working_capital" jsonschema_description:"Change in net working capital (zero value means not known)"`
}

type IncomeStatementSchema struct {
	Datekey           string  `json:"datekey" jsonschema_description:"Date key"`
	FiscalYear        string  `json:"fiscal_year" jsonschema_description:"Fiscal year"`
	FiscalQuarter     string  `json:"fiscal_quarter" jsonschema_description:"Fiscal quarter"`
	Revenue           float64 `json:"revenue" jsonschema_description:"Total revenue (zero value means not known)"`
	RevenueGrowth     float64 `json:"revenue_growth" jsonschema_description:"Revenue growth percentage (zero value means not known)"`
	Cor               float64 `json:"cor" jsonschema_description:"Cost of revenue (zero value means not known)"`
	Gp                float64 `json:"gp" jsonschema_description:"Gross profit (zero value means not known)"`
	Sgna              float64 `json:"sgna" jsonschema_description:"Selling, general and administrative expenses (zero value means not known)"`
	Rnd               float64 `json:"rnd" jsonschema_description:"Research and development expenses (zero value means not known)"`
	Opex              float64 `json:"opex" jsonschema_description:"Operating expenses (zero value means not known)"`
	Opinc             float64 `json:"opinc" jsonschema_description:"Operating income (zero value means not known)"`
	InterestExpense   float64 `json:"interest_expense" jsonschema_description:"Interest expense (zero value means not known)"`
	InterestIncome    float64 `json:"interest_income" jsonschema_description:"Interest income (zero value means not known)"`
	CurrencyGains     float64 `json:"currency_gains" jsonschema_description:"Currency gains/losses (zero value means not known)"`
	OtherNonOperating float64 `json:"other_non_operating" jsonschema_description:"Other non-operating income/expenses (zero value means not known)"`
	EbtExcl           float64 `json:"ebt_excl" jsonschema_description:"Earnings before tax excluding items (zero value means not known)"`
	GainInvestments   float64 `json:"gain_investments" jsonschema_description:"Gain on investments (zero value means not known)"`
	Pretax            float64 `json:"pretax" jsonschema_description:"Pre-tax income (zero value means not known)"`
	Taxexp            float64 `json:"taxexp" jsonschema_description:"Tax expense (zero value means not known)"`
	Netinc            float64 `json:"netinc" jsonschema_description:"Net income (zero value means not known)"`
	Netinccmn         float64 `json:"netinccmn" jsonschema_description:"Net income available to common shareholders (zero value means not known)"`
	NetIncomeGrowth   float64 `json:"net_income_growth" jsonschema_description:"Net income growth percentage (zero value means not known)"`
	SharesBasic       float64 `json:"shares_basic" jsonschema_description:"Basic shares outstanding (zero value means not known)"`
	SharesDiluted     float64 `json:"shares_diluted" jsonschema_description:"Diluted shares outstanding (zero value means not known)"`
	SharesYoY         float64 `json:"shares_yoy" jsonschema_description:"Year-over-year change in shares percentage (zero value means not known)"`
	EpsBasic          float64 `json:"eps_basic" jsonschema_description:"Basic earnings per share (zero value means not known)"`
	EpsDil            float64 `json:"eps_dil" jsonschema_description:"Diluted earnings per share (zero value means not known)"`
	EpsGrowth         float64 `json:"eps_growth" jsonschema_description:"EPS growth percentage (zero value means not known)"`
	Fcf               float64 `json:"fcf" jsonschema_description:"Free cash flow (zero value means not known)"`
	Fcfps             float64 `json:"fcfps" jsonschema_description:"Free cash flow per share (zero value means not known)"`
	Dps               float64 `json:"dps" jsonschema_description:"Dividends per share (zero value means not known)"`
	DividendGrowth    float64 `json:"dividend_growth" jsonschema_description:"Dividend growth percentage (zero value means not known)"`
	GrossMargin       float64 `json:"gross_margin" jsonschema_description:"Gross profit margin percentage (zero value means not known)"`
	OperatingMargin   float64 `json:"operating_margin" jsonschema_description:"Operating margin percentage (zero value means not known)"`
	ProfitMargin      float64 `json:"profit_margin" jsonschema_description:"Net profit margin percentage (zero value means not known)"`
	FcfMargin         float64 `json:"fcf_margin" jsonschema_description:"Free cash flow margin percentage (zero value means not known)"`
	Taxrate           float64 `json:"taxrate" jsonschema_description:"Effective tax rate percentage (zero value means not known)"`
	Ebitda            float64 `json:"ebitda" jsonschema_description:"Earnings before interest, taxes, depreciation, and amortization (zero value means not known)"`
	DepAmorEbitda     float64 `json:"dep_amor_ebitda" jsonschema_description:"Depreciation and amortization from EBITDA (zero value means not known)"`
	EbitdaMargin      float64 `json:"ebitda_margin" jsonschema_description:"EBITDA margin percentage (zero value means not known)"`
	Ebit              float64 `json:"ebit" jsonschema_description:"Earnings before interest and taxes (zero value means not known)"`
	EbitMargin        float64 `json:"ebit_margin" jsonschema_description:"EBIT margin percentage (zero value means not known)"`
	RevenueAsReported float64 `json:"revenue_as_reported" jsonschema_description:"Revenue as reported (zero value means not known)"`
	PayoutRatio       float64 `json:"payout_ratio" jsonschema_description:"Dividend payout ratio percentage (zero value means not known)"`
}

type GetStockFinancialsResponse struct {
	CurrentDate      string                  `json:"current_date"`
	Symbol           string                  `json:"symbol" jsonschema_description:"The symbol of the stock"`
	BalanceSheets    []BalanceSheetSchema    `json:"balance_sheets,omitempty" jsonschema_description:"A list with the latest quarterly balance sheets of the stock company"`
	IncomeStatements []IncomeStatementSchema `json:"income_statements,omitempty" jsonschema_description:"A list with the latest quarterly income statements of the stock company"`
	CashFlows        []CashFlowSchema        `json:"cash_flows,omitempty" jsonschema_description:"A list with the latest quarterly cash flows of the stock company"`
}

type GetStockFinancialsTool struct {
	stockFinancialsService StockFinancialsService
}

func NewGetStockFinancialsTool(stockFinancialsService StockFinancialsService) (*GetStockFinancialsTool, error) {
	return &GetStockFinancialsTool{
		stockFinancialsService: stockFinancialsService,
	}, nil
}

func (t *GetStockFinancialsTool) HandleGetStockFinancials(ctx context.Context, req mcp.CallToolRequest, args GetStockFinancialsRequest) (GetStockFinancialsResponse, error) {
	if args.StockSymbol == "" {
		return GetStockFinancialsResponse{}, fmt.Errorf("stock_symbol is required")
	}
	stockSymbol := strings.ToLower(args.StockSymbol)
	var balanceSheetsResponse []BalanceSheetSchema
	var incomeStatementsResponse []IncomeStatementSchema
	var cashFlowsResponse []CashFlowSchema

	if args.IncludeBalanceSheets {
		balanceSheets, err := t.stockFinancialsService.GetBalanceSheets(stockSymbol)
		if err != nil {
			return GetStockFinancialsResponse{}, err
		}

		var limit int
		if args.Limit > 0 {
			limit = args.Limit
		} else {
			limit = len(balanceSheets)
		}

		balanceSheetsResponse = make([]BalanceSheetSchema, 0, limit)
		for i, balanceSheet := range balanceSheets {
			if i == limit {
				break
			}
			balanceSheetsResponse = append(
				balanceSheetsResponse,
				BalanceSheetSchema{
					Datekey:                        balanceSheet.Datekey,
					FiscalYear:                     balanceSheet.FiscalYear,
					FiscalQuarter:                  balanceSheet.FiscalQuarter,
					Cashneq:                        balanceSheet.Cashneq,
					Investmentsc:                   balanceSheet.Investmentsc,
					Totalcash:                      balanceSheet.Totalcash,
					CashGrowth:                     balanceSheet.CashGrowth,
					AccountsReceivable:             balanceSheet.AccountsReceivable,
					OtherReceivables:               balanceSheet.OtherReceivables,
					Receivables:                    balanceSheet.Receivables,
					Inventory:                      balanceSheet.Inventory,
					RestrictedCash:                 balanceSheet.RestrictedCash,
					Othercurrent:                   balanceSheet.Othercurrent,
					Assetsc:                        balanceSheet.Assetsc,
					NetPPE:                         balanceSheet.NetPPE,
					Investmentsnc:                  balanceSheet.Investmentsnc,
					Goodwill:                       balanceSheet.Goodwill,
					OtherIntangibles:               balanceSheet.OtherIntangibles,
					Othernoncurrent:                balanceSheet.Othernoncurrent,
					Assets:                         balanceSheet.Assets,
					AccountsPayable:                balanceSheet.AccountsPayable,
					AccruedExpenses:                balanceSheet.AccruedExpenses,
					Debtc:                          balanceSheet.Debtc,
					CurrentPortDebt:                balanceSheet.CurrentPortDebt,
					CurrentCapLeases:               balanceSheet.CurrentCapLeases,
					CurrentIncomeTaxesPayable:      balanceSheet.CurrentIncomeTaxesPayable,
					CurrentUnearnedRevenue:         balanceSheet.CurrentUnearnedRevenue,
					OtherCurrentLiabilities:        balanceSheet.OtherCurrentLiabilities,
					CurrentLiabilities:             balanceSheet.CurrentLiabilities,
					Debtnc:                         balanceSheet.Debtnc,
					CapitalLeases:                  balanceSheet.CapitalLeases,
					LongTermUnearnedRevenue:        balanceSheet.LongTermUnearnedRevenue,
					LongTermDeferredTaxLiabilities: balanceSheet.LongTermDeferredTaxLiabilities,
					Otherliabilitiesnoncurrent:     balanceSheet.Otherliabilitiesnoncurrent,
					Liabilities:                    balanceSheet.Liabilities,
					CommonStock:                    balanceSheet.CommonStock,
					Retearn:                        balanceSheet.Retearn,
					OtherEquity:                    balanceSheet.OtherEquity,
					Equity:                         balanceSheet.Equity,
					Liabilitiesequity:              balanceSheet.Liabilitiesequity,
					SharesOutFilingDate:            balanceSheet.SharesOutFilingDate,
					SharesOutTotalCommon:           balanceSheet.SharesOutTotalCommon,
					Bvps:                           balanceSheet.Bvps,
					TangibleBookValue:              balanceSheet.TangibleBookValue,
					TangibleBookValuePerShare:      balanceSheet.TangibleBookValuePerShare,
					Debt:                           balanceSheet.Debt,
					Netcash:                        balanceSheet.Netcash,
					NetCashGrowth:                  balanceSheet.NetCashGrowth,
					Netcashpershare:                balanceSheet.Netcashpershare,
					Workingcapital:                 balanceSheet.Workingcapital,
					Land:                           balanceSheet.Land,
					Machinery:                      balanceSheet.Machinery,
					LeaseholdImprovements:          balanceSheet.LeaseholdImprovements,
					TradingAssetSecurities:         balanceSheet.TradingAssetSecurities,
				},
			)
		}
	}

	if args.IncludeIncomeStatements {
		incomeStatements, err := t.stockFinancialsService.GetIncomeStatements(stockSymbol)
		if err != nil {
			return GetStockFinancialsResponse{}, err
		}

		var limit int
		if args.Limit > 0 {
			limit = args.Limit
		} else {
			limit = len(incomeStatements)
		}

		incomeStatementsResponse = make([]IncomeStatementSchema, 0, limit)
		for i, incomeStatement := range incomeStatements {
			if i == limit {
				break
			}
			incomeStatementsResponse = append(
				incomeStatementsResponse,
				IncomeStatementSchema{
					Datekey:           incomeStatement.Datekey,
					FiscalYear:        incomeStatement.FiscalYear,
					FiscalQuarter:     incomeStatement.FiscalQuarter,
					Revenue:           incomeStatement.Revenue,
					RevenueGrowth:     incomeStatement.RevenueGrowth,
					Cor:               incomeStatement.Cor,
					Gp:                incomeStatement.Gp,
					Sgna:              incomeStatement.Sgna,
					Rnd:               incomeStatement.Rnd,
					Opex:              incomeStatement.Opex,
					Opinc:             incomeStatement.Opinc,
					InterestExpense:   incomeStatement.InterestExpense,
					InterestIncome:    incomeStatement.InterestIncome,
					CurrencyGains:     incomeStatement.CurrencyGains,
					OtherNonOperating: incomeStatement.OtherNonOperating,
					EbtExcl:           incomeStatement.EbtExcl,
					GainInvestments:   incomeStatement.GainInvestments,
					Pretax:            incomeStatement.Pretax,
					Taxexp:            incomeStatement.Taxexp,
					Netinc:            incomeStatement.Netinc,
					Netinccmn:         incomeStatement.Netinccmn,
					NetIncomeGrowth:   incomeStatement.NetIncomeGrowth,
					SharesBasic:       incomeStatement.SharesBasic,
					SharesDiluted:     incomeStatement.SharesDiluted,
					SharesYoY:         incomeStatement.SharesYoY,
					EpsBasic:          incomeStatement.EpsBasic,
					EpsDil:            incomeStatement.EpsDil,
					EpsGrowth:         incomeStatement.EpsGrowth,
					Fcf:               incomeStatement.Fcf,
					Fcfps:             incomeStatement.Fcfps,
					Dps:               incomeStatement.Dps,
					DividendGrowth:    incomeStatement.DividendGrowth,
					GrossMargin:       incomeStatement.GrossMargin,
					OperatingMargin:   incomeStatement.OperatingMargin,
					ProfitMargin:      incomeStatement.ProfitMargin,
					FcfMargin:         incomeStatement.FcfMargin,
					Taxrate:           incomeStatement.Taxrate,
					Ebitda:            incomeStatement.Ebitda,
					DepAmorEbitda:     incomeStatement.DepAmorEbitda,
					EbitdaMargin:      incomeStatement.EbitdaMargin,
					Ebit:              incomeStatement.Ebit,
					EbitMargin:        incomeStatement.EbitMargin,
					RevenueAsReported: incomeStatement.RevenueAsReported,
					PayoutRatio:       incomeStatement.PayoutRatio,
				},
			)
		}
	}

	if args.IncludeCashFlows {
		cashFlows, err := t.stockFinancialsService.GetCashFlows(stockSymbol)
		if err != nil {
			return GetStockFinancialsResponse{}, err
		}

		var limit int
		if args.Limit > 0 {
			limit = args.Limit
		} else {
			limit = len(cashFlows)
		}

		cashFlowsResponse = make([]CashFlowSchema, 0, limit)
		for i, cashFlow := range cashFlows {
			if i == limit {
				break
			}
			cashFlowsResponse = append(
				cashFlowsResponse,
				CashFlowSchema{
					Datekey:                  cashFlow.Datekey,
					FiscalYear:               cashFlow.FiscalYear,
					FiscalQuarter:            cashFlow.FiscalQuarter,
					NetIncomeCF:              cashFlow.NetIncomeCF,
					TotalDepAmorCF:           cashFlow.TotalDepAmorCF,
					Sbcomp:                   cashFlow.Sbcomp,
					ChangeAR:                 cashFlow.ChangeAR,
					ChangeInventory:          cashFlow.ChangeInventory,
					ChangeAP:                 cashFlow.ChangeAP,
					ChangeUnearnedRev:        cashFlow.ChangeUnearnedRev,
					ChangeOtherNetOperAssets: cashFlow.ChangeOtherNetOperAssets,
					OtherOperating:           cashFlow.OtherOperating,
					Ncfo:                     cashFlow.Ncfo,
					OcfGrowth:                cashFlow.OcfGrowth,
					Capex:                    cashFlow.Capex,
					CashAcquisition:          cashFlow.CashAcquisition,
					SalePurchaseIntangibles:  cashFlow.SalePurchaseIntangibles,
					InvestInSecurities:       cashFlow.InvestInSecurities,
					OtherInvesting:           cashFlow.OtherInvesting,
					Ncfi:                     cashFlow.Ncfi,
					DebtIssuedShortTerm:      cashFlow.DebtIssuedShortTerm,
					DebtIssuedLongTerm:       cashFlow.DebtIssuedLongTerm,
					DebtIssuedTotal:          cashFlow.DebtIssuedTotal,
					DebtRepaidShortTerm:      cashFlow.DebtRepaidShortTerm,
					DebtRepaidLongTerm:       cashFlow.DebtRepaidLongTerm,
					DebtRepaidTotal:          cashFlow.DebtRepaidTotal,
					NetDebtIssued:            cashFlow.NetDebtIssued,
					CommonIssued:             cashFlow.CommonIssued,
					CommonRepurchased:        cashFlow.CommonRepurchased,
					CommonDividendCF:         cashFlow.CommonDividendCF,
					OtherFinancing:           cashFlow.OtherFinancing,
					Ncff:                     cashFlow.Ncff,
					Ncf:                      cashFlow.Ncf,
					Fcf:                      cashFlow.Fcf,
					FcfGrowth:                cashFlow.FcfGrowth,
					FcfMargin:                cashFlow.FcfMargin,
					Fcfps:                    cashFlow.Fcfps,
					LeveredFCF:               cashFlow.LeveredFCF,
					UnleveredFCF:             cashFlow.UnleveredFCF,
					CashInterestPaid:         cashFlow.CashInterestPaid,
					CashTaxesPaid:            cashFlow.CashTaxesPaid,
					ChangeNetWorkingCapital:  cashFlow.ChangeNetWorkingCapital,
				},
			)
		}
	}

	return GetStockFinancialsResponse{
		Symbol:           args.StockSymbol,
		CurrentDate:      time.Now().Format("2006-01-02"),
		BalanceSheets:    balanceSheetsResponse,
		IncomeStatements: incomeStatementsResponse,
		CashFlows:        cashFlowsResponse,
	}, nil
}

func (t *GetStockFinancialsTool) GetTool() mcp.Tool {
	return mcp.NewTool("getStockFinancials",
		mcp.WithDescription("Get the financials(balance sheets, income statements, cash flows) of the stock with the given symbol."),
		mcp.WithInputSchema[GetStockFinancialsRequest](),
		mcp.WithOutputSchema[GetStockFinancialsResponse](),
	)
}
