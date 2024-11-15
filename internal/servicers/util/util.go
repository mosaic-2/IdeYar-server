package util

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mosaic-2/IdeYar-server/internal/config"
)

type ProfileIDCtx struct{}

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

func AuthMiddleware(next http.Handler) http.Handler {

	hmacSecret := []byte(config.GetSecretKey())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		allowedMethods := []string{
			"/auth/signup",
			"/auth/login",
			"/auth/code-verification",
		}

		for _, method := range allowedMethods {
			if r.RequestURI == method {
				next.ServeHTTP(w, r)
				return
			}
		}

		authHeader := r.Header.Get("Authorization")
		bearerToken := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			bearerToken = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			log.Printf("No Bearer Token found")
			return
		}

		token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return hmacSecret, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		profileIDStr, err := token.Claims.GetSubject()
		if err != nil {
			return
		}

		profileID, err := strconv.ParseInt(profileIDStr, 10, 64)
		if err != nil {
			return
		}

		ctx := context.WithValue(r.Context(), ProfileIDCtx{}, profileID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
