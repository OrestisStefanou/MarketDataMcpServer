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
	Datekey                             string  `json:"datekey" jsonschema_description:"Period end date (YYYY-MM-DD)"`
	FiscalYear                          string  `json:"fiscal_year" jsonschema_description:"Fiscal year"`
	FiscalQuarter                       string  `json:"fiscal_quarter" jsonschema_description:"Fiscal quarter"`
	Cashneq                             float64 `json:"cashneq" jsonschema_description:"Cash and cash equivalents (zero value means not known)"`
	Investmentsc                        float64 `json:"investmentsc" jsonschema_description:"Short term investments (zero value means not known)"`
	Totalcash                           float64 `json:"totalcash" jsonschema_description:"Total cash and short term investments (zero value means not known)"`
	AccountsReceivable                  float64 `json:"accounts_receivable" jsonschema_description:"Accounts receivable (zero value means not known)"`
	OtherReceivables                    float64 `json:"other_receivables" jsonschema_description:"Other receivables (zero value means not known)"`
	Receivables                         float64 `json:"receivables" jsonschema_description:"Total trade receivables (zero value means not known)"`
	Inventory                           float64 `json:"inventory" jsonschema_description:"Inventory (zero value means not known)"`
	Othercurrent                        float64 `json:"othercurrent" jsonschema_description:"Other current assets (zero value means not known)"`
	Assetsc                             float64 `json:"assetsc" jsonschema_description:"Total current assets (zero value means not known)"`
	NetPPE                              float64 `json:"net_ppe" jsonschema_description:"Net property, plant and equipment (zero value means not known)"`
	Investmentsnc                       float64 `json:"investmentsnc" jsonschema_description:"Long term investments (zero value means not known)"`
	Goodwill                            float64 `json:"goodwill" jsonschema_description:"Goodwill (zero value means not known)"`
	OtherIntangibles                    float64 `json:"other_intangibles" jsonschema_description:"Other intangible assets (zero value means not known)"`
	Othernoncurrent                     float64 `json:"othernoncurrent" jsonschema_description:"Other long term assets (zero value means not known)"`
	Assets                              float64 `json:"assets" jsonschema_description:"Total assets (zero value means not known)"`
	AccountsPayable                     float64 `json:"accounts_payable" jsonschema_description:"Accounts payable (zero value means not known)"`
	AccruedExpenses                     float64 `json:"accrued_expenses" jsonschema_description:"Accrued expenses (zero value means not known)"`
	Debtc                               float64 `json:"debtc" jsonschema_description:"Short term debt (zero value means not known)"`
	CurrentPortDebt                     float64 `json:"current_port_debt" jsonschema_description:"Current portion of long term debt (zero value means not known)"`
	CurrentCapLeases                    float64 `json:"current_cap_leases" jsonschema_description:"Current portion of leases (zero value means not known)"`
	CurrentUnearnedRevenue              float64 `json:"current_unearned_revenue" jsonschema_description:"Unearned revenue (zero value means not known)"`
	OtherCurrentLiabilities             float64 `json:"other_current_liabilities" jsonschema_description:"Other current liabilities (zero value means not known)"`
	CurrentLiabilities                  float64 `json:"current_liabilities" jsonschema_description:"Total current liabilities (zero value means not known)"`
	Debtnc                              float64 `json:"debtnc" jsonschema_description:"Long term debt (zero value means not known)"`
	CapitalLeases                       float64 `json:"capital_leases" jsonschema_description:"Long term leases (zero value means not known)"`
	Otherliabilitiesnoncurrent          float64 `json:"otherliabilitiesnoncurrent" jsonschema_description:"Other long term liabilities (zero value means not known)"`
	TotalLongTermLiabilities            float64 `json:"total_long_term_liabilities" jsonschema_description:"Total long term liabilities (zero value means not known)"`
	Liabilities                         float64 `json:"liabilities" jsonschema_description:"Total liabilities (zero value means not known)"`
	CommonStock                         float64 `json:"common_stock" jsonschema_description:"Common stock (zero value means not known)"`
	AdditionalPaidInCapital             float64 `json:"additional_paid_in_capital" jsonschema_description:"Additional paid-in capital (zero value means not known)"`
	Retearn                             float64 `json:"retearn" jsonschema_description:"Retained earnings (zero value means not known)"`
	TreasuryStock                       float64 `json:"treasury_stock" jsonschema_description:"Treasury stock (zero value means not known)"`
	PreferredStock                      float64 `json:"preferred_stock" jsonschema_description:"Preferred stock (zero value means not known)"`
	AccumulatedOtherComprehensiveIncome float64 `json:"accumulated_other_comprehensive_income" jsonschema_description:"Accumulated other comprehensive income (zero value means not known)"`
	TotalCommonEquity                   float64 `json:"total_common_equity" jsonschema_description:"Total common shareholders equity (zero value means not known)"`
	MinorityInterest                    float64 `json:"minority_interest" jsonschema_description:"Minority interest (zero value means not known)"`
	Equity                              float64 `json:"equity" jsonschema_description:"Total equity (zero value means not known)"`
	Liabilitiesequity                   float64 `json:"liabilitiesequity" jsonschema_description:"Total liabilities and equity (zero value means not known)"`
	Debt                                float64 `json:"debt" jsonschema_description:"Total debt (zero value means not known)"`
	Netcash                             float64 `json:"netcash" jsonschema_description:"Net cash position as reported by the provider: total cash minus total debt, including long term investments where the provider treats them as cash-like (zero value means not known)"`
	Netcashpershare                     float64 `json:"netcashpershare" jsonschema_description:"Net cash per share (zero value means not known)"`
	BookValue                           float64 `json:"book_value" jsonschema_description:"Book value (zero value means not known)"`
	Bvps                                float64 `json:"bvps" jsonschema_description:"Book value per share (zero value means not known)"`
	TangibleBookValue                   float64 `json:"tangible_book_value" jsonschema_description:"Tangible book value (zero value means not known)"`
	TangibleBookValuePerShare           float64 `json:"tangible_book_value_per_share" jsonschema_description:"Tangible book value per share (zero value means not known)"`
	GrossLoans                          float64 `json:"gross_loans" jsonschema_description:"Gross loans (banks only) (zero value means not known)"`
	AllowanceForLoanLosses              float64 `json:"allowance_for_loan_losses" jsonschema_description:"Allowance for loan losses (banks only) (zero value means not known)"`
	NetLoans                            float64 `json:"net_loans" jsonschema_description:"Net loans (banks only) (zero value means not known)"`
	TotalDeposits                       float64 `json:"total_deposits" jsonschema_description:"Total deposits (banks only) (zero value means not known)"`
	InterestBearingDeposits             float64 `json:"interest_bearing_deposits" jsonschema_description:"Interest bearing deposits (banks only) (zero value means not known)"`
	NoninterestBearingDeposits          float64 `json:"noninterest_bearing_deposits" jsonschema_description:"Non-interest bearing deposits (banks only) (zero value means not known)"`
	SecuritiesAndInvestments            float64 `json:"securities_and_investments" jsonschema_description:"Securities and investments (banks only) (zero value means not known)"`
	TradingAssets                       float64 `json:"trading_assets" jsonschema_description:"Trading assets (banks only) (zero value means not known)"`
	TradingLiabilities                  float64 `json:"trading_liabilities" jsonschema_description:"Trading liabilities (banks only) (zero value means not known)"`
	ShortTermBorrowings                 float64 `json:"short_term_borrowings" jsonschema_description:"Short term borrowings (banks only) (zero value means not known)"`
	InterbankBorrowing                  float64 `json:"interbank_borrowing" jsonschema_description:"Short term interbank borrowing and repurchase agreements (banks only) (zero value means not known)"`
	InterbankLending                    float64 `json:"interbank_lending" jsonschema_description:"Short term interbank lending and reverse repurchase agreements (banks only) (zero value means not known)"`
	OtherNonEarningAssets               float64 `json:"other_non_earning_assets" jsonschema_description:"Other non-earning assets (banks only) (zero value means not known)"`
	AccruedInterestAndReceivables       float64 `json:"accrued_interest_and_receivables" jsonschema_description:"Accrued interest and accounts receivable (banks only) (zero value means not known)"`
	OtherLiabilities                    float64 `json:"other_liabilities" jsonschema_description:"Other liabilities (banks only) (zero value means not known)"`
}

