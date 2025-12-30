package service

import (
	"food-service/models"
	"food-service/repository"

	"github.com/google/uuid"
)

type FoodService interface {
	CreateFood(food *models.Food) error
	GetFoodByID(id uuid.UUID) (*models.Food, error)
	GetAllFoods() ([]models.Food, error)
	GetFoodsByRestaurantID(restaurantID uuid.UUID) ([]models.Food, error)
	UpdateFood(food *models.Food) error
	DeleteFood(id uuid.UUID) error
}

type FoodServiceImpl struct {
	foodRepository repository.FoodRepository
}

func NewFoodServiceImpl(foodRepository repository.FoodRepository) FoodService {
	return &FoodServiceImpl{foodRepository: foodRepository}
}

func (s *FoodServiceImpl) CreateFood(food *models.Food) error {
	return s.foodRepository.CreateFood(food)
}

func (s *FoodServiceImpl) GetFoodByID(id uuid.UUID) (*models.Food, error) {
	return s.foodRepository.GetFoodByID(id)
}

func (s *FoodServiceImpl) GetAllFoods() ([]models.Food, error) {
	return s.foodRepository.GetAllFoods()
}

func (s *FoodServiceImpl) GetFoodsByRestaurantID(restaurantID uuid.UUID) ([]models.Food, error) {
	return s.foodRepository.GetFoodsByRestaurantID(restaurantID)
}

func (s *FoodServiceImpl) UpdateFood(food *models.Food) error {
	return s.foodRepository.UpdateFood(food)
}

func (s *FoodServiceImpl) DeleteFood(id uuid.UUID) error {
	return s.foodRepository.DeleteFood(id)
}
