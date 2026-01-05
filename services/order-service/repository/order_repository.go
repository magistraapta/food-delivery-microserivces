package repository

import (
	"order-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) error
	GetOrderById(id uuid.UUID) (*models.Order, error)
	UpdateOrderStatus(id uuid.UUID, status string) error
}

type OrderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepositoryImpl(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{db: db}
}

func (r *OrderRepositoryImpl) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepositoryImpl) GetOrderById(id uuid.UUID) (*models.Order, error) {
	var order models.Order

	if err := r.db.First(&order, id).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// UpdateOrderStatus updates the status of an order
func (r *OrderRepositoryImpl) UpdateOrderStatus(id uuid.UUID, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}
