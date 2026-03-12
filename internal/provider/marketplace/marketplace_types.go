package marketplace

type TokenResponse struct {
	Message string `json:"message"`
	Data    struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	} `json:"data"`
}

type OrderListResponse struct {
	Message string  `json:"message"`
	Data    []Order `json:"data"`
}

type Order struct {
	OrderSN        string      `json:"order_sn"`
	ShopID         string      `json:"shop_id"`
	Status         string      `json:"status"`
	ShippingStatus string      `json:"shipping_status"`
	Items          []OrderItem `json:"items"`
	TotalAmount    float64     `json:"total_amount"`
}

type OrderItem struct {
	SKU      string  `json:"sku"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type ShipResponse struct {
	Message string `json:"message"`
	Data    struct {
		OrderSN        string `json:"order_sn"`
		TrackingNo     string `json:"tracking_no"`
		ShippingStatus string `json:"shipping_status"`
	} `json:"data"`
}
