// Package dropboxsign provides a Go client library for the Dropbox Sign API (formerly HelloSign).
//
// This package provides a type-safe interface to interact with the Dropbox Sign API
// for sending and managing signature requests.
//
// Example:
//
//	client := dropboxsign.NewClient("your-api-key")
//
//	signer := dropboxsign.NewSubSignatureRequestTemplateSigner(
//		"Signer",
//		"John Doe",
//		"john@example.com",
//	)
//
//	request := dropboxsign.NewSendSignatureRequest(
//		[]dropboxsign.SubSignatureRequestTemplateSigner{signer},
//		[]string{"template-id"},
//	)
//
//	response, warnings, err := client.SendWithTemplate(context.Background(), request)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Signature request sent: %s\n", response.SignatureRequestID)
package dropboxsign

import (
	"fmt"
	"net/http"
)

// ResponseWithWarnings wraps API responses that may contain warnings alongside the main data.
//
// The Inner field contains the actual response data, while Warnings contains
// any non-fatal warnings returned by the API.
type ResponseWithWarnings[T any] struct {
	// Inner is the main response data
	Inner T `json:"-"`
	// Warnings contains optional warnings returned by the API
	Warnings []WarningResponse `json:"warnings,omitempty"`
}

// WarningResponse represents a non-fatal warning returned by the Dropbox Sign API.
//
// Warnings indicate potential issues or important information that doesn't
// prevent the operation from completing successfully.
type WarningResponse struct {
	// WarningMsg is the human-readable warning message
	WarningMsg string `json:"warning_msg"`
	// WarningName is the machine-readable warning identifier
	WarningName string `json:"warning_name"`
}

// String returns a formatted string representation of the warning.
func (w WarningResponse) String() string {
	return fmt.Sprintf("%s (%s)", w.WarningMsg, w.WarningName)
}

// ErrorResponse is the top-level error response structure from the Dropbox Sign API.
type ErrorResponse struct {
	// Error contains the detailed error information
	Error ErrorResponseError `json:"error"`
}

// ErrorResponseError contains detailed error information from the Dropbox Sign API.
//
// Contains structured error details including HTTP status codes,
// error messages, and optional path information for field-specific errors.
type ErrorResponseError struct {
	// Status is the HTTP status code
	Status int `json:"-"`
	// ErrorMsg is the human-readable error message
	ErrorMsg string `json:"error_msg"`
	// ErrorPath is the optional path to the field that caused the error
	ErrorPath *string `json:"error_path,omitempty"`
	// ErrorName is the machine-readable error identifier
	ErrorName string `json:"error_name"`
}

// Error implements the error interface for ErrorResponseError.
func (e ErrorResponseError) Error() string {
	if e.ErrorPath != nil && *e.ErrorPath != "" {
		return fmt.Sprintf("%s (%s): %s", e.ErrorName, *e.ErrorPath, e.ErrorMsg)
	}
	return fmt.Sprintf("%s: %s", e.ErrorName, e.ErrorMsg)
}

// ClientError wraps errors that occur when using the Dropbox Sign client.
type ClientError struct {
	// Message is the error message
	Message string
	// StatusCode is the HTTP status code (if applicable)
	StatusCode int
	// Err is the underlying error (if any)
	Err error
}

// Error implements the error interface for ClientError.
func (e *ClientError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("dropboxsign client error (status %d): %s", e.StatusCode, e.Message)
	}
	if e.Err != nil {
		return fmt.Sprintf("dropboxsign client error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("dropboxsign client error: %s", e.Message)
}

// Unwrap returns the underlying error.
func (e *ClientError) Unwrap() error {
	return e.Err
}

// NewClientError creates a new ClientError.
func NewClientError(message string, statusCode int, err error) *ClientError {
	return &ClientError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// IsNotFound returns true if the error is a 404 Not Found error.
func IsNotFound(err error) bool {
	if apiErr, ok := err.(ErrorResponseError); ok {
		return apiErr.Status == http.StatusNotFound
	}
	if clientErr, ok := err.(*ClientError); ok {
		return clientErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsBadRequest returns true if the error is a 400 Bad Request error.
func IsBadRequest(err error) bool {
	if apiErr, ok := err.(ErrorResponseError); ok {
		return apiErr.Status == http.StatusBadRequest
	}
	if clientErr, ok := err.(*ClientError); ok {
		return clientErr.StatusCode == http.StatusBadRequest
	}
	return false
}

// IsUnauthorized returns true if the error is a 401 Unauthorized error.
func IsUnauthorized(err error) bool {
	if apiErr, ok := err.(ErrorResponseError); ok {
		return apiErr.Status == http.StatusUnauthorized
	}
	if clientErr, ok := err.(*ClientError); ok {
		return clientErr.StatusCode == http.StatusUnauthorized
	}
	return false
}
