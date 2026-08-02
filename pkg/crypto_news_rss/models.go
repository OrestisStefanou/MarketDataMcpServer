package cryptonewsrss

import "encoding/xml"

type rssFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Items []rssItem `xml:"item"`
	} `xml:"channel"`
}

type rssItem struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	Categories  []string `xml:"category"`
	// The namespaced element is <media:content>. Its local name is "content", which does
	// not collide with <content:encoded> because that one's local name is "encoded".
	MediaContent []struct {
		Url string `xml:"url,attr"`
	} `xml:"http://search.yahoo.com/mrss/ content"`
}

type feed struct {
	Name string
	Url  string
}

// feeds are public RSS endpoints and need no API key.
var feeds = []feed{
	// The trailing slash on the CoinDesk feed permanently redirects, so it is omitted.
	{Name: "CoinDesk", Url: "https://www.coindesk.com/arc/outboundfeeds/rss"},
	{Name: "CoinTelegraph", Url: "https://cointelegraph.com/rss"},
}
