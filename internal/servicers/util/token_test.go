package util

import (
	"testing"
	"time"
)

func TestCreateAndParseLoginToken(t *testing.T) {
	secretKey := []byte("mySecretKey")
	userID := "12345"
	duration := time.Minute

	token, err := CreateLoginToken(userID, duration, secretKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	parsedUserID, err := ParseLoginToken(token, secretKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("expected userID to be %v, got %v", userID, parsedUserID)
	}
}

func TestCreateAndParseChangeMailToken(t *testing.T) {
	secretKey := []byte("mySecretKey")
	userID := int64(12345)
	newMail := "newemail@example.com"
	duration := time.Minute

	token, err := CreateChangeMailToken(userID, newMail, duration, secretKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	parsedUserID, parsedMail, err := ParseChangeMailToken(token, secretKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("expected userID to be %v, got %v", userID, parsedUserID)
	}

	if parsedMail != newMail {
		t.Errorf("expected mail to be %v, got %v", newMail, parsedMail)
	}
}

func TestCreateAndParseForgetPassToken(t *testing.T) {
	secretKey := []byte("mySecretKey")
	mail := "user@example.com"
	duration := time.Minute

	token, err := CreateForgetPassToken(mail, duration, secretKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	parsedMail, err := ParseForgetPassToken(token, secretKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if parsedMail != mail {
		t.Errorf("expected mail to be %v, got %v", mail, parsedMail)
	}
}
