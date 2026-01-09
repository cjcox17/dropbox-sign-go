package main

import (
	"context"
	"fmt"
	"log"
	"os"

	dropboxsign "github.com/cjcox17/dropbox-sign-go"
)

func main() {
	// Load environment variables
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable must be set")
	}

	signerEmail := os.Getenv("SIGNER_EMAIL")
	if signerEmail == "" {
		log.Fatal("SIGNER_EMAIL environment variable must be set")
	}

	templateID := os.Getenv("TEMPLATE_ID")
	if templateID == "" {
		log.Fatal("TEMPLATE_ID environment variable must be set")
	}

	fmt.Printf("Using API_KEY: %s\n", apiKey)

	// Create a new Dropbox Sign client
	client := dropboxsign.NewClient(apiKey)

	// Create a signer (role should match the template role name)
	signer := dropboxsign.NewSubSignatureRequestTemplateSigner(
		"Client",
		"test name",
		signerEmail,
	)

	// Create custom fields
	customFields := []dropboxsign.SubCustomField{
		dropboxsign.NewSubCustomField("test_field_one").
			WithValue("This is test field one!"),
		dropboxsign.NewSubCustomField("test_field_two").
			WithValue("This is test field two!"),
	}

	// Create the signature request
	signatureRequest := dropboxsign.NewSendSignatureRequest(
		[]dropboxsign.SubSignatureRequestTemplateSigner{signer},
		[]string{templateID},
	).WithCustomFields(customFields)

	// Send the signature request with template
	ctx := context.Background()
	response, warnings, err := client.SendWithTemplate(ctx, signatureRequest)
	if err != nil {
		log.Fatalf("Failed to send signature request: %v", err)
	}

	fmt.Printf("Dropbox response: %+v\n", response)
	if len(warnings) > 0 {
		fmt.Printf("Dropbox warnings: %+v\n", warnings)
	}

	// Get the signature request
	signatureRequestID := response.SignatureRequestID
	getResponse, getWarnings, err := client.GetSignatureRequest(ctx, signatureRequestID)
	if err != nil {
		log.Fatalf("Failed to get signature request: %v", err)
	}

	fmt.Printf("Dropbox get_response: %+v\n", getResponse)
	if len(getWarnings) > 0 {
		fmt.Printf("Dropbox get_warnings: %+v\n", getWarnings)
	}
}
