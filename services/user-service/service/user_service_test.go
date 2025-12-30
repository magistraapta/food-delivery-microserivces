package service

import (
	"errors"
	"testing"
	"user-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of the UserRepository interface
type MockUserRepository struct {
	users       map[string]models.User
	shouldError bool
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:       make(map[string]models.User),
		shouldError: false,
	}
}

func (m *MockUserRepository) CreateUser(user models.User) error {
	if m.shouldError {
		return errors.New("database error")
	}
	user.ID = uuid.New()
	m.users[user.Username] = user
	return nil
}

func (m *MockUserRepository) GetUserByUsername(username string) (models.User, error) {
	if m.shouldError {
		return models.User{}, errors.New("database error")
	}
	user, exists := m.users[username]
	if !exists {
		return models.User{}, gorm.ErrRecordNotFound
	}
	return user, nil
}

func TestNewUserService(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	if service.userRepository == nil {
		t.Error("UserService should have a repository")
	}
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	user := models.User{
		Username: "newuser",
		Password: "password123",
	}

	err := service.CreateUser(user)
	if err != nil {
		t.Errorf("CreateUser failed: %v", err)
	}

	// Verify user was stored
	storedUser, err := mockRepo.GetUserByUsername("newuser")
	if err != nil {
		t.Errorf("User was not stored in repository: %v", err)
	}

	if storedUser.Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", storedUser.Role)
	}
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	// Create first user
	user := models.User{
		Username: "existinguser",
		Password: "password123",
	}
	_ = service.CreateUser(user)

	// Try to create another user with the same username
	duplicateUser := models.User{
		Username: "existinguser",
		Password: "differentpassword",
	}
	err := service.CreateUser(duplicateUser)

	if err == nil {
		t.Error("CreateUser should fail for duplicate username")
	}

	expectedError := "username already exists"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestCreateUser_DatabaseError(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	// Simulate database error
	mockRepo.shouldError = true

	user := models.User{
		Username: "newuser",
		Password: "password123",
	}

	err := service.CreateUser(user)
	if err == nil {
		t.Error("CreateUser should fail when database has error")
	}
}

func TestGetUserByUsername_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	// Add a user to the mock repo
	expectedUser := models.User{
		ID:       uuid.New(),
		Username: "testuser",
		Password: "hashedpassword",
		Role:     "user",
	}
	mockRepo.users["testuser"] = expectedUser

	user, err := service.GetUserByUsername("testuser")
	if err != nil {
		t.Errorf("GetUserByUsername failed: %v", err)
	}

	if user.Username != expectedUser.Username {
		t.Errorf("Expected username '%s', got '%s'", expectedUser.Username, user.Username)
	}
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	_, err := service.GetUserByUsername("nonexistent")
	if err == nil {
		t.Error("GetUserByUsername should fail for non-existent user")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected gorm.ErrRecordNotFound, got %v", err)
	}
}
