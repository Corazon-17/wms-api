package dto

type OrderQuery struct {
	WMSStatus string `query:"wmsStatus"`

	Page     int `query:"page"`
	PageSize int `query:"pageSize"`

	Sort  string `query:"sort"`
	Order string `query:"order"`
}
