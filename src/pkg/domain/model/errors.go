package domain

import "fmt"

// define custom errors
var (
	ErrSimpleExample              = fmt.Errorf("just a simple example error without details")
	ErrAuthorizationHeaderMissing = fmt.Errorf("authorization header missing")
	ErrInvalidTokenFormat         = fmt.Errorf("invalid token format, expected Bearer")
	ErrInvalidToken               = fmt.Errorf("invalid token")
	ErrAccessDenied               = fmt.Errorf("access denied")
)

type ErrorPlaceholderMissing struct {
	MissingPlaceholder string
}

func (e ErrorPlaceholderMissing) Error() string {
	return fmt.Sprintf("placeholder is missing: %s", e.MissingPlaceholder)
}

type ErrorTemplateNotFound struct {
	TemplateName string
}

func (e ErrorTemplateNotFound) Error() string {
	return fmt.Sprintf("template not found: %s", e.TemplateName)
}

type ErrorSendingEmail struct {
	Reason string
}

func (e ErrorSendingEmail) Error() string {
	return fmt.Sprintf("sending email failed: %s", e.Reason)
}

type ErrorRenderingMJML struct {
	Reason string
}

func (e ErrorRenderingMJML) Error() string {
	return fmt.Sprintf("rendering mjml failed: %s", e.Reason)
}

type ErrorAddingTemplate struct {
	Reason string
}

func (e ErrorAddingTemplate) Error() string {
	return fmt.Sprintf("adding template failed: %s", e.Reason)
}

type ErrorTemplateAlreadyExists struct {
	TemplateName string
}

func (e ErrorTemplateAlreadyExists) Error() string {
	return fmt.Sprintf("template already exists: %s", e.TemplateName)
}

type ErrorGettingDownloadURL struct {
	Reason string
}

func (e ErrorGettingDownloadURL) Error() string {
	return fmt.Sprintf("getting download URL failed: %s", e.Reason)
}

type ErrorGettingUploadURL struct {
	Reason string
}

func (e ErrorGettingUploadURL) Error() string {
	return fmt.Sprintf("getting upload URL failed: %s", e.Reason)
}

type ErrorFillingTemplate struct {
	Reason string
}

func (e ErrorFillingTemplate) Error() string {
	return fmt.Sprintf("filling template failed: %s", e.Reason)
}

type ErrorRenderingTypst struct {
	Reason string
}

func (e ErrorRenderingTypst) Error() string {
	return fmt.Sprintf("rendering typst failed: %s", e.Reason)
}

type ErrorCreatingSMSRequest struct {
	Reason string
}

func (e ErrorCreatingSMSRequest) Error() string {
	return fmt.Sprintf("creating sms request failed: %s", e.Reason)
}

type ErrorPerformingSMSRequest struct {
	Reason string
}

func (e ErrorPerformingSMSRequest) Error() string {
	return fmt.Sprintf("performing sms request failed: %s", e.Reason)
}

type ErrorSendingSMS struct {
	StatusCode int
}

func (e ErrorSendingSMS) Error() string {
	return fmt.Sprintf("sending sms failed with status code: %d", e.StatusCode)
}

type ErrorWorkflowAlreadyExists struct {
	WorkflowName string
}

func (e ErrorWorkflowAlreadyExists) Error() string {
	return fmt.Sprintf("workflow already exists: %s", e.WorkflowName)
}

type ErrorWorkflowNotFound struct {
	WorkflowName string
}

func (e ErrorWorkflowNotFound) Error() string {
	return fmt.Sprintf("workflow not found: %s", e.WorkflowName)
}

type ErrorAttachmentNameInvalid struct {
	AttachmentName string
}

func (e ErrorAttachmentNameInvalid) Error() string {
	return fmt.Sprintf("attachment name invalid: %s, expected to contain file extension", e.AttachmentName)
}

type ErrorDownloadingFile struct {
	Reason string
}

func (e ErrorDownloadingFile) Error() string {
	return fmt.Sprintf("downloading file failed: %s", e.Reason)
}
