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
	fmt.Println("Using connection string: ", connectionString)
	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		log.Fatal("Connecting to DB failed", err)
	}

	r.dbConnection = db
}

func (r *Repository) GetTemplateByName(name string) (*domain.Template, error) {
	tx := r.dbConnection.MustBegin()
	getTemplateByNameQuery := "SELECT * FROM templates WHERE name=$1"
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
		Name:       templateDB.Name,
		MJMLString: templateDB.MJMLString,
	}
	return &templateDomain, nil
}

func (r *Repository) AddTemplate(name string, mjmlString string) error {
	tx := r.dbConnection.MustBegin()
	addTemplateQuery := "INSERT INTO templates (name, mjml_string) VALUES ($1, $2)"
	tx.MustExec(addTemplateQuery, name, mjmlString)
	return tx.Commit()

}