type CashFlowSchema struct {
	Datekey                           string  `json:"datekey" jsonschema_description:"Period end date (YYYY-MM-DD)"`
	FiscalYear                        string  `json:"fiscal_year" jsonschema_description:"Fiscal year"`
	FiscalQuarter                     string  `json:"fiscal_quarter" jsonschema_description:"Fiscal quarter"`
	NetIncomeCF                       float64 `json:"net_income_cf" jsonschema_description:"Net income (zero value means not known)"`
	TotalDepAmorCF                    float64 `json:"total_dep_amor_cf" jsonschema_description:"Depreciation and amortization (zero value means not known)"`
	Sbcomp                            float64 `json:"sbcomp" jsonschema_description:"Stock-based compensation (zero value means not known)"`
	ChangeAR                          float64 `json:"change_ar" jsonschema_description:"Change in receivables (zero value means not known)"`
	ChangeInventory                   float64 `json:"change_inventory" jsonschema_description:"Change in inventories (zero value means not known)"`
	ChangeAP                          float64 `json:"change_ap" jsonschema_description:"Change in accounts payable (zero value means not known)"`
	ChangeAccruedExpenses             float64 `json:"change_accrued_expenses" jsonschema_description:"Change in accrued expenses (zero value means not known)"`
	ChangeIncomeTaxesPayable          float64 `json:"change_income_taxes_payable" jsonschema_description:"Change in income taxes payable (zero value means not known)"`
	ChangeUnearnedRev                 float64 `json:"change_unearned_rev" jsonschema_description:"Change in unearned revenue (zero value means not known)"`
	ChangeOtherNetOperAssets          float64 `json:"change_other_net_oper_assets" jsonschema_description:"Change in other operating activities (zero value means not known)"`
	OtherOperating                    float64 `json:"other_operating" jsonschema_description:"Other operating adjustments (zero value means not known)"`
	Ncfo                              float64 `json:"ncfo" jsonschema_description:"Net cash from operating activities (zero value means not known)"`
	Capex                             float64 `json:"capex" jsonschema_description:"Capital expenditures (zero value means not known)"`
	SaleOfPPE                         float64 `json:"sale_of_ppe" jsonschema_description:"Sale of property, plant and equipment (zero value means not known)"`
	CashAcquisition                   float64 `json:"cash_acquisition" jsonschema_description:"Cash used for acquisitions (zero value means not known)"`
	ProceedsFromDivestments           float64 `json:"proceeds_from_divestments" jsonschema_description:"Proceeds from business divestments (zero value means not known)"`
	InvestInSecurities                float64 `json:"invest_in_securities" jsonschema_description:"Purchases of investments (zero value means not known)"`
	SaleOfInvestments                 float64 `json:"sale_of_investments" jsonschema_description:"Proceeds from sale of investments (zero value means not known)"`
	OtherInvesting                    float64 `json:"other_investing" jsonschema_description:"Other investing activities (zero value means not known)"`
	Ncfi                              float64 `json:"ncfi" jsonschema_description:"Net cash from investing activities (zero value means not known)"`
	DebtIssuedShortTerm               float64 `json:"debt_issued_short_term" jsonschema_description:"Short term debt issued (zero value means not known)"`
	DebtIssuedLongTerm                float64 `json:"debt_issued_long_term" jsonschema_description:"Long term debt issued (zero value means not known)"`
	DebtRepaidShortTerm               float64 `json:"debt_repaid_short_term" jsonschema_description:"Short term debt repaid (zero value means not known)"`
	DebtRepaidLongTerm                float64 `json:"debt_repaid_long_term" jsonschema_description:"Long term debt repaid (zero value means not known)"`
	NetDebtIssuedShortTerm            float64 `json:"net_debt_issued_short_term" jsonschema_description:"Net short term debt issued (zero value means not known)"`
	NetDebtIssuedLongTerm             float64 `json:"net_debt_issued_long_term" jsonschema_description:"Net long term debt issued (zero value means not known)"`
	CommonIssued                      float64 `json:"common_issued" jsonschema_description:"Common stock issued (zero value means not known)"`
	CommonRepurchased                 float64 `json:"common_repurchased" jsonschema_description:"Common stock repurchased (zero value means not known)"`
	NetStockIssued                    float64 `json:"net_stock_issued" jsonschema_description:"Net common stock issued (zero value means not known)"`
	PreferredIssued                   float64 `json:"preferred_issued" jsonschema_description:"Preferred stock issued (zero value means not known)"`
	PreferredRepurchased              float64 `json:"preferred_repurchased" jsonschema_description:"Preferred stock repurchased (zero value means not known)"`
	NetPreferredStockIssued           float64 `json:"net_preferred_stock_issued" jsonschema_description:"Net preferred stock issued (zero value means not known)"`
	CommonDividendCF                  float64 `json:"common_dividend_cf" jsonschema_description:"Common stock dividends paid (zero value means not known)"`
	OtherFinancing                    float64 `json:"other_financing" jsonschema_description:"Other financing activities (zero value means not known)"`
	Ncff                              float64 `json:"ncff" jsonschema_description:"Net cash from financing activities (zero value means not known)"`
	ExchangeRateAdjustments           float64 `json:"exchange_rate_adjustments" jsonschema_description:"Exchange rate adjustments (zero value means not known)"`
	Ncf                               float64 `json:"ncf" jsonschema_description:"Net change in cash (zero value means not known)"`
	Fcf                               float64 `json:"fcf" jsonschema_description:"Free cash flow (zero value means not known)"`
	FcfMargin                         float64 `json:"fcf_margin" jsonschema_description:"Free cash flow margin, as a fraction (0.25 means 25%) (zero value means not known)"`
	Fcfps                             float64 `json:"fcfps" jsonschema_description:"Free cash flow per share (zero value means not known)"`
	LeveredFCF                        float64 `json:"levered_fcf" jsonschema_description:"Levered free cash flow (zero value means not known)"`
	UnleveredFCF                      float64 `json:"unlevered_fcf" jsonschema_description:"Unlevered free cash flow (zero value means not known)"`
	ProvisionForCreditLosses          float64 `json:"provision_for_credit_losses" jsonschema_description:"Provision for credit losses (banks only) (zero value means not known)"`
	NetChangeInDeposits               float64 `json:"net_change_in_deposits" jsonschema_description:"Net change in deposits (banks only) (zero value means not known)"`
	NetChangeInLoansHeldForInvestment float64 `json:"net_change_in_loans_held_for_investment" jsonschema_description:"Net change in loans held for investment (banks only) (zero value means not known)"`
	NetChangeInLoansHeldForSale       float64 `json:"net_change_in_loans_held_for_sale" jsonschema_description:"Net change in loans held for sale (banks only) (zero value means not known)"`
	NetChangeInSecurities             float64 `json:"net_change_in_securities" jsonschema_description:"Net change in securities and investments (banks only) (zero value means not known)"`
	ChangeInTradingAssets             float64 `json:"change_in_trading_assets" jsonschema_description:"Change in trading assets (banks only) (zero value means not known)"`
	ChangeInTradingLiabilities        float64 `json:"change_in_trading_liabilities" jsonschema_description:"Change in trading liabilities (banks only) (zero value means not known)"`
	ChangeInSecuritiesBorrowed        float64 `json:"change_in_securities_borrowed" jsonschema_description:"Change in securities borrowed (banks only) (zero value means not known)"`
	ChangeInAccruedInterestReceivable float64 `json:"change_in_accrued_interest_receivable" jsonschema_description:"Change in accrued interest and accounts receivable (banks only) (zero value means not known)"`
	NetChangeInInterbankBorrowing     float64 `json:"net_change_in_interbank_borrowing" jsonschema_description:"Net change in short term interbank borrowing (banks only) (zero value means not known)"`
	NetChangeInInterbankLending       float64 `json:"net_change_in_interbank_lending" jsonschema_description:"Net change in short term interbank lending (banks only) (zero value means not known)"`
}

