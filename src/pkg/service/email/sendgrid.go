package emailservice

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"templify/pkg/domain"

	"github.com/joho/godotenv"
)

type SendgridConfig struct {
	ApiKey       string
	FromEmail    string
	FromName     string
	ReplyToEmail string
	ReplyToName  string
}

type SendGridService struct {
	config SendgridConfig
}

func NewSendGridService(config SendgridConfig) *SendGridService {
	return &SendGridService{
		config: config,
	}
}

const (
	SENDGRID_URL = "https://api.sendgrid.com/v3/mail/send"
)

func (es *SendGridService) SendEmail(emailRequest *domain.EmailRequest) error {
	godotenv.Load()

	if es.config.ApiKey == "" || es.config.FromEmail == "" || es.config.FromName == "" || es.config.ReplyToEmail == "" || es.config.ReplyToName == "" {
		return errors.New("missing environment variables")
	}

	//Create Email Data
	emailData := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{
					{
						"email": emailRequest.ToEmail,
						"name":  emailRequest.ToName,
					},
				},
				"subject": emailRequest.Subject,
			},
		},
		"content": []map[string]string{
			{
				"type":  "text/html",
				"value": emailRequest.MessageBody,
			},
		},
		"from": map[string]string{
			"email": es.config.FromEmail,
			"name":  es.config.FromName,
		},
		"reply_to": map[string]string{
			"email": es.config.ReplyToEmail,
			"name":  es.config.ReplyToName,
		},
	}

	// add attachments if present
	if emailRequest.AttachmentInfo != nil {
		emailData["attachments"] = []map[string]any{
			{
				"content":     base64.StdEncoding.EncodeToString(emailRequest.AttachmentInfo.Content),
				"disposition": "attachment",
				"filename":    emailRequest.AttachmentInfo.FileName,
				"type":        emailRequest.AttachmentInfo.FileExtension,
			},
		}
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
	req.Header.Set("Authorization", "Bearer "+es.config.ApiKey)

	//Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Check HTTP status code
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Error response from SendGrid:", string(body))
		return fmt.Errorf(fmt.Sprintf("SendGrid returned status code %d", resp.StatusCode))
	}

	return nil
}
