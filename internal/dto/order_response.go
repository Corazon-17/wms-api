package dto

import (
	"time"
	"wms-api/internal/domain"
	"wms-api/internal/model"
)

type OrderResponse struct {
	OrderSN           string    `json:"orderSN"`
	WMSStatus         string    `json:"wmsStatus"`
	MarketplaceStatus string    `json:"marketplaceStatus"`
	ShippingStatus    string    `json:"shippingStatus"`
	TrackingNumber    string    `json:"trackingNumber"`
	TotalAmount       float64   `json:"totalAmount"`
	AllowedActions    []string  `json:"allowedActions"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func NewOrderResponse(order model.Order) OrderResponse {

	return OrderResponse{
		OrderSN:           order.OrderSN,
		WMSStatus:         order.WMSStatus,
		MarketplaceStatus: order.MarketplaceStatus,
		ShippingStatus:    order.ShippingStatus,
		TrackingNumber:    order.TrackingNumber,
		TotalAmount:       order.TotalAmount,
		AllowedActions:    domain.AllowedActions(order.WMSStatus),
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
	}
}
