package db

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"

	domain "templify/pkg/domain/model"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	gensqlc "templify/pkg/db/gen-sqlc/generated"
)

type Repository struct {
	config       *RepositoryConfig
	dbConnection *pgxpool.Pool
	log          *slog.Logger
	queries      *gensqlc.Queries
}

type RepositoryConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewRepository(config *RepositoryConfig, log *slog.Logger) *Repository {
	repo := &Repository{
		config: config,
		log:    log,
	}
	repo.ConnectToDB()
	repo.queries = gensqlc.New(repo.dbConnection)
	return repo
}

func (r *Repository) ConnectToDB() error {
	// Connection string in the pgx format
	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		r.config.User, r.config.Password, r.config.Host, r.config.Port, r.config.DBName,
	)

	ctx := context.Background()

	// Parse the connection string and create a connection pool
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Connect to the database using pgxpool
	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Storing the connection pool in the repository struct
	r.dbConnection = dbpool

	if err := r.dbConnection.Ping(ctx); err != nil {
		r.log.With("error", err).Error("Failed to ping DB")
		return err
	}

	// You may want to log successful connection
	log.Println("Connected to DB successfully")
	return nil
}

func (r *Repository) GetEmailTemplateByName(ctx context.Context, name string) (*domain.Template, error) {
	pgx_name := pgtype.Text{
		String: name,
		Valid:  true,
	}
	templateDB, err := r.queries.GetEmailTemplateByName(ctx, pgx_name)
	if err != nil {
		r.log.With("error", err).Error("Failed to get email template by name")
		return nil, err
	}

	// map to domain model
	templateDomain := domain.Template{
		Name:        templateDB.Name.String,
		TemplateStr: templateDB.TemplString.String,
		IsMJML:      templateDB.IsMjml.Bool,
	}
	return &templateDomain, nil
}

func (r *Repository) AddWorkflow(ctx context.Context, workflow *domain.WorkflowCreateRequest) error {
	staticAttachmentNamesStr := strings.Join(workflow.StaticAttachments, ",")
	templatedPDFNamesStr := strings.Join(workflow.TemplatedPDFs, ",")

	//map domain model to db model
	workflowDB := Workflow{
		Name:              workflow.Name,
		EmailTemplateName: workflow.EmailTemplateName,
		EmailSubject:      workflow.EmailSubject,
		StaticAttachments: staticAttachmentNamesStr,
		TemplatedPDFs:     templatedPDFNamesStr,
	}
	r.log.With("workflowDB", workflowDB).Debug("WorkflowDB")
	_, err := r.queries.AddWorkflow(ctx, gensqlc.AddWorkflowParams{
		Name:              pgtype.Text{String: workflowDB.Name, Valid: true},
		EmailTemplateName: pgtype.Text{String: workflowDB.EmailTemplateName, Valid: true},
		EmailSubject:      pgtype.Text{String: workflowDB.EmailSubject, Valid: true},
		StaticAttachments: pgtype.Text{String: workflowDB.StaticAttachments, Valid: true},
		TemplatedPdfs:     pgtype.Text{String: workflowDB.TemplatedPDFs, Valid: true},
	})
	if err != nil {
		r.log.With("error", err).Error("Failed to add workflow")
		return err
	}
	return nil
}

func (r *Repository) GetWorkflowByName(ctx context.Context, workflowName string) (*domain.Workflow, error) {
	pgx_name := pgtype.Text{
		String: workflowName,
		Valid:  true,
	}
	workflowDB, err := r.queries.GetWorkflowByName(ctx, pgx_name)
	if err != nil {
		r.log.With("error", err).Error("Failed to get workflow by name")
		return nil, err
	}

	// map to domain model
	workflowDomain := domain.Workflow{
		Name:              workflowDB.Name.String,
		EmailTemplateName: workflowDB.EmailTemplateName.String,
		PDFTemplateNames:  workflowDB.TemplatedPdfs.String,
		StaticAttachments: workflowDB.StaticAttachments.String,
		EmailSubject:      workflowDB.EmailSubject.String,
	}
	return &workflowDomain, nil
}

func (r *Repository) AddEmailTemplate(ctx context.Context, template *domain.Template) error {
	// map domain model to db model
	templateDB := Template{
		Name:        template.Name,
		TemplString: template.TemplateStr,
		IsMJML:      template.IsMJML,
	}
	_, err := r.queries.AddEmailTemplate(ctx, gensqlc.AddEmailTemplateParams{
		Name:        pgtype.Text{String: templateDB.Name, Valid: true},
		TemplString: pgtype.Text{String: templateDB.TemplString, Valid: true},
		IsMjml:      pgtype.Bool{Bool: templateDB.IsMJML, Valid: true},
	})
	if err != nil {
		r.log.With("error", err, "templateName", template.Name).Error("Failed to add email template")
		return err
	}
	return nil
}

func (r *Repository) AddSMSTemplate(ctx context.Context, template *domain.Template) error {
	// map domain model to db model
	templateDB := Template{
		Name:        template.Name,
		TemplString: template.TemplateStr,
	}
	_, err := r.queries.AddSMSTemplate(ctx, gensqlc.AddSMSTemplateParams{
		Name:        pgtype.Text{String: templateDB.Name, Valid: true},
		TemplString: pgtype.Text{String: templateDB.TemplString, Valid: true},
	})
	if err != nil {
		r.log.With("error", err, "templateName", template.Name).Error("Failed to add sms template")
		return err
	}
	return nil
}

func (r *Repository) AddPDFTemplate(ctx context.Context, template *domain.Template) error {
	// map domain model to db model
	templateDB := Template{
		Name:        template.Name,
		TemplString: template.TemplateStr,
	}
	_, err := r.queries.AddPDFTemplate(ctx, gensqlc.AddPDFTemplateParams{
		Name:        pgtype.Text{String: templateDB.Name, Valid: true},
		TemplString: pgtype.Text{String: templateDB.TemplString, Valid: true},
	})
	if err != nil {
		r.log.With("error", err, "templateName", template.Name).Error("Failed to add pdf template")
		return err
	}
	return nil
}

func (r *Repository) GetPDFTemplateByName(ctx context.Context, name string) (*domain.Template, error) {
	pgx_name := pgtype.Text{
		String: name,
		Valid:  true,
	}
	templateDB, err := r.queries.GetPDFTemplateByName(ctx, pgx_name)
	if err != nil {
		r.log.With("error", err, "templateName", name).Error("Failed to get PDF template by name")
		return nil, err
	}

	// map to domain model
	templateDomain := domain.Template{
		Name:        templateDB.Name.String,
		TemplateStr: templateDB.TemplString.String,
	}
	return &templateDomain, nil
}

func (r *Repository) GetSMSTemplateByName(ctx context.Context, name string) (*domain.Template, error) {
	pgx_name := pgtype.Text{
		String: name,
		Valid:  true,
	}
	templateDB, err := r.queries.GetSMSTemplateByName(ctx, pgx_name)
	if err != nil {
		r.log.With("error", err, "templateName", name).Error("Failed to get SMS template by name")
		return nil, err
	}

	// map to domain model
	templateDomain := domain.Template{
		Name:        templateDB.Name.String,
		TemplateStr: templateDB.TemplString.String,
	}
	return &templateDomain, nil
}
