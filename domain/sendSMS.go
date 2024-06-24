package domain

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func SendSMS(toNumber string, messageBody string) error {
	godotenv.Load()
	accountSID := os.Getenv("ACCOUNT_SID")
	URL := "https://api.twilio.com/2010-04-01/Accounts/" + accountSID + "/Messages.json"
	authToken := os.Getenv("AUTH_TOKEN")

	fromNumber := "+14042366595"

	// Create SMS data
	SMSData := url.Values{
		"To":   {toNumber},
		"From": {fromNumber},
		"Body": {messageBody},
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", URL, strings.NewReader(SMSData.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set the Authorization header with Basic Auth
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
