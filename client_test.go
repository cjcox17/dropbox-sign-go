package dropboxsign

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	if client.apiKey != apiKey {
		t.Errorf("expected apiKey %s, got %s", apiKey, client.apiKey)
	}

	if client.baseURL != APIBaseURL {
		t.Errorf("expected baseURL %s, got %s", APIBaseURL, client.baseURL)
	}

	if client.httpClient == nil {
		t.Error("expected httpClient to be initialized")
	}
}

func TestClientWithTimeout(t *testing.T) {
	client := NewClient("test-api-key").WithTimeout(60 * time.Second)

	if client.httpClient.Timeout != 60*time.Second {
		t.Errorf("expected timeout 60s, got %v", client.httpClient.Timeout)
	}
}

func TestClientWithBaseURL(t *testing.T) {
	customURL := "https://custom.api.com"
	client := NewClient("test-api-key").WithBaseURL(customURL)

	if client.baseURL != customURL {
		t.Errorf("expected baseURL %s, got %s", customURL, client.baseURL)
	}
}

func TestGetSignatureRequest_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/v3/signature_request/test-sig-req-id" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Verify authentication
		username, password, ok := r.BasicAuth()
		if !ok || username != "test-api-key" || password != "" {
			t.Error("invalid basic auth")
		}

		// Send mock response
		response := map[string]interface{}{
			"signature_request": map[string]interface{}{
				"signature_request_id": "test-sig-req-id",
				"title":                "Test Document",
				"original_title":       "Test Document",
				"is_complete":          false,
				"is_declined":          false,
				"has_error":            false,
				"files_url":            "https://example.com/files",
				"details_url":          "https://example.com/details",
				"cc_email_addresses":   []string{},
				"metadata":             map[string]string{},
				"created_at":           1234567890,
				"signatures": []map[string]interface{}{
					{
						"signature_id":         "sig-1",
						"signer_email_address": "test@example.com",
						"status_code":          "awaiting_signature",
						"has_pin":              false,
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClient("test-api-key").WithBaseURL(server.URL + "/v3")

	// Test GetSignatureRequest
	ctx := context.Background()
	sigRequest, warnings, err := client.GetSignatureRequest(ctx, "test-sig-req-id")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sigRequest == nil {
		t.Fatal("expected signature request, got nil")
	}

	if sigRequest.SignatureRequestID != "test-sig-req-id" {
		t.Errorf("expected signature_request_id 'test-sig-req-id', got %s", sigRequest.SignatureRequestID)
	}

	if sigRequest.Title != "Test Document" {
		t.Errorf("expected title 'Test Document', got %s", sigRequest.Title)
	}

	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %d", len(warnings))
	}
}

func TestGetSignatureRequest_NotFound(t *testing.T) {
	// Create a mock server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		response := ErrorResponse{
			Error: ErrorResponseError{
				ErrorMsg:  "Signature request not found",
				ErrorName: "not_found",
			},
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient("test-api-key").WithBaseURL(server.URL + "/v3")

	ctx := context.Background()
	sigRequest, warnings, err := client.GetSignatureRequest(ctx, "nonexistent-id")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if sigRequest != nil {
		t.Errorf("expected nil signature request, got %v", sigRequest)
	}

	if warnings != nil {
		t.Errorf("expected nil warnings, got %v", warnings)
	}

	if !IsNotFound(err) {
		t.Errorf("expected NotFound error, got %v", err)
	}
}

func TestSendWithTemplate_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/v3/signature_request/send_with_template" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Verify Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Parse request body
		var reqBody SendSignatureRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		if len(reqBody.Signers) != 1 {
			t.Errorf("expected 1 signer, got %d", len(reqBody.Signers))
		}

		// Send mock response
		response := map[string]interface{}{
			"signature_request": map[string]interface{}{
				"signature_request_id": "new-sig-req-id",
				"title":                "Test Document",
				"original_title":       "Test Document",
				"is_complete":          false,
				"is_declined":          false,
				"has_error":            false,
				"files_url":            "https://example.com/files",
				"details_url":          "https://example.com/details",
				"cc_email_addresses":   []string{},
				"metadata":             map[string]string{},
				"created_at":           1234567890,
				"signatures": []map[string]interface{}{
					{
						"signature_id":         "sig-1",
						"signer_email_address": "test@example.com",
						"status_code":          "awaiting_signature",
						"has_pin":              false,
					},
				},
			},
			"warnings": []map[string]string{
				{
					"warning_msg":  "Test warning",
					"warning_name": "test_warning",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClient("test-api-key").WithBaseURL(server.URL + "/v3")

	// Create request
	signer := NewSubSignatureRequestTemplateSigner("Signer", "John Doe", "john@example.com")
	request := NewSendSignatureRequest(
		[]SubSignatureRequestTemplateSigner{signer},
		[]string{"template-id"},
	).WithTitle("Test Document").WithTestMode(true)

	// Test SendWithTemplate
	ctx := context.Background()
	sigRequest, warnings, err := client.SendWithTemplate(ctx, request)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sigRequest == nil {
		t.Fatal("expected signature request, got nil")
	}

	if sigRequest.SignatureRequestID != "new-sig-req-id" {
		t.Errorf("expected signature_request_id 'new-sig-req-id', got %s", sigRequest.SignatureRequestID)
	}

	if len(warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(warnings))
	}

	if len(warnings) > 0 && warnings[0].WarningName != "test_warning" {
		t.Errorf("expected warning_name 'test_warning', got %s", warnings[0].WarningName)
	}
}

func TestSendWithTemplate_BadRequest(t *testing.T) {
	// Create a mock server that returns 400
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		response := ErrorResponse{
			Error: ErrorResponseError{
				ErrorMsg:  "Invalid template ID",
				ErrorName: "bad_request",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-api-key").WithBaseURL(server.URL + "/v3")

	signer := NewSubSignatureRequestTemplateSigner("Signer", "John Doe", "john@example.com")
	request := NewSendSignatureRequest(
		[]SubSignatureRequestTemplateSigner{signer},
		[]string{"invalid-template-id"},
	)

	ctx := context.Background()
	sigRequest, warnings, err := client.SendWithTemplate(ctx, request)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if sigRequest != nil {
		t.Errorf("expected nil signature request, got %v", sigRequest)
	}

	if warnings != nil {
		t.Errorf("expected nil warnings, got %v", warnings)
	}

	if !IsBadRequest(err) {
		t.Errorf("expected BadRequest error, got %v", err)
	}
}

func TestCancelIncompleteSignatureRequest_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/v3/signature_request/cancel/test-sig-req-id" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-api-key").WithBaseURL(server.URL + "/v3")

	ctx := context.Background()
	err := client.CancelIncompleteSignatureRequest(ctx, "test-sig-req-id")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseResponse(t *testing.T) {
	jsonData := []byte(`{
		"signature_request": {
			"signature_request_id": "test-id",
			"title": "Test",
			"original_title": "Test",
			"is_complete": false,
			"is_declined": false,
			"has_error": false,
			"files_url": "https://example.com",
			"details_url": "https://example.com",
			"cc_email_addresses": [],
			"metadata": {},
			"created_at": 1234567890,
			"signatures": []
		},
		"warnings": [
			{
				"warning_msg": "Test warning",
				"warning_name": "test_warning"
			}
		]
	}`)

	sigRequest, warnings, err := parseResponse[SignatureRequestResponse](jsonData, "signature_request")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sigRequest == nil {
		t.Fatal("expected signature request, got nil")
	}

	if sigRequest.SignatureRequestID != "test-id" {
		t.Errorf("expected signature_request_id 'test-id', got %s", sigRequest.SignatureRequestID)
	}

	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}

	if warnings[0].WarningName != "test_warning" {
		t.Errorf("expected warning_name 'test_warning', got %s", warnings[0].WarningName)
	}
}

func TestParseResponse_MissingKey(t *testing.T) {
	jsonData := []byte(`{"other_key": {}}`)

	_, _, err := parseResponse[SignatureRequestResponse](jsonData, "signature_request")

	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestErrorResponseError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      ErrorResponseError
		expected string
	}{
		{
			name: "with path",
			err: ErrorResponseError{
				ErrorName: "validation_error",
				ErrorPath: stringPtr("signers[0].email"),
				ErrorMsg:  "Invalid email address",
			},
			expected: "validation_error (signers[0].email): Invalid email address",
		},
		{
			name: "without path",
			err: ErrorResponseError{
				ErrorName: "not_found",
				ErrorMsg:  "Resource not found",
			},
			expected: "not_found: Resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestWarningResponse_String(t *testing.T) {
	warning := WarningResponse{
		WarningMsg:  "This is a test warning",
		WarningName: "test_warning",
	}

	expected := "This is a test warning (test_warning)"
	if got := warning.String(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name: "ErrorResponseError with 404",
			err: ErrorResponseError{
				Status:    http.StatusNotFound,
				ErrorName: "not_found",
				ErrorMsg:  "Not found",
			},
			expected: true,
		},
		{
			name: "ClientError with 404",
			err: &ClientError{
				StatusCode: http.StatusNotFound,
				Message:    "Not found",
			},
			expected: true,
		},
		{
			name: "ErrorResponseError with 400",
			err: ErrorResponseError{
				Status:    http.StatusBadRequest,
				ErrorName: "bad_request",
				ErrorMsg:  "Bad request",
			},
			expected: false,
		},
		{
			name:     "other error",
			err:      http.ErrServerClosed,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFound(tt.err); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
