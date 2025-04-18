package service

import (
	"context"
	"errors"

	"order-service/internal/domain"
	"order-service/internal/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order domain.Order) (domain.Order, error)
	GetOrderByID(ctx context.Context, id string) (domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status domain.OrderStatus) error
	ListOrders(ctx context.Context, filter domain.OrderFilter) ([]domain.Order, int, error)
	GetUserOrders(ctx context.Context, userID string) ([]domain.Order, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	if order.UserID == "" {
		return domain.Order{}, errors.New("user ID is required")
	}

	if len(order.Items) == 0 {
		return domain.Order{}, errors.New("order must contain at least one item")
	}

	var total float64
	for _, item := range order.Items {
		if item.Quantity <= 0 {
			return domain.Order{}, errors.New("item quantity must be greater than zero")
		}
		total += item.Price * float64(item.Quantity)
	}

	order.Total = total
	order.Status = domain.OrderStatusPending

	return s.orderRepo.Create(ctx, order)
}

func (s *orderService) GetOrderByID(ctx context.Context, id string) (domain.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !isValidStatusTransition(order.Status, status) {
		return errors.New("invalid status transition")
	}

	order.Status = status
	return s.orderRepo.Update(ctx, order)
}

func (s *orderService) ListOrders(ctx context.Context, filter domain.OrderFilter) ([]domain.Order, int, error) {
	return s.orderRepo.List(ctx, filter)
}

func (s *orderService) GetUserOrders(ctx context.Context, userID string) ([]domain.Order, error) {
	return s.orderRepo.GetUserOrders(ctx, userID)
}

func isValidStatusTransition(current, next domain.OrderStatus) bool {
	switch current {
	case domain.OrderStatusPending:
		return next == domain.OrderStatusPaid || next == domain.OrderStatusCancelled
	case domain.OrderStatusPaid:
		return next == domain.OrderStatusShipped || next == domain.OrderStatusCancelled
	case domain.OrderStatusShipped:
		return next == domain.OrderStatusDelivered || next == domain.OrderStatusCancelled
	case domain.OrderStatusDelivered, domain.OrderStatusCancelled:
		return false
	default:
		return false
	}
}
