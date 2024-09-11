package emailservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	domain "templify/pkg/domain/model"

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
	config *SendgridConfig
}

func NewSendGridService(config *SendgridConfig) *SendGridService {
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
		// Initialize the attachments slice
		var attachments []map[string]string

		// Range over the attachments and add them to the attachments slice
		for _, attachment := range emailRequest.AttachmentInfo {
			attachmentData := map[string]string{
				"content":     attachment.Base64Content,
				"disposition": "attachment",
				"filename":    attachment.FileName,
				"type":        attachment.FileExtension,
			}
			attachments = append(attachments, attachmentData)
		}

		// Add the attachments slice to the emailData
		emailData["attachments"] = attachments
	}

	jsonData, err := json.Marshal(emailData)
	if err != nil {
		return err
	}

	// Create a new POST request
	r, err := http.NewRequest("POST", SENDGRID_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+es.config.ApiKey)

	//Perform the request
	client := &http.Client{}
	resp, err := client.Do(r)
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
