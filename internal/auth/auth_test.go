package auth

import (
	"net/http"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "ApiKey 12345")

	key, err := GetAPIKey(headers)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if key != "12345" {
		t.Errorf("Expected key '12345', got '%s'", key)
	}
}

func TestGetAPIKeyNoHeader(t *testing.T) {
	headers := http.Header{}

	_, err := GetAPIKey(headers)
	if err != ErrNoAuthHeaderIncluded {
		t.Errorf("Expected ErrNoAuthHeaderIncluded, got %v", err)
	}
}
