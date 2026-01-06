// Package dropboxsign provides data models and types for signature request operations.
package dropboxsign

import (
	"encoding/json"
	"strings"
)

// SendSignatureRequest represents a request structure for sending signature requests with templates.
//
// This struct represents a complete signature request that will be sent to signers.
// It requires signers and template IDs, with many optional configuration parameters
// for customizing the signing experience.
//
// Example:
//
//	signer := dropboxsign.NewSubSignatureRequestTemplateSigner(
//		"Signer",
//		"John Doe",
//		"john@example.com",
//	)
//	request := dropboxsign.NewSendSignatureRequest(
//		[]dropboxsign.SubSignatureRequestTemplateSigner{signer},
//		[]string{"template-id"},
//	).WithTitle("Contract Signature").WithTestMode(true)
type SendSignatureRequest struct {
	// Signers is the list of signers who will receive the signature request
	Signers []SubSignatureRequestTemplateSigner `json:"signers"`
	// TemplateIDs is the list of template IDs to use for this signature request
	TemplateIDs []string `json:"template_ids"`
	// AllowDecline specifies whether signers can decline to sign (default: true)
	AllowDecline *bool `json:"allow_decline,omitempty"`
	// CCs is the list of CC recipients who will receive copies of the signature request
	CCs []SubCC `json:"ccs,omitempty"`
	// ClientID is the client ID for API apps
	ClientID *string `json:"client_id,omitempty"`
	// CustomFields are custom form fields to pre-populate in the document
	CustomFields []SubCustomField `json:"custom_fields,omitempty"`
	// Files is file data as byte arrays (alternative to FileURLs)
	Files [][]byte `json:"files,omitempty"`
	// FileURLs are URLs to files to be signed (alternative to Files)
	FileURLs []string `json:"file_urls,omitempty"`
	// IsEID specifies whether to enable eIDAS compliance (European electronic signatures)
	IsEID *bool `json:"is_eid,omitempty"`
	// Message is the custom message to include in the signature request email
	Message *string `json:"message,omitempty"`
	// Metadata contains key-value pairs for storing custom data with the signature request
	Metadata map[string]string `json:"metadata,omitempty"`
	// SigningOptions is the configuration for signature methods and options
	SigningOptions *SubSigningOptions `json:"signing_options,omitempty"`
	// SigningRedirectURL is the URL to redirect signers to after completing their signature
	SigningRedirectURL *string `json:"signing_redirect_url,omitempty"`
	// TestMode specifies whether to create the signature request in test mode
	TestMode *bool `json:"test_mode,omitempty"`
	// Title is the title for the signature request
	Title *string `json:"title,omitempty"`
}

// NewSendSignatureRequest creates a new signature request with the minimum required fields.
//
// Example:
//
//	signer := dropboxsign.NewSubSignatureRequestTemplateSigner("Signer", "Jane Doe", "jane@example.com")
//	request := dropboxsign.NewSendSignatureRequest(
//		[]dropboxsign.SubSignatureRequestTemplateSigner{signer},
//		[]string{"template-id"},
//	)
func NewSendSignatureRequest(signers []SubSignatureRequestTemplateSigner, templateIDs []string) *SendSignatureRequest {
	return &SendSignatureRequest{
		Signers:     signers,
		TemplateIDs: templateIDs,
	}
}

// WithAllowDecline sets whether signers can decline to sign the document.
func (s *SendSignatureRequest) WithAllowDecline(allowDecline bool) *SendSignatureRequest {
	s.AllowDecline = &allowDecline
	return s
}

// WithCCs sets the list of CC recipients for the signature request.
func (s *SendSignatureRequest) WithCCs(ccs []SubCC) *SendSignatureRequest {
	s.CCs = ccs
	return s
}

// WithClientID sets the client ID for API apps.
func (s *SendSignatureRequest) WithClientID(clientID string) *SendSignatureRequest {
	s.ClientID = &clientID
	return s
}

// WithCustomFields sets custom form fields to pre-populate in the document.
func (s *SendSignatureRequest) WithCustomFields(customFields []SubCustomField) *SendSignatureRequest {
	s.CustomFields = customFields
	return s
}

// WithFiles sets file data as byte arrays for documents to be signed.
func (s *SendSignatureRequest) WithFiles(files [][]byte) *SendSignatureRequest {
	s.Files = files
	return s
}

