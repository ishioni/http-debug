package main

import (
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

	// Test case for a valid status code, which should return a human-readable response.
	t.Run("valid_status_code", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/200", nil)
		req.Header.Set("X-Test-Header", "test-value")
		rr := httptest.NewRecorder()
		statusHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		expectedBody := "Status Code: 200 OK\nHeaders:\n  X-Test-Header: test-value"
		if !strings.Contains(rr.Body.String(), expectedBody) {
			t.Errorf("expected body to contain %q, got %q", expectedBody, rr.Body.String())
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

	// Test case for an unknown status code, which should return a human-readable response.
	t.Run("unknown_status_code", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/999", nil)
		rr := httptest.NewRecorder()
		statusHandler(rr, req)

		if rr.Code != 999 {
			t.Errorf("expected status code 999, got %d", rr.Code)
		}

		expectedBody := "Status Code: 999 Unknown Status Code\nHeaders:\n"
		if !strings.Contains(rr.Body.String(), expectedBody) {
			t.Errorf("expected body to contain %q, got %q", expectedBody, rr.Body.String())
		}
	})
}
