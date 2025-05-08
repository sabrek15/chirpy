package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPasword_Success(t *testing.T) {
	password := "MySecurePassword123"
	hashed, err := HashPasword(password)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if hashed == "" {
		t.Fatal("Expected a hashed password, got an empty string")
	}
}

func TestCheckPassword_Success(t *testing.T) {
	password := "AnotherSecurePass!"
	hashed, err := HashPasword(password)
	if err != nil {
		t.Fatalf("Hashing failed: %v", err)
	}

	err = CheckPassword(hashed, password)
	if err != nil {
		t.Errorf("Expected password to match, got error: %v", err)
	}
}

func TestCheckPassword_Failure(t *testing.T) {
	password := "CorrectPassword"
	wrongPassword := "WrongPassword"
	hashed, err := HashPasword(password)
	if err != nil {
		t.Fatalf("Hashing failed: %v", err)
	}

	err = CheckPassword(hashed, wrongPassword)
	if err == nil {
		t.Error("Expected error for mismatched password, got nil")
	}
}

func TestEmptyPassword(t *testing.T) {
	empty := ""
	hashed, err := HashPasword(empty)
	if err != nil {
		t.Fatalf("Hashing failed for empty password: %v", err)
	}

	err = CheckPassword(hashed, empty)
	if err != nil {
		t.Errorf("Expected empty password to match, got error: %v", err)
	}
}

func TestMakeAndValidateJWT_Success(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := time.Minute

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if parsedID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, parsedID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := -1 * time.Minute // already expired

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	correctSecret := "correct-secret"
	wrongSecret := "wrong-secret"
	expiresIn := time.Minute

	token, err := MakeJWT(userID, correctSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("Expected error for token signed with wrong secret, got nil")
	}
}

func TestValidateJWT_InvalidTokenFormat(t *testing.T) {
	badToken := "this.is.not.a.valid.jwt"
	_, err := ValidateJWT(badToken, "some-secret")
	if err == nil {
		t.Error("Expected error for malformed token, got nil")
	}
}