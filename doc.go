// Package dropboxsign provides a Go client library for the Dropbox Sign API (formerly HelloSign).
//
// This package provides a type-safe interface to interact with the Dropbox Sign API
// for sending and managing signature requests with comprehensive error handling and
// support for various authentication methods.
//
// # Features
//
//   - Type-safe API interactions with comprehensive error handling
//   - Context-aware HTTP requests with timeout support
//   - Support for signature requests with templates
//   - Proper handling of API warnings and errors
//   - Builder patterns for complex request construction
//   - Idiomatic Go design with interfaces and error handling
//
// # Installation
//
// To install the package, run:
//
//	go get github.com/cjcox17/dropbox-sign-go
//
// # Quick Start
//
// Create a client and send a signature request:
//
//	package main
//
//	import (
//		"context"
//		"fmt"
//		"log"
//
//		"github.com/cjcox17/dropbox-sign-go"
//	)
//
//	func main() {
//		// Create a new client with your API key
//		client := dropboxsign.NewClient("your-api-key")
//
//		// Create a signer
//		signer := dropboxsign.NewSubSignatureRequestTemplateSigner(
//			"Signer",
//			"John Doe",
//			"john@example.com",
//		)
//
//		// Create and send a signature request
//		request := dropboxsign.NewSendSignatureRequest(
//			[]dropboxsign.SubSignatureRequestTemplateSigner{signer},
//			[]string{"template-id"},
//		).WithTitle("Contract Signature").WithTestMode(true)
//
//		ctx := context.Background()
//		response, warnings, err := client.SendWithTemplate(ctx, request)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		fmt.Printf("Signature request sent: %s\n", response.SignatureRequestID)
//	}
//
// # Client Configuration
//
// The client can be customized with various options:
//
//	// Set custom timeout
//	client := dropboxsign.NewClient("api-key").
//		WithTimeout(60 * time.Second)
//
//	// Use custom HTTP client
//	httpClient := &http.Client{
//		Timeout: 30 * time.Second,
//		Transport: customTransport,
//	}
//	client := dropboxsign.NewClient("api-key").
//		WithHTTPClient(httpClient)
//
//	// Set custom base URL (for testing)
//	client := dropboxsign.NewClient("api-key").
//		WithBaseURL("https://test.api.com/v3")
//
// # Working with Signature Requests
//
// Create a signature request with custom fields:
//
//	customFields := []dropboxsign.SubCustomField{
//		dropboxsign.NewSubCustomField("company_name").
//			WithValue("Acme Corp").
//			WithRequired(true),
//		dropboxsign.NewSubCustomField("employee_id").
//			WithValue("12345"),
//	}
//
//	request := dropboxsign.NewSendSignatureRequest(signers, templateIDs).
//		WithCustomFields(customFields).
//		WithMessage("Please review and sign this document").
//		WithMetadata(map[string]string{
//			"contract_id": "12345",
//			"department":  "HR",
//		})
//
// Add signers with authentication:
//
//	// Signer with PIN protection
//	signer := dropboxsign.NewSubSignatureRequestTemplateSigner(
//		"Signer",
//		"Jane Doe",
//		"jane@example.com",
//	).WithPin("1234")
//
//	// Signer with SMS authentication
//	signerWithSMS := dropboxsign.NewSubSignatureRequestTemplateSigner(
//		"Signer",
//		"Bob Smith",
//		"bob@example.com",
//	).WithSMSPhoneNumber("+1234567890").
//	  WithSMSPhoneNumberType(dropboxsign.SMSPhoneNumberTypeAuthentication)
//
// # Retrieving Signature Requests
//
// Retrieve a signature request by ID:
//
//	ctx := context.Background()
//	sigRequest, warnings, err := client.GetSignatureRequest(ctx, "signature_request_id")
//	if err != nil {
//		if dropboxsign.IsNotFound(err) {
//			log.Println("Signature request not found")
//			return
//		}
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Status: Complete=%v, Declined=%v\n",
//		sigRequest.IsComplete, sigRequest.IsDeclined)
//
// # Error Handling
//
// The library provides helper functions for common error scenarios:
//
//	sigRequest, _, err := client.GetSignatureRequest(ctx, signatureRequestID)
//	if err != nil {
//		switch {
//		case dropboxsign.IsNotFound(err):
//			// Handle 404 - signature request not found
//		case dropboxsign.IsUnauthorized(err):
//			// Handle 401 - invalid API key
//		case dropboxsign.IsBadRequest(err):
//			// Handle 400 - invalid request parameters
//		default:
//			// Handle other errors
//		}
//	}
//
// Error types:
//
//   - ErrorResponseError: API errors returned by Dropbox Sign with detailed error information
//   - ClientError: Client-side errors (network issues, parsing errors, etc.)
//
// # Context Support
//
// All API methods accept a context.Context parameter for timeout and cancellation:
//
//	// With timeout
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	sigRequest, _, err := client.GetSignatureRequest(ctx, signatureRequestID)
//
//	// With cancellation
//	ctx, cancel := context.WithCancel(context.Background())
//	go func() {
//		// Cancel after some condition
//		cancel()
//	}()
//
//	sigRequest, _, err := client.SendWithTemplate(ctx, request)
//
// # Warnings
//
// The API may return warnings alongside successful responses. These warnings
// indicate potential issues or important information that doesn't prevent
// the operation from completing:
//
//	response, warnings, err := client.SendWithTemplate(ctx, request)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if len(warnings) > 0 {
//		for _, warning := range warnings {
//			log.Printf("Warning: %s\n", warning)
//		}
//	}
//
// # Test Mode
//
// Use test mode to test your integration without sending real emails or
// incurring charges:
//
//	request := dropboxsign.NewSendSignatureRequest(signers, templateIDs).
//		WithTestMode(true)
//
// # Thread Safety
//
// The Client is safe for concurrent use by multiple goroutines. You can
// create a single client and reuse it across your application:
//
//	var client = dropboxsign.NewClient("your-api-key")
//
//	func handler1(w http.ResponseWriter, r *http.Request) {
//		// Safe to use client concurrently
//		sigRequest, _, err := client.GetSignatureRequest(r.Context(), id)
//		// ...
//	}
//
//	func handler2(w http.ResponseWriter, r *http.Request) {
//		// Also safe
//		sigRequest, _, err := client.SendWithTemplate(r.Context(), request)
//		// ...
//	}
//
// # API Reference
//
// For more information about the Dropbox Sign API, visit:
// https://developers.hellosign.com/api/reference/
//
// # License
//
// This project is dual-licensed under MIT OR Apache-2.0.
package dropboxsign
