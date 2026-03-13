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
	AllowedAction     *string   `json:"allowedAction"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func NewOrderResponse(order model.Order) OrderResponse {

	return OrderResponse{
		OrderSN:           order.OrderSN,
		WMSStatus:         order.WMSStatusID,
		MarketplaceStatus: order.MarketplaceStatus,
		ShippingStatus:    order.ShippingStatus,
		TrackingNumber:    order.TrackingNumber,
		TotalAmount:       order.TotalAmount,
		AllowedAction:     domain.AllowedAction(order.WMSStatusID),
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
	}
}

type OrderItemResponse struct {
	SKU      string  `json:"sku"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func NewOrderItemResponse(item model.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		SKU:      item.SKU,
		Quantity: item.Quantity,
		Price:    item.Price,
	}
}

func NewOrderItemResponses(items []model.OrderItem) []OrderItemResponse {

	if items == nil {
		return []OrderItemResponse{}
	}

	res := make([]OrderItemResponse, 0, len(items))

	for _, item := range items {
		res = append(res, NewOrderItemResponse(item))
	}

	return res
}

type OrderDetailResponse struct {
	OrderSN           string              `json:"orderSN"`
	WMSStatusID       string              `json:"wmsStatus"`
	MarketplaceStatus string              `json:"marketplaceStatus"`
	ShippingStatus    string              `json:"shippingStatus"`
	TrackingNumber    string              `json:"trackingNumber"`
	TotalAmount       float64             `json:"totalAmount"`
	AllowedAction     *string             `json:"allowedAction"`
	CreatedAt         time.Time           `json:"createdAt"`
	UpdatedAt         time.Time           `json:"updatedAt"`
	Items             []OrderItemResponse `json:"items"`
}

func NewOrderDetailResponse(order model.OrderDetail) OrderDetailResponse {

	return OrderDetailResponse{
		OrderSN:           order.OrderSN,
		WMSStatusID:       order.WMSStatusID,
		MarketplaceStatus: order.MarketplaceStatus,
		ShippingStatus:    order.ShippingStatus,
		TrackingNumber:    order.TrackingNumber,
		TotalAmount:       order.TotalAmount,
		AllowedAction:     domain.AllowedAction(order.WMSStatusID),
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
		Items:             NewOrderItemResponses(order.Items),
	}
}

type WMSStatusResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewWMSStatusResponse(status []model.WMSStatus) []WMSStatusResponse {

	if status == nil {
		return []WMSStatusResponse{}
	}

	res := make([]WMSStatusResponse, 0, len(status))

	for _, st := range status {
		res = append(res, WMSStatusResponse{
			ID:   st.ID,
			Name: st.Name,
		})
	}

	return res
}
