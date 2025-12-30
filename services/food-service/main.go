package main

import (
	"fmt"
	"food-service/config"
	"food-service/controller"
	"food-service/middleware"
	"food-service/repository"
	"food-service/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Printf("failed to load config: %v", err)
	}

	db, err := config.LoadDatabase()

	if err != nil {
		fmt.Printf("failed to load database: %v", err)
	}

	router := gin.Default()

	restaurantRepository := repository.NewRestaurantRepositoryImpl(db)
	restaurantService := service.NewRestaurantServiceImpl(restaurantRepository)
	restaurantController := controller.NewRestaurantController(restaurantService)

	foodRepository := repository.NewFoodRepositoryImpl(db)
	foodService := service.NewFoodServiceImpl(foodRepository)
	foodController := controller.NewFoodController(foodService)

	// Restaurant routes
	restaurant := router.Group("/restaurant")
	{
		// Public routes
		restaurant.GET("/:id", restaurantController.GetRestaurantByID)
		restaurant.GET("", restaurantController.GetAllRestaurants)

		// Admin-only routes
		adminRestaurant := restaurant.Group("")
		adminRestaurant.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			adminRestaurant.POST("", restaurantController.CreateRestaurant)
		}
	}

	// Food routes
	food := router.Group("/food")
	{
		// Public routes
		food.GET("/:id", foodController.GetFoodByID)
		food.GET("", foodController.GetAllFoods)
		food.GET("/restaurant/:restaurantId", foodController.GetFoodsByRestaurantID)

		// Admin-only routes
		adminFood := food.Group("")
		adminFood.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			adminFood.POST("", foodController.CreateFood)
			adminFood.PUT("/:id", foodController.UpdateFood)
			adminFood.DELETE("/:id", foodController.DeleteFood)
		}
	}

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.Run(os.Getenv("PORT"))
}
