package usecase

import (
	"context"
	_ "embed"
	"log/slog"
	domain "templify/pkg/domain/model"
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
	GetFileUploadURL(fileName string) (*domain.FileUploadResponse, error)
	GetFileDownloadURL(fileName string) (string, error)
	DownloadFile(fileDownloadRequest domain.FileDownloadRequest) ([]byte, error)
	UploadFile(fileUploafRequest domain.FileUploadRequest) error
	ListBuckets() ([]string, error)
	ListFiles(bucketName string) ([]string, error)
}

type Repository interface {
	// Email
	GetEmailTemplateByName(ctx context.Context, name string) (*domain.Template, error)
	AddEmailTemplate(ctx context.Context, template *domain.Template) error
	// PDF
	GetPDFTemplateByName(ctx context.Context, name string) (*domain.Template, error)
	AddPDFTemplate(ctx context.Context, template *domain.Template) error
	// SMS
	AddSMSTemplate(ctx context.Context, template *domain.Template) error
	GetSMSTemplateByName(ctx context.Context, name string) (*domain.Template, error)
	// Workflow
	AddWorkflow(ctx context.Context, workflow *domain.WorkflowCreateRequest) error
	GetWorkflowByName(ctx context.Context, workflowName string) (*domain.Workflow, error)
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
	filemanagerService FileManagerService
	log                *slog.Logger
}

func NewUsecase(emailSender EmailSender, smsSender SMSSender, mjmlService MJMLService, repository Repository, typstService TypstService, filemanagerService FileManagerService, log *slog.Logger) *Usecase {
	return &Usecase{
		emailSender:        emailSender,
		smsSender:          smsSender,
		mjmlService:        mjmlService,
		repository:         repository,
		typstService:       typstService,
		filemanagerService: filemanagerService,
		log:                log,
	}
}
