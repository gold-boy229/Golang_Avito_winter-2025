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

type expirationSetting struct {
	expiresAt    int64
	useGivenTime bool
}

func GenerateTokenFor(user entities.User) (string, error) {
	return generateTokenFor(user, expirationSetting{useGivenTime: false})
}

func generateTokenFor(user entities.User, expirationSetting expirationSetting) (string, error) {
	tokenClaims := JWTClaim{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: calculateExpitarionTime(expirationSetting),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func calculateExpitarionTime(expirationSetting expirationSetting) (expirationTime int64) {
	if expirationSetting.useGivenTime {
		expirationTime = expirationSetting.expiresAt
	} else {
		expirationTime = time.Now().Add(1 * time.Hour).Unix()
	}
	return expirationTime
}

func VerifyToken(tokenString string) (claims JWTClaim, err error) {
	_, err = jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		return claims, errors.New("Some error after jwt.ParseWithClaims\n" + err.Error())
	}

	return claims, nil
}
