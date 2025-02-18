package jwtutil

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestGenerateJWT(t *testing.T) {
	username := "testuser"
	tokenString, err := GenerateJWT(username)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the token's signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorUnverifiable)
		}
		return jwtKey, nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	// Check claims
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		t.Fatalf("Expected claims of type *JWTClaim, got %T", token.Claims)
	}

	if !token.Valid {
		t.Fatalf("Token is invalid. Claims: %+v", claims)
	}

	if claims.Username != username {
		t.Errorf("Expected username %s, got %s", username, claims.Username)
	}

	// Check expiration time
	if claims.ExpiresAt < time.Now().Unix() {
		t.Errorf("Token has expired. ExpiresAt: %d, Current time: %d", claims.ExpiresAt, time.Now().Unix())
	}
}

func TestVerifyToken(t *testing.T) {
	// Helper function to generate a valid token
	generateValidToken := func(username string) string {
		claims := &JWTClaim{
			Username: username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			t.Fatalf("Failed to generate valid token: %v", err)
		}
		return tokenString
	}

	// Test valid token
	t.Run("Valid token", func(t *testing.T) {
		tokenString := generateValidToken("testuser")
		claims, err := VerifyToken(tokenString)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if claims.Username != "testuser" {
			t.Errorf("Expected username 'testuser', got '%s'", claims.Username)
		}
	})

	// Test expired token
	t.Run("Expired token", func(t *testing.T) {
		claims := &JWTClaim{
			Username: "testuser",
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			t.Fatalf("Failed to generate expired token: %v", err)
		}

		_, err = VerifyToken(tokenString)
		if err == nil {
			t.Fatal("Expected error for expired token, got none")
		}
		if err.Error() != "token is either expired or not active yet" {
			t.Errorf("Expected 'token is either expired or not active yet', got '%v'", err)
		}
	})

	// Test malformed token
	t.Run("Malformed token", func(t *testing.T) {
		_, err := VerifyToken("malformed.token.string")
		if err == nil {
			t.Fatal("Expected error for malformed token, got none")
		}
		if err.Error() != "malformed token" {
			t.Errorf("Expected 'malformed token', got '%v'", err)
		}
	})

	// Test invalid key
	t.Run("Invalid key", func(t *testing.T) {
		tokenString := generateValidToken("testuser")
		// Temporarily change the global key to an invalid one
		originalKey := jwtKey
		jwtKey = []byte("wrongkey")
		defer func() { jwtKey = originalKey }() // Restore the original key

		_, err := VerifyToken(tokenString)
		if err == nil {
			t.Fatal("Expected error for invalid key, got none")
		}
	})
}
