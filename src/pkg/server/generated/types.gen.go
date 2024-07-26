// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.2.0 DO NOT EDIT.
package server

// Defines values for ErrorType.
const (
	BadRequest          ErrorType = "Bad Request"
	InternalServerError ErrorType = "Internal Server Error"
	NotFound            ErrorType = "Not Found"
	NotImplemented      ErrorType = "Not Implemented"
)

// Defines values for Status.
const (
	HEALTHY   Status = "HEALTHY"
	UNHEALTHY Status = "UNHEALTHY"
	UNKNOWN   Status = "UNKNOWN"
)

// EmailSendRequest defines model for EmailSendRequest.
type EmailSendRequest struct {
	Message string `json:"message"`
	Subject string `json:"subject"`
	ToEmail string `json:"toEmail"`
	ToName  string `json:"toName"`
}

// EmailSendRequestAttachment defines model for EmailSendRequestAttachment.
type EmailSendRequestAttachment struct {
	// File The file that should be send as attachment
	File    string `json:"file"`
	Message string `json:"message"`
	Subject string `json:"subject"`
	ToEmail string `json:"toEmail"`
	ToName  string `json:"toName"`
}

// Error This object holds the error response data.
type Error struct {
	// ErrorType The error type
	ErrorType ErrorType `json:"ErrorType"`

	// Code The error code
	Code int `json:"code"`

	// Error The error message
	Error string `json:"error"`

	// ErrorId The unique identifier for the error
	ErrorId string `json:"errorId"`

	// Timestamp The time the error occurred
	Timestamp string `json:"timestamp"`
}

// ErrorType The error type
type ErrorType string

// FilledTemplateResponse defines model for FilledTemplateResponse.
type FilledTemplateResponse = string

// MJMLSendRequestAttachment defines model for MJMLSendRequestAttachment.
type MJMLSendRequestAttachment struct {
	// PlaceHolder Each placeholder key-value pair should be in its own form field
	PlaceHolder string `json:"PlaceHolder"`

	// File The file that should be send as attachment
	File string `json:"file"`

	// ShouldBeSent Determines if the email will be sent (true/false)
	ShouldBeSent string `json:"shouldBeSent"`
	Subject      string `json:"subject"`
	ToEmail      string `json:"toEmail"`
	ToName       string `json:"toName"`
}

// Placeholders defines model for Placeholders.
type Placeholders struct {
	Placeholders []struct {
		Data *[]struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"data,omitempty"`
	} `json:"placeholders"`
}

// SMSSendRequest defines model for SMSSendRequest.
type SMSSendRequest struct {
	Message             string `json:"message"`
	ReceiverPhoneNumber string `json:"receiverPhoneNumber"`
}

// Status The status of the API
type Status string

// Template defines model for Template.
type Template struct {
	Name             string `json:"name"`
	TemplateMJMLCode string `json:"templateMJMLCode"`
}

// TemplateFillRequest defines model for TemplateFillRequest.
type TemplateFillRequest struct {
	Placeholders struct {
		Age       string `json:"Age"`
		FirstName string `json:"FirstName"`
		LastName  string `json:"LastName"`
	} `json:"placeholders"`
	ShouldBeSent bool   `json:"shouldBeSent"`
	Subject      string `json:"subject"`
	ToEmail      string `json:"toEmail"`
	ToName       string `json:"toName"`
}

// TemplatePostRequest defines model for TemplatePostRequest.
type TemplatePostRequest struct {
	TemplateMJMLCode string `json:"templateMJMLCode"`
}

// Version This object holds the API version data.
type Version struct {
	// BuildDate The date the code was built
	BuildDate string `json:"buildDate"`

	// CommitDate The date of the commit
	CommitDate string `json:"commitDate"`

	// CommitHash The hash of the commit
	CommitHash string `json:"commitHash"`

	// Details A description of the API
	Details string `json:"details"`

	// Version The version of the API
	Version string `json:"version"`
}

// SendEmailJSONRequestBody defines body for SendEmail for application/json ContentType.
type SendEmailJSONRequestBody = EmailSendRequest

// SendEmailWithAttachmentMultipartRequestBody defines body for SendEmailWithAttachment for multipart/form-data ContentType.
type SendEmailWithAttachmentMultipartRequestBody = EmailSendRequestAttachment

// AddNewTemplateJSONRequestBody defines body for AddNewTemplate for application/json ContentType.
type AddNewTemplateJSONRequestBody = TemplatePostRequest

// FillTemplateJSONRequestBody defines body for FillTemplate for application/json ContentType.
type FillTemplateJSONRequestBody = TemplateFillRequest

// SendMJMLEmailWithAttachmentMultipartRequestBody defines body for SendMJMLEmailWithAttachment for multipart/form-data ContentType.
type SendMJMLEmailWithAttachmentMultipartRequestBody = MJMLSendRequestAttachment

// SendSMSJSONRequestBody defines body for SendSMS for application/json ContentType.
type SendSMSJSONRequestBody = SMSSendRequest