// WithFileURLs sets URLs to files that should be downloaded and used as documents.
func (s *SendSignatureRequest) WithFileURLs(fileURLs []string) *SendSignatureRequest {
	s.FileURLs = fileURLs
	return s
}

// WithIsEID sets whether to enable eIDAS compliance for European electronic signatures.
func (s *SendSignatureRequest) WithIsEID(isEID bool) *SendSignatureRequest {
	s.IsEID = &isEID
	return s
}

// WithMessage sets a custom message to include in signature request emails.
func (s *SendSignatureRequest) WithMessage(message string) *SendSignatureRequest {
	s.Message = &message
	return s
}

// WithMetadata sets custom metadata key-value pairs for the signature request.
func (s *SendSignatureRequest) WithMetadata(metadata map[string]string) *SendSignatureRequest {
	s.Metadata = metadata
	return s
}

// WithSigningOptions sets configuration for available signature methods.
func (s *SendSignatureRequest) WithSigningOptions(signingOptions *SubSigningOptions) *SendSignatureRequest {
	s.SigningOptions = signingOptions
	return s
}

// WithSigningRedirectURL sets the URL to redirect signers to after they complete signing.
func (s *SendSignatureRequest) WithSigningRedirectURL(signingRedirectURL string) *SendSignatureRequest {
	s.SigningRedirectURL = &signingRedirectURL
	return s
}

// WithTestMode sets whether to create the signature request in test mode.
func (s *SendSignatureRequest) WithTestMode(testMode bool) *SendSignatureRequest {
	s.TestMode = &testMode
	return s
}

// WithTitle sets the title for the signature request.
func (s *SendSignatureRequest) WithTitle(title string) *SendSignatureRequest {
	s.Title = &title
	return s
}

// SubSignatureRequestTemplateSigner represents a signer in a template-based signature request.
//
// Each signer must have a role (matching the template), name, and email address.
// Additional authentication options like PIN or SMS can be configured.
type SubSignatureRequestTemplateSigner struct {
	// Role is the role name that matches a role defined in the template
	Role string `json:"role"`
	// Name is the full name of the signer
	Name string `json:"name"`
	// EmailAddress is the email address where the signature request will be sent
	EmailAddress string `json:"email_address"`
	// Pin is an optional PIN for additional security (4-12 digits)
	Pin *string `json:"pin,omitempty"`
	// SMSPhoneNumber is the phone number for SMS authentication or delivery
	SMSPhoneNumber *string `json:"sms_phone_number,omitempty"`
	// SMSPhoneNumberType is the type of SMS usage (authentication or delivery)
	SMSPhoneNumberType *SMSPhoneNumberType `json:"sms_phone_number_type,omitempty"`
}

// NewSubSignatureRequestTemplateSigner creates a new signer with the minimum required information.
//
// Example:
//
//	signer := dropboxsign.NewSubSignatureRequestTemplateSigner(
//		"Signer",
//		"Jane Doe",
//		"jane@example.com",
//	)
func NewSubSignatureRequestTemplateSigner(role, name, emailAddress string) SubSignatureRequestTemplateSigner {
	return SubSignatureRequestTemplateSigner{
		Role:         role,
		Name:         name,
		EmailAddress: emailAddress,
	}
}

// WithPin sets a PIN that the signer must enter before signing.
func (s SubSignatureRequestTemplateSigner) WithPin(pin string) SubSignatureRequestTemplateSigner {
	s.Pin = &pin
	return s
}

// WithSMSPhoneNumber sets the phone number for SMS authentication or delivery.
func (s SubSignatureRequestTemplateSigner) WithSMSPhoneNumber(smsPhoneNumber string) SubSignatureRequestTemplateSigner {
	s.SMSPhoneNumber = &smsPhoneNumber
	return s
}

// WithSMSPhoneNumberType sets how the SMS phone number should be used.
func (s SubSignatureRequestTemplateSigner) WithSMSPhoneNumberType(smsPhoneNumberType SMSPhoneNumberType) SubSignatureRequestTemplateSigner {
	s.SMSPhoneNumberType = &smsPhoneNumberType
	return s
}

// SMSPhoneNumberType specifies how SMS phone numbers are used in signature requests.
type SMSPhoneNumberType string

const (
	// SMSPhoneNumberTypeAuthentication means SMS is used for two-factor authentication
	SMSPhoneNumberTypeAuthentication SMSPhoneNumberType = "authentication"
	// SMSPhoneNumberTypeDelivery means SMS is used for document delivery notifications
	SMSPhoneNumberTypeDelivery SMSPhoneNumberType = "delivery"
)

