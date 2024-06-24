package handler

import (
	"errors"

	"example.SMSService.com/domain"
)

func SenderPostRequest(toNumber string, messageBody string) error {
	if toNumber == "" || messageBody == "" {
		return errors.New("toNumber or messageBody is empty")
	}
	domainError := domain.SendSMS(toNumber, messageBody)
	if domainError != nil {
		return domainError
	}
	return nil
}

//toNumber := "+4915170640522"
//messageBody := "Hello from Go Code!"
