package utils

import (
	"os"
	"testing"
	"time"
	"user-service/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGenerateToken(t *testing.T) {
	user := models.User{
		ID:       uuid.New(),
		Username: "testuser",
		Role:     "user",
	}

	tokenString, err := GenerateToken(user)
	if err != nil {
		t.Errorf("GenerateToken failed: %v", err)
	}

	if tokenString == "" {
		t.Error("GenerateToken returned empty string")
	}
}

func TestValidateToken_Success(t *testing.T) {
	user := models.User{
		ID:       uuid.New(),
		Username: "testuser",
		Role:     "user",
	}

	tokenString, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := ValidateToken(tokenString)
	if err != nil {
		t.Errorf("ValidateToken failed: %v", err)
	}

	if claims["username"] != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, claims["username"])
	}

	if claims["role"] != user.Role {
		t.Errorf("Expected role %s, got %s", user.Role, claims["role"])
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	invalidToken := "invalid.token.string"

	_, err := ValidateToken(invalidToken)
	if err == nil {
		t.Error("ValidateToken should fail for invalid token")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	// Create a token with a different secret
	user := models.User{
		ID:       uuid.New(),
		Username: "testuser",
		Role:     "user",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user.Username,
			"role":     user.Role,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, _ := token.SignedString([]byte("different-secret-key"))

	_, err := ValidateToken(tokenString)
	if err == nil {
		t.Error("ValidateToken should fail for token signed with different secret")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	user := models.User{
		ID:       uuid.New(),
		Username: "testuser",
		Role:     "user",
	}

	// Create an expired token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user.Username,
			"role":     user.Role,
			"exp":      time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
		})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	_, err := ValidateToken(tokenString)
	if err == nil {
		t.Error("ValidateToken should fail for expired token")
	}
}
