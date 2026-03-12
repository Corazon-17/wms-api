package marketplace

import (
	"net/http"
	"time"

	"wms-api/internal/config"
)

type Client struct {
	baseURL    string
	partnerID  string
	partnerKey string

	shopID   string
	redirect string

	http  *http.Client
	token *TokenManager
}

func NewClient(cfg *config.Config) *Client {

	return &Client{
		baseURL:    cfg.MarketplaceBaseURL,
		partnerID:  cfg.MarketplacePartnerId,
		partnerKey: cfg.MarketplacePartnerKey,

		shopID:   cfg.MarketplaceShopID,
		redirect: cfg.MarketplaceRedirect,

		http: &http.Client{
			Timeout: 10 * time.Second,
		},

		token: &TokenManager{},
	}
}