// SubCC represents a carbon copy recipient for signature requests.
//
// CC recipients receive copies of signature request emails and completion notifications
// but are not required to sign the document.
type SubCC struct {
	// Role is the role name for the CC recipient (must match template if using templates)
	Role string `json:"role"`
	// Email is the email address of the CC recipient
	Email string `json:"email"`
}

// NewSubCC creates a new CC recipient.
func NewSubCC(role, email string) SubCC {
	return SubCC{
		Role:  role,
		Email: email,
	}
}

// SubCustomField represents a custom form field that can be pre-populated in signature requests.
//
// Custom fields allow you to set default values for form fields in the document
// before sending it to signers.
type SubCustomField struct {
	// Name is the name of the custom field (must match field name in template)
	Name string `json:"name"`
	// Editor is the email address of the person who can edit this field
	Editor *string `json:"editor,omitempty"`
	// Required specifies whether this field is required to be filled out
	Required *bool `json:"required,omitempty"`
	// Value is the default value for the field
	Value *string `json:"value,omitempty"`
}

// NewSubCustomField creates a new custom field with the specified name.
func NewSubCustomField(name string) SubCustomField {
	return SubCustomField{
		Name: name,
	}
}

// WithEditor sets the email address of the person who can edit this field.
func (s SubCustomField) WithEditor(editor string) SubCustomField {
	s.Editor = &editor
	return s
}

// WithRequired sets whether this field is required to be filled out.
func (s SubCustomField) WithRequired(required bool) SubCustomField {
	s.Required = &required
	return s
}

// WithValue sets the default value for this field.
func (s SubCustomField) WithValue(value string) SubCustomField {
	s.Value = &value
	return s
}

// SubSigningOptions represents configuration for available signature methods.
//
// Defines which signature methods are available to signers and which one
// is the default option.
type SubSigningOptions struct {
	// DefaultType is the default signature method that will be pre-selected
	DefaultType SubSigningOptionsDefaultType `json:"default_type"`
	// Draw specifies whether signers can draw their signature
	Draw *bool `json:"draw,omitempty"`
	// Phone specifies whether signers can use phone-based signatures
	Phone *bool `json:"phone,omitempty"`
	// Type specifies whether signers can type their signature
	Type *bool `json:"type,omitempty"`
	// Upload specifies whether signers can upload an image of their signature
	Upload *bool `json:"upload,omitempty"`
}

// NewSubSigningOptions creates new signing options with the specified default signature method.
func NewSubSigningOptions(defaultType SubSigningOptionsDefaultType) *SubSigningOptions {
	return &SubSigningOptions{
		DefaultType: defaultType,
	}
}

// WithDraw sets whether signers can draw their signature with mouse or finger.
func (s *SubSigningOptions) WithDraw(draw bool) *SubSigningOptions {
	s.Draw = &draw
	return s
}

// WithPhone sets whether signers can use phone-based signature verification.
func (s *SubSigningOptions) WithPhone(phone bool) *SubSigningOptions {
	s.Phone = &phone
	return s
}

// WithType sets whether signers can type their signature using a font.
func (s *SubSigningOptions) WithType(typeEnabled bool) *SubSigningOptions {
	s.Type = &typeEnabled
	return s
}

// WithUpload sets whether signers can upload an image of their signature.
func (s *SubSigningOptions) WithUpload(upload bool) *SubSigningOptions {
	s.Upload = &upload
	return s
}

// SubSigningOptionsDefaultType represents available signature methods for the default signing option.
type SubSigningOptionsDefaultType string

const (
	// SubSigningOptionsDefaultTypeDraw means draw signature with mouse/finger
	SubSigningOptionsDefaultTypeDraw SubSigningOptionsDefaultType = "draw"
	// SubSigningOptionsDefaultTypePhone means phone-based signature verification
	SubSigningOptionsDefaultTypePhone SubSigningOptionsDefaultType = "phone"
	// SubSigningOptionsDefaultTypeType means type signature using a font
	SubSigningOptionsDefaultTypeType SubSigningOptionsDefaultType = "type"
	// SubSigningOptionsDefaultTypeUpload means upload an image of the signature
	SubSigningOptionsDefaultTypeUpload SubSigningOptionsDefaultType = "upload"
)

