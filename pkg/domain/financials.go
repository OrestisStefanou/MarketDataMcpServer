package domain

// Field tags map to the JSON keys served by stockanalysis.com's
// financial statement endpoints. Keep them in sync with that payload:
// an unmatched tag silently yields a zero value rather than an error.

type BalanceSheet struct {
	Datekey                             string  `json:"datekey"`
	FiscalYear                          string  `json:"fiscalYear"`
	FiscalQuarter                       string  `json:"fiscalQuarter"`
	Cashneq                             float64 `json:"cashneq"`
	Investmentsc                        float64 `json:"investmentsc"`
	Totalcash                           float64 `json:"totalcash"`
	AccountsReceivable                  float64 `json:"balance_sheet_accounts_receivable"`
	OtherReceivables                    float64 `json:"balance_sheet_other_receivables"`
	Receivables                         float64 `json:"balance_sheet_total_trade_receivables"`
	Inventory                           float64 `json:"inventory"`
	Othercurrent                        float64 `json:"balance_sheet_other_current_assets"`
	Assetsc                             float64 `json:"assetsc"`
	NetPPE                              float64 `json:"balance_sheet_net_property_plant_and_equipment"`
	Investmentsnc                       float64 `json:"balance_sheet_long_term_investments"`
	Goodwill                            float64 `json:"balance_sheet_goodwill"`
	OtherIntangibles                    float64 `json:"otherIntangibles"`
	Othernoncurrent                     float64 `json:"balance_sheet_other_long_term_assets"`
	Assets                              float64 `json:"assets"`
	AccountsPayable                     float64 `json:"balance_sheet_accounts_payable"`
	AccruedExpenses                     float64 `json:"balance_sheet_accrued_expenses"`
	Debtc                               float64 `json:"balance_sheet_short_term_debt"`
	CurrentPortDebt                     float64 `json:"balance_sheet_current_portion_of_long_term_debt"`
	CurrentCapLeases                    float64 `json:"balance_sheet_current_portion_of_leases"`
	CurrentUnearnedRevenue              float64 `json:"balance_sheet_unearned_revenue"`
	OtherCurrentLiabilities             float64 `json:"balance_sheet_other_current_liabilities"`
	CurrentLiabilities                  float64 `json:"liabilitiesc"`
	Debtnc                              float64 `json:"balance_sheet_long_term_debt"`
	CapitalLeases                       float64 `json:"longTermLeases"`
	Otherliabilitiesnoncurrent          float64 `json:"balance_sheet_other_long_term_liabilities"`
	TotalLongTermLiabilities            float64 `json:"balance_sheet_total_long_term_liabilities"`
	Liabilities                         float64 `json:"liabilities"`
	CommonStock                         float64 `json:"balance_sheet_common_stock"`
	AdditionalPaidInCapital             float64 `json:"balance_sheet_additional_paid_in_capital"`
	Retearn                             float64 `json:"balance_sheet_retained_earnings"`
	TreasuryStock                       float64 `json:"balance_sheet_treasury_stock"`
	PreferredStock                      float64 `json:"balance_sheet_preferred_stock"`
	AccumulatedOtherComprehensiveIncome float64 `json:"balance_sheet_accumulated_other_comprehensive_income"`
	TotalCommonEquity                   float64 `json:"balance_sheet_total_common_shareholders_equity"`
	MinorityInterest                    float64 `json:"minorityinterestbs"`
	Equity                              float64 `json:"equity"`
	Liabilitiesequity                   float64 `json:"liabilitiesequity"`
	Debt                                float64 `json:"debt"`
	Netcash                             float64 `json:"netCash"`
	Netcashpershare                     float64 `json:"netCashPerShare"`
	BookValue                           float64 `json:"bookValue"`
	Bvps                                float64 `json:"bookValuePerShare"`
	TangibleBookValue                   float64 `json:"tangibleBookValue"`
	TangibleBookValuePerShare           float64 `json:"tangibleBookValuePerShare"`
	GrossLoans                          float64 `json:"balance_sheet_gross_loans"`
	AllowanceForLoanLosses              float64 `json:"balance_sheet_less_allowance_for_loan_losses"`
	NetLoans                            float64 `json:"netloans"`
	TotalDeposits                       float64 `json:"balance_sheet_total_deposits"`
	InterestBearingDeposits             float64 `json:"balance_sheet_interest_bearing_deposits"`
	NoninterestBearingDeposits          float64 `json:"balance_sheet_noninterest_bearing_deposits"`
	SecuritiesAndInvestments            float64 `json:"balance_sheet_securities_and_investments"`
	TradingAssets                       float64 `json:"balance_sheet_trading_assets"`
	TradingLiabilities                  float64 `json:"balance_sheet_trading_liabilities"`
	ShortTermBorrowings                 float64 `json:"balance_sheet_short_term_borrowings"`
	InterbankBorrowing                  float64 `json:"balance_sheet_short_term_interbank_borrowing_and_repurchase_agreements"`
	InterbankLending                    float64 `json:"balance_sheet_short_term_interbank_lending_and_reverse_repurchase_agreements"`
	OtherNonEarningAssets               float64 `json:"balance_sheet_other_non_earning_assets"`
	AccruedInterestAndReceivables       float64 `json:"balance_sheet_accrued_interest_and_accounts_receivable"`
	OtherLiabilities                    float64 `json:"balance_sheet_other_liabilities"`
}

