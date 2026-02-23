package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStatusHandler(t *testing.T) {
	// Test case for the root path, which should return a help message.
	t.Run("root_path_help_message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		statusHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		expectedBody := "Welcome to http-debug!"
		if !strings.Contains(rr.Body.String(), expectedBody) {
			t.Errorf("expected body to contain %q, got %q", expectedBody, rr.Body.String())
		}
	})

	// Test case for a valid status code, which should return a JSON response.
	t.Run("valid_status_code", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/200", nil)
		req.Header.Set("X-Test-Header", "test-value")
		rr := httptest.NewRecorder()
		statusHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("expected Content-Type header to be 'application/json', got %q", contentType)
		}

		var data struct {
			StatusCode int         `json:"statusCode"`
			StatusText string      `json:"statusText"`
			Headers    http.Header `json:"headers"`
		}

		body, _ := io.ReadAll(rr.Body)
		if err := json.Unmarshal(body, &data); err != nil {
			t.Fatalf("could not unmarshal json: %v", err)
		}

		if data.StatusCode != http.StatusOK {
			t.Errorf("expected statusCode in json to be %d, got %d", http.StatusOK, data.StatusCode)
		}

		if data.StatusText != http.StatusText(http.StatusOK) {
			t.Errorf("expected statusText in json to be %q, got %q", http.StatusText(http.StatusOK), data.StatusText)
		}

		if data.Headers.Get("X-Test-Header") != "test-value" {
			t.Errorf("expected X-Test-Header to be 'test-value', got %q", data.Headers.Get("X-Test-Header"))
		}
	})

	// Test case for an invalid status code, which should return a 400 Bad Request.
	t.Run("invalid_status_code", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rr := httptest.NewRecorder()
		statusHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	// Test case for an unknown status code, which should still return a JSON response.
	t.Run("unknown_status_code", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/999", nil)
		rr := httptest.NewRecorder()
		statusHandler(rr, req)

		if rr.Code != 999 {
			t.Errorf("expected status code 999, got %d", rr.Code)
		}

		if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("expected Content-Type header to be 'application/json', got %q", contentType)
		}

		var data struct {
			StatusCode int    `json:"statusCode"`
			StatusText string `json:"statusText"`
		}

		body, _ := io.ReadAll(rr.Body)
		if err := json.Unmarshal(body, &data); err != nil {
			t.Fatalf("could not unmarshal json: %v", err)
		}

		if data.StatusCode != 999 {
			t.Errorf("expected statusCode in json to be 999, got %d", data.StatusCode)
		}

		if data.StatusText != "Unknown Status Code" {
			t.Errorf("expected statusText in json to be 'Unknown Status Code', got %q", data.StatusText)
		}
	})
}
