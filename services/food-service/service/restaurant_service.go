package service

import (
	"food-service/models"
	"food-service/repository"

	"github.com/google/uuid"
)

type RestaurantService interface {
	CreateRestaurant(restaurant *models.Restaurant) error
	GetRestaurantByID(id uuid.UUID) (*models.Restaurant, error)
	GetAllRestaurants() ([]models.Restaurant, error)
}

type RestaurantServiceImpl struct {
	restaurantRepository repository.RestaurantRepository
}

func NewRestaurantServiceImpl(restaurantRepository repository.RestaurantRepository) RestaurantService {
	return &RestaurantServiceImpl{restaurantRepository: restaurantRepository}
}

func (s *RestaurantServiceImpl) CreateRestaurant(restaurant *models.Restaurant) error {
	return s.restaurantRepository.CreateRestaurant(restaurant)
}

func (s *RestaurantServiceImpl) GetRestaurantByID(id uuid.UUID) (*models.Restaurant, error) {
	return s.restaurantRepository.GetRestaurantByID(id)
}

func (s *RestaurantServiceImpl) GetAllRestaurants() ([]models.Restaurant, error) {
	return s.restaurantRepository.GetAllRestaurants()
}
