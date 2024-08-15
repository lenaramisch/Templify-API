package domain

type EmailTemplateSendRequest struct {
	ToEmail        string
	ToName         string
	Subject        string
	TemplateName   string
	Placeholders   map[string]string
	AttachmentInfo *AttachmentInfo
}
