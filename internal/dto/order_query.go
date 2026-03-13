package dto

type OrderQuery struct {
	Page     int `query:"page"`
	PageSize int `query:"pageSize"`

	Search       string `query:"search"`
	FilterField  string `query:"filterField"`
	FilterValues string `query:"filterValues"`

	SortField string `query:"sortField"`
	SortDir   string `query:"sortDir"`
}
