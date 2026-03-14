package service

import (
	"context"
	"errors"

	"wms-api/internal/domain"
	"wms-api/internal/dto"
	"wms-api/internal/model"
	"wms-api/internal/provider/marketplace"
	"wms-api/internal/repository"
)

type OrderService struct {
	repo        *repository.OrderRepository
	marketplace *marketplace.Client
}

func NewOrderService(repo *repository.OrderRepository, marketplace *marketplace.Client) *OrderService {

	return &OrderService{
		repo:        repo,
		marketplace: marketplace,
	}
}

func (s *OrderService) GetOrders(ctx context.Context, q dto.OrderQuery) ([]dto.OrderResponse, int, error) {

	orders, total, err := s.repo.FindAll(ctx, q)
	if err != nil {
		return nil, 0, err
	}

	var result []dto.OrderResponse

	for _, o := range orders {
		result = append(result, dto.NewOrderResponse(o))
	}

	return result, total, nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderSN string) (*dto.OrderDetailResponse, error) {

	order, err := s.repo.FindOrderDetailByOrderSN(ctx, orderSN)
	if err != nil {
		return nil, err
	}

	resp := dto.NewOrderDetailResponse(*order)

	return &resp, nil
}

func (s *OrderService) PickOrder(ctx context.Context, orderSN string) error {

	order, err := s.repo.FindByOrderSN(ctx, orderSN)
	if err != nil {
		return err
	}

	if order.WMSStatusID != domain.WMSReadyToPick {
		return errors.New("order cannot be picked")
	}

	order.WMSStatusID = domain.WMSPicking

	return s.repo.Update(ctx, order)
}

func (s *OrderService) PackOrder(ctx context.Context, orderSN string) error {

	order, err := s.repo.FindByOrderSN(ctx, orderSN)
	if err != nil {
		return err
	}

	if order.WMSStatusID != domain.WMSPicking {
		return errors.New("order cannot be packed")
	}

	order.WMSStatusID = domain.WMSPacked

	return s.repo.Update(ctx, order)
}

func (s *OrderService) ShipOrder(ctx context.Context, orderSN string, channelID string) (*model.Order, error) {

	order, err := s.repo.FindByOrderSN(ctx, orderSN)
	if err != nil {
		return nil, err
	}

	if order.WMSStatusID != domain.WMSPacked {
		return nil, errors.New("order must be packed before shipping")
	}

	resp, err := s.marketplace.ShipOrder(orderSN, channelID)
	if err != nil {
		return nil, err
	}

	order.WMSStatusID = domain.WMSShipped
	order.ShippingStatus = resp.Data.ShippingStatus
	order.TrackingNumber = resp.Data.TrackingNo

	err = s.repo.Update(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) SyncOrders(ctx context.Context) error {

	orders, err := s.marketplace.ListOrders()
	if err != nil {
		return err
	}

	return s.repo.SyncOrders(ctx, orders)
}

func (s *OrderService) UpdateMarketplaceStatus(ctx context.Context, orderSN string, status string) error {

	order, err := s.repo.FindByOrderSN(ctx, orderSN)
	if err != nil {
		return err
	}

	order.MarketplaceStatus = status

	return s.repo.Update(ctx, order)
}

func (s *OrderService) UpdateShippingStatus(ctx context.Context, orderSN string, shippingStatus string) error {

	order, err := s.repo.FindByOrderSN(ctx, orderSN)
	if err != nil {
		return err
	}

	order.ShippingStatus = shippingStatus

	return s.repo.Update(ctx, order)
}

func (s *OrderService) GetWMSStatuses(ctx context.Context) ([]dto.WMSStatusResponse, error) {
	status, err := s.repo.GetWMSStatuses(ctx)
	if err != nil {
		return nil, err
	}

	resp := dto.NewWMSStatusResponse(status)

	return resp, nil
}

func (s *OrderService) GetMarketplaceStatuses(ctx context.Context) ([]string, error) {
	return s.repo.GetMarketplaceStatuses(ctx)
}

func (s *OrderService) GetShippingStatuses(ctx context.Context) ([]string, error) {
	return s.repo.GetShippingStatuses(ctx)
}

type OrderSummary struct {
	TotalOrders     int `json:"totalOrders"`
	CancelledOrders int `json:"cancelledOrders"`
}

func (s *OrderService) GetOrderSummary(ctx context.Context) (OrderSummary, error) {

	total, cancelled, err := s.repo.GetOrderCounts(ctx)
	if err != nil {
		return OrderSummary{}, err
	}

	return OrderSummary{
		TotalOrders:     total,
		CancelledOrders: cancelled,
	}, nil
}
