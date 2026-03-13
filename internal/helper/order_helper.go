package helper

import "wms-api/internal/model"

func ExtractOrderSNs(orders []model.Order) []string {

	sns := make([]string, 0, len(orders))

	for _, o := range orders {
		sns = append(sns, o.OrderSN)
	}

	return sns
}

func ExtractIDs(orderMap map[string]model.Order) []int64 {

	ids := make([]int64, 0, len(orderMap))

	for _, o := range orderMap {
		ids = append(ids, o.ID)
	}

	return ids
}
