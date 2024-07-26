package db

import (
	"fmt"
	"log"

	"example.SMSService.com/pkg/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	config       RepositoryConfig
	dbConnection *sqlx.DB
}

type RepositoryConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewRepository(config RepositoryConfig) *Repository {
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
		TemplateStr: templateDB.MJMLString,
	}
	return &templateDomain, nil
}

func (r *Repository) AddEmailTemplate(name string, mjmlString string) error {
	tx := r.dbConnection.MustBegin()
	addTemplateQuery := "INSERT INTO emailtemplates (name, mjml_string) VALUES ($1, $2)"
	tx.MustExec(addTemplateQuery, name, mjmlString)
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
