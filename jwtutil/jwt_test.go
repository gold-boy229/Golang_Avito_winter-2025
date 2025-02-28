package jwtutil

import (
	"MerchShop/entities"
	"fmt"
	"testing"
)

const (
	TIMESTAMP_123        int64 = 123
	TIMESTAMP_123456     int64 = 123456
	TIMESTAMP_1740690120 int64 = 1740690120
	TIMESTAMP_2040690120 int64 = 2040690120

	TOKEN_NICK_1740690120     string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Ik5pY2siLCJleHAiOjE3NDA2OTAxMjB9.uOPMKALxOCm14Wf6QGdykQDMRgHGQEHSFt4WjpOyWyg"
	TOKEN_LEONARDO_1740690120 string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Ikxlb25hcmRvIiwiZXhwIjoxNzQwNjkwMTIwfQ.6K6xtEL2yVbxHtroAPHVd2bB5Da1NMD3uv2oDT1966w"
)

func TestGenerateTokenFor(t *testing.T) {
	tests := []struct {
		user          entities.User
		expiresAt     int64
		expectedToken string
	}{
		{
			entities.User{Username: "Nick"},
			TIMESTAMP_1740690120,
			TOKEN_NICK_1740690120,
		},
		{
			entities.User{Username: "Leonardo"},
			TIMESTAMP_1740690120,
			TOKEN_LEONARDO_1740690120,
		},
		{
			entities.User{Username: "Leonardo", Password: "JWT uses only username"},
			TIMESTAMP_1740690120,
			TOKEN_LEONARDO_1740690120,
		},
	}

	for _, test := range tests {
		gotToken, _ := generateTokenFor(test.user, expirationSetting{test.expiresAt, true})
		if gotToken != test.expectedToken {
			t.Errorf("user = %v\n expectedToken = %q\n gotToken = %q",
				test.user, test.expectedToken, gotToken)
		}
	}
}

func TestCalculateExpitarionTime(t *testing.T) {
	tests := []struct {
		expirationSetting expirationSetting
		expectedResult    int64
		wantError         bool
	}{
		{
			expirationSetting{expiresAt: TIMESTAMP_123, useGivenTime: true},
			TIMESTAMP_123,
			false,
		},
		{
			expirationSetting{expiresAt: TIMESTAMP_123, useGivenTime: false},
			TIMESTAMP_123,
			true,
		},
		{
			expirationSetting{expiresAt: TIMESTAMP_123456, useGivenTime: true},
			TIMESTAMP_123456,
			false,
		},
		{
			expirationSetting{expiresAt: TIMESTAMP_123, useGivenTime: true},
			TIMESTAMP_123456,
			true,
		},
	}

	for _, test := range tests {
		gotExpirationTime := calculateExpitarionTime(test.expirationSetting)
		if (gotExpirationTime != test.expectedResult) != test.wantError {
			t.Errorf("wantError = %v\n expectedResult = %v\n gotExpirationTime = %v",
				test.wantError, test.expectedResult, gotExpirationTime)
		}
	}
}

type verifyTokenStruct struct {
	providedToken        string
	expectUsername       string
	expectExpirationTime int64
	wantError            bool
}

func TestVerifyToken(t *testing.T) {
	tests := []verifyTokenStruct{
		{
			TOKEN_NICK_1740690120,
			"Nick",
			TIMESTAMP_1740690120,
			false,
		},
		{
			TOKEN_LEONARDO_1740690120,
			"Leonardo",
			TIMESTAMP_1740690120,
			false,
		},
		{
			TOKEN_LEONARDO_1740690120,
			"Not Leonardd",
			TIMESTAMP_1740690120,
			true,
		},
	}

	for _, test := range tests {
		gotClaims, tokenErr := VerifyToken(test.providedToken)
		err := areEqualGotAndExpectedClaims(gotClaims, test)
		if (err != nil) != test.wantError {
			t.Errorf("tokenErr = %v\n err = %v\n wantError = %v",
				tokenErr.Error(), err.Error(), test.wantError)
		}
	}
}

func areEqualGotAndExpectedClaims(gotClaims JWTClaim, test verifyTokenStruct) error {
	var (
		sameUsername       bool = (gotClaims.Username == test.expectUsername)
		sameExpirationTime bool = (gotClaims.StandardClaims.ExpiresAt == test.expectExpirationTime)
	)
	if sameUsername && sameExpirationTime {
		return nil
	} else {
		return fmt.Errorf("gotClaims.Username = %q; test.expectUsername = %q\n gotClaims.StandardClaims.ExpiresAt = %v; test.expectExpirationTime = %v",
			gotClaims.Username, test.expectUsername, gotClaims.StandardClaims.ExpiresAt, test.expectExpirationTime)
	}
}
