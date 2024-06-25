package handler

import (
	"encoding/json"
	"net/http"

	"example.SMSService.com/pkg/domain"
)

func SMSPostRequest(res http.ResponseWriter, req *http.Request) {
	var smsRequest SmsRequest

	err := json.NewDecoder(req.Body).Decode(&smsRequest)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if smsRequest.ToNumber == "" || smsRequest.MessageBody == "" {
		http.Error(res, "Empty string content in either ToNumber or MessageBody", http.StatusBadRequest)
	}
	err = domain.SendSMS(smsRequest.ToNumber, smsRequest.MessageBody)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("SMS sent successfully"))
}

func EmailPostRequest(res http.ResponseWriter, req *http.Request) {
	var emailRequest EmailRequest

	err := json.NewDecoder(req.Body).Decode(&emailRequest)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if emailRequest.ToEmail == "" || emailRequest.ToName == "" || emailRequest.Subject == "" || emailRequest.MessageBody == "" {
		http.Error(res, "Empty string content in either ToEmail, ToName, Subject or MessageBody", http.StatusBadRequest)
	}
	err = domain.SendEmail(emailRequest.ToEmail, emailRequest.ToName, emailRequest.Subject, emailRequest.MessageBody)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Email sent successfully"))
}
