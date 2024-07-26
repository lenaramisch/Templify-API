package apihandler

import (
	"net/http"
	server "templify/pkg/server/generated"
	"templify/pkg/server/handler"

	"github.com/go-chi/render"
)

// Send a SMS with custom text
// (POST /sms)
func (ah *APIHandler) SendSMS(w http.ResponseWriter, r *http.Request) {
	// TODO create DTO in api spec
	var smsRequest server.SMSSendRequest

	if err := handler.ReadRequestBody(w, r, &smsRequest); err != nil {
		return
	}

	err := ah.Usecase.SendSMS(smsRequest.ReceiverPhoneNumber, smsRequest.Message)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "SMS sent successfully")
}
