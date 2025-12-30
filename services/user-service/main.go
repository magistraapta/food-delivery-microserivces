package main

import (
	"log"
	"net/http"
	"user-service/common"
	"user-service/controller"
	"user-service/database"
	"user-service/middleware"
	"user-service/repository"
	"user-service/service"

	"github.com/gin-gonic/gin"
)

type Main struct {
	router *gin.Engine
}

func (m *Main) InitServer() error {
	err := common.LoadConfig()
	if err != nil {
		return err
	}

	m.router = gin.Default()

	return nil
}

func main() {
	main := Main{}

	if err := main.InitServer(); err != nil {
		log.Fatal("Failed to init server: ", err)
		return
	}

	database, err := database.ConnectDatabase(common.ConfigData)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
		return
	}

	userRepository := repository.NewUserRepositoryImpl(database)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	v1 := main.router.Group("/v1")
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

	main.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	main.router.Run(common.ConfigData.Port)

}
