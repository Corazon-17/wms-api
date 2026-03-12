package pagination

import (
	"math"
	"wms-api/pkg/utils"
)

type PaginationMeta struct {
	Page        int  `json:"page"`
	PageSize    int  `json:"pageSize"`
	TotalPages  int  `json:"totalPages"`
	TotalItems  int  `json:"totalItems"`
	HasPrevPage bool `json:"hasPrevPage"`
	HasNextPage bool `json:"hasNextPage"`
	PrevPage    *int `json:"prevPage"`
	NextPage    *int `json:"nextPage"`
}

func GenerateOffset(page, pageSize int) int32 {
	return (int32(page) - 1) * int32(pageSize)
}

func GeneratePaginationMeta(currPage, pageSize, totalItems int) PaginationMeta {
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	var prevPage, nextPage *int

	if currPage > 1 {
		prevPage = utils.Pointer(currPage - 1)
	}

	if currPage < totalPages {
		nextPage = utils.Pointer(currPage + 1)
	}

	return PaginationMeta{
		Page:        currPage,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		HasPrevPage: currPage > 1,
		HasNextPage: currPage < totalPages,
		PrevPage:    prevPage,
		NextPage:    nextPage,
	}
}
