package controllers

import (
	"MerchShop/jwtutil"
	"net/http"
	"strings"
)

const (
	USERNAME string = "username"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			tokenString := r.Header.Get("Authorization")
			if len(tokenString) == 0 {
				respondError(w, http.StatusUnauthorized, "Missing Authorization Header")
				return
			}
			tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

			claims, err := jwtutil.VerifyToken(tokenString)
			if err != nil {
				respondError(w, http.StatusUnauthorized, "Error verifying JWT token: "+err.Error())
				return
			}

			r.Header.Set(USERNAME, claims.Username)
			next.ServeHTTP(w, r)
		},
	)
}
