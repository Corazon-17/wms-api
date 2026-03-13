package repository

import (
	"context"
	"encoding/json"
	"strings"

	"wms-api/internal/domain"
	"wms-api/internal/dto"
	"wms-api/internal/helper"
	"wms-api/internal/model"
	"wms-api/internal/provider/marketplace"

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

	if q.Search != "" {
		query.Where("order_sn = ?", q.Search)
	}

	if q.FilterField != "" && q.FilterValues != "" {
		allowedFilter := map[string]bool{
			"marketplace_status": true,
			"shipping_status":    true,
			"wms_status_id":      true,
		}

		filterField := ""
		if allowedFilter[q.FilterField] {
			filterField = q.FilterField
		}

		filterValues := strings.Split(q.FilterValues, ",")

		query.Where("? IN (?)", bun.Ident(filterField), bun.Tuple(filterValues))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if q.Page == 0 {
		q.Page = 1
	}

	if q.PageSize == 0 {
		q.PageSize = 10
	}

	offset := (q.Page - 1) * q.PageSize

	sortField := "updated_at"
	if q.SortField != "" {
		allowedSort := map[string]bool{
			"created_at": true,
			"updated_at": true,
			"order_sn":   true,
		}

		if allowedSort[q.SortField] {
			sortField = q.SortField
		}
	}

	sortDir := "DESC"
	if q.SortDir == "asc" {
		sortDir = "ASC"
	}

	query.
		Order(sortField + " " + sortDir).
		Limit(q.PageSize).
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

func (r *OrderRepository) FindOrderDetailByOrderSN(ctx context.Context, sn string) (*model.OrderDetail, error) {
	order := new(model.OrderDetail)

	err := r.db.NewSelect().
		Model(order).
		Where("order_sn = ?", sn).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	orderItems, err := r.FindOrderItemsByOrderSN(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	order.Items = *orderItems

	return order, nil
}

func (r *OrderRepository) FindOrderItemsByOrderSN(ctx context.Context, orderId int64) (*[]model.OrderItem, error) {

	orderItem := new([]model.OrderItem)

	err := r.db.NewSelect().
		Model(orderItem).
		Where("order_id = ?", orderId).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return orderItem, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {

	_, err := r.db.NewUpdate().
		Model(order).
		WherePK().
		Exec(ctx)

	return err
}

func (r *OrderRepository) UpsertOrdersTx(ctx context.Context, tx bun.Tx, orders []model.Order) error {

	_, err := tx.NewInsert().
		Model(&orders).
		On("CONFLICT (order_sn) DO UPDATE").
		Set("marketplace_status = EXCLUDED.marketplace_status").
		Set("shipping_status = EXCLUDED.shipping_status").
		Set("total_amount = EXCLUDED.total_amount").
		Set("raw_marketplace_payload = EXCLUDED.raw_marketplace_payload").
		Exec(ctx)

	return err
}

func (r *OrderRepository) DeleteItemsByOrderIDsTx(ctx context.Context, tx bun.Tx, orderIDs []int64) error {

	_, err := tx.NewDelete().
		Model((*model.OrderItem)(nil)).
		Where("order_id IN (?)", bun.In(orderIDs)).
		Exec(ctx)

	return err
}

func (r *OrderRepository) InsertItemsTx(ctx context.Context, tx bun.Tx, items []model.OrderItem) error {

	_, err := tx.NewInsert().
		Model(&items).
		Exec(ctx)

	return err
}

func (r *OrderRepository) FindOrdersBySNsTx(ctx context.Context, tx bun.Tx, orders []model.Order) (map[string]model.Order, error) {

	orderSNs := helper.ExtractOrderSNs(orders)

	var dbOrders []model.Order

	err := tx.NewSelect().
		Model(&dbOrders).
		Where("order_sn IN (?)", bun.In(orderSNs)).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	result := make(map[string]model.Order)

	for _, o := range dbOrders {
		result[o.OrderSN] = o
	}

	return result, nil
}

func (r *OrderRepository) SyncOrders(ctx context.Context, marketplaceOrders []marketplace.Order) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {

		var orders []model.Order
		var items []model.OrderItem

		for _, o := range marketplaceOrders {
			raw, err := json.Marshal(o)
			if err != nil {
				return err
			}

			var payload map[string]any
			json.Unmarshal(raw, &payload)

			orders = append(orders, model.Order{
				OrderSN:               o.OrderSN,
				ShopID:                o.ShopID,
				MarketplaceStatus:     o.Status,
				ShippingStatus:        o.ShippingStatus,
				TotalAmount:           o.TotalAmount,
				WMSStatusID:           domain.WMSReadyToPick,
				RawMarketplacePayload: payload,
			})
		}

		err := r.UpsertOrdersTx(ctx, tx, orders)
		if err != nil {
			return err
		}

		dbOrders, err := r.FindOrdersBySNsTx(ctx, tx, orders)
		if err != nil {
			return err
		}

		for _, o := range marketplaceOrders {

			orderID := dbOrders[o.OrderSN].ID

			for _, i := range o.Items {

				items = append(items, model.OrderItem{
					OrderID:  orderID,
					SKU:      i.SKU,
					Quantity: i.Quantity,
					Price:    i.Price,
				})
			}
		}

		err = r.DeleteItemsByOrderIDsTx(ctx, tx, helper.ExtractIDs(dbOrders))
		if err != nil {
			return err
		}

		return r.InsertItemsTx(ctx, tx, items)
	})
}

func (r *OrderRepository) GetMarketplaceStatuses(ctx context.Context) ([]string, error) {

	var statuses []string

	err := r.db.NewSelect().
		Table("orders").
		ColumnExpr("DISTINCT marketplace_status").
		Where("marketplace_status IS NOT NULL").
		Scan(ctx, &statuses)

	return statuses, err
}

func (r *OrderRepository) GetShippingStatuses(ctx context.Context) ([]string, error) {

	var statuses []string

	err := r.db.NewSelect().
		Table("orders").
		ColumnExpr("DISTINCT shipping_status").
		Where("shipping_status IS NOT NULL").
		Scan(ctx, &statuses)

	return statuses, err
}

func (r *OrderRepository) GetWMSStatuses(ctx context.Context) ([]model.WMSStatus, error) {

	var statuses []model.WMSStatus

	err := r.db.NewSelect().
		Model(&statuses).
		Scan(ctx)

	return statuses, err
}
