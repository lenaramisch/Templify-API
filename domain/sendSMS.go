package domain

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	DEFAULT_FROM_NUMBER = "+14042366595"
	TWILIO_BASE_URL     = "https://api.twilio.com/"

	TWILIO_ACCOUNTS_URL = "2010-04-01/Accounts/"
)

func SendSMS(toNumber string, messageBody string) error {
	godotenv.Load()
	accountSID := os.Getenv("ACCOUNT_SID")
	authToken := os.Getenv("AUTH_TOKEN")

	// Create SMS data
	SMSData := url.Values{
		"To":   {toNumber},
		"From": {DEFAULT_FROM_NUMBER},
		"Body": {messageBody},
	}

	// Create a new POST request
	URL := fmt.Sprintf(TWILIO_BASE_URL+TWILIO_ACCOUNTS_URL+"%s/Messages.json", accountSID)
	req, err := http.NewRequest("POST", URL, strings.NewReader(SMSData.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set the Authorization header with Basic Auth base64 encoded
	auth := base64.StdEncoding.EncodeToString([]byte(accountSID + ":" + authToken))
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
