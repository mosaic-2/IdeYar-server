package util

import (
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		username string
		expected bool
	}{
		{"john_doe", true},
		{"john", true},
		{"jo", false},       // too short
		{"john@doe", false}, // invalid character
	}

	for _, test := range tests {
		result := ValidateUsername(test.username)
		if result != test.expected {
			t.Errorf("ValidateUsername(%q) = %v; expected %v", test.username, result, test.expected)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"johndoe@gmail.com", true},
		{"john.doe@domain.co.uk", true},
		{"john.doe@", false},            // missing domain
		{"john.doe@domain.c", false},    // domain TLD too short
		{"john.doe@domain..com", false}, // double dot
	}

	for _, test := range tests {
		result := ValidateEmail(test.email)
		if result != test.expected {
			t.Errorf("ValidateEmail(%q) = %v; expected %v", test.email, result, test.expected)
		}
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"John Doe", true},
		{"Anna", true},
		{"Jo", false},       // too short
		{"J0hn D0e", false}, // numbers not allowed
		{"John Doe John Doe John Doe John D", false}, // too long
	}

	for _, test := range tests {
		result := ValidateName(test.name)
		if result != test.expected {
			t.Errorf("ValidateName(%q) = %v; expected %v", test.name, result, test.expected)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"Password1!", true},
		{"Password", false},    // missing digit and special char
		{"password1!", false},  // missing uppercase letter
		{"PASSWORD1!", false},  // missing lowercase letter
		{"Password123", false}, // missing special char
		{"P@1a", false},        // too short
	}

	for _, test := range tests {
		result := ValidatePassword(test.password)
		if result != test.expected {
			t.Errorf("ValidatePassword(%q) = %v; expected %v", test.password, result, test.expected)
		}
	}
}
