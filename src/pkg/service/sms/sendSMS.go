package smsservice

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
	config TwilioSMSSenderConfig
}

func NewTwilioSMSSender(config TwilioSMSSenderConfig) *TwilioSMSSender {
	return &TwilioSMSSender{
		config: config,
	}
}

func (s *TwilioSMSSender) SendSMS(toNumber string, messageBody string) error {
	// Create SMS data
	SMSData := url.Values{
		"To":   {toNumber},
		"From": {s.config.FromNumber},
		"Body": {messageBody},
	}

	// Create a new POST request
	URL := fmt.Sprintf(TWILIO_BASE_URL+TWILIO_ACCOUNTS_URL+"%s/Messages.json", s.config.AccountSID)
	req, err := http.NewRequest("POST", URL, strings.NewReader(SMSData.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set the Authorization header with Basic Auth base64 encoded
	auth := base64.StdEncoding.EncodeToString([]byte(s.config.AccountSID + ":" + s.config.AuthToken))
	req.Header.Set("Authorization", "Basic "+auth)

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	return nil
}
