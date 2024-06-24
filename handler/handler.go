package handler

import (
	"encoding/json"
	"net/http"

	"example.SMSService.com/domain"
)

func SenderPostRequest(res http.ResponseWriter, req *http.Request) {
	var senderReq SenderRequest

	err := json.NewDecoder(req.Body).Decode(&senderReq)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if senderReq.ToNumber == "" || senderReq.MessageBody == "" {
		http.Error(res, "Empty string content in either ToNumber or MessageBody", http.StatusBadRequest)
	}
	err = domain.SendSMS(senderReq.ToNumber, senderReq.MessageBody)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Message sent successfully"))
}
