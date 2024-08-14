package domain

type SMSTemplateFillRequest struct {
	ReceiverPhoneNumber string            `json:"receiverPhoneNumber"`
	Placeholders        map[string]string `json:"placeholders"`
}
