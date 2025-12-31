package service

import (
	"order-service/models"
	"order-service/repository"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(order *models.Order) error
	GetOrderById(id uuid.UUID) (*models.Order, error)
}

type OrderServiceImpl struct {
	orderRepository repository.OrderRepository
}

func NewOrderServiceImpl(orderRepository repository.OrderRepository) OrderService {
	return &OrderServiceImpl{orderRepository: orderRepository}
}

func (s *OrderServiceImpl) CreateOrder(order *models.Order) error {
	return s.orderRepository.CreateOrder(order)
}

func (s *OrderServiceImpl) GetOrderById(id uuid.UUID) (*models.Order, error) {
	return s.orderRepository.GetOrderById(id)
}
