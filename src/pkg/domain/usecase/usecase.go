package usecase

import (
	_ "embed"
	"log/slog"
	domain "templify/pkg/domain/model"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type EmailSender interface {
	SendEmail(domainEmailReq *domain.EmailRequest) error
}

type SMSSender interface {
	SendSMS(smsRequest domain.SmsRequest) error
}

type MJMLService interface {
	RenderMJML(MJMLString string) (string, error)
}

type FileManagerService interface {
	UploadFile(fileUploadRequest domain.FileUploadRequest) error
	DownloadFile(fileDownloadRequest domain.FileDownloadRequest) ([]byte, error)
	ListBuckets() ([]types.Bucket, error)
}

type Repository interface {
	// Email
	GetEmailTemplateByName(name string) (*domain.Template, error)
	AddEmailTemplate(name string, templateStr string, isMJML bool) error
	// PDF
	GetPDFTemplateByName(name string) (*domain.Template, error)
	AddPDFTemplate(name string, typstString string) error
	// SMS
	AddSMSTemplate(name string, smsTemplString string) error
	GetSMSTemplateByName(name string) (*domain.Template, error)
	// Workflow
	AddWorkflow(workflow *domain.WorkflowCreateRequest) error
	GetWorkflowByName(workflowName string) (*domain.Workflow, error)
}

type TypstService interface {
	RenderTypst(typstString string) ([]byte, error)
}

type Usecase struct {
	emailSender        EmailSender
	smsSender          SMSSender
	mjmlService        MJMLService
	repository         Repository
	typstService       TypstService
	fileManagerService FileManagerService
	log                *slog.Logger
}

func NewUsecase(emailSender EmailSender, smsSender SMSSender, mjmlService MJMLService, repository Repository, typstService TypstService, fileManagerService FileManagerService, log *slog.Logger) *Usecase {
	return &Usecase{
		emailSender:        emailSender,
		smsSender:          smsSender,
		mjmlService:        mjmlService,
		repository:         repository,
		typstService:       typstService,
		fileManagerService: fileManagerService,
		log:                log,
	}
}