// SignatureRequestResponse contains complete response data for a signature request.
//
// Contains all information about a signature request including its status,
// signer information, URLs, and metadata.
type SignatureRequestResponse struct {
	// TestMode indicates whether this signature request was created in test mode
	TestMode *bool `json:"test_mode,omitempty"`
	// SignatureRequestID is the unique identifier for this signature request
	SignatureRequestID string `json:"signature_request_id"`
	// RequesterEmailAddress is the email address of the person who created this request
	RequesterEmailAddress *string `json:"requester_email_address,omitempty"`
	// Title is the current title of the signature request
	Title string `json:"title"`
	// OriginalTitle is the original title of the signature request (before any modifications)
	OriginalTitle string `json:"original_title"`
	// Subject is the subject line used in signature request emails
	Subject *string `json:"subject,omitempty"`
	// Message is the custom message included in signature request emails
	Message *string `json:"message,omitempty"`
	// Metadata contains custom metadata key-value pairs
	Metadata map[string]string `json:"metadata"`
	// CreatedAt is the Unix timestamp when the signature request was created
	CreatedAt int64 `json:"created_at"`
	// ExpiresAt is the Unix timestamp when the signature request expires (if set)
	ExpiresAt *int64 `json:"expires_at,omitempty"`
	// IsComplete indicates whether all required signatures have been completed
	IsComplete bool `json:"is_complete"`
	// IsDeclined indicates whether any signer has declined to sign
	IsDeclined bool `json:"is_declined"`
	// HasError indicates whether there are any errors with this signature request
	HasError bool `json:"has_error"`
	// FilesURL is the URL to download the signed documents
	FilesURL string `json:"files_url"`
	// SigningURL is the URL for signers to access the signing interface
	SigningURL *string `json:"signing_url,omitempty"`
	// DetailsURL is the URL to view signature request details
	DetailsURL string `json:"details_url"`
	// CCEmailAddresses are email addresses that received CC copies of the request
	CCEmailAddresses []string `json:"cc_email_addresses"`
	// SigningRedirectURL is the URL to redirect signers after they complete signing
	SigningRedirectURL *string `json:"signing_redirect_url,omitempty"`
	// FinalCopyURI is the URI to the final signed document copy
	FinalCopyURI *string `json:"final_copy_uri,omitempty"`
	// TemplateIDs are template IDs used to create this signature request
	TemplateIDs []string `json:"template_ids,omitempty"`
	// CustomIDs are custom IDs associated with this signature request
	CustomIDs []string `json:"custom_ids,omitempty"`
	// Attachments are file attachments associated with this signature request
	Attachments *SignatureRequestResponseAttachment `json:"attachments,omitempty"`
	// ResponseData contains form field response data from signers
	ResponseData []SignatureRequestResponseData `json:"response_data,omitempty"`
	// Signatures contains individual signature status for each signer
	Signatures []SignatureRequestResponseSignatures `json:"signatures"`
	// BulkSendJobID is the bulk send job ID if this was part of a bulk operation
	BulkSendJobID *string `json:"bulk_send_job_id,omitempty"`
}

// SignatureRequestResponseCustomFieldBase represents base structure for custom form fields in signature request responses.
//
// Represents form fields that were filled out by signers or pre-populated
// when the signature request was created.
type SignatureRequestResponseCustomFieldBase struct {
	// Type is the type of the form field (text, checkbox, etc.)
	Type SignatureRequestResponseCustomFieldBaseType `json:"type"`
	// Name is the name/identifier of the form field
	Name string `json:"name"`
	// Required specifies whether this field was required to be filled out
	Required *bool `json:"required,omitempty"`
	// APIID is the API identifier for this field
	APIID *string `json:"api_id,omitempty"`
	// Editor is the email of the person who can edit this field
	Editor *string `json:"editor,omitempty"`
	// Value is the current value of the form field
	Value *string `json:"value,omitempty"`
}

// SignatureRequestResponseCustomFieldBaseType represents types of custom form fields available in signature requests.
type SignatureRequestResponseCustomFieldBaseType string

const (
	// SignatureRequestResponseCustomFieldBaseTypeText is a single-line or multi-line text input field
	SignatureRequestResponseCustomFieldBaseTypeText SignatureRequestResponseCustomFieldBaseType = "text"
	// SignatureRequestResponseCustomFieldBaseTypeCheckbox is a checkbox that can be checked or unchecked
	SignatureRequestResponseCustomFieldBaseTypeCheckbox SignatureRequestResponseCustomFieldBaseType = "checkbox"
)

