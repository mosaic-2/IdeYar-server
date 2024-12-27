package util

import (
	"github.com/golang-jwt/jwt/v5"
	"math/rand/v2"
	"strconv"
	"time"
)

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

func CreateChangeMailToken(userID int64, newMail string, duration time.Duration, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		Issuer:    "IdeYar",
		Subject:   strconv.FormatInt(userID, 10) + "$" + newMail,
		Audience:  jwt.ClaimStrings{"ChangeMail"},
	})

	return token.SignedString(key)
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