type CashFlow struct {
	Datekey                           string  `json:"datekey"`
	FiscalYear                        string  `json:"fiscalYear"`
	FiscalQuarter                     string  `json:"fiscalQuarter"`
	NetIncomeCF                       float64 `json:"cash_flow_statement_net_income"`
	TotalDepAmorCF                    float64 `json:"cash_flow_statement_depreciation_and_amortization"`
	Sbcomp                            float64 `json:"sbcomp"`
	ChangeAR                          float64 `json:"changeInReceivables"`
	ChangeInventory                   float64 `json:"cash_flow_statement_changes_in_inventories"`
	ChangeAP                          float64 `json:"cash_flow_statement_changes_in_accounts_payable"`
	ChangeAccruedExpenses             float64 `json:"cash_flow_statement_changes_in_accrued_expenses"`
	ChangeIncomeTaxesPayable          float64 `json:"cash_flow_statement_changes_in_income_taxes_payable"`
	ChangeUnearnedRev                 float64 `json:"cash_flow_statement_changes_in_unearned_revenue"`
	ChangeOtherNetOperAssets          float64 `json:"cash_flow_statement_changes_in_other_operating_activities"`
	OtherOperating                    float64 `json:"cash_flow_statement_other_adjustments"`
	Ncfo                              float64 `json:"ncfo"`
	Capex                             float64 `json:"capex"`
	SaleOfPPE                         float64 `json:"saleofpropertyplantandequipment"`
	CashAcquisition                   float64 `json:"cashAcquisition"`
	ProceedsFromDivestments           float64 `json:"cash_flow_statement_proceeds_from_business_divestments"`
	InvestInSecurities                float64 `json:"purchasesOfInvestments"`
	SaleOfInvestments                 float64 `json:"cash_flow_statement_proceeds_from_sale_of_investments"`
	OtherInvesting                    float64 `json:"cash_flow_statement_other_investing_activities"`
	Ncfi                              float64 `json:"ncfi"`
	DebtIssuedShortTerm               float64 `json:"debtissuedshortterm"`
	DebtIssuedLongTerm                float64 `json:"debtissuedlongterm"`
	DebtRepaidShortTerm               float64 `json:"debtrepaidshortterm"`
	DebtRepaidLongTerm                float64 `json:"debtrepaidlongterm"`
	NetDebtIssuedShortTerm            float64 `json:"netdebtissuedshortterm"`
	NetDebtIssuedLongTerm             float64 `json:"netdebtissuedlongterm"`
	CommonIssued                      float64 `json:"commonissued"`
	CommonRepurchased                 float64 `json:"commonrepurchased"`
	NetStockIssued                    float64 `json:"netstockissued"`
	PreferredIssued                   float64 `json:"preferredissued"`
	PreferredRepurchased              float64 `json:"preferredrepurchased"`
	NetPreferredStockIssued           float64 `json:"netpreferredstockissued"`
	CommonDividendCF                  float64 `json:"commondividendcf"`
	OtherFinancing                    float64 `json:"otherfinancing"`
	Ncff                              float64 `json:"ncff"`
	ExchangeRateAdjustments           float64 `json:"exchangeRateAdjustments"`
	Ncf                               float64 `json:"ncf"`
	Fcf                               float64 `json:"fcf"`
	FcfMargin                         float64 `json:"fcfMargin"`
	Fcfps                             float64 `json:"fcfps"`
	LeveredFCF                        float64 `json:"leveredFCF"`
	UnleveredFCF                      float64 `json:"unleveredFCF"`
	ProvisionForCreditLosses          float64 `json:"cash_flow_statement_provision_for_credit_losses"`
	NetChangeInDeposits               float64 `json:"cash_flow_statement_net_change_in_deposits"`
	NetChangeInLoansHeldForInvestment float64 `json:"cash_flow_statement_net_change_in_loans_held_for_investment"`
	NetChangeInLoansHeldForSale       float64 `json:"cash_flow_statement_net_change_in_loans_held_for_sale"`
	NetChangeInSecurities             float64 `json:"cash_flow_statement_net_change_in_securities_and_investments"`
	ChangeInTradingAssets             float64 `json:"cash_flow_statement_changes_in_trading_assets"`
	ChangeInTradingLiabilities        float64 `json:"cash_flow_statement_changes_in_trading_liabilities"`
	ChangeInSecuritiesBorrowed        float64 `json:"cash_flow_statement_changes_in_securities_borrowed"`
	ChangeInAccruedInterestReceivable float64 `json:"cash_flow_statement_changes_in_accrued_interest_and_accounts_receivable"`
	NetChangeInInterbankBorrowing     float64 `json:"cash_flow_statement_net_change_in_short_term_interbank_borrowing_and_repurchase_agreements"`
	NetChangeInInterbankLending       float64 `json:"cash_flow_statement_net_change_in_short_term_interbank_lending_and_reverse_repurchase_agreements"`
}

