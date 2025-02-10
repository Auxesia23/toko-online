package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/Auxesia23/toko-online/internal/utils"
	"github.com/golang-jwt/jwt"
)

type ContextKey string

const UserEmailContextKey ContextKey = "userEmail"

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := utils.VerifyJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		email := claims["email"]
		ctx := context.WithValue(r.Context(), UserEmailContextKey, email)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func SuperUserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := utils.VerifyJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		is_superuser := claims["is_superuser"]
		log.Println(is_superuser)
		if is_superuser != true {
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}
		email := claims["email"]

		ctx := context.WithValue(r.Context(), UserEmailContextKey, email)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
