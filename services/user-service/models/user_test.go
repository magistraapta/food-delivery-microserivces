package models

import (
	"testing"
)

func TestValidateUser_Success(t *testing.T) {
	user := AddUser{
		Username: "testuser",
		Password: "testpassword",
	}

	err := user.ValidateUser()
	if err != nil {
		t.Errorf("ValidateUser failed for valid user: %v", err)
	}
}

func TestValidateUser_EmptyUsername(t *testing.T) {
	user := AddUser{
		Username: "",
		Password: "testpassword",
	}

	err := user.ValidateUser()
	if err == nil {
		t.Error("ValidateUser should fail for empty username")
	}

	expectedError := "username is required"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestValidateUser_EmptyPassword(t *testing.T) {
	user := AddUser{
		Username: "testuser",
		Password: "",
	}

	err := user.ValidateUser()
	if err == nil {
		t.Error("ValidateUser should fail for empty password")
	}

	expectedError := "password is required"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestValidateUser_BothEmpty(t *testing.T) {
	user := AddUser{
		Username: "",
		Password: "",
	}

	err := user.ValidateUser()
	if err == nil {
		t.Error("ValidateUser should fail when both fields are empty")
	}

	// Should fail on username first
	expectedError := "username is required"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}