type IncomeStatementSchema struct {
	Datekey                     string  `json:"datekey" jsonschema_description:"Period end date (YYYY-MM-DD)"`
	FiscalYear                  string  `json:"fiscal_year" jsonschema_description:"Fiscal year"`
	FiscalQuarter               string  `json:"fiscal_quarter" jsonschema_description:"Fiscal quarter"`
	Revenue                     float64 `json:"revenue" jsonschema_description:"Total revenue (zero value means not known)"`
	Cor                         float64 `json:"cor" jsonschema_description:"Cost of revenue (zero value means not known)"`
	Gp                          float64 `json:"gross_profit" jsonschema_description:"Gross profit (zero value means not known)"`
	Sgna                        float64 `json:"sgna" jsonschema_description:"Selling, general and administrative expenses (zero value means not known)"`
	Rnd                         float64 `json:"rnd" jsonschema_description:"Research and development expenses (zero value means not known)"`
	OtherOpex                   float64 `json:"other_opex" jsonschema_description:"Other operating expenses (zero value means not known)"`
	Opex                        float64 `json:"opex" jsonschema_description:"Total operating expenses (zero value means not known)"`
	Opinc                       float64 `json:"opinc" jsonschema_description:"Operating income (zero value means not known)"`
	InterestExpense             float64 `json:"interest_expense" jsonschema_description:"Interest expense (zero value means not known)"`
	InterestIncome              float64 `json:"interest_income" jsonschema_description:"Interest income (zero value means not known)"`
	OtherNonOperating           float64 `json:"other_non_operating" jsonschema_description:"Other non-operating income/expense (zero value means not known)"`
	TotalNonOperating           float64 `json:"total_non_operating" jsonschema_description:"Total non-operating income/expense (zero value means not known)"`
	Pretax                      float64 `json:"pretax" jsonschema_description:"Pre-tax income (zero value means not known)"`
	Taxexp                      float64 `json:"taxexp" jsonschema_description:"Provision for income taxes (zero value means not known)"`
	EarningsDiscontinued        float64 `json:"earnings_discontinued" jsonschema_description:"Earnings from discontinued operations (zero value means not known)"`
	MinorityInterest            float64 `json:"minority_interest" jsonschema_description:"Minority interest (zero value means not known)"`
	Netinc                      float64 `json:"netinc" jsonschema_description:"Net income attributable to common shareholders, after minority interest (zero value means not known)"`
	NetincCompany               float64 `json:"netinc_company" jsonschema_description:"Net income to company, before minority interest (zero value means not known)"`
	PreferredDividends          float64 `json:"preferred_dividends" jsonschema_description:"Net income attributable to preferred dividends (zero value means not known)"`
	SharesBasic                 float64 `json:"shares_basic" jsonschema_description:"Basic shares outstanding (zero value means not known)"`
	SharesDiluted               float64 `json:"shares_diluted" jsonschema_description:"Diluted shares outstanding (zero value means not known)"`
	EpsBasic                    float64 `json:"eps_basic" jsonschema_description:"Basic earnings per share (zero value means not known)"`
	EpsDil                      float64 `json:"eps_dil" jsonschema_description:"Diluted earnings per share (zero value means not known)"`
	Dps                         float64 `json:"dps" jsonschema_description:"Dividends per share (zero value means not known)"`
	Fcf                         float64 `json:"fcf" jsonschema_description:"Free cash flow (zero value means not known)"`
	Fcfps                       float64 `json:"fcfps" jsonschema_description:"Free cash flow per share (zero value means not known)"`
	GrossMargin                 float64 `json:"gross_margin" jsonschema_description:"Gross profit margin, as a fraction (0.25 means 25%) (zero value means not known)"`
	OperatingMargin             float64 `json:"operating_margin" jsonschema_description:"Operating margin, as a fraction (0.25 means 25%) (zero value means not known)"`
	ProfitMargin                float64 `json:"profit_margin" jsonschema_description:"Net profit margin, as a fraction (0.25 means 25%) (zero value means not known)"`
	FcfMargin                   float64 `json:"fcf_margin" jsonschema_description:"Free cash flow margin, as a fraction (0.25 means 25%) (zero value means not known)"`
	Taxrate                     float64 `json:"taxrate" jsonschema_description:"Effective tax rate, as a fraction (0.25 means 25%) (zero value means not known)"`
	Ebitda                      float64 `json:"ebitda" jsonschema_description:"EBITDA (zero value means not known)"`
	EbitdaMargin                float64 `json:"ebitda_margin" jsonschema_description:"EBITDA margin, as a fraction (0.25 means 25%) (zero value means not known)"`
	Ebit                        float64 `json:"ebit" jsonschema_description:"EBIT (zero value means not known)"`
	EbitMargin                  float64 `json:"ebit_margin" jsonschema_description:"EBIT margin, as a fraction (0.25 means 25%) (zero value means not known)"`
	DepAmor                     float64 `json:"dep_amor" jsonschema_description:"Depreciation and amortization expenses (zero value means not known)"`
	NetInterestIncomeBank       float64 `json:"net_interest_income_bank" jsonschema_description:"Net interest income (banks only) (zero value means not known)"`
	RevenuesBeforeLoanLosses    float64 `json:"revenues_before_loan_losses" jsonschema_description:"Revenue before loan losses (banks only) (zero value means not known)"`
	LoanLosses                  float64 `json:"loan_losses" jsonschema_description:"Provision for loan losses (banks only) (zero value means not known)"`
	CompensationExpenses        float64 `json:"compensation_expenses" jsonschema_description:"Compensation expenses (banks only) (zero value means not known)"`
	TotalNonInterestExpense     float64 `json:"total_non_interest_expense" jsonschema_description:"Total non-interest expense (banks only) (zero value means not known)"`
	OtherNonInterestExpenses    float64 `json:"other_non_interest_expenses" jsonschema_description:"Other non-interest expenses (banks only) (zero value means not known)"`
	NonInterestIncomeGrowthBank float64 `json:"non_interest_income_growth_bank" jsonschema_description:"Non-interest income growth (banks only) (zero value means not known)"`
	ExplorationExpenses         float64 `json:"exploration_expenses" jsonschema_description:"Exploration expenses (energy only) (zero value means not known)"`
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
					Datekey:                             balanceSheet.Datekey,
					FiscalYear:                          balanceSheet.FiscalYear,
					FiscalQuarter:                       balanceSheet.FiscalQuarter,
					Cashneq:                             balanceSheet.Cashneq,
					Investmentsc:                        balanceSheet.Investmentsc,
					Totalcash:                           balanceSheet.Totalcash,
					AccountsReceivable:                  balanceSheet.AccountsReceivable,
					OtherReceivables:                    balanceSheet.OtherReceivables,
					Receivables:                         balanceSheet.Receivables,
					Inventory:                           balanceSheet.Inventory,
					Othercurrent:                        balanceSheet.Othercurrent,
					Assetsc:                             balanceSheet.Assetsc,
					NetPPE:                              balanceSheet.NetPPE,
					Investmentsnc:                       balanceSheet.Investmentsnc,
					Goodwill:                            balanceSheet.Goodwill,
					OtherIntangibles:                    balanceSheet.OtherIntangibles,
					Othernoncurrent:                     balanceSheet.Othernoncurrent,
					Assets:                              balanceSheet.Assets,
					AccountsPayable:                     balanceSheet.AccountsPayable,
					AccruedExpenses:                     balanceSheet.AccruedExpenses,
					Debtc:                               balanceSheet.Debtc,
					CurrentPortDebt:                     balanceSheet.CurrentPortDebt,
					CurrentCapLeases:                    balanceSheet.CurrentCapLeases,
					CurrentUnearnedRevenue:              balanceSheet.CurrentUnearnedRevenue,
					OtherCurrentLiabilities:             balanceSheet.OtherCurrentLiabilities,
					CurrentLiabilities:                  balanceSheet.CurrentLiabilities,
					Debtnc:                              balanceSheet.Debtnc,
					CapitalLeases:                       balanceSheet.CapitalLeases,
					Otherliabilitiesnoncurrent:          balanceSheet.Otherliabilitiesnoncurrent,
					TotalLongTermLiabilities:            balanceSheet.TotalLongTermLiabilities,
					Liabilities:                         balanceSheet.Liabilities,
					CommonStock:                         balanceSheet.CommonStock,
					AdditionalPaidInCapital:             balanceSheet.AdditionalPaidInCapital,
					Retearn:                             balanceSheet.Retearn,
					TreasuryStock:                       balanceSheet.TreasuryStock,
					PreferredStock:                      balanceSheet.PreferredStock,
					AccumulatedOtherComprehensiveIncome: balanceSheet.AccumulatedOtherComprehensiveIncome,
					TotalCommonEquity:                   balanceSheet.TotalCommonEquity,
					MinorityInterest:                    balanceSheet.MinorityInterest,
					Equity:                              balanceSheet.Equity,
					Liabilitiesequity:                   balanceSheet.Liabilitiesequity,
					Debt:                                balanceSheet.Debt,
					Netcash:                             balanceSheet.Netcash,
					Netcashpershare:                     balanceSheet.Netcashpershare,
					BookValue:                           balanceSheet.BookValue,
					Bvps:                                balanceSheet.Bvps,
					TangibleBookValue:                   balanceSheet.TangibleBookValue,
					TangibleBookValuePerShare:           balanceSheet.TangibleBookValuePerShare,
					GrossLoans:                          balanceSheet.GrossLoans,
					AllowanceForLoanLosses:              balanceSheet.AllowanceForLoanLosses,
					NetLoans:                            balanceSheet.NetLoans,
					TotalDeposits:                       balanceSheet.TotalDeposits,
					InterestBearingDeposits:             balanceSheet.InterestBearingDeposits,
					NoninterestBearingDeposits:          balanceSheet.NoninterestBearingDeposits,
					SecuritiesAndInvestments:            balanceSheet.SecuritiesAndInvestments,
					TradingAssets:                       balanceSheet.TradingAssets,
					TradingLiabilities:                  balanceSheet.TradingLiabilities,
					ShortTermBorrowings:                 balanceSheet.ShortTermBorrowings,
					InterbankBorrowing:                  balanceSheet.InterbankBorrowing,
					InterbankLending:                    balanceSheet.InterbankLending,
					OtherNonEarningAssets:               balanceSheet.OtherNonEarningAssets,
					AccruedInterestAndReceivables:       balanceSheet.AccruedInterestAndReceivables,
					OtherLiabilities:                    balanceSheet.OtherLiabilities,
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
					Datekey:                     incomeStatement.Datekey,
					FiscalYear:                  incomeStatement.FiscalYear,
					FiscalQuarter:               incomeStatement.FiscalQuarter,
					Revenue:                     incomeStatement.Revenue,
					Cor:                         incomeStatement.Cor,
					Gp:                          incomeStatement.Gp,
					Sgna:                        incomeStatement.Sgna,
					Rnd:                         incomeStatement.Rnd,
					OtherOpex:                   incomeStatement.OtherOpex,
					Opex:                        incomeStatement.Opex,
					Opinc:                       incomeStatement.Opinc,
					InterestExpense:             incomeStatement.InterestExpense,
					InterestIncome:              incomeStatement.InterestIncome,
					OtherNonOperating:           incomeStatement.OtherNonOperating,
					TotalNonOperating:           incomeStatement.TotalNonOperating,
					Pretax:                      incomeStatement.Pretax,
					Taxexp:                      incomeStatement.Taxexp,
					EarningsDiscontinued:        incomeStatement.EarningsDiscontinued,
					MinorityInterest:            incomeStatement.MinorityInterest,
					Netinc:                      incomeStatement.Netinc,
					NetincCompany:               incomeStatement.NetincCompany,
					PreferredDividends:          incomeStatement.PreferredDividends,
					SharesBasic:                 incomeStatement.SharesBasic,
					SharesDiluted:               incomeStatement.SharesDiluted,
					EpsBasic:                    incomeStatement.EpsBasic,
					EpsDil:                      incomeStatement.EpsDil,
					Dps:                         incomeStatement.Dps,
					Fcf:                         incomeStatement.Fcf,
					Fcfps:                       incomeStatement.Fcfps,
					GrossMargin:                 incomeStatement.GrossMargin,
					OperatingMargin:             incomeStatement.OperatingMargin,
					ProfitMargin:                incomeStatement.ProfitMargin,
					FcfMargin:                   incomeStatement.FcfMargin,
					Taxrate:                     incomeStatement.Taxrate,
					Ebitda:                      incomeStatement.Ebitda,
					EbitdaMargin:                incomeStatement.EbitdaMargin,
					Ebit:                        incomeStatement.Ebit,
					EbitMargin:                  incomeStatement.EbitMargin,
					DepAmor:                     incomeStatement.DepAmor,
					NetInterestIncomeBank:       incomeStatement.NetInterestIncomeBank,
					RevenuesBeforeLoanLosses:    incomeStatement.RevenuesBeforeLoanLosses,
					LoanLosses:                  incomeStatement.LoanLosses,
					CompensationExpenses:        incomeStatement.CompensationExpenses,
					TotalNonInterestExpense:     incomeStatement.TotalNonInterestExpense,
					OtherNonInterestExpenses:    incomeStatement.OtherNonInterestExpenses,
					NonInterestIncomeGrowthBank: incomeStatement.NonInterestIncomeGrowthBank,
					ExplorationExpenses:         incomeStatement.ExplorationExpenses,
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
					Datekey:                           cashFlow.Datekey,
					FiscalYear:                        cashFlow.FiscalYear,
					FiscalQuarter:                     cashFlow.FiscalQuarter,
					NetIncomeCF:                       cashFlow.NetIncomeCF,
					TotalDepAmorCF:                    cashFlow.TotalDepAmorCF,
					Sbcomp:                            cashFlow.Sbcomp,
					ChangeAR:                          cashFlow.ChangeAR,
					ChangeInventory:                   cashFlow.ChangeInventory,
					ChangeAP:                          cashFlow.ChangeAP,
					ChangeAccruedExpenses:             cashFlow.ChangeAccruedExpenses,
					ChangeIncomeTaxesPayable:          cashFlow.ChangeIncomeTaxesPayable,
					ChangeUnearnedRev:                 cashFlow.ChangeUnearnedRev,
					ChangeOtherNetOperAssets:          cashFlow.ChangeOtherNetOperAssets,
					OtherOperating:                    cashFlow.OtherOperating,
					Ncfo:                              cashFlow.Ncfo,
					Capex:                             cashFlow.Capex,
					SaleOfPPE:                         cashFlow.SaleOfPPE,
					CashAcquisition:                   cashFlow.CashAcquisition,
					ProceedsFromDivestments:           cashFlow.ProceedsFromDivestments,
					InvestInSecurities:                cashFlow.InvestInSecurities,
					SaleOfInvestments:                 cashFlow.SaleOfInvestments,
					OtherInvesting:                    cashFlow.OtherInvesting,
					Ncfi:                              cashFlow.Ncfi,
					DebtIssuedShortTerm:               cashFlow.DebtIssuedShortTerm,
					DebtIssuedLongTerm:                cashFlow.DebtIssuedLongTerm,
					DebtRepaidShortTerm:               cashFlow.DebtRepaidShortTerm,
					DebtRepaidLongTerm:                cashFlow.DebtRepaidLongTerm,
					NetDebtIssuedShortTerm:            cashFlow.NetDebtIssuedShortTerm,
					NetDebtIssuedLongTerm:             cashFlow.NetDebtIssuedLongTerm,
					CommonIssued:                      cashFlow.CommonIssued,
					CommonRepurchased:                 cashFlow.CommonRepurchased,
					NetStockIssued:                    cashFlow.NetStockIssued,
					PreferredIssued:                   cashFlow.PreferredIssued,
					PreferredRepurchased:              cashFlow.PreferredRepurchased,
					NetPreferredStockIssued:           cashFlow.NetPreferredStockIssued,
					CommonDividendCF:                  cashFlow.CommonDividendCF,
					OtherFinancing:                    cashFlow.OtherFinancing,
					Ncff:                              cashFlow.Ncff,
					ExchangeRateAdjustments:           cashFlow.ExchangeRateAdjustments,
					Ncf:                               cashFlow.Ncf,
					Fcf:                               cashFlow.Fcf,
					FcfMargin:                         cashFlow.FcfMargin,
					Fcfps:                             cashFlow.Fcfps,
					LeveredFCF:                        cashFlow.LeveredFCF,
					UnleveredFCF:                      cashFlow.UnleveredFCF,
					ProvisionForCreditLosses:          cashFlow.ProvisionForCreditLosses,
					NetChangeInDeposits:               cashFlow.NetChangeInDeposits,
					NetChangeInLoansHeldForInvestment: cashFlow.NetChangeInLoansHeldForInvestment,
					NetChangeInLoansHeldForSale:       cashFlow.NetChangeInLoansHeldForSale,
					NetChangeInSecurities:             cashFlow.NetChangeInSecurities,
					ChangeInTradingAssets:             cashFlow.ChangeInTradingAssets,
					ChangeInTradingLiabilities:        cashFlow.ChangeInTradingLiabilities,
					ChangeInSecuritiesBorrowed:        cashFlow.ChangeInSecuritiesBorrowed,
					ChangeInAccruedInterestReceivable: cashFlow.ChangeInAccruedInterestReceivable,
					NetChangeInInterbankBorrowing:     cashFlow.NetChangeInInterbankBorrowing,
					NetChangeInInterbankLending:       cashFlow.NetChangeInInterbankLending,
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

type GetCompanyKpiMetricsRequest struct {
	StockSymbol string `json:"stock_symbol" jsonschema_description:"Symbol of the stock to get data for"`
}

type GetCompanyKpiMetricsResponse struct {
	Symbol               string              `json:"symbol" jsonschema_description:"Symbol of the stock"`
	KpiMetricsCategories []KpiMetricCategory `json:"kpi_metrics_categories" jsonschema_description:"KPI metrics categories of the stock"`
}

type KpiMetricCategory struct {
	Name    string      `json:"name" jsonschema_description:"Name of the KPI metric category"`
	Metrics []KpiMetric `json:"metrics" jsonschema_description:"Metrics of the KPI metric category"`
}

type KpiMetric struct {
	Title  string  `json:"title" jsonschema_description:"Title of the KPI metric"`
	Values []Value `json:"values" jsonschema_description:"Values of the KPI metric"`
}

type Value struct {
	Year  int     `json:"year" jsonschema_description:"Year of the value"`
	Value float64 `json:"value" jsonschema_description:"Value of the metric"`
}

type KpiMetricService interface {
	GetCompanyKpiMetrics(symbol string) (domain.CompanyKpiMetrics, error)
}

type GetCompanyKpiMetricsTool struct {
	kpiMetricService KpiMetricService
}

func NewGetCompanyKpiMetricsTool(kpiMetricService KpiMetricService) (*GetCompanyKpiMetricsTool, error) {
	if kpiMetricService == nil {
		return nil, fmt.Errorf("kpiMetricService is required")
	}
	return &GetCompanyKpiMetricsTool{kpiMetricService: kpiMetricService}, nil
}

func (t *GetCompanyKpiMetricsTool) HandleGetCompanyKpiMetrics(ctx context.Context, req mcp.CallToolRequest, args GetCompanyKpiMetricsRequest) (GetCompanyKpiMetricsResponse, error) {
	if args.StockSymbol == "" {
		return GetCompanyKpiMetricsResponse{}, fmt.Errorf("stock_symbol is required")
	}

	companyKpiMetrics, err := t.kpiMetricService.GetCompanyKpiMetrics(args.StockSymbol)
	if err != nil {
		return GetCompanyKpiMetricsResponse{}, err
	}

	// Convert domain.CompanyKpiMetrics to GetCompanyKpiMetricsResponse
	companyKpiMetricsResponse := GetCompanyKpiMetricsResponse{
		Symbol: args.StockSymbol,
	}

	for _, category := range companyKpiMetrics.KpiCategories {
		kpiCategory := KpiMetricCategory{
			Name:    category.Name,
			Metrics: make([]KpiMetric, 0, len(category.Metrics)),
		}

		for _, metric := range category.Metrics {
			kpiMetric := KpiMetric{
				Title:  metric.Title,
				Values: make([]Value, 0, len(metric.Values)),
			}

			for _, val := range metric.Values {
				year := val.Year

				var floatVal float64
				switch v := val.Value.(type) {
				case float64:
					floatVal = v
				case int:
					floatVal = float64(v)
				case int64:
					floatVal = float64(v)
				case string:
					// Attempt to parse string as float if needed, though domain values should ideally be numeric
					fmt.Sscanf(v, "%f", &floatVal)
				}

				kpiMetric.Values = append(kpiMetric.Values, Value{
					Year:  year,
					Value: floatVal,
				})
			}
			kpiCategory.Metrics = append(kpiCategory.Metrics, kpiMetric)
		}
		companyKpiMetricsResponse.KpiMetricsCategories = append(companyKpiMetricsResponse.KpiMetricsCategories, kpiCategory)
	}

	return companyKpiMetricsResponse, nil
}

func (t *GetCompanyKpiMetricsTool) GetTool() mcp.Tool {
	return mcp.NewTool("getCompanyKpiMetrics",
		mcp.WithDescription("Get the KPI metrics(revenue breakdown, revenue by geography etc) of the stock with the given symbol."),
		mcp.WithInputSchema[GetCompanyKpiMetricsRequest](),
		mcp.WithOutputSchema[GetCompanyKpiMetricsResponse](),
	)
}
