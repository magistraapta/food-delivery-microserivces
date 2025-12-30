package controller

import (
	"fmt"
	"net/http"
	"user-service/models"
	"user-service/service"
	"user-service/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return UserController{userService: userService}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var userRequest models.CreateUser

	user := models.User{
		Username: userRequest.Username,
		Password: userRequest.Password,
		Role:     "user",
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.userService.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (uc *UserController) LoginUser(c *gin.Context) {
	var userLogin models.UserLogin

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userService.GetUserByUsername(userLogin.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ComparePassword(user.Password, userLogin.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": tokenString})
}

func (uc *UserController) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	// Assuming "Bearer <token>" format, or just "<token>" depending on client.
	// utils.ValidateToken expects the token string.
	// Let's support both raw token and Bearer scheme for robustness if needed,
	// but utils.ValidateToken seems to take the raw string.
	// For simplicity with Traefik which might forward headers, let's pass the header value
	// directly if utils.ValidateToken handles it, or strip "Bearer ".
	// Looking at utils.ValidateToken (step 16), it calls jwt.Parse.
	// jwt.Parse expects the pure token.

	tokenString := authHeader
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Optional: Return user info in headers or body if needed by downstream
	c.Header("X-User-Id", fmt.Sprintf("%v", claims["username"])) // Example
	c.JSON(http.StatusOK, gin.H{"message": "Token is valid", "claims": claims})
}
