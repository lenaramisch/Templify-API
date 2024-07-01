package domain

type EmailSender interface {
	SendEmail(toEmail string, toName string, subject string, message string) error
}

type SMSSender interface {
	SendSMS(toNumber string, messageBody string) error
}

type MJMLService interface {
	GetTemplatePlaceholders(MJMLTemplate string) (string, error)
	TemplatePostRequest(MJMLTemplate string) error
}

type Usecase struct {
	emailSender EmailSender
	smsSender   SMSSender
	mjmlService MJMLService
}

func NewUsecase(emailsender EmailSender, smsSender SMSSender, mjmlService MJMLService) *Usecase {
	return &Usecase{
		emailSender: emailsender,
		smsSender:   smsSender,
		mjmlService: mjmlService,
	}
}

// domain layer functions (usecases with actual business logic)
// Here we currently don't have any logic and forward to our services
// Later we may have sequential service calls or mapping, some logic etc.

func (u *Usecase) SendEmail(toEmail string, toName string, subject string, message string) error {
	return u.emailSender.SendEmail(toEmail, toName, subject, message)
}

func (u *Usecase) SendSMS(toNumber string, messageBody string) error {
	return u.smsSender.SendSMS(toNumber, messageBody)
}

func (u *Usecase) GetTemplatePlaceholders(MJMLTemplate string) (string, error) {
	return u.mjmlService.GetTemplatePlaceholders(MJMLTemplate)
}
