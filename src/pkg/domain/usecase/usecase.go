package usecase

import (
	_ "embed"
	domain "templify/pkg/domain/model"
)

type EmailSender interface {
	SendEmail(domainEmailReq *domain.EmailRequest) error
}

type SMSSender interface {
	SendSMS(toNumber string, messageBody string) error
}

type MJMLService interface {
	RenderMJML(MJMLString string) (string, error)
}

type Repository interface {
	// Email
	GetEmailTemplateByName(name string) (*domain.Template, error)
	AddEmailTemplate(name string, templateStr string, isMJML bool) error
	// PDF
	SavePDF(fileName string, base64Content string) error
	GetPDFTemplateByName(name string) (*domain.Template, error)
	AddPDFTemplate(name string, typstString string) error
	// SMS
	AddSMSTemplate(name string, smsTemplString string) error
	GetSMSTemplateByName(name string) (*domain.Template, error)
	// Workflow
	AddWorkflow(workflow *domain.WorkflowCreateRequest) error
}

type TypstService interface {
	RenderTypst(typstString string) ([]byte, error)
}

type Usecase struct {
	emailSender  EmailSender
	smsSender    SMSSender
	mjmlService  MJMLService
	repository   Repository
	typstService TypstService
}

func NewUsecase(emailSender EmailSender, smsSender SMSSender, mjmlService MJMLService, repository Repository, typstService TypstService) *Usecase {
	return &Usecase{
		emailSender:  emailSender,
		smsSender:    smsSender,
		mjmlService:  mjmlService,
		repository:   repository,
		typstService: typstService,
	}
}
