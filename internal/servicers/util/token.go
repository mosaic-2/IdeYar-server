package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"strings"
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

func ParseChangeMailToken(token string, hmacSecret []byte) (userID int64, newMail string, err error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})
	if err != nil {
		return 0, "", err
	}

	joinedStr, err := jwtToken.Claims.GetSubject()
	if err != nil {
		return 0, "", err
	}

	strArray := strings.Split(joinedStr, "$")

	userID, err = strconv.ParseInt(strArray[0], 10, 64)
	if err != nil {
		return 0, "", err
	}

	return userID, strArray[1], nil
}

func CreateForgetPassToken(mail string, duration time.Duration, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		Issuer:    "IdeYar",
		Subject:   mail,
		Audience:  jwt.ClaimStrings{"ForgetPassword"},
	})

	return token.SignedString(key)

}

func ParseForgetPassToken(token string, key []byte) (string, error) {
	jwToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return "", err
	}

	return jwToken.Claims.GetSubject()

}
