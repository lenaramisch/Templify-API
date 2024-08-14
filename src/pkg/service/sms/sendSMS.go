package smsservice

import (
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
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
}

func NewTwilioSMSSender(config *TwilioSMSSenderConfig) *TwilioSMSSender {
	return &TwilioSMSSender{
		config: config,
	}
}

func (s *TwilioSMSSender) SendSMS(toNumber string, messageBody string) error {
	slog.With(
		"toNumber", toNumber,
		"messageBody", messageBody,
	).Debug("Trying to send SMS")
	// Create SMS data
	SMSData := url.Values{
		"To":   {toNumber},
		"From": {s.config.FromNumber},
		"Body": {messageBody},
	}

	// Create a new POST request
	URL := fmt.Sprintf(TWILIO_BASE_URL+TWILIO_ACCOUNTS_URL+"%s/Messages.json", s.config.AccountSID)
	r, err := http.NewRequest("POST", URL, strings.NewReader(SMSData.Encode()))
	if err != nil {
		slog.Warn("Error creating request")
		return err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set the Authorization header with Basic Auth base64 encoded
	auth := base64.StdEncoding.EncodeToString([]byte(s.config.AccountSID + ":" + s.config.AuthToken))
	r.Header.Set("Authorization", "Basic "+auth)

	slog.With(
		"request", r,
	).Debug("Building request")
	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		slog.Warn("Error performing request")
		return err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Warn("Error reading response body")
		return err
	}

	// Check if the response is not between 200 and 299
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		slog.With(
			"statusCode", resp.StatusCode,
			"response", body,
		).Warn("Error response status code")
		return fmt.Errorf("Error response status code: %d", resp.StatusCode)
	}

	return nil
}
