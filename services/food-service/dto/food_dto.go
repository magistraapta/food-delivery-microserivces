package dto

import "github.com/google/uuid"

type CreateFoodRequest struct {
	Name         string  `json:"name" validate:"required,min=3"`
	Price        float64 `json:"price" validate:"required,gte=0"`
	Description  string  `json:"description" validate:"required"`
	RestaurantID string  `json:"restaurant_id" validate:"required,uuid"`
}

type FoodResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
}
