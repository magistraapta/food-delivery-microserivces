package utils

import (
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	password := "testpassword123"
	hashedPassword := GeneratePassword(password)

	if hashedPassword == "" {
		t.Error("GeneratePassword returned empty string")
	}

	if hashedPassword == password {
		t.Error("GeneratePassword returned unhashed password")
	}
}

func TestGeneratePassword_DifferentHashes(t *testing.T) {
	password := "testpassword123"
	hash1 := GeneratePassword(password)
	hash2 := GeneratePassword(password)

	// bcrypt should generate different hashes for the same password (due to salt)
	if hash1 == hash2 {
		t.Error("GeneratePassword should generate different hashes due to salt")
	}
}

func TestComparePassword_Success(t *testing.T) {
	password := "testpassword123"
	hashedPassword := GeneratePassword(password)

	err := ComparePassword(hashedPassword, password)
	if err != nil {
		t.Errorf("ComparePassword failed for correct password: %v", err)
	}
}

func TestComparePassword_Failure(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"
	hashedPassword := GeneratePassword(password)

	err := ComparePassword(hashedPassword, wrongPassword)
	if err == nil {
		t.Error("ComparePassword should fail for incorrect password")
	}
}

func TestComparePassword_EmptyPassword(t *testing.T) {
	password := "testpassword123"
	hashedPassword := GeneratePassword(password)

	err := ComparePassword(hashedPassword, "")
	if err == nil {
		t.Error("ComparePassword should fail for empty password")
	}
}
