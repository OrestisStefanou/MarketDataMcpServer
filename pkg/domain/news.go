package domain

type NewsArticle struct {
	Url    string
	Image  string
	Title  string
	Text   string
	Source string
	Time   string
}

// CryptocurrencyNews carries news for one cryptocurrency.
//
// MatchedSymbol records whether the articles are actually about the requested symbol. When
// no recent article mentions it we return general crypto news instead of nothing, and this
// flag is what stops that from being read as coverage of the symbol.
type CryptocurrencyNews struct {
	Symbol        string
	MatchedSymbol bool
	Articles      []NewsArticle
}
