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
	ShouldBeSent bool           `json:"shouldBeSent"`
	Subject      string         `json:"subject"`
	ToEmail      string         `json:"toEmail"`
	ToName       string         `json:"toName"`
	Placeholders map[string]any `json:"placeholders"`
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
