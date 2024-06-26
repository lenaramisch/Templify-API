package handler

import (
	"encoding/json"
	"net/http"

	"example.SMSService.com/pkg/domain"
	"github.com/go-chi/render"
)

type APIHandler struct {
	usecase *domain.Usecase
}

func NewAPIHandler(usecase *domain.Usecase) *APIHandler {
	return &APIHandler{
		usecase: usecase,
	}
}

func (ah *APIHandler) SMSPostRequest(res http.ResponseWriter, req *http.Request) {
	var smsRequest SmsRequest

	err := json.NewDecoder(req.Body).Decode(&smsRequest)
	if err != nil {
		render.Status(req, http.StatusBadRequest)
		render.PlainText(res, req, "You messed up: "+err.Error())
		return
	}

	if smsRequest.ToNumber == "" || smsRequest.MessageBody == "" {
		http.Error(res, "Empty string content in either ToNumber or MessageBody", http.StatusBadRequest)
		return
	}
	err = ah.usecase.SendSMS(smsRequest.ToNumber, smsRequest.MessageBody)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	render.Status(req, http.StatusOK)
	render.PlainText(res, req, "SMS sent successfully")
}

func (ah *APIHandler) EmailPostRequest(res http.ResponseWriter, req *http.Request) {
	var emailRequest EmailRequest

	err := json.NewDecoder(req.Body).Decode(&emailRequest)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if emailRequest.ToEmail == "" || emailRequest.ToName == "" || emailRequest.Subject == "" || emailRequest.MessageBody == "" {
		http.Error(res, "Empty string content in either ToEmail, ToName, Subject or MessageBody", http.StatusBadRequest)
		return
	}
	err = ah.usecase.SendEmail(emailRequest.ToEmail, emailRequest.ToName, emailRequest.Subject, emailRequest.MessageBody)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	render.Status(req, http.StatusOK)
	render.PlainText(res, req, "Email sent successfully")
}

//TODO add MJML functions
