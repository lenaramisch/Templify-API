package domain

type SMSTemplateFillRequest struct {
	ShouldBeSent        bool              `json:"shouldBeSent"`
	ReceiverPhoneNumber string            `json:"receiverPhoneNumber"`
	Placeholders        map[string]string `json:"placeholders"`
}
