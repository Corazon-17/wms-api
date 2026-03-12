package repository

import (
	"context"

	"wms-api/internal/dto"
	"wms-api/internal/model"

	"github.com/uptrace/bun"
)

type OrderRepository struct {
	db *bun.DB
}

func NewOrderRepository(db *bun.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) FindAll(ctx context.Context, q dto.OrderQuery) ([]model.Order, int, error) {

	var orders []model.Order

	query := r.db.NewSelect().
		Model(&orders)

	if q.WMSStatus != "" {
		query.Where("wms_status = ?", q.WMSStatus)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if q.Page == 0 {
		q.Page = 1
	}

	if q.Limit == 0 {
		q.Limit = 10
	}

	offset := (q.Page - 1) * q.Limit

	sortField := "updated_at"
	if q.Sort != "" {
		allowedSort := map[string]bool{
			"created_at": true,
			"updated_at": true,
			"order_sn":   true,
		}

		if allowedSort[q.Sort] {
			sortField = q.Sort
		}
	}

	sortOrder := "DESC"
	if q.Order == "asc" {
		sortOrder = "ASC"
	}

	query.
		Order(sortField + " " + sortOrder).
		Limit(q.Limit).
		Offset(offset)

	err = query.Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) FindByOrderSN(ctx context.Context, sn string) (*model.Order, error) {

	order := new(model.Order)

	err := r.db.NewSelect().
		Model(order).
		Where("order_sn = ?", sn).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {

	_, err := r.db.NewUpdate().
		Model(order).
		WherePK().
		Exec(ctx)

	return err
}

func (r *OrderRepository) Upsert(ctx context.Context, order *model.Order) error {

	_, err := r.db.NewInsert().
		Model(order).
		On("CONFLICT (order_sn) DO UPDATE").
		Set("shop_id = EXCLUDED.shop_id").
		Set("marketplace_status = EXCLUDED.marketplace_status").
		Set("shipping_status = EXCLUDED.shipping_status").
		Set("total_amount = EXCLUDED.total_amount").
		Set("updated_at = NOW()").
		Exec(ctx)

	return err
}
