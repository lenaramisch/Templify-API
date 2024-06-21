package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func SendSMS() {
	//TODO accountSID from env file
	URL := "https://api.twilio.com/2010-04-01/Accounts/" + accountSID + "/Messages.json"
	//TODO authToken from env file

	// Define phone numbers and message body
	toNumber := "+4915170640522"
	fromNumber := "+14042366595"
	messageBody := "Hello from Go Code!"

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
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	fmt.Println(result)
}
