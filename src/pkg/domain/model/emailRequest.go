package domain

type EmailRequest struct {
	ToEmail     string
	ToName      string
	Subject     string
	MessageBody string
	// if attachments
	AttachmentInfo *AttachmentInfo
}