// SignatureRequestResponseAttachment represents a file attachment associated with a signature request.
//
// Represents additional documents that signers can upload as part of
// the signing process.
type SignatureRequestResponseAttachment struct {
	// ID is the unique identifier for this attachment
	ID string `json:"id"`
	// Signer is the email address of the signer this attachment is assigned to
	Signer string `json:"signer"`
	// Name is the display name of the attachment
	Name string `json:"name"`
	// Required specifies whether uploading this attachment is required
	Required bool `json:"required"`
	// Instructions contains instructions for the signer about this attachment
	Instructions *string `json:"instructions,omitempty"`
	// UploadedAt is the Unix timestamp when the attachment was uploaded
	UploadedAt *int64 `json:"uploaded_at,omitempty"`
}

// SignatureRequestResponseData represents form field response data from a signature request.
//
// Contains the values that signers entered in form fields, along with
// metadata about each field.
type SignatureRequestResponseData struct {
	// APIID is the API identifier for this form field
	APIID *string `json:"api_id,omitempty"`
	// SignatureID is the ID of the signature this data belongs to
	SignatureID *string `json:"signature_id,omitempty"`
	// Name is the name/label of the form field
	Name *string `json:"name,omitempty"`
	// Required specifies whether this field was required
	Required *bool `json:"required,omitempty"`
	// Type is the type of form field (text, checkbox, dropdown, etc.)
	Type *SignatureRequestResponseDataType `json:"type,omitempty"`
	// Value is the value entered by the signer
	Value *string `json:"value,omitempty"`
}

// SignatureRequestResponseSignatures represents individual signature status and metadata for each signer.
//
// Contains detailed information about each signer's interaction with
// the signature request, including status, timestamps, and authentication details.
type SignatureRequestResponseSignatures struct {
	// SignatureID is the unique identifier for this signature
	SignatureID string `json:"signature_id"`
	// SignerGroupGUID is the group GUID if this signer belongs to a signer group
	SignerGroupGUID *string `json:"signer_group_guid,omitempty"`
	// SignerEmailAddress is the email address of the signer
	SignerEmailAddress string `json:"signer_email_address"`
	// SignerName is the full name of the signer
	SignerName *string `json:"signer_name,omitempty"`
	// SignerRole is the role of the signer in the signature request
	SignerRole *string `json:"signer_role,omitempty"`
	// Order is the signing order (for sequential signing workflows)
	Order *int `json:"order,omitempty"`
	// StatusCode is the current status of this signature (awaiting_signature, signed, declined, etc.)
	StatusCode string `json:"status_code"`
	// DeclineReason is the reason provided if the signer declined to sign
	DeclineReason *string `json:"decline_reason,omitempty"`
	// SignedAt is the Unix timestamp when the signature was completed
	SignedAt *int64 `json:"signed_at,omitempty"`
	// LastViewedAt is the Unix timestamp when the signer last viewed the document
	LastViewedAt *int64 `json:"last_viewed_at,omitempty"`
	// LastRemindedAt is the Unix timestamp when the signer was last sent a reminder
	LastRemindedAt *int64 `json:"last_reminded_at,omitempty"`
	// HasPin indicates whether this signer is required to enter a PIN
	HasPin bool `json:"has_pin"`
	// HasSMSAuth indicates whether SMS authentication is enabled for this signer
	HasSMSAuth *bool `json:"has_sms_auth,omitempty"`
	// HasSMSDelivery indicates whether SMS delivery notifications are enabled for this signer
	HasSMSDelivery *bool `json:"has_sms_delivery,omitempty"`
	// SMSPhoneNumber is the phone number used for SMS authentication or delivery
	SMSPhoneNumber *string `json:"sms_phone_number,omitempty"`
	// ReassignedBy is the email of the person who reassigned this signature
	ReassignedBy *string `json:"reassigned_by,omitempty"`
	// ReassignmentReason is the reason provided for reassigning this signature
	ReassignmentReason *string `json:"reassignment_reason,omitempty"`
	// ReassignedFrom is the email of the original signer before reassignment
	ReassignedFrom *string `json:"reassigned_from,omitempty"`
	// Error is the error message if there was a problem with this signature
	Error *string `json:"error,omitempty"`
}

// SignerStatus represents the status of a signer in a signature request.
type SignerStatus string

