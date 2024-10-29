package domain

import "fmt"

// define custom errors
var (
	ErrSimpleExample = fmt.Errorf("just a simple example error without details")
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

type ErrorSendingEmailFailed struct {
	Reason string
}

func (e ErrorSendingEmailFailed) Error() string {
	return fmt.Sprintf("sending email failed: %s", e.Reason)
}

type ErrorRenderingMJMLFailed struct {
	Reason string
}

func (e ErrorRenderingMJMLFailed) Error() string {
	return fmt.Sprintf("rendering mjml failed: %s", e.Reason)
}

type ErrorAddingTemplateFailed struct {
	Reason string
}

func (e ErrorAddingTemplateFailed) Error() string {
	return fmt.Sprintf("adding template failed: %s", e.Reason)
}

type ErrorTemplateAlreadyExists struct {
	TemplateName string
}

func (e ErrorTemplateAlreadyExists) Error() string {
	return fmt.Sprintf("template already exists: %s", e.TemplateName)
}
