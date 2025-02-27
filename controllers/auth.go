package handlers

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
			authorizationString := r.Header.Get("Authorization")
			if isEmptyString(authorizationString) {
				respondUnauthorized(w, "Missing Authorization Header")
				return
			}

			token := getBearerToken(authorizationString)

			claims, err := jwtutil.VerifyToken(token)
			if err != nil {
				respondUnauthorized(w, "Error verifying JWT token: "+err.Error())
				return
			}

			setUsernameIntoHeader(r, claims.Username)
			next.ServeHTTP(w, r)
		},
	)
}

func isEmptyString(s string) bool {
	return len(s) == 0
}

func getBearerToken(authorizationString string) string {
	return strings.Replace(authorizationString, "Bearer ", "", 1)
}

func setUsernameIntoHeader(r *http.Request, username string) {
	r.Header.Set(USERNAME, username)
}

func getUsernameFromHeader(r *http.Request) (username string) {
	return r.Header.Get(USERNAME)
}
