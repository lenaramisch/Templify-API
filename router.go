package main

import (
	"encoding/json"
	"net/http"

	"example.SMSService.com/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type SenderRequest struct {
	ToNumber    string `json:"receiverPhoneNumber"`
	MessageBody string `json:"message"`
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/sender", senderPostRequest)
	http.ListenAndServe(":8080", router)
}

func senderPostRequest(res http.ResponseWriter, req *http.Request) {
	var senderReq SenderRequest

	err := json.NewDecoder(req.Body).Decode(&senderReq)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	handlerErr := handler.SenderPostRequest(senderReq.ToNumber, senderReq.MessageBody)
	if handlerErr != nil {
		http.Error(res, handlerErr.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Message sent successfully"))
}
