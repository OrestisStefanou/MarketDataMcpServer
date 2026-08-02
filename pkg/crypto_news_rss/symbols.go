package cryptonewsrss

import (
	"regexp"
	"strings"
	"sync"
)

// symbolKeywords maps a ticker to the words a story about it is likely to use. Headlines
// almost always say "Bitcoin" rather than "BTC", so matching on the symbol alone finds very
// little.
var symbolKeywords = map[string][]string{
	"BTC":   {"bitcoin", "btc"},
	"ETH":   {"ethereum", "eth", "ether"},
	"SOL":   {"solana", "sol"},
	"XRP":   {"xrp", "ripple"},
	"ADA":   {"cardano", "ada"},
	"DOGE":  {"dogecoin", "doge"},
	"BNB":   {"binance coin", "bnb"},
	"USDT":  {"tether", "usdt"},
	"USDC":  {"usdc", "circle"},
	"AVAX":  {"avalanche", "avax"},
	"DOT":   {"polkadot", "dot"},
	"LINK":  {"chainlink", "link"},
	"MATIC": {"polygon", "matic"},
	"POL":   {"polygon", "pol"},
	"LTC":   {"litecoin", "ltc"},
	"TRX":   {"tron", "trx"},
	"SHIB":  {"shiba inu", "shib"},
	"TON":   {"toncoin", "ton"},
	"XLM":   {"stellar", "xlm"},
	"UNI":   {"uniswap", "uni"},
}

var (
	patternCacheMu sync.Mutex
	patternCache   = map[string]*regexp.Regexp{}
)

// keywordPattern builds a case-insensitive, word-boundary pattern for a symbol.
//
// Word boundaries matter more than they look: a substring match for "sol" also hits
// "solution" and "console", and "dot" hits "dot com".
func keywordPattern(symbol string) *regexp.Regexp {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	patternCacheMu.Lock()
	defer patternCacheMu.Unlock()

	if pattern, ok := patternCache[symbol]; ok {
		return pattern
	}

	keywords, ok := symbolKeywords[symbol]
	if !ok {
		keywords = []string{strings.ToLower(symbol)}
	}

	quoted := make([]string, 0, len(keywords))
	for _, keyword := range keywords {
		quoted = append(quoted, regexp.QuoteMeta(keyword))
	}

	pattern := regexp.MustCompile(`(?i)\b(` + strings.Join(quoted, "|") + `)\b`)
	patternCache[symbol] = pattern

	return pattern
}
