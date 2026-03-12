package dto

type OrderQuery struct {
	WMSStatus string `query:"wms_status"`

	Page  int `query:"page"`
	Limit int `query:"limit"`

	Sort  string `query:"sort"`
	Order string `query:"order"`
}
