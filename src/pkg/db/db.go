package db

import (
	"fmt"
	"log"

	"example.SMSService.com/pkg/domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	return &Repository{
		config: config,
	}
}

func (r *Repository) ConnectToDB() {
	connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=disable", r.config.User, r.config.DBName)
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatal("Connecting to DB failed", err)
	}

	r.dbConnection = db
}

func (r *Repository) GetTemplateByName(name string) (*domain.Template, error) {
	db := r.dbConnection.MustBegin()
	getTemplateByNameQuery := "SELECT * FROM templates WHERE name=$1"
	templateDB := Template{}
	err := db.Get(&templateDB, getTemplateByNameQuery, name)
	if err != nil {
		return nil, err
	}
	templateDomain := domain.Template{
		Name:       templateDB.Name,
		MJMLString: templateDB.MJMLString,
	}
	return &templateDomain, nil
}

func (r *Repository) AddTemplate(name string, mjmlString string) (int64, error) {
	db := r.dbConnection.MustBegin()
	addTemplateQuery := "INSERT INTO templates (name, mjml_string) VALUES ($1, $2)"
	sqlResult := db.MustExec(addTemplateQuery, name, mjmlString)
	return sqlResult.LastInsertId()
}
