package domain

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	SENDGRID_URL = "https://api.sendgrid.com/v3/mail/send"
)

func SendEmail(toEmail string, toName string, subject string, message string) error {
	godotenv.Load()
	apiKey := os.Getenv("API_KEY")
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")
	replyToEmail := os.Getenv("REPLY_TO_EMAIL")
	replyToName := os.Getenv("REPLY_TO_NAME")

	if apiKey == "" || fromEmail == "" || fromName == "" || replyToEmail == "" || replyToName == "" {
		return errors.New("missing environment variables")
	}

	//Create Email Data
	emailData := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{
					{
						"email": toEmail,
						"name":  toName,
					},
				},
				"subject": subject,
			},
		},
		"content": []map[string]string{
			{
				"type":  "text/plain",
				"value": message,
			},
		},
		"from": map[string]string{
			"email": fromEmail,
			"name":  fromName,
		},
		"reply_to": map[string]string{
			"email": replyToEmail,
			"name":  replyToName,
		},
	}

	jsonData, err := json.Marshal(emailData)
	if err != nil {
		return err
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", SENDGRID_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	//Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	return nil
}
