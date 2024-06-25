package domain

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
