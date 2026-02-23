package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	CacheTtl int // The ttl for the cache in seconds

	// Alpha Vantage configs
	AlphaVantageApiKey   string
	AlphaVantageCacheTtl int // The ttl for the alpha vantage cache in seconds

	// CoinGecko configs
	CoinGeckoApiKey   string
	CoinGeckoCacheTtl int // The ttl for the coin gecko cache in seconds

	// Investing ideas configs
	InvestingIdeasDataPath string
}

func LoadConfig() (Config, error) {
	// Load .env file if it exists, but don't fail if it's missing
	_ = godotenv.Load()

	cacheTtl, err := strconv.Atoi(getEnv("CACHE_TTL", "3600"))
	if err != nil {
		cacheTtl = 3600
	}

	alphaVantageCacheTtl, err := strconv.Atoi(getEnv("ALPHA_VANTAGE_CACHE_TTL", "3600"))
	if err != nil {
		alphaVantageCacheTtl = 3600
	}

	coinGeckoCacheTtl, err := strconv.Atoi(getEnv("COIN_GECKO_CACHE_TTL", "3600"))
	if err != nil {
		coinGeckoCacheTtl = 3600
	}

	return Config{
		CacheTtl:               cacheTtl,
		AlphaVantageApiKey:     getEnv("ALPHA_VANTAGE_API_KEY", ""),
		AlphaVantageCacheTtl:   alphaVantageCacheTtl,
		CoinGeckoApiKey:        getEnv("COIN_GECKO_API_KEY", ""),
		CoinGeckoCacheTtl:      coinGeckoCacheTtl,
		InvestingIdeasDataPath: getEnv("INVESTING_IDEAS_DATA_PATH", "static_data/investing_ideas.json"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvFloat32(key string, fallback float32) float32 {
	if value, exists := os.LookupEnv(key); exists {
		if floatValue, err := strconv.ParseFloat(value, 32); err == nil {
			return float32(floatValue)
		}
	}
	return fallback
}
