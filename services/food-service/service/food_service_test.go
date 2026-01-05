package service

import (
	"food-service/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFoodRepository struct {
	mock.Mock
}

func (m *MockFoodRepository) CreateFood(food *models.Food) error {
	args := m.Called(food)
	return args.Error(0)
}

func (m *MockFoodRepository) GetFoodByID(id uuid.UUID) (*models.Food, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Food), args.Error(1)
}

func (m *MockFoodRepository) GetAllFoods() ([]models.Food, error) {
	args := m.Called()
	return args.Get(0).([]models.Food), args.Error(1)
}

func (m *MockFoodRepository) GetFoodsByRestaurantID(restaurantID uuid.UUID) ([]models.Food, error) {
	args := m.Called(restaurantID)
	return args.Get(0).([]models.Food), args.Error(1)
}

func (m *MockFoodRepository) UpdateFood(food *models.Food) error {
	args := m.Called(food)
	return args.Error(0)
}

func (m *MockFoodRepository) DeleteFood(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateFood(t *testing.T) {
	mockFoodRepository := &MockFoodRepository{}
	foodService := NewFoodServiceImpl(mockFoodRepository)

	food := &models.Food{
		ID:           uuid.New(),
		Name:         "Test Food",
		Description:  "Test Description",
		Price:        10.0,
		RestaurantID: uuid.New(),
	}

	mockFoodRepository.On("CreateFood", food).Return(nil)

	err := foodService.CreateFood(food)

	if err != nil {
		t.Errorf("CreateFood() error = %v, want nil", err)
	}

	mockFoodRepository.AssertExpectations(t)
	mockFoodRepository.AssertCalled(t, "CreateFood", food)
	assert.NoError(t, err)

	// check if the food is created
	assert.Equal(t, food.Name, "Test Food")
	assert.Equal(t, food.Description, "Test Description")
	assert.Equal(t, food.Price, 10.0)
}

func TestGetFoodByID(t *testing.T) {
	mockFoodRepository := &MockFoodRepository{}
	foodService := NewFoodServiceImpl(mockFoodRepository)

	food := &models.Food{
		ID:           uuid.New(),
		Name:         "Test Food",
		Description:  "Test Description",
		Price:        10.0,
		RestaurantID: uuid.New(),
	}

	mockFoodRepository.On("GetFoodByID", food.ID).Return(food, nil)

	f, err := foodService.GetFoodByID(food.ID)

	if err != nil {
		t.Errorf("GetFoodByID() error = %v, want nil", err)
	}

	mockFoodRepository.AssertExpectations(t)
	mockFoodRepository.AssertCalled(t, "GetFoodByID", food.ID)
	assert.NoError(t, err)
	assert.Equal(t, f, food)
	assert.Equal(t, f.ID, food.ID)
}

func TestGetAllFoods(t *testing.T) {
	mockFoodRepository := &MockFoodRepository{}
	foodService := NewFoodServiceImpl(mockFoodRepository)

	foods := []models.Food{
		{
			ID:           uuid.New(),
			Name:         "Test Food",
			Description:  "Test Description",
			Price:        10.0,
			RestaurantID: uuid.New(),
		},
		{
			ID:           uuid.New(),
			Name:         "Test Food 2",
			Description:  "Test Description 2",
			Price:        20.0,
			RestaurantID: uuid.New(),
		},
	}

	mockFoodRepository.On("GetAllFoods").Return(foods, nil)

	f, err := foodService.GetAllFoods()

	if err != nil {
		t.Errorf("GetAllFoods() error = %v, want nil", err)
	}

	mockFoodRepository.AssertExpectations(t)
	mockFoodRepository.AssertCalled(t, "GetAllFoods")
	assert.NoError(t, err)
	assert.Equal(t, f, foods)
}

func TestUpdateFood(t *testing.T) {
	mockFoodRepository := &MockFoodRepository{}
	foodService := NewFoodServiceImpl(mockFoodRepository)

	food := &models.Food{
		ID:           uuid.New(),
		Name:         "Test Food",
		Description:  "Test Description",
		Price:        10.0,
		RestaurantID: uuid.New(),
	}

	mockFoodRepository.On("UpdateFood", food).Return(nil)

	err := foodService.UpdateFood(food)

	if err != nil {
		t.Errorf("UpdateFood() error = %v, want nil", err)
	}

	mockFoodRepository.AssertExpectations(t)
	mockFoodRepository.AssertCalled(t, "UpdateFood", food)
	assert.NoError(t, err)

	assert.Equal(t, food.Name, "Test Food")
	assert.Equal(t, food.Description, "Test Description")
	assert.Equal(t, food.Price, 10.0)
}

func TestDeleteFoodById(t *testing.T) {
	mockRepository := &MockFoodRepository{}

	foodService := NewFoodServiceImpl(mockRepository)

	foodId := uuid.New()

	mockRepository.On("DeleteFood", foodId).Return(nil)

	err := foodService.DeleteFood(foodId)

	if err != nil {
		t.Errorf("DeleteFood() error = %v, want nil", err)
	}

	mockRepository.AssertExpectations(t)
	mockRepository.AssertCalled(t, "DeleteFood", foodId)
	assert.NoError(t, err)
}
