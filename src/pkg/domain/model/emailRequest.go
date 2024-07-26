package domain

type EmailRequest struct {
	ToEmail      string
	ToName       string
	Subject      string
	MessageBody  string
	ShouldBeSent bool
	// if attachments
	AttachmentInfo *AttachmentInfo
}
