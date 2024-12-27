package util

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type UserIDCtxKey struct{}

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
		Issuer:    "IdeYar",
		Subject:   userID,
		Audience:  jwt.ClaimStrings{"Login"},
	})

	return token.SignedString(key)
}

func CreateRefreshToken(userID string, duration time.Duration, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		Issuer:    "IdeYar",
		Subject:   userID,
		Audience:  jwt.ClaimStrings{"Refresh"},
	})

	return token.SignedString(key)
}

func CreateChangeMailToken(userID string, old string, new string, duration time.Duration, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		Issuer:    "IdeYar",
		Subject:   userID,
		Audience:  jwt.ClaimStrings{"Refresh"},
	})

	return token.SignedString(key)
}

func SendSignUpEmail(email string, signeUpToken string, code string) {
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
	message, err := verificationEmail(fmt.Sprintf("localhost:3000/code-veification/%s/%s", signeUpToken, code))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("SignUp link:\n localhost:3000/code-veification/%s/%s", signeUpToken, code)

	mimeHeaders := "MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n"

	// Email subject.
	header := fmt.Sprintf("From: no-reply@ideyar-app.ir\r\nSubject: Email Verification\r\nTo: %s\r\n", email)

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

func verificationEmail(code string) (string, error) {

	tmp, err := template.ParseFiles("./internal/servicers/util/templates/verification.gohtml")
	if err != nil {
		return "", err
	}

	b := new(bytes.Buffer)

	err = tmp.Execute(b, code)
	if err != nil {
		return "", err
	}

	message := string(b.Bytes())

	return message, nil
}

func GenerateFileName() string {
	return uuid.New().String()
}

func GetUserIDFromCtx(ctx context.Context) int64 {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0
	}

	userID, _ := strconv.ParseInt(md["user-id"][0], 10, 64)
	return int64(userID)
}

func SendForgetPasswordEmail(ctx context.Context, email string, token string) {
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
	message, err := verificationEmail(fmt.Sprintf("localhost:3000/forget-password/%s", token))
	if err != nil {
		fmt.Println(err)
		return
	}

	mimeHeaders := "MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n"

	// Email subject.
	header := fmt.Sprintf("From: no-reply@ideyar-app.ir\r\nSubject: Email Verification\r\nTo: %s\r\n", email)

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

func GenerateForgetPassToken() string {
	const charset = `ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
	//const charset = `0123456789`
	//for front test
	const codeLen = 20
	b := make([]byte, codeLen)
	for i := 0; i < codeLen; i++ {
		b[i] = charset[rand.Int()%len(charset)]
	}
	return string(b)
}
