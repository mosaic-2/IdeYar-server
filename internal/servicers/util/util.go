package util

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateVerificationCode() string {
	const charset = `ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
	//const charset = `0123456789`
	//for front test
	const codeLen = 6
	b := make([]byte, codeLen)
	for i := 0; i < codeLen; i++ {
		b[i] = charset[rand.Int()%len(charset)]
	}
	return string(b)
}

func CreateLoginToken(userID string, duration time.Duration, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		Issuer:    "KhanWeb",
		Subject:   userID,
		Audience:  jwt.ClaimStrings{"Login"},
	})

	return token.SignedString(key)
}

func CreateRefreshToken(userID string, duration time.Duration, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		Issuer:    "KhanWeb",
		Subject:   userID,
		Audience:  jwt.ClaimStrings{"Refresh"},
	})

	return token.SignedString(key)
}

func SendSignUpEmail(email string, code string) {
	username := "mmdhossein.haghdadi@gmail.com"
	password := "xsmtpsib-eb6248a76b82480199faf72cd07e43092f9d8c6ed89357698b5ac6a362171213-3aXDxUdMpNI01h9r"

	from := "no-reply@ideyar-app.ir"

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp-relay.brevo.com"
	smtpPort := "587"

	// Message.
	message, err := verificationEmail(code)
	if err != nil {
		fmt.Println(err)
		return
	}

	mimeHeaders := "MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n"

	// Email subject.
	header := fmt.Sprintf("From: no-reply@khanmedia.ir\r\nSubject: Email Verification\r\nTo: %s\r\n", email)

	// Putting together the email message with headers and body content.
	emailMessage := []byte(header + mimeHeaders + "\r\n" + message)

	// Authentication.
	auth := smtp.PlainAuth("", username, password, smtpHost)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, emailMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
}
