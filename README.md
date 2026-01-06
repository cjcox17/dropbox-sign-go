# Dropbox Sign Go Client

A modern Go client library for the Dropbox Sign API (formerly HelloSign).

## Features

- Type-safe API interactions with comprehensive error handling
- Context-aware HTTP requests with timeout support
- Support for signature requests with templates
- Proper handling of API warnings and errors
- Builder patterns for complex request construction
- Idiomatic Go design with interfaces and error handling

## Installation

```bash
go get github.com/cjcox17/dropboxsign-go-client
```

## Usage

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/cjcox17/dropboxsign-go-client"
)

func main() {
    // Create a new client with your API key
    client := dropboxsign.NewClient("your-api-key")

    // Create a signer (role should match the template role name)
    signer := dropboxsign.NewSubSignatureRequestTemplateSigner(
        "Signer",
        "John Doe",
        "john@example.com",
    )

    // Create a signature request
    request := dropboxsign.NewSendSignatureRequest(
        []dropboxsign.SubSignatureRequestTemplateSigner{signer},
        []string{"template-id"},
    ).WithTitle("Contract Signature").WithTestMode(true)

    // Send the signature request
    ctx := context.Background()
    response, warnings, err := client.SendWithTemplate(ctx, request)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Signature request sent: %s\n", response.SignatureRequestID)

    // Check for warnings
    if len(warnings) > 0 {
        for _, warning := range warnings {
            fmt.Printf("Warning: %s\n", warning)
        }
    }
}
```

### Advanced Configuration

```go
// Create a client with custom timeout
client := dropboxsign.NewClient("your-api-key").
    WithTimeout(60 * time.Second)

// Or use a custom HTTP client
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: customTransport,
}
client := dropboxsign.NewClient("your-api-key").
    WithHTTPClient(httpClient)
```

### Working with Custom Fields

```go
// Create custom fields with default values
customFields := []dropboxsign.SubCustomField{
    dropboxsign.NewSubCustomField("company_name").
        WithValue("Acme Corp").
        WithRequired(true),
    dropboxsign.NewSubCustomField("employee_id").
        WithValue("12345"),
}

request := dropboxsign.NewSendSignatureRequest(signers, templateIDs).
    WithCustomFields(customFields)
```

### Adding Signers with Authentication

```go
// Add a signer with PIN protection
signer := dropboxsign.NewSubSignatureRequestTemplateSigner(
    "Signer",
    "Jane Doe",
    "jane@example.com",
).WithPin("1234")

// Add a signer with SMS authentication
signerWithSMS := dropboxsign.NewSubSignatureRequestTemplateSigner(
    "Signer",
    "Bob Smith",
    "bob@example.com",
).WithSMSPhoneNumber("+1234567890").
  WithSMSPhoneNumberType(dropboxsign.SMSPhoneNumberTypeAuthentication)
```

### Retrieving Signature Requests

```go
ctx := context.Background()
sigRequest, warnings, err := client.GetSignatureRequest(ctx, "signature_request_id")
if err != nil {
    if dropboxsign.IsNotFound(err) {
        log.Println("Signature request not found")
    } else {
        log.Fatal(err)
    }
}

fmt.Printf("Status: Complete=%v, Declined=%v\n", 
    sigRequest.IsComplete, sigRequest.IsDeclined)
```

### Canceling Signature Requests

```go
ctx := context.Background()
err := client.CancelIncompleteSignatureRequest(ctx, "signature_request_id")
if err != nil {
    log.Fatal(err)
}
```

### Error Handling

The library provides helper functions for common error scenarios:

```go
sigRequest, _, err := client.GetSignatureRequest(ctx, signatureRequestID)
if err != nil {
    if dropboxsign.IsNotFound(err) {
        // Handle 404 - signature request not found
    } else if dropboxsign.IsUnauthorized(err) {
        // Handle 401 - invalid API key
    } else if dropboxsign.IsBadRequest(err) {
        // Handle 400 - invalid request parameters
    } else {
        // Handle other errors
    }
}
```

## Environment Variables

For the example application, set these environment variables:

```bash
export API_KEY="your-dropbox-sign-api-key"
export SIGNER_EMAIL="signer@example.com"
export TEMPLATE_ID="your-template-id"
```

## Running Examples

```bash
cd examples
go run basic_usage.go
```

## Error Types

### ErrorResponseError

API errors returned by Dropbox Sign:

```go
type ErrorResponseError struct {
    Status    int
    ErrorMsg  string
    ErrorPath *string
    ErrorName string
}
```

### ClientError

Client-side errors (network issues, parsing errors, etc.):

```go
type ClientError struct {
    Message    string
    StatusCode int
    Err        error
}
```

## License

This project is dual-licensed under MIT OR Apache-2.0, matching the original Rust implementation.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Related Projects

- [Original Rust Client](https://github.com/cjcox17/dropboxsign-client)

## Acknowledgments

This Go client is a port of the Rust Dropbox Sign client library, maintaining API compatibility while following Go idioms and best practices.
