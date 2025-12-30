package repository

import (
	"food-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FoodRepository interface {
	CreateFood(food *models.Food) error
	GetFoodByID(id uuid.UUID) (*models.Food, error)
	GetAllFoods() ([]models.Food, error)
	GetFoodsByRestaurantID(restaurantID uuid.UUID) ([]models.Food, error)
	UpdateFood(food *models.Food) error
	DeleteFood(id uuid.UUID) error
}

type FoodRepositoryImpl struct {
	db *gorm.DB
}

func NewFoodRepositoryImpl(db *gorm.DB) FoodRepository {
	return &FoodRepositoryImpl{db: db}
}

func (r *FoodRepositoryImpl) CreateFood(food *models.Food) error {
	return r.db.Create(food).Error
}

func (r *FoodRepositoryImpl) GetFoodByID(id uuid.UUID) (*models.Food, error) {
	var food models.Food
	if err := r.db.First(&food, id).Error; err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *FoodRepositoryImpl) GetAllFoods() ([]models.Food, error) {
	var foods []models.Food
	if err := r.db.Find(&foods).Error; err != nil {
		return nil, err
	}
	return foods, nil
}

func (r *FoodRepositoryImpl) GetFoodsByRestaurantID(restaurantID uuid.UUID) ([]models.Food, error) {
	var foods []models.Food
	if err := r.db.Where("restaurant_id = ?", restaurantID).Find(&foods).Error; err != nil {
		return nil, err
	}
	return foods, nil
}

func (r *FoodRepositoryImpl) UpdateFood(food *models.Food) error {
	return r.db.Save(food).Error
}

func (r *FoodRepositoryImpl) DeleteFood(id uuid.UUID) error {
	return r.db.Delete(&models.Food{}, id).Error
}
