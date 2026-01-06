// Package dropboxsign provides HTTP client implementation for the Dropbox Sign API.
package dropboxsign

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// APIBaseURL is the base URL for the Dropbox Sign API (v3)
	APIBaseURL = "https://api.hellosign.com/v3"
	// DefaultTimeout is the default request timeout
	DefaultTimeout = 30 * time.Second
)

// Client is an HTTP client for interacting with the Dropbox Sign API.
//
// This client handles authentication, request/response processing, and error handling
// for Dropbox Sign API operations.
//
// Example:
//
//	client := dropboxsign.NewClient("your-api-key").
//		WithTimeout(60 * time.Second)
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Dropbox Sign client with the specified API key.
//
// The client is configured with sensible defaults including a 30-second timeout
// and connection pooling.
//
// Example:
//
//	client := dropboxsign.NewClient("your-api-key")
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: APIBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 5,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// WithTimeout sets a custom timeout for HTTP requests.
//
// Returns the client instance for method chaining.
//
// Example:
//
//	client := dropboxsign.NewClient("api-key").WithTimeout(60 * time.Second)
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.httpClient.Timeout = timeout
	return c
}

// WithHTTPClient sets a custom HTTP client.
//
// This allows for advanced configuration of the underlying HTTP transport.
//
// Returns the client instance for method chaining.
func (c *Client) WithHTTPClient(httpClient *http.Client) *Client {
	c.httpClient = httpClient
	return c
}

// WithBaseURL sets a custom base URL for the API.
//
// This is primarily useful for testing against mock servers.
//
// Returns the client instance for method chaining.
func (c *Client) WithBaseURL(baseURL string) *Client {
	c.baseURL = baseURL
	return c
}

// GetSignatureRequest retrieves a signature request by its ID.
//
// Returns the signature request data and any warnings, or an error
// if the request fails or the signature request is not found.
//
// Example:
//
//	ctx := context.Background()
//	sigRequest, warnings, err := client.GetSignatureRequest(ctx, "signature_request_id")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Title: %s\n", sigRequest.Title)
func (c *Client) GetSignatureRequest(ctx context.Context, signatureRequestID string) (*SignatureRequestResponse, []WarningResponse, error) {
	url := fmt.Sprintf("%s/signature_request/%s", c.baseURL, signatureRequestID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, NewClientError("failed to create request", 0, err)
	}

	req.SetBasicAuth(c.apiKey, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, NewClientError("failed to execute request", 0, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, NewClientError("failed to read response body", resp.StatusCode, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, c.parseErrorResponse(body, resp.StatusCode)
	}

	sigRequest, warnings, err := parseResponse[SignatureRequestResponse](body, "signature_request")
	if err != nil {
		return nil, nil, NewClientError("failed to parse response", resp.StatusCode, err)
	}

	return sigRequest, warnings, nil
}

// SendWithTemplate sends a signature request using a template.
//
// This method creates and sends a signature request based on a pre-existing
// template. The template defines the document layout and form fields, while
// the request specifies the signers and other dynamic parameters.
//
// Returns the created signature request data and any warnings, or an error
// if the request fails.
//
// Example:
//
//	ctx := context.Background()
//	signer := dropboxsign.NewSubSignatureRequestTemplateSigner(
//		"Signer",
//		"John Doe",
//		"john@example.com",
//	)
//	request := dropboxsign.NewSendSignatureRequest(
//		[]dropboxsign.SubSignatureRequestTemplateSigner{signer},
//		[]string{"template-id"},
//	).WithTitle("Contract Signature").WithTestMode(true)
//
//	sigRequest, warnings, err := client.SendWithTemplate(ctx, request)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sent: %s\n", sigRequest.SignatureRequestID)
func (c *Client) SendWithTemplate(ctx context.Context, request *SendSignatureRequest) (*SignatureRequestResponse, []WarningResponse, error) {
	url := fmt.Sprintf("%s/signature_request/send_with_template", c.baseURL)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, nil, NewClientError("failed to marshal request", 0, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, NewClientError("failed to create request", 0, err)
	}

	req.SetBasicAuth(c.apiKey, "")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, NewClientError("failed to execute request", 0, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, NewClientError("failed to read response body", resp.StatusCode, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, c.parseErrorResponse(body, resp.StatusCode)
	}

	sigRequest, warnings, err := parseResponse[SignatureRequestResponse](body, "signature_request")
	if err != nil {
		return nil, nil, NewClientError("failed to parse response", resp.StatusCode, err)
	}

	return sigRequest, warnings, nil
}

// CancelIncompleteSignatureRequest cancels an incomplete signature request.
//
// This can only be used on signature requests that have not been completed
// by all signers.
//
// Returns an error if the request fails or the signature request cannot be cancelled.
//
// Example:
//
//	ctx := context.Background()
//	err := client.CancelIncompleteSignatureRequest(ctx, "signature_request_id")
//	if err != nil {
//		log.Fatal(err)
//	}
func (c *Client) CancelIncompleteSignatureRequest(ctx context.Context, signatureRequestID string) error {
	url := fmt.Sprintf("%s/signature_request/cancel/%s", c.baseURL, signatureRequestID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return NewClientError("failed to create request", 0, err)
	}

	req.SetBasicAuth(c.apiKey, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return NewClientError("failed to execute request", 0, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return NewClientError("failed to read error response body", resp.StatusCode, err)
		}
		return c.parseErrorResponse(body, resp.StatusCode)
	}

	return nil
}

// parseResponse parses a JSON response from the Dropbox Sign API, extracting the main payload and any warnings.
//
// This utility function handles the common pattern of Dropbox Sign API responses which
// contain the main data under a specific key (e.g., "signature_request") and optional
// warnings at the top level.
func parseResponse[T any](body []byte, key string) (*T, []WarningResponse, error) {
	var rawResponse map[string]json.RawMessage
	if err := json.Unmarshal(body, &rawResponse); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract main payload by key
	payload, ok := rawResponse[key]
	if !ok {
		return nil, nil, fmt.Errorf("missing key '%s' in response", key)
	}

	// Deserialize the payload into T
	var result T
	if err := json.Unmarshal(payload, &result); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Extract warnings if present
	var warnings []WarningResponse
	if warningsData, ok := rawResponse["warnings"]; ok {
		if err := json.Unmarshal(warningsData, &warnings); err != nil {
			// Non-fatal: we can continue without warnings
			warnings = nil
		}
	}

	return &result, warnings, nil
}

// parseErrorResponse parses an error response from the Dropbox Sign API.
func (c *Client) parseErrorResponse(body []byte, statusCode int) error {
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return NewClientError(fmt.Sprintf("failed to parse error response: %s", string(body)), statusCode, err)
	}

	errResp.Error.Status = statusCode
	return errResp.Error
}
