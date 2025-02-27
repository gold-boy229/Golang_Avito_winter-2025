package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsEmptyString(t *testing.T) {
	tests := []struct {
		providedString string
		expectedResult bool
	}{
		{"", true},
		{"not_empty", false},
		{"1", false},
	}

	for _, test := range tests {
		gotResult := isEmptyString(test.providedString)
		if gotResult != test.expectedResult {
			t.Errorf("providedString = %q; expectedResult = %v; gotResult = %v",
				test.providedString, test.expectedResult, gotResult)
		}
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		providedString string
		expectedResult string
	}{
		{"Bearer token", "token"},
		{"Bearer xx.yy.zz", "xx.yy.zz"},
		{"Carrier t", "Carrier t"},
	}

	for _, test := range tests {
		gotResult := getBearerToken(test.providedString)
		if gotResult != test.expectedResult {
			t.Errorf("providedString = %q; expectedResult = %q; gotResult = %q",
				test.providedString, test.expectedResult, gotResult)
		}
	}
}

func TestSetUsernameIntoHeader(t *testing.T) {
	request := createRequest()

	tests := []struct {
		providedUsername string
	}{
		{"Nick"},
		{"123"},
		{""},
	}

	for _, test := range tests {
		setUsernameIntoHeader(request, test.providedUsername)
		gotResult := getUsernameFromHeader(request)
		if gotResult != test.providedUsername {
			t.Errorf("providedUsername = %q; expectedResult = %q; gotResult = %q",
				test.providedUsername, test.providedUsername, gotResult)
		}
	}
}

func createRequest() *http.Request {
	return httptest.NewRequest("GET", "http://example.com", nil)
}
