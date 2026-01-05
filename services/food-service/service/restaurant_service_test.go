package service

import (
	"food-service/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRestaurantRepository struct {
	mock.Mock
}

func (m *MockRestaurantRepository) CreateRestaurant(restaurant *models.Restaurant) error {
	args := m.Called(restaurant)
	return args.Error(0)
}

func (m *MockRestaurantRepository) GetRestaurantByID(id uuid.UUID) (*models.Restaurant, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) GetAllRestaurants() ([]models.Restaurant, error) {
	args := m.Called()
	return args.Get(0).([]models.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) UpdateRestaurant(restaurant *models.Restaurant) error {
	args := m.Called(restaurant)
	return args.Error(0)
}

func (m *MockRestaurantRepository) DeleteRestaurant(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateRestaurant(t *testing.T) {
	mockRestaurantRepository := &MockRestaurantRepository{}

	restaurantService := NewRestaurantServiceImpl(mockRestaurantRepository)

	restaurant := &models.Restaurant{
		ID:      uuid.New(),
		Name:    "Test Restaurant",
		Address: "Test Address",
		Foods: []models.Food{
			{
				ID:           uuid.New(),
				Name:         "Test Food",
				Description:  "Test Description",
				Price:        10.0,
				RestaurantID: uuid.New(),
			},
		},
	}

	mockRestaurantRepository.On("CreateRestaurant", restaurant).Return(nil)

	err := restaurantService.CreateRestaurant(restaurant)

	if err != nil {
		t.Errorf("CreateRestaurant() error = %v, want nil", err)
	}

	mockRestaurantRepository.AssertExpectations(t)
	mockRestaurantRepository.AssertCalled(t, "CreateRestaurant", restaurant)
	assert.NoError(t, err)

	assert.Equal(t, restaurant.Name, "Test Restaurant")
	assert.Equal(t, restaurant.Address, "Test Address")
	assert.Equal(t, restaurant.Foods[0].Name, "Test Food")
	assert.Equal(t, restaurant.Foods[0].Description, "Test Description")
	assert.Equal(t, restaurant.Foods[0].Price, 10.0)
}

func TestGetRestaurantByID(t *testing.T) {
	mockRestaurantRepository := &MockRestaurantRepository{}

	restaurantService := NewRestaurantServiceImpl(mockRestaurantRepository)

	restaurant := &models.Restaurant{
		ID:      uuid.New(),
		Name:    "Test Restaurant",
		Address: "Test Address",
		Foods: []models.Food{
			{
				ID:           uuid.New(),
				Name:         "Test Food",
				Description:  "Test Description",
				Price:        10.0,
				RestaurantID: uuid.New(),
			},
		},
	}

	mockRestaurantRepository.On("GetRestaurantByID", restaurant.ID).Return(restaurant, nil)

	f, err := restaurantService.GetRestaurantByID(restaurant.ID)

	if err != nil {
		t.Errorf("GetRestaurantByID() error = %v, want nil", err)
	}

	mockRestaurantRepository.AssertExpectations(t)
	mockRestaurantRepository.AssertCalled(t, "GetRestaurantByID", restaurant.ID)
	assert.NoError(t, err)
	assert.Equal(t, f, restaurant)
	assert.Equal(t, f.ID, restaurant.ID)
}
