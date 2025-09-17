package auth

import (
        "net/http"
        "testing"
)

func TestGetAPIKey_ValidHeader(t *testing.T) {
        // Arrange: create a header with a valid API key
        headers := http.Header{}
        headers.Set("Authorization", "ApiKey supersecret")

        // Act: call the function
        got, err := GetAPIKey(headers)

        // Assert: expect no error and the correct key
        if err != nil {
                t.Fatalf("unexpected error: %v", err)
        }
        want := "supersecret"
        if got != want {
                t.Errorf("got %q, want %q", got, want)
        }
}

func TestGetAPIKey_NoHeader(t *testing.T) {
        // Arrange: create an empty header
        headers := http.Header{}

        // Act: call the function
        _, err := GetAPIKey(headers)
