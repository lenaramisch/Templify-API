package domain

type SmsRequest struct {
	ToNumber    string
	MessageBody string
}

type EmailRequest struct {
	ToEmail      string
	ToName       string
	Subject      string
	MessageBody  string
	ShouldBeSent bool
	// if attachments
	AttachmentInfo *AttachmentInfo
}

type Template struct {
	Name        string
	TemplateStr string
}

type AttachmentInfo struct {
	FileName      string
	FileExtension string
	Content       []byte
}
