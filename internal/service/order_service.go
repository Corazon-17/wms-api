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

func (s *OrderService) GetOrder(ctx context.Context, orderSN string) (*dto.OrderResponse, error) {

	order, err := s.repo.FindByOrderSN(ctx, orderSN)
	if err != nil {
		return nil, err
	}

	resp := dto.NewOrderResponse(*order)

	return &resp, nil
}

func (s *OrderService) PickOrder(ctx context.Context, orderSN string) error {

	order, err := s.repo.FindByOrderSN(ctx, orderSN)
	if err != nil {
		return err
	}

	if order.WMSStatus != domain.WMSReadyToPick {
		return errors.New("order cannot be picked")
	}

	order.WMSStatus = domain.WMSPicking

	return s.repo.Update(ctx, order)
}

func (s *OrderService) PackOrder(ctx context.Context, orderSN string) error {

	order, err := s.repo.FindByOrderSN(ctx, orderSN)
	if err != nil {
		return err
	}

	if order.WMSStatus != domain.WMSPicking {
		return errors.New("order cannot be packed")
	}

	order.WMSStatus = domain.WMSPacked

	return s.repo.Update(ctx, order)
}

func (s *OrderService) ShipOrder(ctx context.Context, orderSN string, channelID string) (*model.Order, error) {

	order, err := s.repo.FindByOrderSN(ctx, orderSN)
	if err != nil {
		return nil, err
	}

	if order.WMSStatus != domain.WMSPacked {
		return nil, errors.New("order must be packed before shipping")
	}

	resp, err := s.marketplace.ShipOrder(orderSN, channelID)
	if err != nil {
		return nil, err
	}

	order.WMSStatus = domain.WMSShipped
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

	for _, o := range orders {

		order := model.Order{
			OrderSN:           o.OrderSN,
			ShopID:            o.ShopID,
			MarketplaceStatus: o.Status,
			ShippingStatus:    o.ShippingStatus,
			WMSStatus:         domain.WMSReadyToPick,
			TotalAmount:       o.TotalAmount,
		}

		s.repo.Upsert(ctx, &order)
	}

	return nil
}

func (s *OrderService) UpdateMarketplaceStatus(ctx context.Context, orderSN string, status string) error {

	order, err := s.repo.FindByOrderSN(ctx, orderSN)
	if err != nil {
		return err
	}

	order.MarketplaceStatus = status

	return s.repo.Update(ctx, order)
}

func (s *OrderService) UpdateShippingStatus(ctx context.Context, orderSN string, shippingState string) error {

	order, err := s.repo.FindByOrderSN(ctx, orderSN)
	if err != nil {
		return err
	}

	order.ShippingStatus = shippingState

	return s.repo.Update(ctx, order)
}
