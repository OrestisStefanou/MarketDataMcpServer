package cryptonewsrss

import (
	"testing"

	"market_data_mcp_server/pkg/domain"
)

// Substring matching would make "sol" hit "solution" and "dot" hit "dot com", so the
// keyword patterns are anchored on word boundaries.
func TestKeywordPatternRespectsWordBoundaries(t *testing.T) {
	tests := []struct {
		symbol   string
		text     string
		expected bool
	}{
		{"BTC", "Bitcoin rallies past 90k", true},
		{"BTC", "Traders moved 39,600 BTC today", true},
		{"BTC", "Ethereum staking hits a record", false},
		{"SOL", "Solana Foundation names a new CISO", true},
		{"SOL", "The solution was consoled by the console", false},
		{"DOT", "Polkadot upgrade ships", true},
		{"DOT", "Visit example dot com for details", true},
		{"DOT", "A dotted line", false},
		{"LINK", "Chainlink oracle expands", true},
		{"LINK", "The article linked above", false},
		{"ETH", "Ether outperforms this week", true},
		{"ETH", "Something ethereal happened", false},
	}

	for _, test := range tests {
		if got := keywordPattern(test.symbol).MatchString(test.text); got != test.expected {
			t.Errorf("%s against %q: expected %v, got %v", test.symbol, test.text, test.expected, got)
		}
	}
}

// An unknown symbol still has to match on the symbol itself.
func TestKeywordPatternForUnknownSymbol(t *testing.T) {
	pattern := keywordPattern("FOOBAR")
	if !pattern.MatchString("FOOBAR announces a fork") {
		t.Error("expected an unknown symbol to match on its own name")
	}
	if pattern.MatchString("Bitcoin rallies") {
		t.Error("expected an unknown symbol not to match unrelated news")
	}
}

func TestSelectForSymbolMarksAMatch(t *testing.T) {
	articles := []domain.NewsArticle{
		{Title: "Bitcoin rallies", Text: "btc is up"},
		{Title: "Solana upgrade", Text: "solana ships"},
	}

	news := selectForSymbol("BTC", articles)
	if !news.MatchedSymbol {
		t.Fatal("expected MatchedSymbol to be true")
	}
	if len(news.Articles) != 1 {
		t.Fatalf("expected 1 article, got %d", len(news.Articles))
	}
	if news.Symbol != "BTC" {
		t.Errorf("expected the symbol to be echoed back, got %q", news.Symbol)
	}
}

// When nothing matches we return general crypto news, and MatchedSymbol is what stops that
// from being reported as coverage of the requested coin.
func TestSelectForSymbolFallsBackAndSaysSo(t *testing.T) {
	articles := []domain.NewsArticle{
		{Title: "Bitcoin rallies", Text: "btc is up"},
		{Title: "Solana upgrade", Text: "solana ships"},
	}

	news := selectForSymbol("ZZZQQQ", articles)
	if news.MatchedSymbol {
		t.Fatal("expected MatchedSymbol to be false for an unmatched symbol")
	}
	if len(news.Articles) != 2 {
		t.Fatalf("expected the general feed back, got %d articles", len(news.Articles))
	}
}

func TestStripHtml(t *testing.T) {
	input := `<p style="float:right"><img src="https://example.com/a.jpg"></p><p>Bitcoin users moved 39,600 BTC &amp; more</p>`
	expected := "Bitcoin users moved 39,600 BTC & more"

	if got := stripHtml(input); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestItemImagePrefersMediaContent(t *testing.T) {
	item := rssItem{Description: `<img src="https://example.com/inline.jpg">`}
	item.MediaContent = append(item.MediaContent, struct {
		Url string `xml:"url,attr"`
	}{Url: "https://example.com/media.jpg"})

	if got := itemImage(item); got != "https://example.com/media.jpg" {
		t.Errorf("expected the media:content image, got %q", got)
	}

	// CoinTelegraph has no media:content, so the inline image is the fallback.
	inlineOnly := rssItem{Description: `<p><img src="https://example.com/inline.jpg" alt="x"></p>`}
	if got := itemImage(inlineOnly); got != "https://example.com/inline.jpg" {
		t.Errorf("expected the inline image, got %q", got)
	}

	if got := itemImage(rssItem{Description: "no image here"}); got != "" {
		t.Errorf("expected no image, got %q", got)
	}
}

// Articles from different feeds are merged and sorted lexically, which only works if their
// timestamps are normalised to a single format.
func TestParsePubDateNormalisesToRfc3339(t *testing.T) {
	if got := parsePubDate("Sat, 01 Aug 2026 22:01:57 +0000"); got != "2026-08-01T22:01:57Z" {
		t.Errorf("expected 2026-08-01T22:01:57Z, got %q", got)
	}
	if got := parsePubDate("Sun, 02 Aug 2026 08:04:56 +0200"); got != "2026-08-02T06:04:56Z" {
		t.Errorf("expected the time converted to UTC, got %q", got)
	}

	// An unparseable date is passed through rather than dropped.
	if got := parsePubDate("whenever"); got != "whenever" {
		t.Errorf("expected the raw value back, got %q", got)
	}
}
