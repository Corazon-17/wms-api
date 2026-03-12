package model

import "github.com/uptrace/bun"

type OrderItem struct {
	bun.BaseModel `bun:"table:order_items"`

	ID       int64   `bun:",pk,autoincrement"`
	OrderID  int64   `bun:"order_id"`
	SKU      string  `bun:"sku"`
	Quantity int     `bun:"quantity"`
	Price    float64 `bun:"price"`
}
