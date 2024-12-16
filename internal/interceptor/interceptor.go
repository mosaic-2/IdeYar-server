package interceptor

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mosaic-2/IdeYar-server/internal/config"
	"github.com/mosaic-2/IdeYar-server/internal/servicers/util"
)

func AuthMiddleware(next http.Handler) http.Handler {

	hmacSecret := []byte(config.GetSecretKey())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		allowedMethods := []string{
			"/auth/signup",
			"/auth/login",
			"/auth/code-verification",
			"/liveness/checkliveness",
			"/api/image",
			"/api/landing-posts",
			"/api/search-post",
		}

		for _, method := range allowedMethods {
			if strings.HasPrefix(r.RequestURI, method) {
				next.ServeHTTP(w, r)
				return
			}
		}

		authHeader := r.Header.Get("Authorization")
		bearerToken := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			bearerToken = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			w.WriteHeader(http.StatusUnauthorized)
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

		userIDStr, err := token.Claims.GetSubject()
		if err != nil {
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return
		}

		ctx := context.WithValue(r.Context(), util.UserIDCtxKey{}, userID)
		r = r.WithContext(ctx)

		r.Header.Set("x-user-id", userIDStr)

		next.ServeHTTP(w, r)
	})
}