type IncomeStatement struct {
	Datekey                     string  `json:"datekey"`
	FiscalYear                  string  `json:"fiscalYear"`
	FiscalQuarter               string  `json:"fiscalQuarter"`
	Revenue                     float64 `json:"revenue"`
	Cor                         float64 `json:"cor"`
	Gp                          float64 `json:"grossProfit"`
	Sgna                        float64 `json:"sgna"`
	Rnd                         float64 `json:"rnd"`
	OtherOpex                   float64 `json:"otheropex"`
	Opex                        float64 `json:"totalOperatingExpenses"`
	Opinc                       float64 `json:"operatingIncome"`
	InterestExpense             float64 `json:"income_statement_interest_expense"`
	InterestIncome              float64 `json:"interestIncome"`
	OtherNonOperating           float64 `json:"otherNonOperatingIncome"`
	TotalNonOperating           float64 `json:"totalNonOperatingIncome"`
	Pretax                      float64 `json:"pretax"`
	Taxexp                      float64 `json:"income_statement_provision_for_income_taxes"`
	EarningsDiscontinued        float64 `json:"earningsDiscontinued"`
	MinorityInterest            float64 `json:"minorityInterest"`
	Netinc                      float64 `json:"netIncome"`
	NetincCompany               float64 `json:"netincCompany"`
	PreferredDividends          float64 `json:"income_statement_net_income_attributable_to_preferred_dividends"`
	SharesBasic                 float64 `json:"sharesBasic"`
	SharesDiluted               float64 `json:"sharesDiluted"`
	EpsBasic                    float64 `json:"epsBasic"`
	EpsDil                      float64 `json:"epsDiluted"`
	Dps                         float64 `json:"dps"`
	Fcf                         float64 `json:"fcf"`
	Fcfps                       float64 `json:"fcfps"`
	GrossMargin                 float64 `json:"grossMargin"`
	OperatingMargin             float64 `json:"operatingMargin"`
	ProfitMargin                float64 `json:"profitMargin"`
	FcfMargin                   float64 `json:"fcfMargin"`
	Taxrate                     float64 `json:"effectiveTaxRate"`
	Ebitda                      float64 `json:"ebitda"`
	EbitdaMargin                float64 `json:"ebitdaMargin"`
	Ebit                        float64 `json:"ebit"`
	EbitMargin                  float64 `json:"ebitMargin"`
	DepAmor                     float64 `json:"income_statement_depreciation_and_amortization_expenses"`
	NetInterestIncomeBank       float64 `json:"netInterestIncomeBank"`
	RevenuesBeforeLoanLosses    float64 `json:"revenuesbeforeloanlosses"`
	LoanLosses                  float64 `json:"loanlosses"`
	CompensationExpenses        float64 `json:"income_statement_compensation_expenses"`
	TotalNonInterestExpense     float64 `json:"income_statement_total_non_interest_expense"`
	OtherNonInterestExpenses    float64 `json:"income_statement_other_non_interest_expenses"`
	NonInterestIncomeGrowthBank float64 `json:"nonInterestIncomeGrowthBank"`
	ExplorationExpenses         float64 `json:"exploration_expenses"`
}

