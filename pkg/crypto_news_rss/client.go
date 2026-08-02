package cryptonewsrss

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/errors"
)

const (
	// maxArticles bounds the merged feed, and maxFallbackArticles the general crypto news
	// returned when a symbol matches nothing.
	maxArticles         = 40
	maxFallbackArticles = 20
)

var (
	httpClient   = &http.Client{Timeout: 10 * time.Second}
	htmlTagsRe   = regexp.MustCompile(`<[^>]*>`)
	imgSrcRe     = regexp.MustCompile(`(?i)<img[^>]+src="([^"]+)"`)
	whitespaceRe = regexp.MustCompile(`\s+`)
)

// CryptoNewsRssClient aggregates public crypto news RSS feeds. No API key is involved.
type CryptoNewsRssClient struct{}

func NewCryptoNewsRssClient() (*CryptoNewsRssClient, error) {
	return &CryptoNewsRssClient{}, nil
}

// GetCryptocurrencyNews returns recent articles about the given symbol.
//
// When nothing matches the symbol it falls back to general crypto news rather than an empty
// list, since on a quiet news day an empty response is indistinguishable from a broken tool.
// The fallback is reported through MatchedSymbol so it cannot be mistaken for coverage of
// the symbol itself.
func (c *CryptoNewsRssClient) GetCryptocurrencyNews(symbol string) (domain.CryptocurrencyNews, error) {
	articles, err := c.GetLatestNews()
	if err != nil {
		return domain.CryptocurrencyNews{}, err
	}

	return selectForSymbol(symbol, articles), nil
}

// selectForSymbol filters the merged feed down to one symbol, falling back to general news.
func selectForSymbol(symbol string, articles []domain.NewsArticle) domain.CryptocurrencyNews {
	pattern := keywordPattern(symbol)

	matched := make([]domain.NewsArticle, 0, len(articles))
	for _, article := range articles {
		if pattern.MatchString(article.Title) || pattern.MatchString(article.Text) {
			matched = append(matched, article)
		}
	}

	if len(matched) > 0 {
		return domain.CryptocurrencyNews{Symbol: symbol, MatchedSymbol: true, Articles: matched}
	}

	if len(articles) > maxFallbackArticles {
		articles = articles[:maxFallbackArticles]
	}

	return domain.CryptocurrencyNews{Symbol: symbol, MatchedSymbol: false, Articles: articles}
}

// GetLatestNews fetches every configured feed concurrently and merges them newest first.
// A feed that fails is skipped; the call only fails if every feed fails.
func (c *CryptoNewsRssClient) GetLatestNews() ([]domain.NewsArticle, error) {
	results := make([][]domain.NewsArticle, len(feeds))
	failures := make([]error, len(feeds))

	var waitGroup sync.WaitGroup
	for i, source := range feeds {
		waitGroup.Add(1)
		go func(index int, source feed) {
			defer waitGroup.Done()
			articles, err := fetchFeed(source)
			if err != nil {
				failures[index] = err
				return
			}
			results[index] = articles
		}(i, source)
	}
	waitGroup.Wait()

	articles := make([]domain.NewsArticle, 0, maxArticles*2)
	succeeded := 0
	for i := range feeds {
		if failures[i] != nil {
			continue
		}
		succeeded++
		articles = append(articles, results[i]...)
	}

	if succeeded == 0 {
		return nil, fmt.Errorf("every crypto news feed failed, most recently: %v", failures[len(failures)-1])
	}

	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].Time > articles[j].Time
	})

	if len(articles) > maxArticles {
		articles = articles[:maxArticles]
	}

	return articles, nil
}

func fetchFeed(source feed) ([]domain.NewsArticle, error) {
	resp, err := httpClient.Get(source.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("Call to get the %s news feed failed", source.Name),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.StreamError{
			Message: fmt.Sprintf("failed to read the %s news feed", source.Name),
			Err:     err,
		}
	}

	var parsed rssFeed
	if err := xml.Unmarshal(body, &parsed); err != nil {
		return nil, errors.StreamError{
			Message: fmt.Sprintf("failed to parse the %s news feed", source.Name),
			Err:     err,
		}
	}

	articles := make([]domain.NewsArticle, 0, len(parsed.Channel.Items))
	for _, item := range parsed.Channel.Items {
		articles = append(articles, domain.NewsArticle{
			Url:    strings.TrimSpace(item.Link),
			Image:  itemImage(item),
			Title:  strings.TrimSpace(item.Title),
			Text:   stripHtml(item.Description),
			Source: source.Name,
			Time:   parsePubDate(item.PubDate),
		})
	}

	return articles, nil
}

// itemImage prefers the media:content element, falling back to the first inline image in
// the description. CoinDesk uses the former, CoinTelegraph the latter.
func itemImage(item rssItem) string {
	for _, media := range item.MediaContent {
		if url := strings.TrimSpace(media.Url); url != "" {
			return url
		}
	}

	if match := imgSrcRe.FindStringSubmatch(item.Description); len(match) == 2 {
		return match[1]
	}

	return ""
}

// stripHtml reduces an RSS description to plain text. Descriptions arrive as CDATA-wrapped
// markup, which is noise in an LLM context.
func stripHtml(value string) string {
	value = htmlTagsRe.ReplaceAllString(value, " ")
	value = strings.NewReplacer("&amp;", "&", "&lt;", "<", "&gt;", ">", "&quot;", `"`, "&#39;", "'", "&nbsp;", " ").Replace(value)
	return strings.TrimSpace(whitespaceRe.ReplaceAllString(value, " "))
}

// parsePubDate normalises the RFC 1123 dates RSS uses into RFC 3339, so that articles from
// different feeds sort against each other lexically.
func parsePubDate(value string) string {
	value = strings.TrimSpace(value)

	for _, layout := range []string{time.RFC1123Z, time.RFC1123, time.RFC822Z, time.RFC822} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.UTC().Format(time.RFC3339)
		}
	}

	return value
}
