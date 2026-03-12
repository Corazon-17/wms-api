package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"table:orders"`

	ID                int64   `bun:",pk,autoincrement"`
	OrderSN           string  `bun:"order_sn,unique"`
	ShopID            string  `bun:"shop_id"`
	MarketplaceStatus string  `bun:"marketplace_status"`
	ShippingStatus    string  `bun:"shipping_status"`
	WMSStatus         string  `bun:"wms_status"`
	TrackingNumber    string  `bun:"tracking_number"`
	TotalAmount       float64 `bun:"total_amount"`

	RawMarketplacePayload string `bun:"raw_marketplace_payload,type:jsonb"`

	CreatedAt time.Time `bun:"created_at,nullzero,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,default:current_timestamp"`
}
