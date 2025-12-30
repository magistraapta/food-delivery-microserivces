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

type FoodController struct {
	foodService service.FoodService
}

func NewFoodController(foodService service.FoodService) FoodController {
	return FoodController{foodService: foodService}
}

func (fc *FoodController) CreateFood(c *gin.Context) {
	var request dto.CreateFoodRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := utils.ValidateStruct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	food := models.Food{
		ID:           uuid.New(),
		Name:         request.Name,
		Price:        request.Price,
		RestaurantID: uuid.MustParse(request.RestaurantID),
	}
	err := fc.foodService.CreateFood(&food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.FoodResponse{
		ID:          food.ID,
		Name:        food.Name,
		Price:       food.Price,
		Description: food.Description,
	}

	c.JSON(http.StatusCreated, response)
}

func (fc *FoodController) GetFoodByID(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid food ID"})
		return
	}

	food, err := fc.foodService.GetFoodByID(parsedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "food not found"})
		return
	}

	c.JSON(http.StatusOK, food)
}

func (fc *FoodController) GetAllFoods(c *gin.Context) {
	foods, err := fc.foodService.GetAllFoods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, foods)
}

func (fc *FoodController) GetFoodsByRestaurantID(c *gin.Context) {
	restaurantID := c.Param("restaurantId")
	parsedID, err := uuid.Parse(restaurantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	foods, err := fc.foodService.GetFoodsByRestaurantID(parsedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, foods)
}

func (fc *FoodController) UpdateFood(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid food ID"})
		return
	}

	var food models.Food
	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	food.ID = parsedID

	err = fc.foodService.UpdateFood(&food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, food)
}

func (fc *FoodController) DeleteFood(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid food ID"})
		return
	}

	err = fc.foodService.DeleteFood(parsedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "food deleted successfully"})
}
