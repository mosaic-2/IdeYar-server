package util

import "regexp"

func ValidateUsername(username string) bool {
	var usernameRegex = regexp.MustCompile(`^[a-zA-z0-9_-]{3,32}$`)
	return usernameRegex.MatchString(username)
}

func ValidateEmail(mail string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(mail)
}

func ValidateName(name string) bool {
	var nameRegex = regexp.MustCompile(`^[a-zA-Z ]{3,40}$`)
	return nameRegex.MatchString(name)
}

func ValidatePassword(password string) bool {
	var lowerChar = regexp.MustCompile(`[a-z]`)
	var upperChar = regexp.MustCompile(`[A-Z]`)
	var digit = regexp.MustCompile(`\d`)
	var specialChar = regexp.MustCompile(`[!@#$%^&*_]`)
	var length = regexp.MustCompile(`^.{8,72}$`)
	return lowerChar.MatchString(password) && upperChar.MatchString(password) && digit.MatchString(password) && specialChar.MatchString(password) && length.MatchString(password)
}
