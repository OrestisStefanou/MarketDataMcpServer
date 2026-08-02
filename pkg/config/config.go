package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// defaultSecEdgarUserAgent identifies this server to the SEC. EDGAR replies 403 to a request
// whose User-Agent does not name a contact address, and also to one containing a URL, so the
// default has to be a working email-shaped string. Operators should replace it with their own.
const defaultSecEdgarUserAgent = "MarketDataMcpServer/1.0 (contact@example.com)"

type Config struct {
	Port     string
	CacheTtl int // The ttl for the cache in seconds

	// CoinGecko configs. The api key is optional; without it CoinGecko serves the same
	// endpoints anonymously at a lower rate limit.
	CoinGeckoApiKey   string
	CoinGeckoCacheTtl int // The ttl for the coin gecko cache in seconds

	// Polymarket configs
	PolymarketCacheTtl int // The ttl for the polymarket cache in seconds

	// FRED configs
	FredCacheTtl int // The ttl for the FRED cache in seconds

	// Frankfurter configs
	FrankfurterCacheTtl int // The ttl for the currency exchange rate cache in seconds

	// SEC EDGAR configs
	SecEdgarUserAgent string
	SecEdgarCacheTtl  int // The ttl for the SEC EDGAR cache in seconds

	// Crypto news configs
	CryptoNewsCacheTtl int // The ttl for the crypto news cache in seconds
}

func LoadConfig() (Config, error) {
	// Load .env file if it exists, but don't fail if it's missing
	_ = godotenv.Load()

	cacheTtl, err := strconv.Atoi(getEnv("CACHE_TTL", "3600"))
	if err != nil {
		cacheTtl = 3600
	}

	coinGeckoCacheTtl, err := strconv.Atoi(getEnv("COIN_GECKO_CACHE_TTL", "3600"))
	if err != nil {
		coinGeckoCacheTtl = 3600
	}

	polymarketCacheTtl, err := strconv.Atoi(getEnv("POLYMARKET_CACHE_TTL", "300"))
	if err != nil {
		polymarketCacheTtl = 300
	}

	// Most FRED series publish once a day at most, and several only monthly.
	fredCacheTtl, err := strconv.Atoi(getEnv("FRED_CACHE_TTL", "86400"))
	if err != nil {
		fredCacheTtl = 86400
	}

	// The ECB publishes reference rates once per working day.
	frankfurterCacheTtl, err := strconv.Atoi(getEnv("FRANKFURTER_CACHE_TTL", "3600"))
	if err != nil {
		frankfurterCacheTtl = 3600
	}

	// A filed Form 4 never changes.
	secEdgarCacheTtl, err := strconv.Atoi(getEnv("SEC_EDGAR_CACHE_TTL", "86400"))
	if err != nil {
		secEdgarCacheTtl = 86400
	}

	// News turns over fast, so this is deliberately short.
	cryptoNewsCacheTtl, err := strconv.Atoi(getEnv("CRYPTO_NEWS_CACHE_TTL", "900"))
	if err != nil {
		cryptoNewsCacheTtl = 900
	}

	return Config{
		Port:                getEnv("PORT", "8080"),
		CacheTtl:            cacheTtl,
		CoinGeckoApiKey:     getEnv("COIN_GECKO_API_KEY", ""),
		CoinGeckoCacheTtl:   coinGeckoCacheTtl,
		PolymarketCacheTtl:  polymarketCacheTtl,
		FredCacheTtl:        fredCacheTtl,
		FrankfurterCacheTtl: frankfurterCacheTtl,
		SecEdgarUserAgent:   getEnv("SEC_EDGAR_USER_AGENT", defaultSecEdgarUserAgent),
		SecEdgarCacheTtl:    secEdgarCacheTtl,
		CryptoNewsCacheTtl:  cryptoNewsCacheTtl,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
