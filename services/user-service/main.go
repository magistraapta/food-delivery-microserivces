package main

import (
	"log"
	"net/http"
	"os"
	"user-service/config"
	"user-service/controller"
	"user-service/database"
	"user-service/middleware"
	"user-service/repository"
	"user-service/service"

	"github.com/gin-gonic/gin"
)

func main() {

	err := config.LoadEnv()

	if err != nil {
		log.Fatal("Failed to load config: ", err)
		return
	}

	database, err := database.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
		return
	}

	userRepository := repository.NewUserRepositoryImpl(database)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		admin := v1.Group("/admin")
		{
			admin.Use(middleware.AdminMiddleware())
			admin.GET("/hello", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Hello Admin"})
			})
		}
		user := v1.Group("/user")
		{
			user.POST("/create", userController.CreateUser)
			user.POST("/login", userController.LoginUser)
			user.GET("/validate", userController.ValidateToken)
			user.GET("/hello", middleware.AuthMiddleware(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Hello User"})
			})
		}
	}

	log.Println("User Service is running at", os.Getenv("PORT"))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	router.Run(os.Getenv("PORT"))

}
