package jwtutil

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("supersecretkey")

type JWTClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: " + err.Error())
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*JWTClaim, error) {
	claims := &JWTClaim{}

	// Parse the token with claims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC-SHA256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		// Handle specific JWT validation errors
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("malformed token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, errors.New("token is either expired or not active yet")
			} else {
				return nil, fmt.Errorf("token validation error: %v", ve)
			}
		}
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Ensure the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Ensure the claims are of the correct type
	if _, ok := token.Claims.(*JWTClaim); !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
