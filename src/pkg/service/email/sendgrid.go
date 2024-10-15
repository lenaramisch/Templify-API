package emailservice

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	domain "templify/pkg/domain/model"
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
	log    *slog.Logger
}

func NewSendGridService(config *SendgridConfig, log *slog.Logger) *SendGridService {
	return &SendGridService{
		config: config,
		log:    log,
	}
}

const (
	SENDGRID_URL = "https://api.sendgrid.com/v3/mail/send"
)

func (es *SendGridService) CreateEmailData(emailRequest *domain.EmailRequest) map[string]any {
	//Create Email Data
	emailData := map[string]any{
		"personalizations": []map[string]any{
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
			base64AttachmentStr := base64.StdEncoding.EncodeToString(attachment.FileBytes)
			var attachmentType string
			// switch on file extension for type
			switch attachment.FileExtension {
			case "html":
				attachmentType = "text/html"
			case "txt":
				attachmentType = "text/plain"
			case "csv":
				attachmentType = "text/csv"
			case "pdf":
				attachmentType = "application/pdf"
			case "png":
				attachmentType = "image/png"
			case "jpg", "jpeg":
				attachmentType = "image/jpeg"
			default:
				attachmentType = "application/octet-stream"
			}
			attachmentData := map[string]string{
				"content":     base64AttachmentStr,
				"disposition": "attachment",
				"filename":    attachment.FileName + "." + attachment.FileExtension,
				"type":        attachmentType,
			}
			attachments = append(attachments, attachmentData)
		}

		// Add the attachments slice to the emailData
		emailData["attachments"] = attachments
	}

	return emailData
}

func (es *SendGridService) SendEmail(emailRequest *domain.EmailRequest) error {
	if es.config.ApiKey == "" || es.config.FromEmail == "" || es.config.FromName == "" || es.config.ReplyToEmail == "" || es.config.ReplyToName == "" {
		return errors.New("missing environment variables")
	}

	// Create Email Data
	emailData := es.CreateEmailData(emailRequest)

	jsonBytes, err := json.Marshal(emailData)
	if err != nil {
		return err
	}

	// log the json bytes as string
	es.log.With("emailData", string(jsonBytes)).Debug("Request Data")

	// Create a new POST request
	r, err := http.NewRequest("POST", SENDGRID_URL, bytes.NewBuffer(jsonBytes))
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
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Error response from SendGrid:", string(body))
		es.log.With("statusCode", resp.StatusCode, "errorResponse", body).Debug("Error response from SendGrid")
		return fmt.Errorf("error response from SendGrid")
	}

	return nil
}
