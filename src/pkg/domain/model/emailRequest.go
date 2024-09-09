package domain

type EmailRequest struct {
	ToEmail        string
	ToName         string
	Subject        string
	MessageBody    string
	AttachmentInfo []AttachmentInfo
}
