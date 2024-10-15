package smptservice

import (
	"encoding/base64"
	"log/slog"

	domain "templify/pkg/domain/model"
)

// TODO - Fill config with the correct SMTP service details
type SMTPServiceConfig struct {
	ApiKey       string
	FromEmail    string
	FromName     string
	ReplyToEmail string
	ReplyToName  string
}

type SMTPService struct {
	config *SMTPServiceConfig
	log    *slog.Logger
}

func NewSMTPService(config *SMTPServiceConfig, log *slog.Logger) *SMTPService {
	return &SMTPService{
		config: config,
		log:    log,
	}
}

// TODO - Create email data as needed for the SMTP service
func (es *SMTPService) CreateEmailData(emailRequest *domain.EmailRequest) map[string]any {
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

// TODO - Implement function to send email using the SMTP service
func (es *SMTPService) SendEmail(emailRequest *domain.EmailRequest) error {
	return nil
}
