package controller

import (
	"food-service/dto"
	"food-service/models"
	"food-service/service"
	"food-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RestaurantController struct {
	restaurantService service.RestaurantService
}

func NewRestaurantController(restaurantService service.RestaurantService) RestaurantController {
	return RestaurantController{restaurantService: restaurantService}
}

func (rc *RestaurantController) CreateRestaurant(c *gin.Context) {
	var request dto.CreateRestaurantRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := utils.ValidateStruct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurant := models.Restaurant{
		ID:      uuid.New(),
		Name:    request.Name,
		Address: request.Address,
	}

	err := rc.restaurantService.CreateRestaurant(&restaurant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.RestaurantResponse{
		ID:      restaurant.ID.String(),
		Name:    restaurant.Name,
		Address: restaurant.Address,
		Foods:   []dto.FoodResponse{},
	}

	c.JSON(http.StatusOK, response)
}

func (rc *RestaurantController) GetRestaurantByID(c *gin.Context) {
	id := c.Param("id")
	restaurant, err := rc.restaurantService.GetRestaurantByID(uuid.Must(uuid.Parse(id)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	foodResponse := make([]dto.FoodResponse, len(restaurant.Foods))
	for i, food := range restaurant.Foods {
		foodResponse[i] = dto.FoodResponse{
			ID:          food.ID,
			Name:        food.Name,
			Price:       food.Price,
			Description: food.Description,
		}
	}

	response := dto.RestaurantResponse{
		ID:      restaurant.ID.String(),
		Name:    restaurant.Name,
		Address: restaurant.Address,
		Foods:   foodResponse,
	}

	c.JSON(http.StatusOK, response)
}

func (rc *RestaurantController) GetAllRestaurants(c *gin.Context) {
	restaurants, err := rc.restaurantService.GetAllRestaurants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurants)
}
