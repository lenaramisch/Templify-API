package smsservice

import (
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	domain "templify/pkg/domain/model"
)

const (
	TWILIO_BASE_URL = "https://api.twilio.com/"

	TWILIO_ACCOUNTS_URL = "2010-04-01/Accounts/"
)

type TwilioSMSSenderConfig struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

type TwilioSMSSender struct {
	config *TwilioSMSSenderConfig
	log    *slog.Logger
}

func NewTwilioSMSSender(config *TwilioSMSSenderConfig, log *slog.Logger) *TwilioSMSSender {
	return &TwilioSMSSender{
		config: config,
		log:    log,
	}
}

func (s *TwilioSMSSender) SendSMS(smsRequest domain.SmsRequest) error {
	s.log.With(
		"toNumber", smsRequest.ToNumber,
		"messageBody", smsRequest.MessageBody,
	).Debug("Trying to send SMS")
	// Create SMS data
	SMSData := url.Values{}
	SMSData.Set("To", smsRequest.ToNumber)
	SMSData.Set("From", s.config.FromNumber)
	SMSData.Set("Body", smsRequest.MessageBody)

	s.log.With("SMSData", SMSData).Debug("SMS data")

	// Create a new POST request
	URL := fmt.Sprintf(TWILIO_BASE_URL+TWILIO_ACCOUNTS_URL+"%s/Messages.json", s.config.AccountSID)
	r, err := http.NewRequest("POST", URL, strings.NewReader(SMSData.Encode()))
	if err != nil {
		s.log.Warn("Error creating request")
		return domain.ErrorCreatingSMSRequest{Reason: err.Error()}
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set the Authorization header with Basic Auth base64 encoded
	auth := base64.StdEncoding.EncodeToString([]byte(s.config.AccountSID + ":" + s.config.AuthToken))
	r.Header.Set("Authorization", "Basic "+auth)

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		s.log.Warn("Error performing request")
		return domain.ErrorPerformingSMSRequest{Reason: err.Error()}
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.log.Warn("Error reading response body")
		return err
	}

	// Check if the response is not between 200 and 299
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		s.log.With(
			"statusCode", resp.StatusCode,
			"response", body,
		).Warn("Error response status code")
		return domain.ErrorSendingSMS{StatusCode: resp.StatusCode}
	}
	return nil
}
