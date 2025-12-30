package repository

import (
	"food-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RestaurantRepository interface {
	CreateRestaurant(restaurant *models.Restaurant) error
	GetRestaurantByID(id uuid.UUID) (*models.Restaurant, error)
	GetAllRestaurants() ([]models.Restaurant, error)
}

type RestaurantRepositoryImpl struct {
	db *gorm.DB
}

func NewRestaurantRepositoryImpl(db *gorm.DB) RestaurantRepository {
	return &RestaurantRepositoryImpl{db: db}
}

func (r *RestaurantRepositoryImpl) CreateRestaurant(restaurant *models.Restaurant) error {
	return r.db.Create(restaurant).Error
}

func (r *RestaurantRepositoryImpl) GetRestaurantByID(id uuid.UUID) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	if err := r.db.Preload("Foods").First(&restaurant, id).Error; err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func (r *RestaurantRepositoryImpl) GetAllRestaurants() ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	if err := r.db.Find(&restaurants).Error; err != nil {
		return nil, err
	}
	return restaurants, nil
}
