package dto

import (
	"wms-api/internal/domain"
	"wms-api/internal/model"
)

type OrderResponse struct {
	OrderSN           string  `json:"order_sn"`
	WMSStatus         string  `json:"wms_status"`
	MarketplaceStatus string  `json:"marketplace_status"`
	ShippingStatus    string  `json:"shipping_status"`
	TrackingNumber    string  `json:"tracking_number"`
	TotalAmount       float64 `json:"total_amount"`

	AllowedActions []string `json:"allowed_actions"`
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
	}
}