const (
	// SignerStatusSuccess indicates the signature was successful
	SignerStatusSuccess SignerStatus = "success"
	// SignerStatusOnHold indicates the signature is on hold
	SignerStatusOnHold SignerStatus = "on_hold"
	// SignerStatusSigned indicates the document has been signed
	SignerStatusSigned SignerStatus = "signed"
	// SignerStatusAwaitingSignature indicates waiting for signature
	SignerStatusAwaitingSignature SignerStatus = "awaiting_signature"
	// SignerStatusDeclined indicates the signer declined to sign
	SignerStatusDeclined SignerStatus = "declined"
	// SignerStatusErrorUnknown indicates an unknown error occurred
	SignerStatusErrorUnknown SignerStatus = "error_unknown"
	// SignerStatusErrorFile indicates a file error occurred
	SignerStatusErrorFile SignerStatus = "error_file"
	// SignerStatusErrorComponentPosition indicates a component position error occurred
	SignerStatusErrorComponentPosition SignerStatus = "error_component_position"
	// SignerStatusErrorTextTag indicates a text tag error occurred
	SignerStatusErrorTextTag SignerStatus = "error_text_tag"
	// SignerStatusOnHoldByRequester indicates the signature is on hold by the requester
	SignerStatusOnHoldByRequester SignerStatus = "on_hold_by_requester"
	// SignerStatusErrorInvalidEmail indicates an invalid email error occurred
	SignerStatusErrorInvalidEmail SignerStatus = "error_invalid_email"
	// SignerStatusExpired indicates the signature request has expired
	SignerStatusExpired SignerStatus = "expired"
	// SignerStatusUnknownEnum indicates an unknown status value
	SignerStatusUnknownEnum SignerStatus = "unknown_enum"
)

// UnmarshalJSON implements custom unmarshaling for SignerStatus.
func (s *SignerStatus) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*s = ParseSignerStatus(str)
	return nil
}

// ParseSignerStatus parses a string into a SignerStatus.
func ParseSignerStatus(s string) SignerStatus {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "success":
		return SignerStatusSuccess
	case "on_hold":
		return SignerStatusOnHold
	case "signed":
		return SignerStatusSigned
	case "awaiting_signature":
		return SignerStatusAwaitingSignature
	case "declined":
		return SignerStatusDeclined
	case "error_unknown":
		return SignerStatusErrorUnknown
	case "error_file":
		return SignerStatusErrorFile
	case "error_component_position":
		return SignerStatusErrorComponentPosition
	case "error_text_tag":
		return SignerStatusErrorTextTag
	case "on_hold_by_requester":
		return SignerStatusOnHoldByRequester
	case "error_invalid_email":
		return SignerStatusErrorInvalidEmail
	case "expired":
		return SignerStatusExpired
	default:
		return SignerStatusUnknownEnum
	}
}

// SignatureRequestResponseDataType represents types of form fields that can appear in signature request responses.
//
// Covers all possible field types that signers can interact with in documents.
type SignatureRequestResponseDataType string

const (
	// SignatureRequestResponseDataTypeText is a single-line or multi-line text input
	SignatureRequestResponseDataTypeText SignatureRequestResponseDataType = "text"
	// SignatureRequestResponseDataTypeCheckbox is a checkbox that can be checked or unchecked
	SignatureRequestResponseDataTypeCheckbox SignatureRequestResponseDataType = "checkbox"
	// SignatureRequestResponseDataTypeDropdown is a dropdown menu with predefined options
	SignatureRequestResponseDataTypeDropdown SignatureRequestResponseDataType = "dropdown"
	// SignatureRequestResponseDataTypeRadio is a radio button group (single selection)
	SignatureRequestResponseDataTypeRadio SignatureRequestResponseDataType = "radio"
	// SignatureRequestResponseDataTypeSignature is an electronic signature field
	SignatureRequestResponseDataTypeSignature SignatureRequestResponseDataType = "signature"
	// SignatureRequestResponseDataTypeDateSigned is automatically filled date when document was signed
	SignatureRequestResponseDataTypeDateSigned SignatureRequestResponseDataType = "date_signed"
	// SignatureRequestResponseDataTypeInitials is an initial signature field
	SignatureRequestResponseDataTypeInitials SignatureRequestResponseDataType = "initials"
	// SignatureRequestResponseDataTypeTextMerge is a text field merged from template data
	SignatureRequestResponseDataTypeTextMerge SignatureRequestResponseDataType = "text-merge"
	// SignatureRequestResponseDataTypeCheckboxMerge is a checkbox field merged from template data
	SignatureRequestResponseDataTypeCheckboxMerge SignatureRequestResponseDataType = "checkbox-merge"
)
