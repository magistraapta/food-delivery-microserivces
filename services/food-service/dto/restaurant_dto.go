package dto

type CreateRestaurantRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
}

type RestaurantResponse struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Address string         `json:"address"`
	Foods   []FoodResponse `json:"foods"`
}
