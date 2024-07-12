package handler

type SmsRequest struct {
	ToNumber    string `json:"receiverPhoneNumber"`
	MessageBody string `json:"message"`
}

type EmailRequest struct {
	ToEmail     string `json:"toEmail"`
	ToName      string `json:"toName"`
	Subject     string `json:"subject"`
	MessageBody string `json:"message"`
}

type Template struct {
	Name       string `json:"name"`
	MJMLString string `json:"mjml_string"`
}

type TemplateFillRequest struct {
	ShouldBeSent bool              `json:"shouldBeSent"`
	Subject      string            `json:"subject"`
	ToEmail      string            `json:"toEmail"`
	ToName       string            `json:"toName"`
	Placeholders map[string]string `json:"placeholders"`
}

type TemplateFillRequestAttm struct {
	ShouldBeSent bool
	Subject      string
	ToEmail      string
	ToName       string
	Placeholders map[string]string
	AttmContent  string
	FileName     string
	FileType     string
}

type EmailAttachmentRequest struct {
	ToEmail           string
	ToName            string
	Subject           string
	MessageBody       string
	AttachmentContent string
	FileName          string
	FileType          string
}
