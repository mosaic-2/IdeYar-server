package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand/v2"
	"net/http"
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

func ParseLoginToken(bearerToken string, hmacSecret []byte) (string, error) {
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})
	if err != nil {
		return "", err
	}

	return token.Claims.GetSubject()
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

func ParseChangeMailToken(token string) (userID int64, newMail string, err error) {

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
