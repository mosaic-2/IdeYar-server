package util

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestGenerateVerificationCode(t *testing.T) {
	code := GenerateVerificationCode()
	codeLen := 6

	if len(code) != codeLen {
		t.Errorf("Expected length %d, got %d", codeLen, len(code))
	}

	for _, char := range code {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			t.Errorf("Unexpected character %c in code", char)
		}
	}
}

// func TestLoadVerificationEmail(t *testing.T) {
// 	code := "TESTCODE"
// 	expectedContents := "verification" // part of the template content for checking

// 	message, err := LoadVerificationEmail(code)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	if !bytes.Contains([]byte(message), []byte(expectedContents)) {
// 		t.Errorf("Expected message to contain %q, got %s", expectedContents, message)
// 	}
// }

// func TestLoadChangeMailEmail(t *testing.T) {
// 	code := "TESTCODE"
// 	expectedContents := "changeMail"

// 	message, err := LoadChangeMailEmail(code)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	if !bytes.Contains([]byte(message), []byte(expectedContents)) {
// 		t.Errorf("Expected message to contain %q, got %s", expectedContents, message)
// 	}
// }

// func TestLoadForgetPasswordEmail(t *testing.T) {
// 	code := "TESTCODE"
// 	expectedContents := "forgetPass"

// 	message, err := LoadForgetPasswordEmail(code)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	if !bytes.Contains([]byte(message), []byte(expectedContents)) {
// 		t.Errorf("Expected message to contain %q, got %s", expectedContents, message)
// 	}
// }

func TestGenerateFileName(t *testing.T) {
	name := GenerateFileName()
	if len(name) == 0 {
		t.Error("Expected non-empty file name")
	}
}

func TestGetUserIDFromCtx(t *testing.T) {
	md := metadata.New(map[string]string{"user-id": "12345"})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	userID := GetUserIDFromCtx(ctx)
	expectedUserID := int64(12345)

	if userID != expectedUserID {
		t.Errorf("Expected user ID %v, got %v", expectedUserID, userID)
	}
}

func TestGetUserIDFromCtx_NoMetadata(t *testing.T) {
	ctx := context.Background()

	userID := GetUserIDFromCtx(ctx)
	expectedUserID := int64(0)

	if userID != expectedUserID {
		t.Errorf("Expected user ID %v, got %v", expectedUserID, userID)
	}
}
