package util

import (
	"context"
	"github.com/google/uuid"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func TestGenerateVerificationCode(t *testing.T) {
	code := GenerateVerificationCode()
	if len(code) != 6 {
		t.Errorf("Expected code of length 6, got %d", len(code))
	}
}

func TestCreateLoginToken(t *testing.T) {
	key := []byte("test_secret_key")
	_, err := CreateLoginToken("12345", time.Hour, key)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateRefreshToken(t *testing.T) {
	key := []byte("test_secret_key")
	_, err := CreateRefreshToken("12345", time.Hour, key)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGenerateFileName(t *testing.T) {
	fileName := GenerateFileName()
	if _, err := uuid.Parse(fileName); err != nil {
		t.Errorf("Expected valid UUID string, got error %v", err)
	}
}

func TestGetUserIDFromCtx(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("user-id", "42"))
	userID := GetUserIDFromCtx(ctx)
	expectedUserID := int64(42)
	if userID != expectedUserID {
		t.Errorf("Expected %d, got %d", expectedUserID, userID)
	}
}

func TestGenerateForgetPassToken(t *testing.T) {
	token := GenerateForgetPassToken()
	if len(token) != 20 {
		t.Errorf("Expected token of length 20, got %d", len(token))
	}
}