type FinancialRatios struct {
	Datekey           string  `json:"datekey"`
	FiscalYear        string  `json:"fiscalYear"`
	FiscalQuarter     string  `json:"fiscalQuarter"`
	Marketcap         float64 `json:"marketCap"`
	Ev                float64 `json:"ev"`
	LastCloseRatios   float64 `json:"lastClosePrice"`
	Pe                float64 `json:"pe"`
	PeForward         float64 `json:"peForward"`
	PegRatio          float64 `json:"pegRatio"`
	Ps                float64 `json:"ps"`
	Pb                float64 `json:"pb"`
	PtbvRatio         float64 `json:"ptbvRatio"`
	Pfcf              float64 `json:"pfcf"`
	Pocf              float64 `json:"pocf"`
	EvRevenue         float64 `json:"evrevenue"`
	EvEbitda          float64 `json:"evebitda"`
	EvEbit            float64 `json:"evebit"`
	EvFcf             float64 `json:"evfcf"`
	DebtEquity        float64 `json:"debtequity"`
	DebtEbitda        float64 `json:"debtebitda"`
	DebtFcf           float64 `json:"debtfcf"`
	NetDebtEquity     float64 `json:"netdebtequity"`
	NetDebtEbitda     float64 `json:"netdebtebitda"`
	NetDebtFcf        float64 `json:"netdebtfcf"`
	AssetTurnover     float64 `json:"assetturnover"`
	InventoryTurnover float64 `json:"inventoryTurnover"`
	QuickRatio        float64 `json:"quickRatio"`
	CurrentRatio      float64 `json:"currentratio"`
	Roe               float64 `json:"roe"`
	Roa               float64 `json:"roa"`
	Roic              float64 `json:"roic"`
	Roce              float64 `json:"roce"`
	EarningsYield     float64 `json:"earningsyield"`
	FcfYield          float64 `json:"fcfyield"`
	DividendYield     float64 `json:"dividendyield"`
	PayoutRatio       float64 `json:"payoutratio"`
	BuybackYield      float64 `json:"buybackyield"`
	TotalReturn       float64 `json:"totalreturn"`
}

type KpiMetricValue struct {
	Year  int
	Value any
}

type KpiMetric struct {
	Title  string
	Values []KpiMetricValue
}

type KpiCategory struct {
	Name    string
	Metrics []KpiMetric
}

type CompanyKpiMetrics struct {
	StockSymbol   string
	KpiCategories []KpiCategory
	Metadata      map[string]any
}
