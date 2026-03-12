package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DatabaseURL string

	WMSEmail       string
	WMSPassword    string
	InternalAPIKey string

	MarketplaceBaseURL    string
	MarketplacePartnerId  string
	MarketplacePartnerKey string
	MarketplaceShopID     string
	MarketplaceRedirect   string
}

var AppConfig *Config

func Load() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}

	cfg := &Config{
		AppPort: getEnv("APP_PORT", "3000"),

		DatabaseURL: getEnv("DATABASE_URL", ""),

		WMSEmail:       getEnv("WMS_EMAIL", ""),
		WMSPassword:    getEnv("WMS_PASSWORD", ""),
		InternalAPIKey: getEnv("INTERNAL_API_KEY", ""),

		MarketplaceBaseURL:    getEnv("MARKETPLACE_BASE_URL", ""),
		MarketplacePartnerId:  getEnv("MARKETPLACE_PARTNER_ID", ""),
		MarketplacePartnerKey: getEnv("MARKETPLACE_PARTNER_KEY", ""),
		MarketplaceShopID:     getEnv("MARKETPLACE_SHOP_ID", ""),
		MarketplaceRedirect:   getEnv("MARKETPLACE_REDIRECT", ""),
	}

	AppConfig = cfg
	return cfg
}

func getEnv(key string, fallback string) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
