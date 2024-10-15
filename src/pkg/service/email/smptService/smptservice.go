package smptservice

import (
	"io"
	"log/slog"

	gomail "gopkg.in/mail.v2"

	domain "templify/pkg/domain/model"
)

// TODO - Fill config with the correct SMTP service details
type SMTPServiceConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
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

func (es *SMTPService) CreateEmailData(emailRequest *domain.EmailRequest) *gomail.Message {
	//Create Email Data
	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", es.config.FromEmail)
	message.SetHeader("To", emailRequest.ToEmail)
	message.SetHeader("Subject", emailRequest.Subject)

	// Set email body
	message.SetBody("text/plain", emailRequest.MessageBody)

	//Add attachments if present
	if emailRequest.AttachmentInfo != nil {
		for _, attachment := range emailRequest.AttachmentInfo {
			message.Attach(attachment.FileName+"."+attachment.FileExtension, gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(attachment.FileBytes)
				return err
			}))
		}
	}
	return message
}

func (es *SMTPService) SendEmail(emailRequest *domain.EmailRequest) {
	// Create email data
	message := es.CreateEmailData(emailRequest)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer(es.config.Host, es.config.Port, es.config.Username, es.config.Password)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		es.log.With("Error", err.Error()).Error("Error sending email")
	} else {
		es.log.Debug("Email sent successfully")
	}
}
