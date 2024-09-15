package db

import (
	"fmt"
	"log"

	domain "templify/pkg/domain/model"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	config       *RepositoryConfig
	dbConnection *sqlx.DB
}

type RepositoryConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewRepository(config *RepositoryConfig) *Repository {
	repo := &Repository{
		config: config,
	}
	repo.ConnectToDB()
	return repo
}

func (r *Repository) ConnectToDB() {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		r.config.Host, r.config.Port, r.config.User, r.config.Password, r.config.DBName,
	)
	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		log.Fatal("Connecting to DB failed", err)
	}

	r.dbConnection = db
}

func (r *Repository) GetEmailTemplateByName(name string) (*domain.Template, error) {
	tx := r.dbConnection.MustBegin()
	getTemplateByNameQuery := "SELECT * FROM emailtemplates WHERE name=$1"
	templateDB := Template{}
	err := tx.Get(&templateDB, getTemplateByNameQuery, name)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// map to domain model
	templateDomain := domain.Template{
		Name:        templateDB.Name,
		TemplateStr: templateDB.TemplString,
		IsMJML:      templateDB.IsMJML,
	}
	return &templateDomain, nil
}

func (r *Repository) AddWorkflow(workflow *domain.Workflow) error {

	var staticAttachments string
	for _, attachment := range workflow.StaticAttachments {
		staticAttachments = staticAttachments + attachment.FileName + ","
	}

	var templatedPDFs string
	for _, templatedPDF := range workflow.TemplatedPDFs {
		templatedPDFs = templatedPDFs + templatedPDF.TemplateName + ","
	}

	//map domain model to db model
	workflowDB := Workflow{
		Name:              workflow.Name,
		EmailTemplateName: workflow.EmailTemplateName,
		EmailSubject:      workflow.EmailSubject,
		StaticAttachments: staticAttachments,
		TemplatedPDFs:     templatedPDFs,
	}
	tx := r.dbConnection.MustBegin()
	addWorkflowQuery := "INSERT INTO workflows (name, email_template_name, email_subject, static_attachments, templated_pdfs) VALUES ($1, $2, $3, $4, $5)"
	tx.MustExec(addWorkflowQuery, workflowDB.Name, workflowDB.EmailTemplateName, workflowDB.EmailSubject, workflowDB.StaticAttachments, workflowDB.TemplatedPDFs)
	//tx.MustExec(addWorkflowQuery, workflowDB)
	return tx.Commit()
}

func (r *Repository) GetWorkflowByName(workflowName string) (*domain.Workflow, error) {
	tx := r.dbConnection.MustBegin()
	getWorkflowByNameQuery := "SELECT * FROM workflows WHERE name=$1"
	workflowDB := Workflow{}
	err := tx.Get(&workflowDB, getWorkflowByNameQuery, workflowName)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// map to domain model
	workflowDomain := domain.Workflow{
		Name:              workflowDB.Name,
		EmailTemplateName: workflowDB.EmailTemplateName,
		EmailSubject:      workflowDB.EmailSubject,
	}
	return &workflowDomain, nil
}

func (r *Repository) AddEmailTemplate(name string, templStr string, isMJML bool) error {
	tx := r.dbConnection.MustBegin()
	addTemplateQuery := "INSERT INTO emailtemplates (name, templ_string, is_mjml) VALUES ($1, $2, $3)"
	tx.MustExec(addTemplateQuery, name, templStr, isMJML)
	return tx.Commit()
}

func (r *Repository) AddSMSTemplate(name string, smsTemplString string) error {
	tx := r.dbConnection.MustBegin()
	addSMSTemplQuery := "INSERT INTO smstemplates (name, templ_string) VALUES ($1, $2)"
	tx.MustExec(addSMSTemplQuery, name, smsTemplString)
	return tx.Commit()
}

func (r *Repository) SavePDF(fileName string, base64Content string) error {
	tx := r.dbConnection.MustBegin()
	savePDFQuery := "INSERT INTO pdfs (name, content) VALUES ($1, $2)"
	tx.MustExec(savePDFQuery, fileName, base64Content)
	return tx.Commit()
}

func (r *Repository) AddPDFTemplate(name string, typstString string) error {
	tx := r.dbConnection.MustBegin()
	addPDFTemplateQuery := "INSERT INTO pdftemplates (name, typst_string) VALUES ($1, $2)"
	tx.MustExec(addPDFTemplateQuery, name, typstString)
	return tx.Commit()
}

func (r *Repository) GetPDFTemplateByName(name string) (*domain.Template, error) {
	tx := r.dbConnection.MustBegin()
	getPDFTemplateByNameQuery := "SELECT * FROM pdftemplates WHERE name=$1"
	templateDB := PDFTemplate{}
	err := tx.Get(&templateDB, getPDFTemplateByNameQuery, name)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// map to domain model
	templateDomain := domain.Template{
		Name:        templateDB.Name,
		TemplateStr: templateDB.TypstString,
	}
	return &templateDomain, nil
}

func (r *Repository) GetSMSTemplateByName(name string) (*domain.Template, error) {
	tx := r.dbConnection.MustBegin()
	getSMSTemplateByNameQuery := "SELECT * FROM smstemplates WHERE name=$1"
	templateDB := SMSTemplate{}
	err := tx.Get(&templateDB, getSMSTemplateByNameQuery, name)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// map to domain model
	templateDomain := domain.Template{
		Name:        templateDB.Name,
		TemplateStr: templateDB.TemplString,
	}
	return &templateDomain, nil
}
