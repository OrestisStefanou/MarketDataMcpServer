package domain

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
	CHF Currency = "CHF"
	CAD Currency = "CAD"
	AUD Currency = "AUD"
)

type CurrencyName string

const (
	UnitedStatesDollar CurrencyName = "United States Dollar"
	Euro               CurrencyName = "Euro"
	BritishPound       CurrencyName = "British Pound"
	JapaneseYen        CurrencyName = "Japanese Yen"
	SwissFranc         CurrencyName = "Swiss Franc"
	CanadianDollar     CurrencyName = "Canadian Dollar"
	AustralianDollar   CurrencyName = "Australian Dollar"
)

var CurrencyCodeToNameMap = map[Currency]CurrencyName{
	USD: UnitedStatesDollar,
	EUR: Euro,
	GBP: BritishPound,
	JPY: JapaneseYen,
	CHF: SwissFranc,
	CAD: CanadianDollar,
	AUD: AustralianDollar,
}

type CurrencyExchangeRate struct {
	FromCurrency     Currency     `json:"from_currency"`
	FromCurrencyName CurrencyName `json:"from_currency_name"`
	ToCurrency       Currency     `json:"to_currency"`
	ToCurrencyName   CurrencyName `json:"to_currency_name"`
	Rate             float64      `json:"rate"`
}
