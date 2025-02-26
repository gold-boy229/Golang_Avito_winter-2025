package jwtutil

import (
	"MerchShop/entities"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("supersecretkey")

type JWTClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateTokenFor(user entities.User) (string, error) {
	signedToken, err := getSignedToken(user)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func getSignedToken(user entities.User) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	tokenClaims := JWTClaim{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	return token.SignedString(jwtKey)
}

func VerifyToken(tokenString string) (claims JWTClaim, err error) {
	token, err := parseSignedToken(tokenString)
	if err != nil {
		return claims, err
	}

	if !isValid(token) {
		return claims, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(JWTClaim)
	if !ok {
		return claims, errors.New("incorrect token claims type")
	}

	return claims, nil
}

func parseSignedToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	return token, err
}

func isValid(token *jwt.Token) bool {
	return token.Valid
}
