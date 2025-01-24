package interceptor

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

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
			"/auth/forget-password",
			"/auth/forget-password-finalize",
			"/api/post",
			""
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

		userIDStr, err := util.ParseLoginToken(bearerToken, hmacSecret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
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
